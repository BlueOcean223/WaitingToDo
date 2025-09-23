package repository

import (
	"backend/internal/models"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTaskRepo 初始化任务仓库测试环境
func setupTaskRepo(t *testing.T) (TaskRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskRepository(gormDB)
	return repo, mock
}

// 测试用例使用的示例任务数据
var testTask = models.Task{
	Id:          1,
	UserId:      1,
	Title:       "Test Task",
	Description: "Test Description",
	Ddl:         "2024-12-31 23:59:59",
	Type:        0,
	Status:      0,
	CreateTime:  time.Now(),
	UpdateTime:  time.Now(),
}

// TestTaskRepository_GetList 测试分页查询任务列表的功能
func TestTaskRepository_GetList(t *testing.T) {
	repo, mock := setupTaskRepo(t)

	countQueryPattern := `SELECT count\(\*\) FROM .tasks. WHERE \(user_id = \? and type = \?\)`
	listQueryPattern := `SELECT \* FROM .tasks. WHERE \(user_id = \? and type = \?\)`

	t.Run("found_with_status", func(t *testing.T) {
		status := 0
		// 模拟count查询
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery(countQueryPattern+` AND status = \? `).
			WithArgs(testTask.UserId, testTask.Type, status).
			WillReturnRows(countRows)

		// 模拟分页查询
		listRows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "ddl", "type", "status", "create_time", "update_time"}).
			AddRow(testTask.Id, testTask.UserId, testTask.Title, testTask.Description, testTask.Ddl, testTask.Type, testTask.Status, testTask.CreateTime, testTask.UpdateTime)
		mock.ExpectQuery(listQueryPattern+` AND status = \? `+`ORDER BY ddl desc, id LIMIT ?`).
			WithArgs(testTask.UserId, testTask.Type, status, 10).
			WillReturnRows(listRows)

		tasks, count, err := repo.GetList(testTask.UserId, 1, 10, testTask.Type, &status)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
		assert.Len(t, tasks, 1)
		assert.Equal(t, testTask.Id, tasks[0].Id)
	})
}

// TestTaskRepository_Create 测试创建任务的功能
func TestTaskRepository_Create(t *testing.T) {
	repo, mock := setupTaskRepo(t)

	query := "INSERT INTO `tasks` (`user_id`,`title`,`description`,`ddl`,`type`,`status`,`create_time`,`update_time`,`id`) VALUES (?,?,?,?,?,?,?,?,?)"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testTask.UserId, testTask.Title, testTask.Description, testTask.Ddl, testTask.Type, testTask.Status, AnyTime{}, AnyTime{}, testTask.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(&testTask, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testTask.UserId, testTask.Title, testTask.Description, testTask.Ddl, testTask.Type, testTask.Status, AnyTime{}, AnyTime{}, testTask.Id).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		err := repo.Create(&testTask, nil)
		assert.Error(t, err)
	})
}

// TestTaskRepository_Update 测试更新任务的功能
func TestTaskRepository_Update(t *testing.T) {
	repo, mock := setupTaskRepo(t)

	queryPattern := `UPDATE .tasks. SET .id.=\?,.user_id.=\?,.title.=\?,.description.=\?,.ddl.=\?,.create_time.=\?,.update_time.=\? WHERE .id. = \?`

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(testTask.Id, testTask.UserId, testTask.Title, testTask.Description, testTask.Ddl, AnyTime{}, sqlmock.AnyArg(), testTask.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(testTask, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(testTask.Id, testTask.UserId, testTask.Title, testTask.Description, testTask.Ddl, AnyTime{}, sqlmock.AnyArg(), testTask.Id).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		err := repo.Update(testTask, nil)
		assert.Error(t, err)
	})
}

// TestTaskRepository_Delete 测试删除任务的功能
func TestTaskRepository_Delete(t *testing.T) {
	repo, mock := setupTaskRepo(t)

	query := "DELETE FROM `tasks` WHERE `tasks`.`id` = ?"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testTask.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(testTask.Id, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testTask.Id).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.Delete(testTask.Id, nil)
		assert.Error(t, err)
	})
}

// TestTaskRepository_GetUrgentList 测试获取紧急任务列表的功能
func TestTaskRepository_GetUrgentList(t *testing.T) {
	repo, mock := setupTaskRepo(t)

	queryPattern := `SELECT \* FROM .tasks. WHERE user_id = \? AND status = \? AND ddl >= \? AND ddl <= \?`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "ddl", "type", "status", "create_time", "update_time"}).
			AddRow(testTask.Id, testTask.UserId, testTask.Title, testTask.Description, testTask.Ddl, testTask.Type, testTask.Status, testTask.CreateTime, testTask.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(testTask.UserId, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		tasks, err := repo.GetUrgentList(testTask.UserId)
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, testTask.Id, tasks[0].Id)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(999, 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "description", "ddl", "type", "status", "create_time", "update_time"}))

		tasks, err := repo.GetUrgentList(999)
		assert.NoError(t, err)
		assert.Len(t, tasks, 0)
	})
}

// TestTaskRepository_GetOneDayDDLTaskList 测试获取一天内到期任务列表的功能
func TestTaskRepository_GetOneDayDDLTaskList(t *testing.T) {
	repo, mock := setupTaskRepo(t)

	queryPattern := `SELECT \* FROM .tasks. WHERE status = \? AND ddl >= \? AND ddl <= \?`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "ddl", "type", "status", "create_time", "update_time"}).
			AddRow(testTask.Id, testTask.UserId, testTask.Title, testTask.Description, testTask.Ddl, testTask.Type, testTask.Status, testTask.CreateTime, testTask.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		tasks, err := repo.GetOneDayDDLTaskList()
		assert.NoError(t, err)
		assert.Len(t, tasks, 1)
		assert.Equal(t, testTask.Id, tasks[0].Id)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(0, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "description", "ddl", "type", "status", "create_time", "update_time"}))

		tasks, err := repo.GetOneDayDDLTaskList()
		assert.NoError(t, err)
		assert.Len(t, tasks, 0)
	})
}

// TestTaskRepository_GetTaskListByIds 测试根据ID列表获取任务列表的功能
func TestTaskRepository_GetTaskListByIds(t *testing.T) {
	repo, mock := setupTaskRepo(t)

	queryPattern := `SELECT \* FROM .tasks. WHERE id IN \(\?,\?\) ORDER BY ddl desc, id`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "title", "description", "ddl", "type", "status", "create_time", "update_time"}).
			AddRow(1, testTask.UserId, "Task 1", "Description 1", testTask.Ddl, testTask.Type, 0, testTask.CreateTime, testTask.UpdateTime).
			AddRow(2, testTask.UserId, "Task 2", "Description 2", testTask.Ddl, testTask.Type, 1, testTask.CreateTime, testTask.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(1, 2).
			WillReturnRows(rows)

		tasks, err := repo.GetTaskListByIds([]int{1, 2})
		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
		assert.Equal(t, 1, tasks[0].Id)
		assert.Equal(t, 2, tasks[1].Id)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(3, 4).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "description", "ddl", "type", "status", "create_time", "update_time"}))

		tasks, err := repo.GetTaskListByIds([]int{3, 4})
		assert.NoError(t, err)
		assert.Len(t, tasks, 0)
	})
}
