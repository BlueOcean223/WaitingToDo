package repository

import (
	"backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetList(userId, page, pageSize, myType int, status *int) ([]models.Task, int64, error)
	Create(task *models.Task, tx *gorm.DB) error
	Update(task models.Task, tx *gorm.DB) error
	Delete(id int, tx *gorm.DB) error
	GetUrgentList(userId int) ([]models.Task, error)
	GetOneDayDDLTaskList() ([]models.Task, error)
	GetTaskListByIds(ids []int) ([]models.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// GetList 分页查询
func (s *taskRepository) GetList(userId, page, pageSize, myType int, status *int) ([]models.Task, int64, error) {
	// 构建基础查询
	query := s.db.Model(&models.Task{}).Where("user_id = ? and type = ?", userId, myType)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 先获取总数
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 如果没有数据，直接返回
	if count == 0 {
		return []models.Task{}, 0, nil
	}

	// 分页查询
	offset := (page - 1) * pageSize
	var taskList []models.Task
	err = query.Order("ddl desc, id").Offset(offset).Limit(pageSize).Find(&taskList).Error
	return taskList, count, err
}

// Create 新增任务
func (s *taskRepository) Create(task *models.Task, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	return db.Create(&task).Error
}

// Update 修改任务
func (s *taskRepository) Update(task models.Task, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	return db.Model(&task).Updates(task).Error
}

// Delete 删除任务
func (s *taskRepository) Delete(id int, tx *gorm.DB) error {
	db := s.db
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.Task{}, id).Error
}

// GetUrgentList 获取紧急任务
func (s *taskRepository) GetUrgentList(userId int) ([]models.Task, error) {
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

// GetOneDayDDLTaskList 获取还有不到一天就过期的未完成的任务
func (s *taskRepository) GetOneDayDDLTaskList() ([]models.Task, error) {
	// 获取当前时间
	now := time.Now()
	// 计算一天后的时间
	oneDayLater := now.Add(24 * time.Hour)
	// 获取所有任务
	var tasks []models.Task
	err := s.db.Where("status = ? AND ddl >= ? AND ddl <= ?", 0, now, oneDayLater).Find(&tasks).Error
	return tasks, err
}

// GetTaskListByIds 根据任务id列表批量获取任务
func (s *taskRepository) GetTaskListByIds(ids []int) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.Where("id IN (?)", ids).Order("ddl desc").Find(&tasks).Error
	return tasks, err
}
