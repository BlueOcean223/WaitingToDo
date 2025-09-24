package ticker

import (
	"backend/internal/configs"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/services/logAnalyzer"
	"backend/pkg/logger"
	"fmt"
	"github.com/robfig/cron/v3"
	"gopkg.in/gomail.v2"
	"sync"
)

// cronScheduler 全局cron调度器
var cronScheduler *cron.Cron

// StartAllTimers 启动所有定时任务
func StartAllTimers() {
	// 创建cron调度器，支持秒级精度
	cronScheduler = cron.New(cron.WithSeconds())

	// 启动日志分析任务 - 每天凌晨1点执行
	startLogAnalysisTask()

	// 启动任务通知任务 - 每小时执行一次
	startTaskNotifyTask()

	// 启动cron调度器
	cronScheduler.Start()
	logger.Info("所有定时任务已启动")
}

// StopAllTimers 停止所有定时任务
func StopAllTimers() {
	if cronScheduler != nil {
		cronScheduler.Stop()
		logger.Info("所有定时任务已停止")
	}
}

// startLogAnalysisTask 启动日志分析定时任务
func startLogAnalysisTask() {
	// 创建日志分析器
	analyzer := logAnalyzer.NewLogAnalyzer()

	if !analyzer.IsEnabled() {
		logger.Info("日志分析功能未启用")
		return
	}

	// 使用cron表达式：每天凌晨1点执行
	// 格式：秒 分 时 日 月 周
	// "0 0 1 * * *" 表示每天凌晨1点0分0秒执行
	_, err := cronScheduler.AddFunc("0 0 1 * * *", func() {
		logger.Info("开始执行日志分析任务")

		// 执行日志分析
		if err := analyzer.AnalyzeLogs(); err != nil {
			logger.Error("日志分析失败", logger.Err(err))
			return
		}

		logger.Info("日志分析完成")
	})

	if err != nil {
		logger.Error("添加日志分析定时任务失败", logger.Err(err))
		return
	}

	logger.Info("日志分析定时任务已添加：每天凌晨1点执行")
}

// startTaskNotifyTask 启动任务通知定时任务
func startTaskNotifyTask() {
	// 使用cron表达式：每小时执行一次
	// "0 0 * * * *" 表示每小时的0分0秒执行
	_, err := cronScheduler.AddFunc("0 0 * * * *", func() {
		executeTaskNotify()
	})

	if err != nil {
		logger.Error("添加任务通知定时任务失败", logger.Err(err))
		return
	}

	logger.Info("任务通知定时任务已添加：每小时执行一次")
}

// executeTaskNotify 执行任务通知逻辑
func executeTaskNotify() {
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

	// 1. 获取即将过期的任务列表
	tasks, err := taskRepo.GetOneDayDDLTaskList()
	if err != nil {
		logger.Error("获取任务列表失败", logger.Err(err))
		return
	}

	if len(tasks) == 0 {
		return
	}

	// 2. 批量获取需要通知的任务(过滤已通知的)
	taskIds := make([]int, len(tasks))
	for i, task := range tasks {
		taskIds[i] = task.Id
	}

	notifiedTasks, err := taskNoticeHistoryRepo.GetHistoriesByTaskIds(taskIds)
	if err != nil {
		logger.Error("获取通知历史失败", logger.Err(err))
		return
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
			return
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
		return
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
