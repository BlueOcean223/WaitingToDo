package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type FileRepository interface {
	Insert(file models.File, tx *gorm.DB) error
	Update(file models.File, tx *gorm.DB) error
	Delete(id int, tx *gorm.DB) error
	DeleteByIds(ids []int, tx *gorm.DB) error
	DeleteByTaskId(taskId int, tx *gorm.DB) error
	GetFileByTaskId(taskId int) ([]models.File, error)
	GetFileByTaskIds(ids []int) ([]models.File, error)
	GetFileByIds(ids []int) ([]models.File, error)
}

type fileRepository struct {
	Db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{Db: db}
}

func (s *fileRepository) Insert(file models.File, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Create(&file).Error
}

func (s *fileRepository) Update(file models.File, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Model(&file).Updates(file).Error
}

func (s *fileRepository) Delete(id int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.File{}, id).Error
}

func (s *fileRepository) DeleteByIds(ids []int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("id IN ?", ids).Delete(&models.File{}).Error
}

func (s *fileRepository) DeleteByTaskId(taskId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("task_id = ?", taskId).Delete(&models.File{}).Error
}

func (s *fileRepository) GetFileByTaskId(taskId int) ([]models.File, error) {
	var files []models.File
	err := s.Db.Where("task_id = ?", taskId).Find(&files).Error
	return files, err
}

func (s *fileRepository) GetFileByTaskIds(ids []int) ([]models.File, error) {
	var files []models.File
	err := s.Db.Where("task_id IN ?", ids).Find(&files).Error
	return files, err
}

func (s *fileRepository) GetFileByIds(ids []int) ([]models.File, error) {
	var files []models.File
	err := s.Db.Where("id IN ?", ids).Find(&files).Error
	return files, err
}
