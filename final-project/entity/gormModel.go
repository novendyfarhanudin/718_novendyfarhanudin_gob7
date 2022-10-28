package entity

import "time"

type GormModel struct {
	Id        uint       `gorm:"primary_key" json:"id" example:"1"`
	CreatedAt *time.Time `json:"created_at,omitempty" example:"2022-11-11T21:21:46+00:00"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" example:"2022-11-11T21:21:46+00:00"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty" example:"2022-11-11T21:21:46+00:00"`
}
