package dto

import (
	"time"
)

type SocialRequest struct {
	Name           string `json:"name"`
	SocialMediaUrl string `json:"social_media_url"`
}

type SocialCreateResponse struct {
	Id             uint       `json:"id"`
	Name           string     `json:"name"`
	SocialMediaUrl string     `json:"social_media_url"`
	UserId         uint       `json:"user_id"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type SocialGetResponse struct {
	Socials []SocialData `json:"social_medias"`
}

type SocialData struct {
	Id             uint               `json:"id"`
	Name           string             `json:"name"`
	SocialMediaUrl string             `json:"social_media_url"`
	CreatedAt      *time.Time         `json:"created_at"`
	UpdatedAt      *time.Time         `json:"updated_at"`
	User           UserSocialResponse `json:"user"`
}

type UserSocialResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}
