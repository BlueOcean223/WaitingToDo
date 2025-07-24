package repository

import (
	"back/models"
	"errors"
	"gorm.io/gorm"
)

type TaskNoticeHistoryRepository struct {
	Db *gorm.DB
}

func NewTaskNoticeHistoryRepository(db *gorm.DB) *TaskNoticeHistoryRepository {
	return &TaskNoticeHistoryRepository{
		Db: db,
	}
}

// Insert 插入任务通知历史
func (s *TaskNoticeHistoryRepository) Insert(taskNoticeHistory models.TaskNoticeHistory, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Create(&taskNoticeHistory).Error
}

// GetHistoryByTaskId 根据任务ID获取任务通知历史
func (s *TaskNoticeHistoryRepository) GetHistoryByTaskId(taskId int) (models.TaskNoticeHistory, error) {
	var taskNoticeHistory models.TaskNoticeHistory
	err := s.Db.Where("task_id = ?", taskId).First(&taskNoticeHistory).Error

	if errors.As(err, &gorm.ErrRecordNotFound) {
		return models.TaskNoticeHistory{}, nil
	}

	return taskNoticeHistory, err
}

// GetHistoriesByTaskIds 根据任务ID列表获取任务通知历史
func (s *TaskNoticeHistoryRepository) GetHistoriesByTaskIds(taskIds []int) ([]models.TaskNoticeHistory, error) {
	var taskNoticeHistories []models.TaskNoticeHistory
	err := s.Db.Where("task_id in (?)", taskIds).Find(&taskNoticeHistories).Error

	return taskNoticeHistories, err
}

// DeleteHistoryByTaskId 根据任务ID删除任务通知历史
func (s *TaskNoticeHistoryRepository) DeleteHistoryByTaskId(taskId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("task_id = ?", taskId).Delete(&models.TaskNoticeHistory{}).Error
}
