package dto

import "time"

type CommentRequest struct {
	Message string `json:"message"`
	PhotoId uint   `json:"photo_id"`
}

type CommentCreateResponse struct {
	Id        uint       `json:"id"`
	Message   string     `json:"message"`
	PhotoId   uint       `json:"photo_id"`
	UserId    uint       `json:"user_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CommentGetResponse struct {
	Comments []CommentData `json:"photos"`
}

type CommentData struct {
	Id        uint                 `json:"id"`
	Message   string               `json:"message"`
	PhotoId   uint                 `json:"photo_id"`
	User      UserCommentResponse  `json:"user"`
	Photos    PhotoCommentResponse `json:"photos"`
	CreatedAt *time.Time           `json:"created_at"`
	UpdatedAt *time.Time           `json:"updated_at"`
}

type UserCommentResponse struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type PhotoCommentResponse struct {
	Id       uint   `json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url"`
	UserId   uint   `json:"user_id"`
}

type CommentUpdateRequest struct {
	Message string `json:"message"`
}
