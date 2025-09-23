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

// setupFriendRepo 初始化好友仓库测试环境
func setupFriendRepo(t *testing.T) (FriendRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewFriendRepository(gormDB)
	return repo, mock
}

// 测试用例使用的示例好友数据
var testFriend = models.Friend{
	Id:         1,
	UserId:     1,
	FriendId:   2,
	Status:     1,
	CreateTime: time.Now(),
	UpdateTime: time.Now(),
}

// TestFriendRepository_GetFriendList 测试获取好友列表的功能
func TestFriendRepository_GetFriendList(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	// 查询模式，匹配 GORM 生成的子查询
	queryPattern := `SELECT \* FROM .users. WHERE id in \(SELECT friend_id FROM .friends. WHERE user_id = \? and status = \?\)`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "name", "create_time", "update_time"}).
			AddRow(testUser.Id, testUser.Email, testUser.Name, testUser.CreateTime, testUser.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(1, 1).
			WillReturnRows(rows)

		friends, err := repo.GetFriendList(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(friends))
		assert.Equal(t, testUser.Id, friends[0].Id)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(999, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "create_time", "update_time"}))

		friends, err := repo.GetFriendList(999)
		assert.NoError(t, err)
		assert.Len(t, friends, 0)
	})
}

// TestFriendRepository_GetFriendRelation 测试获取好友关系的功能
func TestFriendRepository_GetFriendRelation(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	queryPattern := `SELECT \* FROM .friends. WHERE user_id = \? and friend_id = \? ORDER BY .friends.\..id. LIMIT \?`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "friend_id", "status", "create_time", "update_time"}).
			AddRow(testFriend.Id, testFriend.UserId, testFriend.FriendId, testFriend.Status, testFriend.CreateTime, testFriend.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(testFriend.UserId, testFriend.FriendId, 1).
			WillReturnRows(rows)

		friend, err := repo.GetFriendRelation(testFriend.UserId, testFriend.FriendId)
		assert.NoError(t, err)
		assert.Equal(t, testFriend.UserId, friend.UserId)
		assert.Equal(t, testFriend.FriendId, friend.FriendId)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(999, 888, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		friend, err := repo.GetFriendRelation(999, 888)
		assert.NoError(t, err)
		assert.Empty(t, friend)
	})
}

// TestFriendRepository_AddFriendRequest 测试添加好友请求的功能
func TestFriendRepository_AddFriendRequest(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	query := "INSERT INTO `friends` (`user_id`,`friend_id`,`status`,`create_time`,`update_time`,`id`) VALUES (?,?,?,?,?,?)"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.UserId, testFriend.FriendId, testFriend.Status, AnyTime{}, AnyTime{}, testFriend.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.AddFriendRequest(&testFriend, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.UserId, testFriend.FriendId, testFriend.Status, AnyTime{}, AnyTime{}, testFriend.Id).
			WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		err := repo.AddFriendRequest(&testFriend, nil)
		assert.Error(t, err)
	})
}

// TestFriendRepository_UpdateFriend 测试更新好友关系的功能
func TestFriendRepository_UpdateFriend(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	queryPattern := `UPDATE .friends. SET .user_id.=\?,.friend_id.=\?,.status.=\?,.create_time.=\?,.update_time.=\? WHERE .id. = \?`

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(testFriend.UserId, testFriend.FriendId, testFriend.Status, AnyTime{}, AnyTime{}, testFriend.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.UpdateFriend(testFriend, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(queryPattern).
			WithArgs(testFriend.UserId, testFriend.FriendId, testFriend.Status, AnyTime{}, AnyTime{}, testFriend.Id).
			WillReturnError(errors.New("update error"))
		mock.ExpectRollback()

		err := repo.UpdateFriend(testFriend, nil)
		assert.Error(t, err)
	})
}

// TestFriendRepository_DeleteFriend 测试删除好友关系的功能
func TestFriendRepository_DeleteFriend(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	query := "DELETE FROM `friends` WHERE id = ?"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteFriend(testFriend.Id, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.Id).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.DeleteFriend(testFriend.Id, nil)
		assert.Error(t, err)
	})
}

// TestFriendRepository_DeleteByUserIdAndFriendId 测试根据用户ID和好友ID删除好友关系的功能
func TestFriendRepository_DeleteByUserIdAndFriendId(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	query := "DELETE FROM `friends` WHERE user_id = ? AND friend_id = ?"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.UserId, testFriend.FriendId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteByUserIdAndFriendId(testFriend.UserId, testFriend.FriendId, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.UserId, testFriend.FriendId).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.DeleteByUserIdAndFriendId(testFriend.UserId, testFriend.FriendId, nil)
		assert.Error(t, err)
	})
}

// TestFriendRepository_GetIsFriend 测试查询是否已经是好友关系的功能
func TestFriendRepository_GetIsFriend(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	queryPattern := `SELECT \* FROM .friends. WHERE user_id = \? AND friend_id = \?`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "user_id", "friend_id", "status", "create_time", "update_time"}).
			AddRow(testFriend.Id, testFriend.UserId, testFriend.FriendId, testFriend.Status, testFriend.CreateTime, testFriend.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(testFriend.UserId, testFriend.FriendId).
			WillReturnRows(rows)

		friends, err := repo.GetIsFriend(testFriend.UserId, testFriend.FriendId)
		assert.NoError(t, err)
		assert.Len(t, friends, 1)
		assert.Equal(t, testFriend.UserId, friends[0].UserId)
		assert.Equal(t, testFriend.FriendId, friends[0].FriendId)
	})

	t.Run("not_found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(999, 888).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "friend_id", "status", "create_time", "update_time"}))

		friends, err := repo.GetIsFriend(999, 888)
		assert.NoError(t, err)
		assert.Len(t, friends, 0)
	})
}

// TestFriendRepository_DeleteIsFriend 测试删除已经是好友的多余好友关系的功能
func TestFriendRepository_DeleteIsFriend(t *testing.T) {
	repo, mock := setupFriendRepo(t)

	query := "DELETE FROM `friends` WHERE user_id = ? AND friend_id = ? AND status = ?"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.UserId, testFriend.FriendId, 0).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteIsFriend(testFriend.UserId, testFriend.FriendId, nil)
		assert.NoError(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testFriend.UserId, testFriend.FriendId, 0).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		err := repo.DeleteIsFriend(testFriend.UserId, testFriend.FriendId, nil)
		assert.Error(t, err)
	})
}
