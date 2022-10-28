package entity

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Social struct {
	GormModel
	Name           string `gorm:"not null" json:"name" valid:"required~Your name is required"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url"  valid:"required~Social media url is required"`
	UserId         uint   `gorm:"not null" json:"user_id"`
	User           *User  `json:"user"`
}

func (social *Social) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(social)
	if errCreate != nil {
		return errCreate
	}
	return
}
