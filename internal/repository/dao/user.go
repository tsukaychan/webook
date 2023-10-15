/**
 * @author tsukiyo
 * @date 2023-08-11 1:29
 */

package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("parameter duplicate")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

var _ UserDao = (*UserGormDao)(nil)

type UserDao interface {
	Create(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx *gin.Context, uid int64) (User, error)
	UpdateNonZeroFields(ctx *gin.Context, user User) error
}

type UserGormDao struct {
	db *gorm.DB
}

func (dao *UserGormDao) Create(ctx context.Context, u User) error {
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func (dao *UserGormDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (dao *UserGormDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (dao *UserGormDao) FindById(ctx *gin.Context, uid int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", uid).First(&user).Error
	return user, err
}

func (dao *UserGormDao) UpdateNonZeroFields(ctx *gin.Context, user User) error {
	timeZeroMilli := time.Time{}.UnixMilli()
	if user.Birthday.Int64 == (timeZeroMilli) {
		user.Birthday.Int64 = 0
		user.Birthday.Valid = false
	}
	if user.CreateAt == (timeZeroMilli) {
		user.CreateAt = 0
	}
	if user.UpdateAt == (timeZeroMilli) {
		user.UpdateAt = 0
	}
	return dao.db.Updates(&user).Error
}

type User struct {
	Id       int64 `gorm:"primaryKey,autoIncrement"`
	Birthday sql.NullInt64
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	NickName sql.NullString
	Password string         `gorm:"not null"`
	AboutMe  sql.NullString `gorm:"default:这个用户很懒什么都没有留下;type=varchar(1024)"`
	CreateAt int64
	UpdateAt int64
}

func NewUserGormDao(db *gorm.DB) UserDao {
	return &UserGormDao{
		db: db,
	}
}
