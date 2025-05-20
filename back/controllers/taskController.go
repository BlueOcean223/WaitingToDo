package controllers

import (
	"back/models"
	"back/models/vo"
	"back/service"
	"back/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskController struct {
	taskService *service.TaskService
}

func NewTaskController(taskService *service.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

// GetTaskList 分页查询
func (s *TaskController) GetTaskList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("pageSize"))

	if err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	// 防止传入较大参数
	if page > 10000 || pageSize > 1000 {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数错误", nil))
		return
	}

	// 获取当前用户的邮箱
	email, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", "令牌无效", nil))
		return
	}

	taskList, err := s.taskService.GetTaskList(email, page, pageSize)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, models.Success("", "查询成功", taskList))
}

// AddTask 新增任务
func (s *TaskController) AddTask(c *gin.Context) {
	var taskVo vo.TaskVo
	if err := c.ShouldBind(&taskVo); err != nil {
		c.JSON(http.StatusBadRequest, models.Fail("", "参数异常", nil))
		return
	}
	// 获取用户邮箱
	email, err := utils.GetUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Fail("", "令牌无效", nil))
		return
	}
	// 添加任务
	err = s.taskService.AddTask(email, taskVo)
	if err != nil {
		if utils.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, models.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, models.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, models.Success("", "添加成功", nil))
}
