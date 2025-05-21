package repository

import (
	"back/models"
	"gorm.io/gorm"
	"time"
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

// Update 修改任务
func (s *TaskRepository) Update(task models.Task) error {
	return s.db.Model(&task).Updates(task).Error
}

// Delete 删除任务
func (s *TaskRepository) Delete(id int) error {
	return s.db.Delete(&models.Task{}, id).Error
}

// GetUrgentList 获取紧急任务
func (s *TaskRepository) GetUrgentList(userId int) ([]models.Task, error) {
	// 获取当前时间
	now := time.Now()
	// 计算三天后的时间
	threeDaysLater := now.Add(3 * 24 * time.Hour)
	// 查询数据库
	var tasks []models.Task
	err := s.db.Where("user_id = ? AND status = ? AND ddl >= ? AND ddl <= ?",
		userId, 0, now, threeDaysLater).Find(&tasks).Error

	return tasks, err
}
