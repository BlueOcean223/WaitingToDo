package repository

import (
	"backend/internal/models"
	"gorm.io/gorm"
)

type TeamTaskRepository interface {
	Insert(teamTask models.TeamTask, tx *gorm.DB) error
	Update(teamTask models.TeamTask, tx *gorm.DB) error
	GetList(page, pageSize, userId int) ([]models.TeamTask, error)
	GetTeamMembers(taskId int) ([]models.User, error)
	GetTeamTaskShipByTaskIds(taskIds []int) ([]models.TeamTask, error)
	Delete(taskId, userId int, tx *gorm.DB) error
}

type teamTaskRepository struct {
	Db *gorm.DB
}

func NewTeamTaskRepository(db *gorm.DB) TeamTaskRepository {
	return &teamTaskRepository{Db: db}
}

// Insert 插入
func (s *teamTaskRepository) Insert(teamTask models.TeamTask, tx *gorm.DB) error {
	Db := s.Db
	if tx != nil {
		Db = tx
	}
	return Db.Create(&teamTask).Error
}

// Update 更新
func (s *teamTaskRepository) Update(teamTask models.TeamTask, tx *gorm.DB) error {
	Db := s.Db
	if tx != nil {
		Db = tx
	}
	return Db.Model(&teamTask).Updates(teamTask).Error
}

// GetList 分页查询
func (s *teamTaskRepository) GetList(page, pageSize, userId int) ([]models.TeamTask, error) {
	var teamTaskList []models.TeamTask
	offset := (page - 1) * pageSize
	err := s.Db.Where("user_id = ?", userId).Offset(offset).Limit(pageSize).Find(&teamTaskList).Error

	return teamTaskList, err
}

// GetTeamMembers 获取小组成员
func (s *teamTaskRepository) GetTeamMembers(taskId int) ([]models.User, error) {
	var users []models.User
	subQuery := s.Db.Table("team_task").Select("user_id").Where("task_id = ?", taskId)
	err := s.Db.Table("users").Where("id IN (?)", subQuery).Find(&users).Error

	return users, err
}

// GetTeamTaskShipByTaskIds 根据任务id获取小组任务关系
func (s *teamTaskRepository) GetTeamTaskShipByTaskIds(taskIds []int) ([]models.TeamTask, error) {
	var teamTaskShip []models.TeamTask
	err := s.Db.Where("task_id IN (?)", taskIds).Find(&teamTaskShip).Error

	return teamTaskShip, err
}

// Delete 删除
func (s *teamTaskRepository) Delete(taskId, userId int, tx *gorm.DB) error {
	Db := s.Db
	if tx != nil {
		Db = tx
	}
	return Db.Where("task_id = ? AND user_id = ?", taskId, userId).Delete(&models.TeamTask{}).Error
}
