package services

import (
	"backend/internal/configs"
	"backend/internal/models"
	"backend/internal/models/dto"
	"backend/internal/models/vo"
	"backend/internal/repositories"
	"backend/pkg/captcha"
	"backend/pkg/logger"
	"backend/pkg/minioContent"
	"backend/pkg/myError"
	"backend/pkg/redisContent"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"strconv"
	"time"
)

type TaskService interface {
	GetTaskList(email string, page, pageSize int, status *int) ([]dto.TaskDto, error)
	AddTask(email string, taskVo vo.TaskVo) (models.Task, error)
	UpdateTask(taskVo vo.TaskVo) error
	DeleteTask(id int) error
	GetUrgentTaskList(email string) ([]dto.TaskDto, error)
	GetTeamTaskList(userId, page, pageSize int) ([]dto.TeamTaskDto, error)
	DeleteTeamTask(taskId, userId int) error
	AddTeamTask(task models.Task) error
	GetInviteCodeByTaskId(taskId int) (string, error)
	DeleteFromMinio(bucketName string, taskId int) error
	JoinTeamTaskByInviteCode(email, inviteCode string) error
	GenerateInviteCode(taskId int)
	CompleteTeamTask(teamTask models.TeamTask) error
	InviteTeamMember(email string, teamTask models.TeamTask) error
}

type taskService struct {
	authRepository     repository.AuthRepository
	messageRepository  repository.MessageRepository
	taskRepository     repository.TaskRepository
	teamTaskRepository repository.TeamTaskRepository
	fileRepository     repository.FileRepository
	inviteCodeRepo     repository.InviteCodeRepository
}

func NewTaskService(authRepository repository.AuthRepository,
	messageRepository repository.MessageRepository,
	taskRepository repository.TaskRepository,
	teamTaskRepository repository.TeamTaskRepository,
	fileRepository repository.FileRepository,
	inviteCodeRepo repository.InviteCodeRepository) TaskService {
	return &taskService{
		authRepository:     authRepository,
		messageRepository:  messageRepository,
		taskRepository:     taskRepository,
		teamTaskRepository: teamTaskRepository,
		fileRepository:     fileRepository,
		inviteCodeRepo:     inviteCodeRepo,
	}
}

// GetTaskList 获取任务列表
func (s *taskService) GetTaskList(email string, page, pageSize int, status *int) ([]dto.TaskDto, error) {
	// 根据邮箱查询用户
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		logger.Error("查询用户失败",
			logger.String("email", email),
			logger.Err(err))
		return nil, err
	}
	if user == (models.User{}) {
		logger.Warn("用户不存在",
			logger.String("email", email))
		return nil, errors.New("用户不存在")
	}

	taskList, count, err := s.taskRepository.GetList(user.Id, page, pageSize, 0, status)
	if err != nil {
		logger.Error("查询任务列表失败",
			logger.String("email", email),
			logger.String("user_id", fmt.Sprintf("%d", user.Id)),
			logger.Err(err))
		return nil, err
	}

	// 查询任务的附件
	var taskIds []int
	for _, task := range taskList {
		taskIds = append(taskIds, task.Id)
	}
	attachments, err := s.fileRepository.GetFileByTaskIds(taskIds)
	if err != nil {
		logger.Error("查询任务附件失败",
			logger.String("email", email),
			logger.String("user_id", fmt.Sprintf("%d", user.Id)),
			logger.Err(err))
		return nil, err
	}

	// 收集任务与附件关系
	taskAttachments := make(map[int][]models.File)
	for _, attachment := range attachments {
		taskAttachments[attachment.TaskId] = append(taskAttachments[attachment.TaskId], attachment)
	}

	// 封装Dto列表
	var taskDtoList []dto.TaskDto
	for _, task := range taskList {
		taskDtoList = append(taskDtoList, dto.TaskDto{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Ddl:         task.Ddl,
			Status:      task.Status,
			Count:       count,
			Attachments: taskAttachments[task.Id],
		})
	}

	return taskDtoList, nil
}

