package service

import (
	"back/configs"
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/repository"
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"sync"
	"time"
)

type TaskService struct {
	authRepository     *repository.AuthRepository
	taskRepository     *repository.TaskRepository
	teamTaskRepository *repository.TeamTaskRepository
}

func NewTaskService(authRepository *repository.AuthRepository,
	taskRepository *repository.TaskRepository, teamTaskRepository *repository.TeamTaskRepository) *TaskService {
	return &TaskService{
		authRepository:     authRepository,
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
	taskUserShipMap := make(map[int][]int)
	for _, teamTaskShip := range teamTaskShipList {
		tId := teamTaskShip.TaskId
		uId := teamTaskShip.UserId
		userIds = append(userIds, uId)
		taskUserShipMap[tId] = append(taskUserShipMap[tId], uId)
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
	teamUsersMap := make(map[int][]models.User)
	for _, taskId := range taskIds {
		for _, uId := range taskUserShipMap[taskId] {
			teamUsersMap[taskId] = append(teamUsersMap[taskId], userInfoMap[uId])
		}
	}

	// 封装teamTaskDto列表
	var teamTaskDtoList []dto.TeamTaskDto
	for _, task := range taskList {
		// 获取同小组的成员信息
		users := teamUsersMap[task.Id]
		// 封装userDto
		var userDtoList []dto.UserDto
		for _, user := range users {
			userDtoList = append(userDtoList, dto.UserDto{
				Id:          user.Id,
				Name:        user.Name,
				Email:       user.Email,
				Description: user.Description,
				Pic:         user.Pic,
			})
		}

		// 封装dto
		teamTaskDtoList = append(teamTaskDtoList, dto.TeamTaskDto{
			Id:          task.Id,
			UserId:      userId,
			Title:       task.Title,
			Description: task.Description,
			Ddl:         task.Ddl,
			Status:      task.Status,
			Users:       userDtoList,
		})
	}

	return teamTaskDtoList, nil
}
