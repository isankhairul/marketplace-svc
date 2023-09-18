package base

import (
	"fmt"
	model_jwt "marketplace-svc/app/model/jwt"
	"time"

	"github.com/go-kit/kit/auth/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"

	"marketplace-svc/pkg/util"
)

type BaseIDModel struct {
	// ID of as primary key
	// in: int64
	ID uint64 `gorm:"primary_key:auto_increment" json:"-"`

	// UID of the product
	// in: int64
	UID string `gorm:"uniqueIndex;type:varchar(255)" json:"uid"`

	IsDeleted    bool      `json:"is_deleted"`
	CreatedByUid string    `json:"-" gorm:"type:varchar(255)"`
	CreatedBy    string    `json:"-" gorm:"type:varchar(255)"`
	CreatedAt    time.Time `json:"-"`
	UpdatedByUid string    `json:"-" gorm:"type:varchar(255)"`
	UpdatedBy    string    `json:"-" gorm:"type:varchar(255)"`
	UpdatedAt    time.Time `json:"-"`
}

func (base *BaseIDModel) BeforeCreate(tx *gorm.DB) error {
	uid, _ := gonanoid.New()
	payload := tx.Statement.Context.Value(jwt.JWTClaimsContextKey)
	if user, ok := payload.(*model_jwt.ClaimsJWT); ok && user != nil {
		tx.Statement.SetColumn("CreatedBy", user.Data.Name)
		tx.Statement.SetColumn("CreatedByUid", user.Data.UserUid)
		tx.Statement.SetColumn("UpdatedBy", user.Data.Name)
		tx.Statement.SetColumn("UpdatedByUid", user.Data.UserUid)
	}

	tx.Statement.SetColumn("UID", uid)
	tx.Statement.SetColumn("IsDeleted", false)
	tx.Statement.SetColumn("CreatedAt", util.TimeNow())
	tx.Statement.SetColumn("UpdatedAt", util.TimeNow())
	return nil
}

func (base *BaseIDModel) BeforeUpdate(tx *gorm.DB) error {
	payload := tx.Statement.Context.Value(jwt.JWTClaimsContextKey)
	if user, ok := payload.(*model_jwt.ClaimsJWT); ok && user != nil {
		tx.Statement.SetColumn("UpdatedBy", user.Data.Name)
		tx.Statement.SetColumn("UpdatedByUid", user.Data.UserUid)
	}

	tx.Statement.SetColumn("UpdatedAt", util.TimeNow())
	return nil
}

type BaseIDModelEpoch struct {
	// ID of as primary key
	// in: int64
	ID uint64 `gorm:"primary_key:auto_increment" json:"-"`

	// UID of the product
	// in: int64
	UID string `gorm:"uniqueIndex;type:varchar(255)" json:"uid"`

	IsDeleted    bool   `json:"is_deleted"`
	CreatedByUid string `json:"-" gorm:"type:varchar(255)"`
	CreatedBy    string `json:"-" gorm:"type:varchar(255)"`
	CreatedAt    int64  `json:"-"`
	UpdatedByUid string `json:"-" gorm:"type:varchar(255)"`
	UpdatedBy    string `json:"-" gorm:"type:varchar(255)"`
	UpdatedAt    int64  `json:"-"`
}

func (base *BaseIDModelEpoch) BeforeCreate(tx *gorm.DB) error {
	uid, _ := gonanoid.New()
	payload := tx.Statement.Context.Value(jwt.JWTClaimsContextKey)
	if user, ok := payload.(*model_jwt.ClaimsJWT); ok && user != nil {
		tx.Statement.SetColumn("CreatedBy", user.Data.Name)
		tx.Statement.SetColumn("CreatedByUid", fmt.Sprint(user.Data.UserUid))
		tx.Statement.SetColumn("UpdatedBy", user.Data.Name)
		tx.Statement.SetColumn("UpdatedByUid", fmt.Sprint(user.Data.UserUid))
	}

	tx.Statement.SetColumn("UID", uid)
	tx.Statement.SetColumn("IsDeleted", false)
	tx.Statement.SetColumn("CreatedAt", util.TimeNow().Unix())
	tx.Statement.SetColumn("UpdatedAt", util.TimeNow().Unix())
	return nil
}

func (base *BaseIDModelEpoch) BeforeUpdate(tx *gorm.DB) error {
	payload := tx.Statement.Context.Value(jwt.JWTClaimsContextKey)
	if user, ok := payload.(*model_jwt.ClaimsJWT); ok && user != nil {
		tx.Statement.SetColumn("UpdatedBy", user.Data.Name)
		tx.Statement.SetColumn("UpdatedByUid", fmt.Sprint(user.Data.UserUid))
	}

	tx.Statement.SetColumn("UpdatedAt", util.TimeNow().Unix())
	return nil
}
