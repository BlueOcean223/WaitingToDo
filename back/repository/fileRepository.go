package repository

import (
	"back/models"
	"gorm.io/gorm"
)

type FileRepository struct {
	Db *gorm.DB
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{Db: db}
}

func (s *FileRepository) Insert(file models.File, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Create(&file).Error
}

func (s *FileRepository) Update(file models.File, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Model(&file).Updates(file).Error
}

func (s *FileRepository) Delete(id int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Delete(&models.File{}, id).Error
}

func (s *FileRepository) DeleteByIds(ids []int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("id IN ?", ids).Delete(&models.File{}).Error
}

func (s *FileRepository) DeleteByTaskId(taskId int, tx *gorm.DB) error {
	db := s.Db
	if tx != nil {
		db = tx
	}
	return db.Where("task_id = ?", taskId).Delete(&models.File{}).Error
}

func (s *FileRepository) GetFileByTaskId(taskId int) ([]models.File, error) {
	var files []models.File
	err := s.Db.Where("task_id = ?", taskId).Find(&files).Error
	return files, err
}

func (s *FileRepository) GetFileByTaskIds(ids []int) ([]models.File, error) {
	var files []models.File
	err := s.Db.Where("task_id IN ?", ids).Find(&files).Error
	return files, err
}

func (s *FileRepository) GetFileByIds(ids []int) ([]models.File, error) {
	var files []models.File
	err := s.Db.Where("id IN ?", ids).Find(&files).Error
	return files, err
}
