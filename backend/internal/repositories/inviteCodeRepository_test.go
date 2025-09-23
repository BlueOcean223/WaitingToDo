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

func setupInviteCodeTest(t *testing.T) (InviteCodeRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewInviteCodeRepository(gormDB)

	return repo, mock
}

func TestInviteCodeRepository_Insert(t *testing.T) {
	repo, mock := setupInviteCodeTest(t)

	pattern := "INSERT INTO `invite_codes` (`task_id`,`invite_code`,`create_time`,`update_time`) VALUES (?,?,?,?)"

	inviteCode := &models.InviteCode{
		TaskId:     1,
		InviteCode: "ABC123",
	}

	t.Run("成功创建邀请码", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(inviteCode.TaskId, inviteCode.InviteCode, AnyTime{}, AnyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Insert(inviteCode, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("创建邀请码失败", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `invite_codes` (`task_id`,`invite_code`,`create_time`,`update_time`,`id`) VALUES (?,?,?,?,?)")).
			WithArgs(inviteCode.TaskId, inviteCode.InviteCode, AnyTime{}, AnyTime{}, 1).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.Insert(inviteCode, nil)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestInviteCodeRepository_GetByTaskId(t *testing.T) {
	repo, mock := setupInviteCodeTest(t)

	pattern := "SELECT * FROM `invite_codes` WHERE task_id = ? ORDER BY `invite_codes`.`id` LIMIT ?"

	t.Run("成功根据任务ID获取邀请码", func(t *testing.T) {
		taskId := 1
		inviteCode := "ABC123"
		createTime := time.Now()
		updateTime := time.Now()

		expectedInviteCode := models.InviteCode{
			Id:         1,
			TaskId:     taskId,
			InviteCode: inviteCode,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}

		rows := sqlmock.NewRows([]string{"id", "task_id", "invite_code", "create_time", "update_time"}).
			AddRow(1, taskId, inviteCode, createTime, updateTime)

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(taskId, 1).
			WillReturnRows(rows)

		result, err := repo.GetByTaskId(taskId)

		assert.NoError(t, err)
		assert.Equal(t, expectedInviteCode.Id, result.Id)
		assert.Equal(t, expectedInviteCode.TaskId, result.TaskId)
		assert.Equal(t, expectedInviteCode.InviteCode, result.InviteCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("任务ID对应的邀请码不存在", func(t *testing.T) {
		taskId := 999

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(taskId, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetByTaskId(taskId)

		assert.NoError(t, err)
		assert.Equal(t, models.InviteCode{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestInviteCodeRepository_GetByInviteCode(t *testing.T) {
	repo, mock := setupInviteCodeTest(t)

	pattern := "SELECT * FROM `invite_codes` WHERE invite_code = ? ORDER BY `invite_codes`.`id` LIMIT ?"

	t.Run("成功根据邀请码获取记录", func(t *testing.T) {
		code := "ABC123"
		taskId := 1
		createTime := time.Now()
		updateTime := time.Now()

		expectedInviteCode := models.InviteCode{
			Id:         1,
			TaskId:     taskId,
			InviteCode: code,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}

		rows := sqlmock.NewRows([]string{"id", "task_id", "invite_code", "create_time", "update_time"}).
			AddRow(1, taskId, code, createTime, updateTime)

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(code, 1).
			WillReturnRows(rows)

		result, err := repo.GetByInviteCode(code)

		assert.NoError(t, err)
		assert.Equal(t, expectedInviteCode.Id, result.Id)
		assert.Equal(t, expectedInviteCode.TaskId, result.TaskId)
		assert.Equal(t, expectedInviteCode.InviteCode, result.InviteCode)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("邀请码不存在", func(t *testing.T) {
		code := "INVALID"

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(code, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.GetByInviteCode(code)

		assert.NoError(t, err)
		assert.Equal(t, models.InviteCode{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		code := "ABC123"

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(code, 1).
			WillReturnError(gorm.ErrInvalidDB)

		result, err := repo.GetByInviteCode(code)

		assert.Error(t, err)
		assert.Equal(t, models.InviteCode{}, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestInviteCodeRepository_DeleteByTaskId(t *testing.T) {
	repo, mock := setupInviteCodeTest(t)

	parttern := "DELETE FROM `invite_codes` WHERE task_id = ?"

	t.Run("成功根据任务ID删除邀请码", func(t *testing.T) {
		taskId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(parttern)).
			WithArgs(taskId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.DeleteByTaskId(taskId, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("删除邀请码失败", func(t *testing.T) {
		taskId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(parttern)).
			WithArgs(taskId).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.DeleteByTaskId(taskId, nil)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
