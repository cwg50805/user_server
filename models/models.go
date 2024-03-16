package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RecommendationItem struct {
	ItemName string  `json:"item_name"`
	Price    float64 `json:"price"`
}
