package routers

import (
	"final-project/common/database"
	"final-project/common/middleware"
	"final-project/controller"
	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {

	db := database.ConnectDB()
	router := gin.Default()
	user := controller.NewUserController(db)
	social := controller.NewSocialController(db)
	photo := controller.NewPhotoController(db)
	comment := controller.NewCommentController(db)

	userGroup := router.Group("/users")
	{
		userGroup.POST("/login", user.UserLogin)
		userGroup.POST("/register", user.CreateUser)
		userGroup.PUT("/", middleware.Auth(), user.UpdateUser)
		userGroup.DELETE("/", middleware.Auth(), user.DeleteUser)
	}

	socialGroup := router.Group("/socials")
	{
		socialGroup.GET("/", middleware.Auth(), social.FindAllSocial)
		socialGroup.POST("/", middleware.Auth(), social.CreateSocial)
		socialGroup.PUT("/:socialMediaId", middleware.Auth(), social.UpdateSocial)
		socialGroup.DELETE("/:socialMediaId", middleware.Auth(), social.DeleteSocial)
	}

	photoGroup := router.Group("/photos")
	{
		photoGroup.GET("/", middleware.Auth(), photo.FindAllPhoto)
		photoGroup.POST("/", middleware.Auth(), photo.CreatePhoto)
		photoGroup.PUT("/:photoId", middleware.Auth(), photo.UpdatePhoto)
		photoGroup.DELETE("/:socialMediaId", middleware.Auth(), photo.DeletePhoto)
	}

	commentGroup := router.Group("/comments")
	{
		commentGroup.GET("/", middleware.Auth(), comment.FindAllComment)
		commentGroup.POST("/", middleware.Auth(), comment.CreateComment)
		commentGroup.PUT("/:commentId", middleware.Auth(), comment.UpdateComment)
		commentGroup.DELETE("/:commentId", middleware.Auth(), comment.DeleteComment)
	}

	//router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
