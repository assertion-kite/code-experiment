package model

type Order struct {
	Id        int64      `json:"id"`
	OrderNo   string     `json:"order_no"`
	OrderInfo *OrderInfo `gorm:"foreignKey:OrderId" json:"order_info"`
}

func (o *Order) TableName() string {
	return "order"
}

type OrderInfo struct {
	Id        int64  `json:"id"`
	OrderId   int64  `json:"order_id"`
	GoodsName string `json:"goods_name"`
}

func (i *OrderInfo) TableName() string {
	return "order_info"
}
