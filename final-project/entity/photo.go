package entity

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Photo struct {
	GormModel
	Title    string    `gorm:"not null" json:"username"  valid:"required~Title is required"`
	Caption  string    `gorm:"not null" json:"email"  valid:"required~Caption is required"`
	PhotoUrl string    `gorm:"not null" json:"photo_url"  valid:"required~Photo url is required"`
	UserId   uint      `gorm:"not null" json:"user_id"`
	Comment  []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments"`
	User     *User     `json:"user"`
}

func (photo *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(photo)
	if errCreate != nil {
		return errCreate
	}
	return
}
