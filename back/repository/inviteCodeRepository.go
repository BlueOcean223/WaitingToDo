package repository

import (
	"back/models"
	"errors"
	"gorm.io/gorm"
)

type InviteCodeRepository struct {
	Db *gorm.DB
}

func NewInviteCodeRepository(db *gorm.DB) *InviteCodeRepository {
	return &InviteCodeRepository{
		Db: db,
	}
}

// Insert 插入
func (r *InviteCodeRepository) Insert(inviteCode *models.InviteCode, tx *gorm.DB) error {
	db := r.Db
	if tx != nil {
		db = tx
	}
	return db.Create(inviteCode).Error
}

// GetByTaskId 根据任务ID获取邀请码
func (r *InviteCodeRepository) GetByTaskId(taskId int) (models.InviteCode, error) {
	var inviteCode models.InviteCode
	err := r.Db.Where("task_id = ?", taskId).First(&inviteCode).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.InviteCode{}, nil // 如果没有找到记录，返回nil
		}
		return models.InviteCode{}, err // 返回其他错误
	}

	return inviteCode, nil
}

// GetByInviteCode 根据邀请码获取记录
func (r *InviteCodeRepository) GetByInviteCode(code string) (models.InviteCode, error) {
	var inviteCode models.InviteCode
	err := r.Db.Where("invite_code = ?", code).First(&inviteCode).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.InviteCode{}, nil // 如果没有找到记录，返回nil
		}
		return models.InviteCode{}, err // 返回其他错误
	}

	return inviteCode, nil
}

// DeleteByTaskId 根据任务ID删除邀请码
func (r *InviteCodeRepository) DeleteByTaskId(taskId int, tx *gorm.DB) error {
	db := r.Db
	if tx != nil {
		db = tx
	}
	return db.Where("task_id = ?", taskId).Delete(&models.InviteCode{}).Error
}