// AddTask 添加任务
func (s *taskService) AddTask(email string, taskVo vo.TaskVo) (models.Task, error) {
	// 根据邮箱查询用户ID
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		logger.Error("查询用户失败",
			logger.String("email", email),
			logger.Err(err))
		return models.Task{}, err
	}
	if user == (models.User{}) {
		logger.Warn("用户不存在",
			logger.String("email", email))
		return models.Task{}, errors.New("用户不存在")
	}

	task := models.Task{
		UserId:      user.Id,
		Title:       taskVo.Title,
		Description: taskVo.Description,
		Ddl:         taskVo.Ddl,
		Type:        taskVo.Type,
		Status:      0,
	}

	// 插入数据库
	err = s.taskRepository.Create(&task, nil)
	if err != nil {
		logger.Error("创建任务失败",
			logger.String("email", email),
			logger.String("title", taskVo.Title),
			logger.Err(err))
		return task, err
	}

	return task, nil
}

// UpdateTask 更新任务
func (s *taskService) UpdateTask(taskVo vo.TaskVo) error {
	// 开启事务
	tx := configs.MysqlDb.Begin()

	// 更新任务
	err := s.taskRepository.Update(models.Task{
		Id:          taskVo.Id,
		Title:       taskVo.Title,
		Description: taskVo.Description,
		Ddl:         taskVo.Ddl,
		Type:        taskVo.Type,
		Status:      taskVo.Status,
	}, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("更新任务失败",
			logger.String("task_id", fmt.Sprintf("%d", taskVo.Id)),
			logger.Err(err))
		return err
	}

	// 更新任务通知表
	layout := "2006-01-02T15:04:05.000Z"
	ddl, err := time.Parse(layout, taskVo.Ddl)
	if err != nil {
		tx.Rollback()
		logger.Error("格式化ddl失败",
			logger.Int("task_id", taskVo.Id),
			logger.Err(err))
		return err
	}
	now := time.Now()

	// 如果ddl比当前时间大于等于一天，则删除通知历史,需要重新通知
	if ddl.Sub(now).Hours() >= 24 {
		// 查询该任务是否有通知过
		taskNoticeHistoryRepo := repository.NewTaskNoticeHistoryRepository(configs.MysqlDb)
		history, err := taskNoticeHistoryRepo.GetHistoryByTaskId(taskVo.Id)
		if err != nil {
			tx.Rollback()
			logger.Error("查询任务通知记录失败",
				logger.Int("task_id", taskVo.Id),
				logger.Err(err))
			return err
		}
		// 如果有通知历史，则删除
		if history != (models.TaskNoticeHistory{}) {
			err = taskNoticeHistoryRepo.DeleteHistoryByTaskId(taskVo.Id, tx)
			if err != nil {
				tx.Rollback()
				logger.Error("删除任务通知记录失败",
					logger.Int("task_id", taskVo.Id),
					logger.Err(err))
				return err
			}
		}
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("提交更新任务事务失败",
			logger.String("task_id", fmt.Sprintf("%d", taskVo.Id)),
			logger.Err(err))
		return err
	}

	return nil
}

// DeleteTask 删除任务
func (s *taskService) DeleteTask(id int) error {
	// 删除minio中的附件
	err := s.DeleteFromMinio(minioContent.FilesBucket, id)
	if err != nil {
		logger.Error("删除Minio附件失败",
			logger.String("task_id", fmt.Sprintf("%d", id)),
			logger.Err(err))
		return err
	}

	// 开启事务
	tx := configs.MysqlDb.Begin()
	// 删除任务附件
	err = s.fileRepository.DeleteByTaskId(id, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("删除任务附件失败",
			logger.String("task_id", fmt.Sprintf("%d", id)),
			logger.Err(err))
		return err
	}
	// 删除任务
	err = s.taskRepository.Delete(id, tx)
	if err != nil {
		tx.Rollback()
		logger.Error("删除任务失败",
			logger.String("task_id", fmt.Sprintf("%d", id)),
			logger.Err(err))
		return err
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		logger.Error("提交删除任务事务失败",
			logger.String("task_id", fmt.Sprintf("%d", id)),
			logger.Err(err))
		return err
	}

	return nil
}

