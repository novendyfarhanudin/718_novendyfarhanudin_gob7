package entity

import (
	"final-project/common/auth"
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type User struct {
	GormModel
	Username string    `gorm:"not null;uniqueIndex" json:"username,omitempty" form:"username" valid:"required~Your username is required"`
	Email    string    `gorm:"not null;uniqueIndex" json:"email,omitempty" form:"email" valid:"required~Your email is required, email~Invalid email format,email~Invalid format email"`
	Password string    `gorm:"not null" json:"password,omitempty" form:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
	Age      int       `gorm:"not null" json:"age,omitempty" form:"age" valid:"required~Your age is required,numeric~Fill age with number,range(8|99)~minimum 8 years old"`
	Photos   []Photo   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos,omitempty"`
	Comments []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comments,omitempty"`
	Socials  []Social  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"socials,omitempty"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(user)
	if errCreate != nil {
		err = errCreate
		return err
	}
	hash, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return
}
