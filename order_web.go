package main

import (
	"code/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type OrderWeb struct {
	db *gorm.DB
}

func NewOrderWeb(db *gorm.DB) *OrderWeb {
	return &OrderWeb{db: db}
}

func (u *OrderWeb) RegisterOrderRoute(en *gin.Engine) {
	g := en.Group("/order")
	g.POST("/create", u.Create)
	g.POST("/update", u.Update)
}

type OrderReq struct {
	Id        int64  `json:"id"`
	OrderNo   string `json:"order_no"`
	GoodsName string `json:"goods_name"`
}

func (u *OrderWeb) Create(ctx *gin.Context) {
	var req OrderReq
	if err := ctx.ShouldBindJSON(&req); err == nil {
		// user.Name 现在指向一个字符串（如果 JSON 中有 name 字段）
	}
	u.db.Create(&model.Order{
		OrderNo: req.OrderNo,
		OrderInfo: &model.OrderInfo{
			GoodsName: req.GoodsName,
		},
	})
	ctx.JSON(http.StatusOK, nil)
}

func (u *OrderWeb) Update(ctx *gin.Context) {
	var req OrderReq
	if err := ctx.ShouldBindJSON(&req); err == nil {
		// user.Name 现在指向一个字符串（如果 JSON 中有 name 字段）
	}
	u.db.Updates(&model.Order{
		Id:      req.Id,
		OrderNo: req.OrderNo,
		OrderInfo: &model.OrderInfo{
			Id:        1,
			OrderId:   req.Id,
			GoodsName: req.GoodsName,
		},
	})
	u.db.Select("order_no", "goods_name").Updates(&model.OrderInfo{
		Id:        1,
		OrderId:   1,
		GoodsName: req.GoodsName,
	})
	ctx.JSON(http.StatusOK, nil)
}
