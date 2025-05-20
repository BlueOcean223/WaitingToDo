package service

import (
	"back/models"
	"back/models/dto"
	"back/models/vo"
	"back/repository"
	"errors"
	"sort"
)

type TaskService struct {
	authRepository *repository.AuthRepository
	taskRepository *repository.TaskRepository
}

func NewTaskService(authRepository *repository.AuthRepository,
	taskRepository *repository.TaskRepository) *TaskService {
	return &TaskService{
		authRepository: authRepository,
		taskRepository: taskRepository,
	}
}

// GetTaskList 获取任务列表
func (s *TaskService) GetTaskList(email string, page, pageSize int) ([]dto.TaskDto, error) {
	// 根据邮箱查询用户
	user, err := s.authRepository.SelectUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if user == (models.User{}) {
		return nil, errors.New("用户不存在")
	}

	taskList, count, err := s.taskRepository.GetList(user.Id, page, pageSize, 0)
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
	// 将任务列表按ddl降序排序
	sort.Slice(taskDtoList, func(i, j int) bool {
		return taskList[i].Ddl > taskList[j].Ddl
	})
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