// DeleteFromMinio 删除Minio中的附件
func (s *taskService) DeleteFromMinio(bucketName string, taskId int) error {
	minioClient := configs.MinioClient
	ctx := context.Background()

	taskFiles, err := s.fileRepository.GetFileByTaskId(taskId)
	if err != nil {
		return err
	}

	var objectNames []string
	for _, taskFile := range taskFiles {
		objectNames = append(objectNames, strconv.Itoa(taskId)+"/"+taskFile.Name)
	}

	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, objectName := range objectNames {
			objectsCh <- minio.ObjectInfo{Key: objectName}
		}
	}()

	var errs []error
	for err := range minioClient.RemoveObjects(ctx, bucketName, objectsCh, minio.RemoveObjectsOptions{}) {
		errs = append(errs, err.Err)
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// GetUrgentTaskList 获取紧急任务列表
func (s *taskService) GetUrgentTaskList(email string) ([]dto.TaskDto, error) {
	// 根据邮箱查询用户
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		logger.Error("查询用户失败",
			logger.String("email", email),
			logger.Err(err))
		return nil, err
	}
	if user == (models.User{}) {
		logger.Warn("用户不存在",
			logger.String("email", email))
		return nil, errors.New("用户不存在")
	}

	userId := user.Id

	// 查询紧急任务列表
	taskList, err := s.taskRepository.GetUrgentList(userId)
	if err != nil {
		logger.Error("查询紧急任务列表失败",
			logger.String("email", email),
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.Err(err))
		return nil, err
	}
	// 封装Dto列表
	var taskDtoList []dto.TaskDto
	for _, task := range taskList {
		taskDtoList = append(taskDtoList, dto.TaskDto{
			Id:          task.Id,
			Title:       task.Title,
			Description: task.Description,
			Ddl:         task.Ddl,
			Status:      task.Status,
		})
	}
	return taskDtoList, nil
}

