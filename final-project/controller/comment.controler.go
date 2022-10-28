package controller

import (
	"final-project/common/helper"
	"final-project/dto"
	"final-project/entity"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type CommentController struct {
	db *gorm.DB
}

func NewCommentController(db *gorm.DB) *CommentController {
	return &CommentController{
		db: db,
	}
}

func (controller *CommentController) CreateComment(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentRequest := dto.CommentRequest{}

	err := ctx.ShouldBindJSON(&commentRequest)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	comment := entity.Comment{
		Message: commentRequest.Message,
		PhotoId: commentRequest.PhotoId,
		UserId:  uint(userId.(float64)),
	}

	_, err = govalidator.ValidateStruct(&comment)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Create(&comment).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, err.Error())
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helper.WriteJsonResponse(ctx, http.StatusCreated, dto.CommentCreateResponse{
		Id:        comment.Id,
		Message:   comment.Message,
		PhotoId:   comment.PhotoId,
		UserId:    comment.UserId,
		CreatedAt: comment.CreatedAt,
	})
}

func (controller *CommentController) FindAllComment(ctx *gin.Context) {
	var comments []entity.Comment

	err := controller.db.Preload("User").Preload("Photo").Find(&comments).Error

	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, err.Error())
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	var response dto.CommentGetResponse
	for _, comment := range comments {
		var userData dto.UserCommentResponse
		if comment.User != nil {
			userData = dto.UserCommentResponse{
				Id:       comment.User.Id,
				Username: comment.User.Username,
				Email:    comment.User.Email,
			}
		}

		var photoData dto.PhotoCommentResponse
		if comment.Photo != nil {
			photoData = dto.PhotoCommentResponse{
				Id:       comment.Photo.Id,
				Title:    comment.Photo.Title,
				Caption:  comment.Photo.Caption,
				PhotoUrl: comment.Photo.PhotoUrl,
				UserId:   comment.Photo.UserId,
			}
		}
		commentsResponse := dto.CommentData{
			Id:        comment.Id,
			Message:   comment.Message,
			PhotoId:   comment.PhotoId,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			User:      userData,
			Photos:    photoData,
		}
		response.Comments = append(response.Comments, commentsResponse)

	}
	helper.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (controller *CommentController) UpdateComment(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentId := ctx.Param("commentId")
	var commentUpdateRequest dto.CommentUpdateRequest
	var comment entity.Comment

	err := ctx.ShouldBindJSON(&commentUpdateRequest)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	updatePhoto := entity.Comment{
		Message: commentUpdateRequest.Message,
	}

	err = controller.db.First(&comment, commentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, "data not found")
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if comment.UserId != uint(userId.(float64)) {
		helper.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to update or edit this comment",
		})
		return
	}

	err = controller.db.Model(&comment).Updates(updatePhoto).Error
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	response := dto.CommentCreateResponse{
		Id:        comment.Id,
		Message:   comment.Message,
		PhotoId:   comment.PhotoId,
		UserId:    comment.UserId,
		UpdatedAt: comment.UpdatedAt,
	}

	helper.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (controller *CommentController) DeleteComment(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	commentId := ctx.Param("commentId")
	var comment entity.Comment

	err := controller.db.First(&comment, commentId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, "data not found")
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if comment.UserId != uint(userId.(float64)) {
		helper.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to delete this comment",
		})
		return
	}

	err = controller.db.Delete(&comment).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, err.Error())
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helper.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"error":   false,
		"message": "Your comment has been successfully deleted",
	})
}
