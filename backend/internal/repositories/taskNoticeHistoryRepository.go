package repository

import (
	"backend/internal/models"
	"errors"
	"gorm.io/gorm"
)

type TaskNoticeHistoryRepository interface {
	Insert(models.TaskNoticeHistory, *gorm.DB) error
	GetHistoryByTaskId(int) (models.TaskNoticeHistory, error)
	GetHistoriesByTaskIds([]int) ([]models.TaskNoticeHistory, error)
	DeleteHistoryByTaskId(int, *gorm.DB) error
	BatchInsert([]models.TaskNoticeHistory, *gorm.DB) error
}

type taskNoticeHistoryRepository struct {
	Db *gorm.DB
}

func NewTaskNoticeHistoryRepository(db *gorm.DB) TaskNoticeHistoryRepository {
	return &taskNoticeHistoryRepository{
		Db: db,
	}
}

// Insert 插入任务通知历史
func (s *taskNoticeHistoryRepository) Insert(taskNoticeHistory models.TaskNoticeHistory, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Create(&taskNoticeHistory).Error
}

// GetHistoryByTaskId 根据任务ID获取任务通知历史
func (s *taskNoticeHistoryRepository) GetHistoryByTaskId(taskId int) (models.TaskNoticeHistory, error) {
	var taskNoticeHistory models.TaskNoticeHistory
	err := s.Db.Where("task_id = ?", taskId).First(&taskNoticeHistory).Error

	if errors.As(err, &gorm.ErrRecordNotFound) {
		return models.TaskNoticeHistory{}, nil
	}

	return taskNoticeHistory, err
}

// GetHistoriesByTaskIds 根据任务ID列表获取任务通知历史
func (s *taskNoticeHistoryRepository) GetHistoriesByTaskIds(taskIds []int) ([]models.TaskNoticeHistory, error) {
	var taskNoticeHistories []models.TaskNoticeHistory
	err := s.Db.Where("task_id in (?)", taskIds).Find(&taskNoticeHistories).Error

	return taskNoticeHistories, err
}

// DeleteHistoryByTaskId 根据任务ID删除任务通知历史
func (s *taskNoticeHistoryRepository) DeleteHistoryByTaskId(taskId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("task_id = ?", taskId).Delete(&models.TaskNoticeHistory{}).Error
}

// BatchInsert 批量插入任务通知历史
func (s *taskNoticeHistoryRepository) BatchInsert(taskNoticeHistories []models.TaskNoticeHistory, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}

	// 批量插入，每次最大一百条
	return db.CreateInBatches(taskNoticeHistories, 100).Error
}
