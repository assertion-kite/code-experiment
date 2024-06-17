package main

import (
	"code/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserWeb struct {
	db *gorm.DB
}

func NewUserWeb(db *gorm.DB) *UserWeb {
	return &UserWeb{db: db}
}

func (u *UserWeb) RegisterUserRoute(en *gin.Engine) {
	g := en.Group("/user")
	g.POST("/home", u.Home)
}

type UserReq struct {
	Name *string `json:"name"`
}

func (u *UserWeb) Home(ctx *gin.Context) {
	var userReq UserReq
	if err := ctx.ShouldBindJSON(&userReq); err == nil {
		// user.Name 现在指向一个字符串（如果 JSON 中有 name 字段）
	}
	u.db.Model(&model.User{}).Create(&model.User{
		Name:     *userReq.Name,
		DateTime: time.Now(),
	})
	ctx.JSON(http.StatusOK, userReq)
}
