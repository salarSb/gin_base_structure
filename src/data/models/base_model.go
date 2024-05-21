package models

import (
	"base_structure/src/constants"
	"database/sql"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	CreatedBy int            `gorm:"not null"`
	UpdatedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy *sql.NullInt64 `gorm:"null"`
}

func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = -1
	if value != nil {
		userId = int(value.(float64))
	}
	m.CreatedBy = userId
	return
}

func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{
			Valid: true,
			Int64: int64(value.(float64)),
		}
	}
	m.UpdatedBy = userId
	return
}

func (m *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value(constants.UserIdKey)
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId = &sql.NullInt64{
			Valid: true,
			Int64: int64(value.(float64)),
		}
	}
	m.DeletedBy = userId
	return
}
