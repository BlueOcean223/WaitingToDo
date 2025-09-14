package handlers

import (
	"backend/internal/models"
	"backend/internal/models/vo"
	"backend/internal/services"
	"backend/pkg/myError"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskService services.TaskService
}

func NewTaskHandler(taskService services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// GetTaskList 分页查询
func (s *TaskHandler) GetTaskList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	statusStr := c.Query("status")

	// 前端是否需要筛选状态
	var status *int
	if statusStr != "" {
		temp, err := strconv.Atoi(statusStr)
		if err != nil {
			c.JSON(http.StatusOK, response.Fail("", "状态参数异常", nil))
			return
		}
		status = &temp
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	// 防止传入较大参数
	if page > 10000 || pageSize > 1000 {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	// 获取当前用户的邮箱
	email := c.GetString("user")

	taskList, err := s.taskService.GetTaskList(email, page, pageSize, status)
	if err != nil {
		if myError.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", taskList))
}

// AddTask 新增任务
func (s *TaskHandler) AddTask(c *gin.Context) {
	var taskVo vo.TaskVo
	if err := c.ShouldBind(&taskVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数异常", nil))
		return
	}
	// 获取用户邮箱
	email := c.GetString("user")
	// 添加任务
	task, err := s.taskService.AddTask(email, taskVo)
	if err != nil {
		if myError.IsMyError(err) {
			// 自定义错误
			c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		} else {
			// 系统错误
			c.JSON(http.StatusInternalServerError, response.Fail("", err.Error(), nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("", "添加成功", task))
}

// UpdateTask 修改任务
func (s *TaskHandler) UpdateTask(c *gin.Context) {
	var taskVo vo.TaskVo
	if err := c.ShouldBind(&taskVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数异常", nil))
		return
	}

	err := s.taskService.UpdateTask(taskVo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "修改成功", nil))
}

// DeleteTask 删除任务
func (s *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err = s.taskService.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "删除成功", nil))
}

// GetUrgentTaskList 获取紧急任务列表
func (s *TaskHandler) GetUrgentTaskList(c *gin.Context) {
	// 获取用户邮箱
	email := c.GetString("user")
	// 获取紧急任务列表
	taskList, err := s.taskService.GetUrgentTaskList(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", taskList))
}

// GetTeamTaskList 获取小组任务列表
func (s *TaskHandler) GetTeamTaskList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	userId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	// 防止传入过大的分页参数
	if page > 10000 || pageSize > 1000 {
		c.JSON(http.StatusBadRequest, response.Fail("", "分页参数过大", nil))
		return
	}

	taskList, err := s.taskService.GetTeamTaskList(userId, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", taskList))
}

// DeleteTeamTask 删除小组任务
func (s *TaskHandler) DeleteTeamTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Query("taskId"))
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err = s.taskService.DeleteTeamTask(taskId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "删除成功", nil))
}

// AddTeamTask 添加小组任务
func (s *TaskHandler) AddTeamTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBind(&task); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err := s.taskService.AddTeamTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "添加成功", nil))
}

// CompleteTeamTask 小组成员完成小组任务
func (s *TaskHandler) CompleteTeamTask(c *gin.Context) {
	var teamTask models.TeamTask
	if err := c.ShouldBind(&teamTask); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	err := s.taskService.CompleteTeamTask(teamTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "完成成功", nil))
}

// InviteTeamMember 邀请成员
func (s *TaskHandler) InviteTeamMember(c *gin.Context) {
	var teamTask models.TeamTask
	if err := c.ShouldBind(&teamTask); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}
	email := c.GetString("user")

	err := s.taskService.InviteTeamMember(email, teamTask)
	if err != nil {
		if myError.IsMyError(err) {
			c.JSON(http.StatusOK, response.Success("", err.Error(), nil))
		} else {
			c.JSON(http.StatusInternalServerError, response.Fail("", "系统错误", nil))
		}
		return
	}

	c.JSON(http.StatusOK, response.Success("", "发送邀请成功", nil))
}

// GetInviteCode 获取小组任务邀请码
func (s *TaskHandler) GetInviteCode(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Query("taskId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	inviteCode, err := s.taskService.GetInviteCodeByTaskId(taskId)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "查询成功", inviteCode))
}

// JoinTeamByInviteCode 通过邀请码加入小组任务
func (s *TaskHandler) JoinTeamByInviteCode(c *gin.Context) {
	type temp struct {
		InviteCode string `json:"inviteCode"`
	}
	var inviteCodeVo temp
	if err := c.ShouldBind(&inviteCodeVo); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail("", "参数错误", nil))
		return
	}

	email := c.GetString("user")

	err := s.taskService.JoinTeamTaskByInviteCode(email, inviteCodeVo.InviteCode)
	if err != nil {
		c.JSON(http.StatusOK, response.Fail("", err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, response.Success("", "加入小组任务成功", nil))
}
