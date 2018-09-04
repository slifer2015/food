package model

import (
	"time"

	"test.com/mine/services/mysql"
)

type tableInit struct {
}

type Manager struct {
	mysql.Manager
}

const (
	RestaurantTableName = "restaurants"
	FoodTableName       = "foods"
	OrderTableName      = "orders"
)

// Initialize orm package
func (m *Manager) Initialize() {
	m.GetConn().AddTableWithName(Restaurant{}, RestaurantTableName).SetKeys(true, "ID")
	m.GetConn().AddTableWithName(Food{}, FoodTableName).SetKeys(true, "ID")
	m.GetConn().AddTableWithName(Order{}, OrderTableName).SetKeys(true, "ID")

}
func init() {
	mysql.Register(NewFoodManager())
}

type Restaurant struct {
	ID        int64     `json:"id" db:"id, primarykey, autoincrement"`
	Name      string    `json:"name" db:"name"`
	Tax       int       `json:"tax" db:"tax"`
	Send      int       `json:"send" db:"send"`
	Packing   int       `json:"packing" db:"packing"`
	Locations string    `json:"locations" db:"locations"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Food struct {
	ID           int64  `json:"id" db:"id, primarykey, autoincrement"`
	RestaurantID int64  `json:"resturant_id" db:"resturant_id"`
	Title        string `json:"title" db:"title"`
	Price        int64  `json:"price" db:"price"`
}

type Order struct {
	ID           int64     `json:"id" db:"id, primarykey, autoincrement"`
	FoodIDs      string    `json:"food_ids" db:"food_ids"`
	RestaurantID int64     `json:"restaurant_id" db:"restaurant_id"`
	Price        int64     `json:"price" db:"price"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func NewFoodManager() *Manager {
	return &Manager{}
}
