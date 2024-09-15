package controllers

import (
	"errors"
	"net/http"

	"github.com/Ferdinand-work/PetalPix/models"
	"github.com/Ferdinand-work/PetalPix/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func New(userService services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if user.UserId == "" || user.Name == "" || user.Email == "" || user.Password == "" || user.ContactNo == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errors.New("invalid Creds")})
		return
	}
	if user.Following == nil {
		user.Following = make([]string, 0)
	}
	if user.Followers == nil {
		user.Followers = make([]string, 0)
	}
	msg, err := uc.UserService.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": msg})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	userId := ctx.Params.ByName("id")
	user, err := uc.UserService.GetUser(&userId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if user.UserId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errors.New("ID cannot be empty")})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Params.ByName("id")
	err := uc.UserService.DeleteUser(&userId)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (uc *UserController) Follow(ctx *gin.Context) {
	userId := ctx.Param("id")
	var followUsers models.FollowRequest
	if err := ctx.ShouldBindJSON(&followUsers); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var usersToFollow interface{}
	if followUsers.FollowUser != "" {
		usersToFollow = followUsers.FollowUser
	} else {
		usersToFollow = followUsers.FollowUsers
	}

	resFollow, err := uc.UserService.Follow(usersToFollow, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
	ctx.JSON(http.StatusOK, *resFollow)
}

func (uc *UserController) GetFollowing(ctx *gin.Context) {
	id := ctx.Param("id")
	followList, err := uc.UserService.GetFollowing(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, *followList)
}

func (uc *UserController) Unfollow(ctx *gin.Context) {
	userId := ctx.Param("id")
	var unfollowUsers models.UnfollowRequest
	if err := ctx.ShouldBindJSON(&unfollowUsers); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var usersToUnfollow interface{}
	if unfollowUsers.UnfollowUser != "" {
		usersToUnfollow = unfollowUsers.UnfollowUser
	} else {
		usersToUnfollow = unfollowUsers.UnfollowUsers
	}

	resFollow, err := uc.UserService.Unfollow(usersToUnfollow, userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
	}
	ctx.JSON(http.StatusOK, *resFollow)
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userRoute := rg.Group("/user")
	userRoute.POST("/create", uc.CreateUser)
	userRoute.GET("/get/:id", uc.GetUser)
	userRoute.GET("/getall", uc.GetAll)
	userRoute.PATCH("/update", uc.UpdateUser)
	userRoute.DELETE("/delete/:id", uc.DeleteUser)
	userRoute.POST("/follow/:id", uc.Follow)
	userRoute.GET("/getFollowing/:id", uc.GetFollowing)
	userRoute.POST("/unfollow/:id", uc.Unfollow)
}
