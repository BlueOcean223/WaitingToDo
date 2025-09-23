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

// setupFileRepo 初始化文件仓库测试环境
func setupFileRepo(t *testing.T) (FileRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewFileRepository(gormDB)
	return repo, mock
}

// 测试用例使用的示例文件数据
var testFile = models.File{
	Id:         1,
	TaskId:     1,
	Name:       "test.txt",
	Url:        "http://example.com/test.txt",
	CreateTime: time.Now(),
	UpdateTime: time.Now(),
}

// TestFileRepository_Insert 测试插入文件的功能
func TestFileRepository_Insert(t *testing.T) {
	repo, mock := setupFileRepo(t)

	query := "INSERT INTO `files` (`task_id`,`name`,`url`,`create_time`,`update_time`,`id`) VALUES (?,?,?,?,?,?)"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFile.TaskId, testFile.Name, testFile.Url, AnyTime{}, AnyTime{}, testFile.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Insert(testFile, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFile.TaskId, testFile.Name, testFile.Url, AnyTime{}, AnyTime{}, testFile.Id).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		err := repo.Insert(testFile, nil)
		assert.Error(t, err)
	})
}

// TestFileRepository_Update 测试更新文件的功能
func TestFileRepository_Update(t *testing.T) {
	repo, mock := setupFileRepo(t)

	queryPattern := `UPDATE .files. SET .id.=\?,.task_id.=\?,.name.=\?,.url.=\?,.create_time.=\?,.update_time.=\? WHERE .id. = \?`

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(testFile.Id, testFile.TaskId, testFile.Name, testFile.Url, AnyTime{}, AnyTime{}, testFile.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(testFile, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(testFile.Id, testFile.TaskId, testFile.Name, testFile.Url, AnyTime{}, AnyTime{}, testFile.Id).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		err := repo.Update(testFile, nil)
		assert.Error(t, err)
	})
}

// TestFileRepository_Delete 测试删除文件的功能
func TestFileRepository_Delete(t *testing.T) {
	repo, mock := setupFileRepo(t)

	query := "DELETE FROM `files` WHERE `files`.`id` = ?"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFile.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(testFile.Id, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFile.Id).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.Delete(testFile.Id, nil)
		assert.Error(t, err)
	})
}

// TestFileRepository_DeleteByIds 测试批量删除文件的功能
func TestFileRepository_DeleteByIds(t *testing.T) {
	repo, mock := setupFileRepo(t)

	queryPattern := `DELETE FROM .files. WHERE id IN \(\?(?:,\?)*\)`

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(1, 2).
			WillReturnResult(sqlmock.NewResult(2, 2))
		mock.ExpectCommit()

		err := repo.DeleteByIds([]int{1, 2}, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(1, 2).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.DeleteByIds([]int{1, 2}, nil)
		assert.Error(t, err)
	})
}

// TestFileRepository_DeleteByTaskId 测试根据任务ID删除文件的功能
func TestFileRepository_DeleteByTaskId(t *testing.T) {
	repo, mock := setupFileRepo(t)

	query := "DELETE FROM `files` WHERE task_id = ?"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFile.TaskId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteByTaskId(testFile.TaskId, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFile.TaskId).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.DeleteByTaskId(testFile.TaskId, nil)
		assert.Error(t, err)
	})
}

// TestFileRepository_GetFileByTaskId 测试根据任务ID查询文件的功能
func TestFileRepository_GetFileByTaskId(t *testing.T) {
	repo, mock := setupFileRepo(t)

	queryPattern := `SELECT \* FROM .files. WHERE task_id = \?`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "task_id", "name", "url", "create_time", "update_time"}).
			AddRow(testFile.Id, testFile.TaskId, testFile.Name, testFile.Url, testFile.CreateTime, testFile.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(testFile.TaskId).
			WillReturnRows(rows)

		files, err := repo.GetFileByTaskId(testFile.TaskId)
		assert.NoError(t, err)
		assert.Len(t, files, 1)
		assert.Equal(t, testFile.TaskId, files[0].TaskId)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(999).
			WillReturnRows(sqlmock.NewRows([]string{"id", "task_id", "name", "url", "create_time", "update_time"}))

		files, err := repo.GetFileByTaskId(999)
		assert.NoError(t, err)
		assert.Len(t, files, 0)
	})
}

// TestFileRepository_GetFileByTaskIds 测试根据任务ID列表查询文件的功能
func TestFileRepository_GetFileByTaskIds(t *testing.T) {
	repo, mock := setupFileRepo(t)

	queryPattern := `SELECT \* FROM .files. WHERE task_id IN \(\?(?:,\?)*\)`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "task_id", "name", "url", "create_time", "update_time"}).
			AddRow(1, 1, "file1.txt", "http://example.com/file1.txt", time.Now(), time.Now()).
			AddRow(2, 2, "file2.txt", "http://example.com/file2.txt", time.Now(), time.Now())

		mock.ExpectQuery(queryPattern).
			WithArgs(1, 2).
			WillReturnRows(rows)

		files, err := repo.GetFileByTaskIds([]int{1, 2})
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Equal(t, 1, files[0].TaskId)
		assert.Equal(t, 2, files[1].TaskId)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(3, 4).
			WillReturnRows(sqlmock.NewRows([]string{"id", "task_id", "name", "url", "create_time", "update_time"}))

		files, err := repo.GetFileByTaskIds([]int{3, 4})
		assert.NoError(t, err)
		assert.Len(t, files, 0)
	})
}

// TestFileRepository_GetFileByIds 测试根据文件ID列表查询文件的功能
func TestFileRepository_GetFileByIds(t *testing.T) {
	repo, mock := setupFileRepo(t)

	queryPattern := `SELECT \* FROM .files. WHERE id IN \(\?(?:,\?)*\)`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "task_id", "name", "url", "create_time", "update_time"}).
			AddRow(1, 1, "file1.txt", "http://example.com/file1.txt", time.Now(), time.Now()).
			AddRow(2, 1, "file2.txt", "http://example.com/file2.txt", time.Now(), time.Now())

		mock.ExpectQuery(queryPattern).
			WithArgs(1, 2).
			WillReturnRows(rows)

		files, err := repo.GetFileByIds([]int{1, 2})
		assert.NoError(t, err)
		assert.Len(t, files, 2)
		assert.Equal(t, 1, files[0].Id)
		assert.Equal(t, 2, files[1].Id)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(3, 4).
			WillReturnRows(sqlmock.NewRows([]string{"id", "task_id", "name", "url", "create_time", "update_time"}))

		files, err := repo.GetFileByIds([]int{3, 4})
		assert.NoError(t, err)
		assert.Len(t, files, 0)
	})
}
