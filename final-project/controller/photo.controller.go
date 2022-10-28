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

type PhotoController struct {
	db *gorm.DB
}

func NewPhotoController(db *gorm.DB) *PhotoController {
	return &PhotoController{
		db: db,
	}
}

func (controller *PhotoController) CreatePhoto(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoRequest := dto.PhotoRequest{}

	err := ctx.ShouldBindJSON(&photoRequest)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	photo := entity.Photo{
		Title:    photoRequest.Title,
		Caption:  photoRequest.Caption,
		PhotoUrl: photoRequest.PhotoUrl,
		UserId:   uint(userId.(float64)),
	}

	_, err = govalidator.ValidateStruct(&photo)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Create(&photo).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, err.Error())
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helper.WriteJsonResponse(ctx, http.StatusCreated, dto.PhotoCreateResponse{
		Id:        photo.Id,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserId:    photo.UserId,
		CreatedAt: photo.CreatedAt,
	})
}

func (controller *PhotoController) FindAllPhoto(ctx *gin.Context) {
	var photos []entity.Photo

	err := controller.db.Preload("User").Find(&photos).Error

	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, err.Error())
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	var response dto.PhotoGetResponse
	for _, photo := range photos {
		var userData dto.UserPhotoResponse
		if photo.User != nil {
			userData = dto.UserPhotoResponse{
				Username: photo.User.Username,
				Email:    photo.User.Email,
			}
		}
		photosResponse := dto.PhotoData{
			Id:        photo.Id,
			Title:     photo.Title,
			Caption:   photo.Caption,
			PhotoUrl:  photo.PhotoUrl,
			CreatedAt: photo.CreatedAt,
			UpdatedAt: photo.UpdatedAt,
			User:      userData,
		}
		response.Photos = append(response.Photos, photosResponse)

	}
	helper.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (controller *PhotoController) UpdatePhoto(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoMediaId := ctx.Param("photoId")
	var photoRequest dto.PhotoRequest
	var photo entity.Photo

	err := ctx.ShouldBindJSON(&photoRequest)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	updatePhoto := entity.Photo{
		Title:    photoRequest.Title,
		Caption:  photoRequest.Caption,
		PhotoUrl: photoRequest.PhotoUrl,
	}

	err = controller.db.First(&photo, photoMediaId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, "data not found")
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if photo.UserId != uint(userId.(float64)) {
		helper.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to update or edit this photo",
		})
		return
	}

	err = controller.db.Model(&photo).Updates(updatePhoto).Error
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	response := dto.PhotoCreateResponse{
		Id:        photo.Id,
		Title:     photo.Title,
		Caption:   photo.Caption,
		PhotoUrl:  photo.PhotoUrl,
		UserId:    photo.UserId,
		UpdatedAt: photo.UpdatedAt,
	}

	helper.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (controller *PhotoController) DeletePhoto(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	photoId := ctx.Param("photoId")
	var photo entity.Photo

	err := controller.db.First(&photo, photoId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, "data not found")
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if photo.UserId != uint(userId.(float64)) {
		helper.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to delete this photo",
		})
		return
	}

	err = controller.db.Delete(&photo).Error
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
		"message": "Your photo has been successfully deleted",
	})
}
