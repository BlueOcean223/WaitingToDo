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

func setupMessageTest(t *testing.T) (MessageRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewMessageRepository(gormDB)

	return repo, mock
}

func TestMessageRepository_GetMessageList(t *testing.T) {
	repo, mock := setupMessageTest(t)

	pattern := "SELECT * FROM `messages` WHERE to_id = ? ORDER BY send_time desc, id LIMIT ?"

	t.Run("成功获取消息列表", func(t *testing.T) {
		userId := 1
		page := 1
		pageSize := 10

		expectedMessages := []models.Message{
			{
				Id:          1,
				Title:       "Hello",
				Description: "Hello message",
				FromId:      2,
				ToId:        userId,
				Type:        0,
				SendTime:    "2024-01-01 10:00:00",
				OutId:       0,
				IsRead:      0,
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			},
			{
				Id:          2,
				Title:       "Hi there",
				Description: "Hi there message",
				FromId:      userId,
				ToId:        2,
				Type:        0,
				SendTime:    "2024-01-01 10:05:00",
				OutId:       0,
				IsRead:      1,
				CreateTime:  time.Now(),
				UpdateTime:  time.Now(),
			},
		}

		rows := sqlmock.NewRows([]string{"id", "title", "description", "from_id", "to_id", "type", "send_time", "out_id", "is_read", "create_time", "update_time"}).
			AddRow(1, "Hello", "Hello message", 2, userId, 0, "2024-01-01 10:00:00", 0, 0, time.Now(), time.Now()).
			AddRow(2, "Hi there", "Hi there message", userId, 2, 0, "2024-01-01 10:05:00", 0, 1, time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(userId, pageSize).
			WillReturnRows(rows)

		messages, err := repo.GetMessageList(page, pageSize, userId)

		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, expectedMessages[0].Title, messages[0].Title)
		assert.Equal(t, expectedMessages[1].Title, messages[1].Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		userId := 1
		page := 1
		pageSize := 10

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(userId, pageSize).
			WillReturnError(gorm.ErrInvalidDB)

		messages, err := repo.GetMessageList(page, pageSize, userId)

		assert.Error(t, err)
		assert.Nil(t, messages)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMessageRepository_InsertMessage(t *testing.T) {
	repo, mock := setupMessageTest(t)

	pattern := "INSERT INTO `messages` (`title`,`description`,`from_id`,`to_id`,`type`,`send_time`,`out_id`,`is_read`,`create_time`,`update_time`) VALUES (?,?,?,?,?,?,?,?,?,?)"

	message := models.Message{
		Title:       "Hello",
		Description: "Hello World",
		FromId:      1,
		ToId:        2,
		Type:        0,
		SendTime:    "2024-01-01 10:00:00",
		OutId:       0,
		IsRead:      0,
	}

	t.Run("成功插入消息", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(message.Title, message.Description, message.FromId, message.ToId, message.Type, message.SendTime, message.OutId, message.IsRead, AnyTime{}, AnyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.InsertMessage(message, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("插入消息失败", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(message.Title, message.Description, message.FromId, message.ToId, message.Type, message.SendTime, message.OutId, message.IsRead, AnyTime{}, AnyTime{}).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.InsertMessage(message, nil)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMessageRepository_GetUnreadMessageCount(t *testing.T) {
	repo, mock := setupMessageTest(t)

	pattern := "SELECT count(*) FROM `messages` WHERE to_id = ? and is_read = ?"

	t.Run("成功获取未读消息数量", func(t *testing.T) {
		userId := 1
		expectedCount := int64(5)

		rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(userId, 0).
			WillReturnRows(rows)

		count, err := repo.GetUnreadMessageCount(userId)

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		userId := 1

		mock.ExpectQuery(regexp.QuoteMeta(pattern)).
			WithArgs(userId, 0).
			WillReturnError(gorm.ErrInvalidDB)

		count, err := repo.GetUnreadMessageCount(userId)

		assert.Error(t, err)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMessageRepository_ReadAllMessage(t *testing.T) {
	repo, mock := setupMessageTest(t)

	pattern := "UPDATE `messages` SET `is_read`=?,`update_time`=? WHERE to_id = ?"

	t.Run("成功标记所有消息为已读", func(t *testing.T) {
		userId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(1, AnyTime{}, userId).
			WillReturnResult(sqlmock.NewResult(0, 3))
		mock.ExpectCommit()

		err := repo.ReadAllMessage(userId, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("标记所有消息为已读失败", func(t *testing.T) {
		userId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(1, AnyTime{}, userId).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.ReadAllMessage(userId, nil)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMessageRepository_Delete(t *testing.T) {
	repo, mock := setupMessageTest(t)

	pattern := "DELETE FROM `messages` WHERE `messages`.`id` = ?"

	t.Run("成功删除消息", func(t *testing.T) {
		messageId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(messageId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(messageId, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("删除消息失败", func(t *testing.T) {
		messageId := 1

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(messageId).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.Delete(messageId, nil)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMessageRepository_Update(t *testing.T) {
	repo, mock := setupMessageTest(t)

	pattern := "UPDATE `messages` SET `title`=?,`description`=?,`is_read`=?,`update_time`=? WHERE `id` = ?"

	message := models.Message{
		Id:          1,
		Title:       "Updated Title",
		Description: "Updated Description",
		IsRead:      1,
	}

	t.Run("成功更新消息", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(message.Title, message.Description, message.IsRead, AnyTime{}, message.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Update(message, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("更新消息失败", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(message.Title, message.Description, message.IsRead, AnyTime{}, message.Id).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.Update(message, nil)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
