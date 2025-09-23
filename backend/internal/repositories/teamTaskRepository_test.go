package repository

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"backend/internal/models"
)

func initTeamTaskRepositoryTest(t *testing.T) (TeamTaskRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTeamTaskRepository(gormDB)
	return repo, mock
}

func TestTeamTaskRepository_GetList(t *testing.T) {
	repo, mock := initTeamTaskRepositoryTest(t)

	queryPattern := "SELECT * FROM `team_task` WHERE user_id = ? LIMIT ?"

	page := 1
	pageSize := 10
	userId := 1

	t.Run("成功获取团队任务列表", func(t *testing.T) {
		expectedTasks := []models.TeamTask{
			{
				Id:         1,
				TaskId:     1,
				UserId:     userId,
				Status:     1,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			},
			{
				Id:         2,
				TaskId:     2,
				UserId:     userId,
				Status:     0,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			},
		}

		rows := sqlmock.NewRows([]string{"id", "task_id", "user_id", "status", "create_time", "update_time"}).
			AddRow(1, 1, userId, 1, time.Now(), time.Now()).
			AddRow(2, 2, userId, 0, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(queryPattern)).
			WithArgs(userId, pageSize).
			WillReturnRows(rows)

		tasks, err := repo.GetList(page, pageSize, userId)

		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
		assert.Equal(t, expectedTasks[0].Status, tasks[0].Status)
		assert.Equal(t, expectedTasks[1].Status, tasks[1].Status)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(queryPattern)).
			WithArgs(userId, pageSize).
			WillReturnError(gorm.ErrInvalidDB)

		tasks, err := repo.GetList(page, pageSize, userId)

		assert.Error(t, err)
		assert.Nil(t, tasks)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeamTaskRepository_Insert(t *testing.T) {
	repo, mock := initTeamTaskRepositoryTest(t)

	pattern := "INSERT INTO `team_task` (`task_id`,`user_id`,`status`,`create_time`,`update_time`) VALUES (?,?,?,?,?)"

	teamTask := models.TeamTask{
		TaskId: 1,
		UserId: 1,
	}

	t.Run("成功创建团队任务", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(teamTask.TaskId, teamTask.UserId, teamTask.Status, AnyTime{}, AnyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Insert(teamTask, nil)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库插入失败", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(teamTask.TaskId, teamTask.UserId, teamTask.Status, AnyTime{}, AnyTime{}).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.Insert(teamTask, nil)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeamTaskRepository_Update(t *testing.T) {
	repo, mock := initTeamTaskRepositoryTest(t)

	pattern := "UPDATE `team_task` SET `id`=?,`task_id`=?,`user_id`=?,`update_time`=? WHERE `id` = ?"

	teamTask := models.TeamTask{
		Id:     1,
		TaskId: 1,
		UserId: 1,
	}

	t.Run("成功更新团队任务", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(teamTask.Id, teamTask.TaskId, teamTask.UserId, AnyTime{}, teamTask.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(teamTask, nil)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库更新失败", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(teamTask.Id, teamTask.TaskId, teamTask.UserId, AnyTime{}, teamTask.Id).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.Update(teamTask, nil)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeamTaskRepository_Delete(t *testing.T) {
	repo, mock := initTeamTaskRepositoryTest(t)

	t.Run("成功删除团队任务", func(t *testing.T) {
		taskId := 1
		userId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `team_task` WHERE task_id = ? AND user_id = ?")).
			WithArgs(taskId, userId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(taskId, userId, nil)

		assert.NoError(t, err)
	})

	t.Run("数据库删除失败", func(t *testing.T) {
		taskId := 1
		userId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `team_task` WHERE task_id = ? AND user_id = ?")).
			WithArgs(taskId, userId).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.Delete(taskId, userId, nil)

		assert.Error(t, err)
	})
}

func TestTeamTaskRepository_GetTeamMembers(t *testing.T) {
	repo, mock := initTeamTaskRepositoryTest(t)

	pattern := "SELECT * FROM `users` WHERE id IN (SELECT user_id FROM `team_task` WHERE task_id = ?)"

	taskId := 1
	t.Run("成功获取团队成员", func(t *testing.T) {
		expectedUsers := []models.User{
			{Id: 1, Name: "user1", Email: "user1@example.com"},
			{Id: 2, Name: "user2", Email: "user2@example.com"},
		}

		rows := sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "user1", "user1@example.com").
			AddRow(2, "user2", "user2@example.com")

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(taskId).
			WillReturnRows(rows)

		users, err := repo.GetTeamMembers(taskId)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(taskId).
			WillReturnError(gorm.ErrInvalidDB)

		users, err := repo.GetTeamMembers(taskId)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestTeamTaskRepository_GetTeamTaskShipByTaskIds(t *testing.T) {
	repo, mock := initTeamTaskRepositoryTest(t)

	pattern := "SELECT * FROM `team_task` WHERE task_id IN (?,?)"

	taskIds := []int{1, 2}

	t.Run("成功获取团队任务关系", func(t *testing.T) {
		expectedTeamTasks := []models.TeamTask{
			{Id: 1, TaskId: 1, UserId: 1, CreateTime: time.Now(), UpdateTime: time.Now()},
			{Id: 2, TaskId: 2, UserId: 2, CreateTime: time.Now(), UpdateTime: time.Now()},
		}

		rows := sqlmock.NewRows([]string{"id", "task_id", "user_id", "create_time", "update_time"}).
			AddRow(1, 1, 1, time.Now(), time.Now()).
			AddRow(2, 2, 2, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(1, 2).
			WillReturnRows(rows)

		teamTasks, err := repo.GetTeamTaskShipByTaskIds(taskIds)

		assert.NoError(t, err)
		assert.Equal(t, expectedTeamTasks, teamTasks)
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(1, 2).
			WillReturnError(gorm.ErrInvalidDB)

		teamTasks, err := repo.GetTeamTaskShipByTaskIds(taskIds)

		assert.Error(t, err)
		assert.Nil(t, teamTasks)
	})
}
