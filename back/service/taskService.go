package service

import (
	"back/configs"
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/repository"
	"back/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"sync"
	"time"
)

type TaskService struct {
	authRepository     *repository.AuthRepository
	messageRepository  *repository.MessageRepository
	taskRepository     *repository.TaskRepository
	teamTaskRepository *repository.TeamTaskRepository
}

func NewTaskService(authRepository *repository.AuthRepository,
	messageRepository *repository.MessageRepository,
	taskRepository *repository.TaskRepository, teamTaskRepository *repository.TeamTaskRepository) *TaskService {
	return &TaskService{
		authRepository:     authRepository,
		messageRepository:  messageRepository,
		taskRepository:     taskRepository,
		teamTaskRepository: teamTaskRepository,
	}
}

// GetTaskList 获取任务列表
func (s *TaskService) GetTaskList(email string, page, pageSize int, status *int) ([]dto.TaskDto, error) {
	// 根据邮箱查询用户
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == (models.User{}) {
		return nil, errors.New("用户不存在")
	}

	taskList, count, err := s.taskRepository.GetList(user.Id, page, pageSize, 0, status)
	if err != nil {
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
			Count:       count,
		})
	}

	return taskDtoList, nil
}

// AddTask 添加任务
func (s *TaskService) AddTask(email string, taskVo vo.TaskVo) error {
	// 根据邮箱查询用户ID
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		return err
	}
	if user == (models.User{}) {
		return errors.New("用户不存在")
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
	return s.taskRepository.Create(task)
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(taskVo vo.TaskVo) error {
	return s.taskRepository.Update(models.Task{
		Id:          taskVo.Id,
		Title:       taskVo.Title,
		Description: taskVo.Description,
		Ddl:         taskVo.Ddl,
		Type:        taskVo.Type,
		Status:      taskVo.Status,
	})
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(id int) error {
	return s.taskRepository.Delete(id)
}

// GetUrgentTaskList 获取紧急任务列表
func (s *TaskService) GetUrgentTaskList(email string) ([]dto.TaskDto, error) {
	// 根据邮箱查询用户
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == (models.User{}) {
		return nil, errors.New("用户不存在")
	}

	userId := user.Id

	// 查询紧急任务列表
	taskList, err := s.taskRepository.GetUrgentList(userId)
	if err != nil {
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

// TickerNotify 定时任务，当用户的任务结束时间小于一天时，发邮件提醒用户
func TickerNotify() {
	// 初始化时创建所有repository
	taskNoticeHistoryRepo := repository.NewTaskNoticeHistoryRepository(configs.MysqlDb)
	taskRepo := repository.NewTaskRepository(configs.MysqlDb)
	authRepo := repository.NewAuthRepository(configs.MysqlDb)

	// 初始化邮件dialer
	mailConfig := configs.AppConfigs.MailConfig
	d := gomail.NewDialer(
		mailConfig.SMTPHost,
		mailConfig.SMTPPort,
		mailConfig.From,
		mailConfig.Password,
	)

	// 测试使用一分钟
	//ticker := time.NewTicker(1 * time.Minute)
	// 使用一小时的定时器
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		// 1. 获取即将过期的任务列表
		tasks, err := taskRepo.GetOneDayDDLTaskList()
		if err != nil {
			log.Printf("获取任务列表失败: %v", err)
			continue
		}

		if len(tasks) == 0 {
			continue
		}

		// 2. 批量获取需要通知的任务(过滤已通知的)
		taskIds := make([]int, 0, len(tasks))
		for _, task := range tasks {
			taskIds = append(taskIds, task.Id)
		}

		notifiedTasks, err := taskNoticeHistoryRepo.GetHistoriesByTaskIds(taskIds)
		if err != nil {
			log.Printf("获取通知历史失败: %v", err)
			continue
		}

		notifiedMap := make(map[int]bool)
		for _, history := range notifiedTasks {
			notifiedMap[history.TaskId] = true
		}

		// 3. 批量获取用户信息
		userIds := make([]int, 0, len(tasks))
		tasksToNotify := make([]models.Task, 0, len(tasks))
		for _, task := range tasks {
			if !notifiedMap[task.Id] {
				userIds = append(userIds, task.UserId)
				tasksToNotify = append(tasksToNotify, task)
			}
		}

		users, err := authRepo.SelectUsersByIds(userIds)
		if err != nil {
			log.Printf("获取用户信息失败: %v", err)
			continue
		}

		userMap := make(map[int]models.User)
		for _, user := range users {
			userMap[user.Id] = user
		}

		// 4. 并发发送邮件并记录历史
		var wg sync.WaitGroup
		for _, task := range tasksToNotify {
			user, ok := userMap[task.UserId]
			if !ok {
				log.Printf("用户%d不存在", task.UserId)
				continue
			}

			wg.Add(1)
			// 使用协程发送邮件
			go func(task models.Task, user models.User) {
				defer wg.Done()

				mail := models.Mail{
					To:      []string{user.Email},
					Subject: "您有一个即将到达ddl的任务",
					Body: fmt.Sprintf(
						`您标题为 <strong>%s</strong> 的任务即将到达ddl，请尽快完成！`, task.Title),
				}

				m := gomail.NewMessage()
				m.SetHeader("From", mailConfig.From)
				m.SetHeader("To", mail.To...)
				m.SetHeader("Subject", mail.Subject)
				m.SetBody("text/html", mail.Body)

				if err := d.DialAndSend(m); err != nil {
					log.Printf("发送邮件失败: %v", err)
					return
				}

				// 记录通知历史
				if err := taskNoticeHistoryRepo.Insert(models.TaskNoticeHistory{
					TaskId: task.Id,
				}); err != nil {
					log.Printf("记录通知历史失败: %v", err)
				}
			}(task, user)
		}
		// 等待所有任务完成
		wg.Wait()
	}
}

// GetTeamTaskList 获取小组任务列表
func (s *TaskService) GetTeamTaskList(userId, page, pageSize int) ([]dto.TeamTaskDto, error) {
	// 分页查询小组任务表中，用户的小组任务id
	teamTaskList, err := s.teamTaskRepository.GetList(page, pageSize, userId)
	if err != nil {
		return nil, err
	}
	var taskIds []int
	for _, teamTask := range teamTaskList {
		taskIds = append(taskIds, teamTask.TaskId)
	}

	// 查询任务列表
	taskList, err := s.taskRepository.GetTaskListByIds(taskIds)
	if err != nil {
		return nil, err
	}

	// 根据任务id列表获取所有任务关系
	teamTaskShipList, err := s.teamTaskRepository.GetTeamTaskShipByTaskIds(taskIds)
	if err != nil {
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
func (s *TaskService) DeleteTeamTask(taskId, userId int) error {
	// 开启事务
	tx := s.teamTaskRepository.Db.Begin()
	// 删除任务关系表数据
	err := tx.Where("task_id = ? AND user_id = ?", taskId, userId).Delete(&models.TeamTask{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 如果小组成员全部删除，则删除任务表数据
	teamTaskList, err := s.teamTaskRepository.GetTeamTaskShipByTaskIds([]int{taskId})
	if err != nil {
		tx.Rollback()
		return err
	}
	if len(teamTaskList) == 1 {
		err = tx.Delete(&models.Task{}, taskId).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// AddTeamTask 添加小组任务
func (s *TaskService) AddTeamTask(task models.Task) error {
	// 开启事务
	tx := s.teamTaskRepository.Db.Begin()
	// 向任务表写数据
	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
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
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// CompleteTeamTask 小组成员完成小组任务
func (s *TaskService) CompleteTeamTask(teamTask models.TeamTask) error {
	// 开启事务
	tx := s.teamTaskRepository.Db.Begin()
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
		return err
	}

	return nil
}

// InviteTeamMember 邀请成员,向被邀请的好友发送信息
func (s *TaskService) InviteTeamMember(email string, teamTask models.TeamTask) error {
	// 获取当前用户信息
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		return err
	}

	// 查询redis，判断24h内是否已经邀请过
	redisClient := configs.RedisClient
	key := fmt.Sprintf(utils.InviteTeamMemberKey+"%d:%d:%d", user.Id, teamTask.UserId, teamTask.TaskId)
	if redisClient.Exists(context.Background(), key).Val() == 1 {
		return utils.NewMyError("已邀请过该用户，请等待用户同意")
	}

	// 查询小组任务信息
	tasks, err := s.taskRepository.GetTaskListByIds([]int{teamTask.TaskId})
	if err != nil {
		return err
	}

	// 填写发送信息
	message := models.Message{
		Title:       "小组任务邀请",
		Description: fmt.Sprintf("好友 %s 邀请你加入名为 %s 的小组任务", user.Name, tasks[0].Title),
		FromId:      user.Id,
		ToId:        teamTask.UserId,
		Type:        2,
		SendTime:    time.Now().Format("2006-01-02 15:04:05"),
		OutId:       teamTask.TaskId,
		IsRead:      0,
	}
	// 写入信息表
	err = s.messageRepository.InsertMessage(message)
	if err != nil {
		return err
	}

	// 写入redis，记录24h内邀请信息
	redisClient.Set(context.Background(), key, 1, time.Hour*24)
	return nil
}

// StartTeamConsumer 监听消息队列，处理同意加入小组任务
func StartTeamConsumer() {
	for {
		MQConn := configs.RabbitMQConn
		channel, err := MQConn.Channel()
		if err != nil {
			log.Printf("打开channel失败: %v, 重新尝试...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		defer channel.Close()

		// 声明队列
		queue, err := channel.QueueDeclare(
			configs.AppConfigs.RabbitMQConfig.Queues["team_request"].Name,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("声明队列失败：%v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 绑定队列
		err = channel.QueueBind(
			queue.Name,
			configs.AppConfigs.RabbitMQConfig.Queues["team_request"].RoutingKey,
			"social",
			false,
			nil,
		)
		if err != nil {
			log.Printf("绑定队列失败：%v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 消费信息
		msgs, err := channel.Consume(
			queue.Name,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Printf("消费信息失败：%v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for msg := range msgs {
			var mqMessage models.MQMessage
			if err = json.Unmarshal(msg.Body, &mqMessage); err != nil {
				log.Printf("反序列化消息失败：%v", err)
				msg.Nack(false, true) // 重新入队
				continue
			}

			// 开启事务
			tx := configs.MysqlDb.Begin()

			// 向任务关系表插入数据
			teamTask := models.TeamTask{
				TaskId: mqMessage.RelationID,
				UserId: mqMessage.ReceiverID,
				Status: 0,
			}
			if err = tx.Create(&teamTask).Error; err != nil {
				tx.Rollback()
				log.Printf("向任务关系表插入数据异常：%v", err)
				msg.Nack(false, true) // 处理失败，重新入队
				continue
			}

			// 将任务表相应任务状态改为未完成
			if err = tx.Model(&models.Task{}).Where("id = ?", mqMessage.RelationID).
				Update("status", 0).Error; err != nil {
				tx.Rollback()
				log.Printf("更新任务状态异常：%v", err)
				msg.Nack(false, true) // 处理失败，重新入队
				continue
			}

			// 提交事务
			if err = tx.Commit().Error; err != nil {
				tx.Rollback()
				log.Printf("提交事务异常：%v", err)
				msg.Nack(false, true) // 处理失败，重新入队
				continue
			}

			// 消息消费完成
			if err = msg.Ack(false); err != nil {
				log.Printf("消息消费确认失败：%v", err)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
