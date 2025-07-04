package routers

import (
	"back/configs"
	"back/controllers"
	"back/repository"
	"back/service"
	"github.com/gin-gonic/gin"
)

func SetTaskRoutes(r *gin.Engine) {
	// 初始化依赖
	authRepository := repository.NewAuthRepository(configs.MysqlDb)
	messageRepository := repository.NewMessageRepository(configs.MysqlDb)
	taskRepository := repository.NewTaskRepository(configs.MysqlDb)
	teamTaskRepo := repository.NewTeamTaskRepository(configs.MysqlDb)
	fileRepository := repository.NewFileRepository(configs.MysqlDb)

	taskService := service.NewTaskService(authRepository, messageRepository, taskRepository, teamTaskRepo, fileRepository)

	taskController := controllers.NewTaskController(taskService)

	task := r.Group("/task")
	{
		// 查询任务列表
		task.GET("/list", taskController.GetTaskList)
		// 新增任务
		task.POST("/add", taskController.AddTask)
		// 删除任务
		task.DELETE("/delete", taskController.DeleteTask)
		// 修改任务
		task.PUT("/update", taskController.UpdateTask)
		// 获取紧急任务列表
		task.GET("/urgent", taskController.GetUrgentTaskList)
		// 查询小组任务列表
		task.GET("/teamList", taskController.GetTeamTaskList)
		// 删除小组任务
		task.DELETE("/team/delete", taskController.DeleteTeamTask)
		// 添加小组任务
		task.POST("/team/add", taskController.AddTeamTask)
		// 小组成员完成任务
		task.PUT("/team/complete", taskController.CompleteTeamTask)
		// 邀请成员
		task.POST("/team/invite", taskController.InviteTeamMember)
	}
}
