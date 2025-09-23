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

func setupImageTest(t *testing.T) (ImageRepository, sqlmock.Sqlmock, *gorm.DB) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	repo := NewImageRepository(gormDB)
	return repo, mock, gormDB
}

func TestImageRepository_InsertImage(t *testing.T) {
	repo, mock, gormDB := setupImageTest(t)

	pattern := "INSERT INTO `images` (`md5`,`url`,`create_time`,`update_time`) VALUES (?,?,?,?)"

	image := models.Image{
		Md5: "abc123def456",
		Url: "/uploads/test.jpg",
	}

	t.Run("成功插入图片", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(image.Md5, image.Url, AnyTime{}, AnyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.InsertImage(image, nil)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("插入图片失败", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(image.Md5, image.Url, AnyTime{}, AnyTime{}).
			WillReturnError(gorm.ErrInvalidDB)
		mock.ExpectRollback()

		err := repo.InsertImage(image, nil)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("使用事务插入图片", func(t *testing.T) {
		// 模拟事务
		mock.ExpectBegin()
		tx := gormDB.Begin()

		mock.ExpectExec(regexp.QuoteMeta(pattern)).
			WithArgs(image.Md5, image.Url, AnyTime{}, AnyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertImage(image, tx)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestImageRepository_GetImageByMD5(t *testing.T) {
	repo, mock, _ := setupImageTest(t)

	parttern := "SELECT * FROM `images` WHERE md5 = ? ORDER BY `images`.`id` LIMIT ?"

	t.Run("成功根据MD5获取图片", func(t *testing.T) {
		md5Hash := "abc123def456"

		expectedImage := models.Image{
			ID:         1,
			Md5:        md5Hash,
			Url:        "/uploads/test.jpg",
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}

		rows := sqlmock.NewRows([]string{"id", "md5", "url", "create_time", "update_time"}).
			AddRow(1, md5Hash, "/uploads/test.jpg", time.Now(), time.Now())

		mock.ExpectQuery(regexp.QuoteMeta(parttern)).
			WithArgs(md5Hash, 1).
			WillReturnRows(rows)

		image, err := repo.GetImageByMD5(md5Hash)

		assert.NoError(t, err)
		assert.Equal(t, expectedImage.Md5, image.Md5)
		assert.Equal(t, expectedImage.Url, image.Url)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("图片不存在时返回空对象", func(t *testing.T) {
		md5Hash := "notfound"

		mock.ExpectQuery(regexp.QuoteMeta(parttern)).
			WithArgs(md5Hash, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		image, err := repo.GetImageByMD5(md5Hash)

		assert.NoError(t, err)
		assert.Equal(t, models.Image{}, image)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("数据库查询失败", func(t *testing.T) {
		md5Hash := "abc123def456"

		mock.ExpectQuery(regexp.QuoteMeta(parttern)).
			WithArgs(md5Hash, 1).
			WillReturnError(gorm.ErrInvalidDB)

		image, err := repo.GetImageByMD5(md5Hash)

		assert.Error(t, err)
		assert.Equal(t, models.Image{}, image)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
