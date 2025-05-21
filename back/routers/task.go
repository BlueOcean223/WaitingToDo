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
	taskRepository := repository.NewTaskRepository(configs.MysqlDb)
	taskService := service.NewTaskService(authRepository, taskRepository)
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
	}
}
