package entity

type Comment struct {
	GormModel
	UserId  uint   `gorm:"not null" json:"user_id"`
	PhotoId uint   `gorm:"not null" json:"photo_id"`
	Message string `gorm:"not null" json:"message" form:"message" valid:"required~Message is required"`
	User    *User
	Photo   *Photo
}
