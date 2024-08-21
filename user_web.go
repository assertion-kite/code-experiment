package main

import (
	"code/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
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
	g.POST("/test", u.Test)
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

type Req struct {
	CorpWxid string `json:"corp_wxid"`
	RelWxid  string `json:"rel_wxid"`
	Type     string `json:"type"`
	Data     []struct {
		ClientId string `json:"client_id"`
		Content  struct {
			Items []struct {
				Text string `json:"text"`
				Type string `json:"type"`
			} `json:"items"`
		} `json:"content"`
		Extra struct {
			CallbackType       int    `json:"callback_type"`
			ClientCallbackType int    `json:"client_callback_type"`
			Index              int    `json:"index"`
			Module             string `json:"module"`
			SenderWxid         string `json:"sender_wxid"`
			Source             int    `json:"source"`
			SubTaskId          string `json:"sub_task_id"`
			SubUserId          string `json:"sub_user_id"`
			TargetId           string `json:"target_id"`
			TargetWxid         string `json:"target_wxid"`
			TaskId             string `json:"task_id"`
			UserId             string `json:"user_id"`
		} `json:"extra"`
		FromWxid  string `json:"from_wxid"`
		MsgHash   string `json:"msg_hash"`
		MsgTime   int64  `json:"msg_time"`
		MsgType   string `json:"msg_type"`
		ToWxid    string `json:"to_wxid"`
		WxMsgid   string `json:"wx_msgid"`
		WxSvrid   string `json:"wx_svrid"`
		UserId    string `json:"user_id"`
		SubUserId int    `json:"sub_user_id"`
		CorpId    int    `json:"corp_id"`
		Wxid      string `json:"wxid"`
		Svrid     string `json:"svrid"`
		ChatWxid  string `json:"chat_wxid"`
		ChatType  int    `json:"chat_type"`
	} `json:"data"`
}

func (u *UserWeb) Test(ctx *gin.Context) {
	// 读取POST请求的body
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		http.Error(ctx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	// 记得关闭Body，因为它实现了io.ReadCloser接口
	defer ctx.Request.Body.Close()

	// 你可以将bodyBytes转换为string，以便于查看或处理
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	var req []Req
	if err := ctx.BindJSON(&req); err == nil {
	}
	ctx.JSON(http.StatusOK, req)
}
