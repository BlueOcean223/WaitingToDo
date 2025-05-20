package repository

import (
	"back/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// GetList 分页查询
func (s *TaskRepository) GetList(userId, page, pageSize, myType int) ([]models.Task, int64, error) {
	// 计算偏移量
	offset := (page - 1) * pageSize
	// 查询总数
	var count int64
	err := s.db.Model(&models.Task{}).Where("user_id = ? and type = ?", userId, myType).Count(&count).Error

	var taskList []models.Task
	// 分页查询
	err = s.db.Where("user_id = ? and type = ?", userId, myType).Offset(offset).Limit(pageSize).Find(&taskList).Error
	if err != nil {
		return nil, 0, err
	}
	return taskList, count, nil
}

// Create 新增任务
func (s *TaskRepository) Create(task models.Task) error {
	err := s.db.Create(&task).Error
	return err
}
