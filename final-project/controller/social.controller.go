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

type SocialController struct {
	db *gorm.DB
}

func NewSocialController(db *gorm.DB) *SocialController {
	return &SocialController{
		db: db,
	}
}

func (controller *SocialController) CreateSocial(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialRequest := dto.SocialRequest{}

	err := ctx.ShouldBindJSON(&socialRequest)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	social := entity.Social{
		Name:           socialRequest.Name,
		SocialMediaUrl: socialRequest.SocialMediaUrl,
		UserId:         uint(userId.(float64)),
	}

	_, err = govalidator.ValidateStruct(&social)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	err = controller.db.Create(&social).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, err.Error())
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	helper.WriteJsonResponse(ctx, http.StatusCreated, dto.SocialCreateResponse{
		Id:             social.Id,
		Name:           social.Name,
		SocialMediaUrl: social.SocialMediaUrl,
		UserId:         social.UserId,
		CreatedAt:      social.CreatedAt,
	})
}

func (controller *SocialController) FindAllSocial(ctx *gin.Context) {
	var socials []entity.Social

	err := controller.db.Preload("User").Find(&socials).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, err.Error())
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	var response dto.SocialGetResponse
	for _, social := range socials {
		var userData dto.UserSocialResponse
		if social.User != nil {
			userData = dto.UserSocialResponse{
				Id:       social.User.Id,
				Username: social.User.Username,
			}
		}
		socialMediasResponse := dto.SocialData{
			Id:             social.Id,
			Name:           social.Name,
			SocialMediaUrl: social.SocialMediaUrl,
			CreatedAt:      social.CreatedAt,
			UpdatedAt:      social.UpdatedAt,
			User:           userData,
		}
		response.Socials = append(response.Socials, socialMediasResponse)

	}
	helper.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (controller *SocialController) UpdateSocial(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialMediaId := ctx.Param("socialMediaId")
	var socialRequest dto.SocialRequest
	var social entity.Social

	err := ctx.ShouldBindJSON(&socialRequest)
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	updatedSocial := entity.Social{
		Name:           socialRequest.Name,
		SocialMediaUrl: socialRequest.SocialMediaUrl,
		UserId:         uint(userId.(float64)),
	}

	err = controller.db.First(&social, socialMediaId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, "data not found")
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if social.UserId != uint(userId.(float64)) {
		helper.WriteJsonResponse(ctx, http.StatusUnauthorized, "you're not allowed to update or edit this social media")
		return
	}

	err = controller.db.Model(&social).Updates(updatedSocial).Error
	if err != nil {
		helper.BadRequestResponse(ctx, err.Error())
		return
	}

	response := dto.SocialCreateResponse{
		Id:             social.Id,
		Name:           social.Name,
		SocialMediaUrl: social.SocialMediaUrl,
		UserId:         social.UserId,
		UpdatedAt:      social.UpdatedAt,
	}

	helper.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (controller *SocialController) DeleteSocial(ctx *gin.Context) {
	userId, _ := ctx.Get("id")
	socialMediaId := ctx.Param("socialMediaId")
	var social entity.Social

	err := controller.db.First(&social, socialMediaId).Error
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			helper.NotFoundResponse(ctx, "data not found")
			return
		}
		helper.InternalServerJsonResponse(ctx, err.Error())
		return
	}

	if social.UserId != uint(userId.(float64)) {
		helper.WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "you're not allowed to delete this social media",
		})
		return
	}

	err = controller.db.Delete(&social).Error
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
		"message": "Your social media has been successfully deleted",
	})
}