// GetTeamTaskList 获取小组任务列表
func (s *taskService) GetTeamTaskList(userId, page, pageSize int) ([]dto.TeamTaskDto, error) {
	// 分页查询小组任务表中，用户的小组任务id
	teamTaskList, err := s.teamTaskRepository.GetList(page, pageSize, userId)
	if err != nil {
		logger.Error("查询小组任务列表失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.Err(err))
		return nil, err
	}

	taskIds := make([]int, len(teamTaskList))
	for i, teamTask := range teamTaskList {
		taskIds[i] = teamTask.TaskId
	}

	// 查询任务列表
	taskList, err := s.taskRepository.GetTaskListByIds(taskIds)
	if err != nil {
		logger.Error("查询任务列表失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.Err(err))
		return nil, err
	}

	// 根据任务id列表获取所有任务关系
	teamTaskShipList, err := s.teamTaskRepository.GetTeamTaskShipByTaskIds(taskIds)
	if err != nil {
		logger.Error("查询任务关系失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.Err(err))
		return nil, err
	}
	// 收集用户id，以及任务与用户关联关系
	var userIds []int
	taskUserShipMap := make(map[int][]models.TeamTask)
	for _, teamTaskShip := range teamTaskShipList {
		tId := teamTaskShip.TaskId
		uId := teamTaskShip.UserId
		userIds = append(userIds, uId)
		taskUserShipMap[tId] = append(taskUserShipMap[tId], teamTaskShip)
	}

	// 根据用户id获取用户信息
	userInfoList, err := s.authRepository.SelectUsersByIds(userIds)
	if err != nil {
		logger.Error("查询用户信息失败",
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.Err(err))
		return nil, err
	}
	// 将用户信息根据id映射成map
	userInfoMap := make(map[int]models.User)
	for _, user := range userInfoList {
		userInfoMap[user.Id] = user
	}

	// 收集同一小组的成员信息
	teamUsersMap := make(map[int][]dto.TeamUserDto)
	for _, taskId := range taskIds {
		for _, taskUserShip := range taskUserShipMap[taskId] {
			uId := taskUserShip.UserId
			userInfo := userInfoMap[uId]
			// 封装对应任务的用户信息
			teamUserDto := dto.TeamUserDto{
				Id:     userInfo.Id,
				Name:   userInfo.Name,
				Pic:    userInfo.Pic,
				Status: taskUserShip.Status,
			}
			teamUsersMap[taskId] = append(teamUsersMap[taskId], teamUserDto)
		}
	}

	// 封装teamTaskDto列表
	var teamTaskDtoList []dto.TeamTaskDto
	for _, task := range taskList {
		// 获取同小组的成员信息
		users := teamUsersMap[task.Id]

		// 封装dto
		teamTaskDtoList = append(teamTaskDtoList, dto.TeamTaskDto{
			Id:          task.Id,
			UserId:      userId,
			Title:       task.Title,
			Description: task.Description,
			Ddl:         task.Ddl,
			Status:      task.Status,
			Users:       users,
		})
	}

	return teamTaskDtoList, nil
}

// DeleteTeamTask 删除小组任务
func (s *taskService) DeleteTeamTask(taskId, userId int) error {
	// 开启事务
	tx := configs.MysqlDb.Begin()
	// 删除任务关系表数据
	err := tx.Where("task_id = ? AND user_id = ?", taskId, userId).Delete(&models.TeamTask{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 查询小组任务成员情况
	teamTaskList, err := s.teamTaskRepository.GetTeamTaskShipByTaskIds([]int{taskId})
	if err != nil {
		tx.Rollback()
		return err
	}
	// 如果小组任务组长为自己且还有其它组员，则转交给其它小组成员
	tasks, err := s.taskRepository.GetTaskListByIds([]int{taskId})
	task := tasks[0]
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(teamTaskList) > 1 && task.UserId == userId {
		for _, teamTask := range teamTaskList {
			if teamTask.UserId != userId {
				task.UserId = teamTask.UserId
				break
			}
		}
		err = s.taskRepository.Update(task, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// 如果小组成员全部删除，则删除任务表数据
	if len(teamTaskList) == 1 {
		err = tx.Delete(&models.Task{}, taskId).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 删除小组的邀请码
		err = s.inviteCodeRepo.DeleteByTaskId(taskId, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("提交删除小组任务事务失败",
			logger.String("task_id", fmt.Sprintf("%d", taskId)),
			logger.String("user_id", fmt.Sprintf("%d", userId)),
			logger.Err(err))
		return err
	}

	return nil
}

// AddTeamTask 添加小组任务
func (s *taskService) AddTeamTask(task models.Task) error {
	// 开启事务
	tx := configs.MysqlDb.Begin()
	// 向任务表写数据
	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		logger.Error("创建任务失败",
			logger.String("title", task.Title),
			logger.String("user_id", fmt.Sprintf("%d", task.UserId)),
			logger.Err(err))
		return err
	}
	// 向任务关系表写数据
	teamTask := models.TeamTask{
		TaskId: task.Id,
		UserId: task.UserId,
		Status: 0,
	}
	if err := tx.Create(&teamTask).Error; err != nil {
		tx.Rollback()
		logger.Error("创建任务关系失败",
			logger.String("task_id", fmt.Sprintf("%d", task.Id)),
			logger.String("user_id", fmt.Sprintf("%d", task.UserId)),
			logger.Err(err))
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("提交添加小组任务事务失败",
			logger.String("task_id", fmt.Sprintf("%d", task.Id)),
			logger.String("user_id", fmt.Sprintf("%d", task.UserId)),
			logger.Err(err))
		return err
	}

	// 生成邀请码，使用goroutine完成
	go s.GenerateInviteCode(task.Id)

	return nil
}

// GenerateInviteCode 生成邀请码
func (s *taskService) GenerateInviteCode(taskId int) {
	// 邀请码
	code := captcha.GenerateInviteCode()

retry:
	for {
		// 检查邀请码是否已存在
		inviteCode, err := s.inviteCodeRepo.GetByInviteCode(code)
		if err != nil {
			log.Println("检查验证码失败:", err)
			continue
		}

		// 验证码不存在，则结束循环
		if inviteCode == (models.InviteCode{}) {
			break
		} else {
			// 验证码存在，重新生成
			code = captcha.GenerateInviteCode()
		}
	}

	// 插入数据库
	inviteCode := models.NewInviteCode(taskId, code)
	err := s.inviteCodeRepo.Insert(inviteCode, nil)
	if err != nil {
		log.Println("插入邀请码失败:", err)

		// 休眠，然后重新尝试
		time.Sleep(30 * time.Second)
		goto retry
	}
}

// GetInviteCodeByTaskId 根据任务ID获取邀请码
func (s *taskService) GetInviteCodeByTaskId(taskId int) (string, error) {
	// 查询邀请码
	inviteCode, err := s.inviteCodeRepo.GetByTaskId(taskId)
	if err != nil {
		return "", err
	}

	// 如果没有邀请码
	if inviteCode == (models.InviteCode{}) {
		return "暂无邀请码", nil
	}

	return inviteCode.InviteCode, nil
}

// JoinTeamTaskByInviteCode 根据邀请码加入小组任务
func (s *taskService) JoinTeamTaskByInviteCode(email, inviteCode string) error {
	// 根据邮箱查询用户
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		logger.Error("查询用户失败",
			logger.String("email", email),
			logger.String("invite_code", inviteCode),
			logger.Err(err))
		return err
	}

	// 根据邀请码查询小组任务ID
	inviteCodeRecord, err := s.inviteCodeRepo.GetByInviteCode(inviteCode)
	if err != nil {
		return err
	}
	// 如果没有查询到邀请码
	if inviteCodeRecord == (models.InviteCode{}) {
		return myError.NewMyError("邀请码不存在")
	}
	// 查询小组任务是否存在
	task, err := s.taskRepository.GetTaskListByIds([]int{inviteCodeRecord.TaskId})
	if err != nil {
		logger.Error("查询小组任务失败",
			logger.String("email", email),
			logger.String("invite_code", inviteCode),
			logger.Int("taskId", inviteCodeRecord.TaskId),
			logger.Err(err))
		return err
	}
	if len(task) == 0 {
		logger.Warn("小组任务不存在",
			logger.String("email", email),
			logger.String("invite_code", inviteCode),
			logger.Int("taskId", inviteCodeRecord.TaskId))
		return myError.NewMyError("小组任务不存在")
	}

	// 检查该用户是否已经加入了该小组任务
	teamTaskMembers, err := s.teamTaskRepository.GetTeamMembers(inviteCodeRecord.TaskId)
	if err != nil {
		logger.Error("查询小组成员失败",
			logger.String("email", email),
			logger.String("task_id", fmt.Sprintf("%d", inviteCodeRecord.TaskId)),
			logger.Err(err))
		return err
	}
	for _, member := range teamTaskMembers {
		if member.Id == user.Id {
			logger.Warn("用户已加入小组",
				logger.String("email", email),
				logger.String("task_id", fmt.Sprintf("%d", inviteCodeRecord.TaskId)))
			return myError.NewMyError("您已经加入了该小组任务")
		}
	}

	// 插入小组任务关系表
	teamTask := models.TeamTask{
		TaskId: inviteCodeRecord.TaskId,
		UserId: user.Id,
		Status: 0,
	}
	err = s.teamTaskRepository.Insert(teamTask, nil)
	if err != nil {
		logger.Error("插入小组任务关系失败",
			logger.String("email", email),
			logger.String("task_id", fmt.Sprintf("%d", inviteCodeRecord.TaskId)),
			logger.Err(err))
		return err
	}

	return nil
}

// CompleteTeamTask 小组成员完成小组任务
func (s *taskService) CompleteTeamTask(teamTask models.TeamTask) error {
	// 开启事务
	tx := configs.MysqlDb.Begin()
	// 更新任务关系表
	err := tx.Where("task_id = ? AND user_id = ?", teamTask.TaskId, teamTask.UserId).Updates(teamTask).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 如果小组成员已经全部完成任务，则更新任务表
	teamTaskList, err := s.teamTaskRepository.GetTeamTaskShipByTaskIds([]int{teamTask.TaskId})
	if err != nil {
		tx.Rollback()
		return err
	}
	count := 0
	for _, taskShip := range teamTaskList {
		count += taskShip.Status
	}

	if count == len(teamTaskList)-1 {
		err = tx.Where("id = ?", teamTask.TaskId).Updates(models.Task{Status: 1}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("提交完成小组任务事务失败",
			logger.String("task_id", fmt.Sprintf("%d", teamTask.TaskId)),
			logger.String("user_id", fmt.Sprintf("%d", teamTask.UserId)),
			logger.Err(err))
		return err
	}
	return nil
}

// InviteTeamMember 邀请成员,向被邀请的好友发送信息
func (s *taskService) InviteTeamMember(email string, teamTask models.TeamTask) error {
	// 获取当前用户信息
	redisClient := configs.RedisClient
	userInfoKey := redisContent.UserInfoKey + email
	var user models.User
	var err error

	// 先查询redis
	if redisClient.Exists(context.Background(), userInfoKey).Val() == 1 {
		val, err := redisClient.Get(context.Background(), userInfoKey).Bytes()
		if err != nil {
			logger.Error("从Redis获取用户信息失败",
				logger.String("email", email),
				logger.Err(err))
			return err
		}
		err = json.Unmarshal(val, &user)
		if err != nil {
			logger.Error("反序列化用户信息失败",
				logger.String("email", email),
				logger.Err(err))
			return err
		}
	} else {
		// Redis中不存在，查询数据库
		user, err = s.authRepository.SelectUserByEmail(email)
		if err != nil {
			logger.Error("查询用户信息失败",
				logger.String("email", email),
				logger.Err(err))
			return err
		}

		// 将用户信息写入redis
		userJson, err := json.Marshal(user)
		if err != nil {
			return err
		}
		redisClient.Set(context.Background(), userInfoKey, userJson, 24*time.Hour)
	}

	// 查询redis，判断24h内是否已经邀请过
	key := fmt.Sprintf(redisContent.InviteTeamMemberKey+"%d:%d:%d", user.Id, teamTask.UserId, teamTask.TaskId)
	if redisClient.Exists(context.Background(), key).Val() == 1 {
		return myError.NewMyError("已邀请过该用户，请等待用户同意")
	}

	// 查询小组任务信息
	tasks, err := s.taskRepository.GetTaskListByIds([]int{teamTask.TaskId})
	if err != nil {
		return err
	}

	// 填写发送信息
	message := models.Message{
		Title:       "小组任务邀请",
		Description: fmt.Sprintf("好友 %%s 邀请你加入名为 %s 的小组任务", tasks[0].Title),
		FromId:      user.Id,
		ToId:        teamTask.UserId,
		Type:        2,
		SendTime:    time.Now().Format("2006-01-02 15:04:05"),
		OutId:       teamTask.TaskId,
		IsRead:      0,
	}
	// 写入信息表
	err = s.messageRepository.InsertMessage(message, nil)
	if err != nil {
		logger.Error("写入邀请消息失败",
			logger.String("email", email),
			logger.String("task_id", fmt.Sprintf("%d", teamTask.TaskId)),
			logger.String("invited_user_id", fmt.Sprintf("%d", teamTask.UserId)),
			logger.Err(err))
		return err
	}

	// 写入redis，记录24h内邀请信息
	redisClient.Set(context.Background(), key, 1, time.Hour*24)
	return nil
}
