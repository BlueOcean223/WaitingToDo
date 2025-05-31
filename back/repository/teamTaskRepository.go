package repository

import (
	"back/models"
	"gorm.io/gorm"
)

type TeamTaskRepository struct {
	Db *gorm.DB
}

func NewTeamTaskRepository(db *gorm.DB) *TeamTaskRepository {
	return &TeamTaskRepository{Db: db}
}

// Insert 插入
func (s *TeamTaskRepository) Insert(teamTask models.TeamTask) error {
	return s.Db.Create(&teamTask).Error
}

// Update 更新
func (s *TeamTaskRepository) Update(teamTask models.TeamTask) error {
	return s.Db.Save(&teamTask).Error
}

// GetList 分页查询
func (s *TeamTaskRepository) GetList(page, pageSize, userId int) ([]models.TeamTask, error) {
	var teamTaskList []models.TeamTask
	offset := (page - 1) * pageSize
	err := s.Db.Where("user_id = ?", userId).Offset(offset).Limit(pageSize).Find(&teamTaskList).Error

	return teamTaskList, err
}

// GetTeamMembers 获取小组成员
func (s *TeamTaskRepository) GetTeamMembers(taskId int) ([]models.User, error) {
	var users []models.User
	subQuery := s.Db.Table("team_task").Select("user_id").Where("task_id = ?", taskId)
	err := s.Db.Table("users").Where("id IN (?)", subQuery).Find(&users).Error

	return users, err
}

// GetTeamTaskShipByTaskIds 根据任务id获取小组任务关系
func (s *TeamTaskRepository) GetTeamTaskShipByTaskIds(taskIds []int) ([]models.TeamTask, error) {
	var teamTaskShip []models.TeamTask
	err := s.Db.Where("task_id IN (?)", taskIds).Find(&teamTaskShip).Error

	return teamTaskShip, err
}

// Delete 删除
func (s *TeamTaskRepository) Delete(taskId, userId int) error {
	return s.Db.Where("task_id = ? AND user_id = ?", taskId, userId).Delete(&models.TeamTask{}).Error
}
