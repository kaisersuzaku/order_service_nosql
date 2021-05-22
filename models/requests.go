package models

type OrderProductReq struct {
	ID  string `json:"id" valid:"required"`
	Qty int    `json:"qty" valid:"required"`
}
