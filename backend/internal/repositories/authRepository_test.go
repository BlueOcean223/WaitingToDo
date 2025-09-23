package repository

import (
	"backend/internal/models"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// AnyTime 用于在sqlmock中匹配任意time.Time值的自定义类型
type AnyTime struct{}

// Match 实现sqlmock.Argument接口，用于验证时间类型参数
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// setup 初始化测试环境
// 返回值:
// - AuthRepository: 认证仓库实例
// - sqlmock.Sqlmock: SQL mock实例，用于模拟数据库操作
func setup(t *testing.T) (AuthRepository, sqlmock.Sqlmock) {
	// 创建一个新的SQL mock数据库连接
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// 使用mock的数据库连接初始化GORM
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewAuthRepository(gormDB)
	return repo, mock
}

// 测试用例使用的示例用户数据
var testUser = models.User{
	Id:         1,
	Email:      "test@example.com",
	Password:   "password",
	CreateTime: time.Now(),
	UpdateTime: time.Now(),
}

// TestAuthRepository_SelectUserByEmail 测试通过邮箱查询用户的功能
func TestAuthRepository_SelectUserByEmail(t *testing.T) {
	repo, mock := setup(t)

	// GORM生成的查询语句模式（包含LIMIT占位符）
	queryPattern := `SELECT \* FROM .users. WHERE email = \? ORDER BY .users.\..id. LIMIT \?`

	t.Run("found", func(t *testing.T) {
		// 模拟查询结果
		rows := sqlmock.NewRows([]string{"id", "email", "password", "create_time", "update_time"}).
			AddRow(testUser.Id, testUser.Email, testUser.Password, testUser.CreateTime, testUser.UpdateTime)

		// 设置预期的查询行为和返回结果
		mock.ExpectQuery(queryPattern).
			WithArgs(testUser.Email, 1).
			WillReturnRows(rows)

		user, err := repo.SelectUserByEmail(testUser.Email)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Email, user.Email)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs("nonexistent@example.com", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.SelectUserByEmail("nonexistent@example.com")
		assert.NoError(t, err)
		assert.Empty(t, user)
	})
}

// TestAuthRepository_SelectUserById 测试通过用户ID查询用户的功能
func TestAuthRepository_SelectUserById(t *testing.T) {
	repo, mock := setup(t)

	// GORM生成的查询语句模式（包含LIMIT占位符）
	queryPattern := `SELECT \* FROM .users. WHERE id = \? ORDER BY .users.\..id. LIMIT \?`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password", "create_time", "update_time"}).
			AddRow(testUser.Id, testUser.Email, testUser.Password, testUser.CreateTime, testUser.UpdateTime)

		mock.ExpectQuery(queryPattern).
			WithArgs(testUser.Id, 1).
			WillReturnRows(rows)

		user, err := repo.SelectUserById(testUser.Id)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Id, user.Id)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.SelectUserById(999)
		assert.NoError(t, err)
		assert.Empty(t, user)
	})
}

// TestAuthRepository_SelectUsersByIds 测试通过用户ID列表查询用户的功能
func TestAuthRepository_SelectUsersByIds(t *testing.T) {
	repo, mock := setup(t)

	// GORM生成的查询语句模式（用于IN查询）
	queryPattern := `SELECT \* FROM .users. WHERE id in \(\?(?:,\?)*\)`

	t.Run("found", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email"}).
			AddRow(1, "test1@example.com").
			AddRow(2, "test2@example.com")

		mock.ExpectQuery(queryPattern).
			WithArgs(1, 2).
			WillReturnRows(rows)

		users, err := repo.SelectUsersByIds([]int{1, 2})
		assert.NoError(t, err)
		assert.Len(t, users, 2)
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(queryPattern).
			WithArgs(3, 4).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

		users, err := repo.SelectUsersByIds([]int{3, 4})
		assert.NoError(t, err)
		assert.Len(t, users, 0)
	})
}

// TestAuthRepository_InsertUser 测试插入用户的功能
func TestAuthRepository_InsertUser(t *testing.T) {
	repo, mock := setup(t)

	// 插入用户的SQL语句模板
	query := "INSERT INTO `users` (`email`,`password`,`name`,`pic`,`description`,`create_time`,`update_time`,`id`) VALUES (?,?,?,?,?,?,?,?)"

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin() // 开始事务
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testUser.Email, testUser.Password, testUser.Name, testUser.Pic, testUser.Description, AnyTime{}, AnyTime{}, testUser.Id).
			WillReturnResult(sqlmock.NewResult(1, 1)) // 模拟返回插入结果
		mock.ExpectCommit() // 提交事务

		err := repo.InsertUser(testUser, nil)
		assert.NoError(t, err)
	})

	// 测试插入失败的场景
	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin() // 开始事务
		mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(testUser.Email, testUser.Password, testUser.Name, testUser.Pic, testUser.Description, AnyTime{}, AnyTime{}, testUser.Id).
			WillReturnError(errors.New("insert error")) // 模拟返回插入错误
		mock.ExpectRollback() // 回滚事务

		err := repo.InsertUser(testUser, nil)
		assert.Error(t, err)
	})
}

// TestAuthRepository_UpdateUser 测试更新用户的功能
func TestAuthRepository_UpdateUser(t *testing.T) {
	repo, mock := setup(t)

	// GORM Updates()生成的更新语句模式
	queryPattern := `UPDATE .users. SET .id.=\?,.email.=\?,.password.=\?,.create_time.=\?,.update_time.=\? WHERE .id. = \?`

	// 测试更新成功的场景
	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin() // 开始事务
		mock.ExpectExec(queryPattern).
			WithArgs(testUser.Id, testUser.Email, testUser.Password, AnyTime{}, AnyTime{}, testUser.Id).
			WillReturnResult(sqlmock.NewResult(1, 1)) // 模拟返回更新结果
		mock.ExpectCommit() // 提交事务

		err := repo.UpdateUser(testUser, nil)
		assert.NoError(t, err)
	})

	// 测试更新失败的场景
	t.Run("fail", func(t *testing.T) {
		mock.ExpectBegin() // 开始事务
		mock.ExpectExec(queryPattern).
			WithArgs(testUser.Id, testUser.Email, testUser.Password, AnyTime{}, AnyTime{}, testUser.Id).
			WillReturnError(errors.New("update error")) // 模拟返回更新错误
		mock.ExpectRollback() // 回滚事务

		err := repo.UpdateUser(testUser, nil)
		assert.Error(t, err)
	})
}
