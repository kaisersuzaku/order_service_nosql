package models

type Product struct {
	ID    string `json:"id"`
	Code  string `json:"code"`
	Stock int    `json:"stock"`
}
