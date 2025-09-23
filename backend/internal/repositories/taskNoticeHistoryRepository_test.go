package repository

import (
	"backend/internal/models"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func initTaskNoticeHistoryRepo(t *testing.T) (TaskNoticeHistoryRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewTaskNoticeHistoryRepository(gormDB)
	return repo, mock
}

var testTaskNoticeHistory = models.TaskNoticeHistory{
	Id:         1,
	TaskId:     1,
	CreateTime: time.Now(),
	UpdateTime: time.Now(),
}

func TestTaskNoticeHistoryRepository_Insert(t *testing.T) {
	repo, mock := initTaskNoticeHistoryRepo(t)

	query := "INSERT INTO `task_notice_history` (`task_id`,`create_time`,`update_time`,`id`) VALUES (?,?,?,?)"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)). // 使用 regexp.QuoteMeta 转义特殊字符，使其在正则表达式中被正确识别
								WithArgs(testTaskNoticeHistory.TaskId, sqlmock.AnyArg(), sqlmock.AnyArg(), testTaskNoticeHistory.Id).
								WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Insert(testTaskNoticeHistory, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testTaskNoticeHistory.TaskId, sqlmock.AnyArg(), sqlmock.AnyArg(), testTaskNoticeHistory.Id).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err := repo.Insert(testTaskNoticeHistory, nil)
		assert.Error(t, err)
	})
}

func TestTaskNoticeHistoryRepository_GetHistoryByTaskId(t *testing.T) {
	repo, mock := initTaskNoticeHistoryRepo(t)

	queryPattern := `SELECT \* FROM .task_notice_history. WHERE task_id = \? ORDER BY .task_notice_history.\..id. LIMIT \?`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "task_id", "create_time", "update_time"}).
			AddRow(testTaskNoticeHistory.Id, testTaskNoticeHistory.TaskId, testTaskNoticeHistory.CreateTime, testTaskNoticeHistory.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(testTaskNoticeHistory.TaskId, 1).
			WillReturnRows(rows)

		history, err := repo.GetHistoryByTaskId(testTaskNoticeHistory.TaskId)
		assert.NoError(t, err)
		assert.Equal(t, testTaskNoticeHistory.TaskId, history.TaskId)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(111, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		history, err := repo.GetHistoryByTaskId(111)
		assert.NoError(t, err)
		assert.Empty(t, history)
	})
}

func TestTaskNoticeHistoryRepository_GetHistoriesByTaskIds(t *testing.T) {
	repo, mock := initTaskNoticeHistoryRepo(t)

	queryPattern := `SELECT \* FROM .task_notice_history. WHERE task_id in \(\?(?:,\?)*\)`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "task_id", "create_time", "update_time"}).
			AddRow(1, 1, time.Now(), time.Now()).
			AddRow(2, 2, time.Now(), time.Now())

		mock.ExpectQuery(queryPattern).
			WithArgs(1, 2).
			WillReturnRows(rows)

		histories, err := repo.GetHistoriesByTaskIds([]int{1, 2})
		assert.NoError(t, err)
		assert.Len(t, histories, 2)
		assert.Equal(t, 1, histories[0].TaskId)
		assert.Equal(t, 2, histories[1].TaskId)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(3, 4).
			WillReturnRows(sqlmock.NewRows([]string{"id", "task_id", "create_time", "update_time"}))

		histories, err := repo.GetHistoriesByTaskIds([]int{3, 4})
		assert.NoError(t, err)
		assert.Len(t, histories, 0)
	})
}

func TestTaskNoticeHistoryRepository_DeleteHistoryByTaskId(t *testing.T) {
	repo, mock := initTaskNoticeHistoryRepo(t)

	query := "DELETE FROM `task_notice_history` WHERE task_id = ?"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testTaskNoticeHistory.TaskId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteHistoryByTaskId(testTaskNoticeHistory.TaskId, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(111).
			WillReturnError(errors.New("fail"))
		mock.ExpectRollback()

		err := repo.DeleteHistoryByTaskId(111, nil)
		assert.Error(t, err)
	})
}

func TestTaskNoticeHistoryRepository_BatchInsert(t *testing.T) {
	repo, mock := initTaskNoticeHistoryRepo(t)

	query := "INSERT INTO `task_notice_history` (`task_id`,`create_time`,`update_time`,`id`) VALUES (?,?,?,?),(?,?,?,?)"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 2, sqlmock.AnyArg(), sqlmock.AnyArg(), 2).
			WillReturnResult(sqlmock.NewResult(2, 2))
		mock.ExpectCommit()

		err := repo.BatchInsert([]models.TaskNoticeHistory{
			{Id: 1, TaskId: 1, CreateTime: time.Now(), UpdateTime: time.Now()},
			{Id: 2, TaskId: 2, CreateTime: time.Now(), UpdateTime: time.Now()},
		}, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(1, sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 2, sqlmock.AnyArg(), sqlmock.AnyArg(), 2).
			WillReturnError(errors.New("fail"))
		mock.ExpectRollback()

		err := repo.BatchInsert([]models.TaskNoticeHistory{
			{Id: 1, TaskId: 1, CreateTime: time.Now(), UpdateTime: time.Now()},
			{Id: 2, TaskId: 2, CreateTime: time.Now(), UpdateTime: time.Now()},
		}, nil)
		assert.Error(t, err)
	})
}
