package ticker

import (
	"backend/internal/configs"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/pkg/logger"
	"fmt"
	"gopkg.in/gomail.v2"
	"sync"
	"time"
)

// TaskNotify 定时任务，当用户的任务结束时间小于一天时，发邮件提醒用户
func TaskNotify() {
	db := configs.MysqlDb
	// 初始化时创建所有repository
	taskNoticeHistoryRepo := repository.NewTaskNoticeHistoryRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	teamTaskRepo := repository.NewTeamTaskRepository(db)
	authRepo := repository.NewAuthRepository(db)

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
			logger.Error("获取任务列表失败",
				logger.Err(err))
			continue
		}

		if len(tasks) == 0 {
			continue
		}

		// 2. 批量获取需要通知的任务(过滤已通知的)
		taskIds := make([]int, len(tasks))
		for i, task := range tasks {
			taskIds[i] = task.Id
		}

		notifiedTasks, err := taskNoticeHistoryRepo.GetHistoriesByTaskIds(taskIds)
		if err != nil {
			logger.Error("获取通知历史失败",
				logger.Err(err))
			continue
		}

		notifiedMap := make(map[int]bool)
		for _, history := range notifiedTasks {
			notifiedMap[history.TaskId] = true
		}

		// 3. 批量获取用户信息
		userIds := make([]int, 0, len(tasks))
		tasksToNotify := make([]models.Task, 0, len(tasks))
		var teamTaskIds []int
		leaderIds := make(map[int]int)
		teamTaskMap := make(map[int]models.Task)

		for _, task := range tasks {
			if !notifiedMap[task.Id] {
				userIds = append(userIds, task.UserId)
				tasksToNotify = append(tasksToNotify, task)
				// 小组任务
				if task.Type == 1 {
					teamTaskIds = append(teamTaskIds, task.Id)
					leaderIds[task.UserId] = 1
					teamTaskMap[task.Id] = task
				}
			}
		}

		// 查询小组任务成员
		var teamMembers []models.TeamTask
		if len(teamTaskIds) > 0 {
			teamMembers, err = teamTaskRepo.GetTeamTaskShipByTaskIds(teamTaskIds)
			if err != nil {
				logger.Error("获取小组任务成员失败", logger.Err(err))
				continue
			}
		}

		// 添加小组成员的id到用户id列表中，并向要通知的任务列表中添加给小组成员的通知
		for _, teamTask := range teamMembers {
			userIds = append(userIds, teamTask.UserId)
			// 如果不是小组组长，则添加通知（小组组长的通知已经包含在了通知任务列表中）
			if _, exist := leaderIds[teamTask.UserId]; !exist {
				newTask := teamTaskMap[teamTask.TaskId]
				newTask.UserId = teamTask.UserId
				tasksToNotify = append(tasksToNotify, newTask)
			}
		}

		users, err := authRepo.SelectUsersByIds(userIds)
		if err != nil {
			logger.Error("获取用户信息失败", logger.Err(err))
			continue
		}

		userMap := make(map[int]models.User)
		for _, user := range users {
			userMap[user.Id] = user
		}

		// 记录已经通知的任务
		var notifiedTaskIds []int

		// 4. 并发发送邮件并记录历史
		var wg sync.WaitGroup
		var mu sync.Mutex
		for _, task := range tasksToNotify {
			user, ok := userMap[task.UserId]
			if !ok {
				logger.Warn("用户不存在", logger.Int("id", user.Id))
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
					logger.Error("发送邮件失败", logger.Err(err))
					return
				}

				mu.Lock()
				// 记录已经通知的任务
				notifiedTaskIds = append(notifiedTaskIds, task.Id)
				mu.Unlock()
			}(task, user)
		}
		// 等待所有任务完成
		wg.Wait()

		// 批量插入已经通知的任务记录，避免 N+1 问题
		if len(notifiedTaskIds) > 0 {
			histories := make([]models.TaskNoticeHistory, len(notifiedTaskIds))
			for i, taskId := range notifiedTaskIds {
				histories[i] = models.TaskNoticeHistory{TaskId: taskId}
			}

			if err := taskNoticeHistoryRepo.BatchInsert(histories, nil); err != nil {
				logger.Error("记录通知历史失败", logger.Err(err))
			}
		}
	}
}
