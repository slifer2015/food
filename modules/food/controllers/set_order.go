package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"fmt"
	"strings"
	"time"

	"test.com/mine/modules/food/model"
	"test.com/mine/services/framework"
)

type setOrderPayload struct {
	FoodIDs      []int64 `json:"food_ids"`
	RestaurantID int64   `json:"restaurant_id"`
}

type OrderFoods struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

type setOrderResponse struct {
	RestaurantID   int64  `json:"restaurant_id"`
	RestaurantName string `json:"restaurant_name"`

	Foods []OrderFoods `json:"foods"`

	Total int64 `json:"total"`
}

func setOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}
	var income setOrderPayload
	err = json.Unmarshal(data, &income)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}
	m := model.NewFoodManager()

	// find restaurant by id
	restaurant, err := m.FindRestaurant(income.RestaurantID)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	// find ordered food by ids and restaurant id
	foods := m.FindFoods(restaurant.ID, income.FoodIDs)
	if len(foods) == 0 {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: "no food found",
		})
		return
	}

	var orderFoods = make([]OrderFoods, 0)
	var foodIDs = make([]string, 0)
	for i := range foods {
		foodIDs = append(foodIDs, fmt.Sprint(foods[i].ID))
		orderFoods = append(orderFoods, OrderFoods{
			Title: foods[i].Title,
			Price: int(foods[i].Price),
		})
	}

	total := calculateTotalPrice(restaurant, foods)
	var order = &model.Order{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Price:        total,
		RestaurantID: restaurant.ID,
		FoodIDs:      strings.Join(foodIDs, ","),
	}
	err = m.CreateOrder(order)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	var finalRes = setOrderResponse{
		RestaurantID:   restaurant.ID,
		RestaurantName: restaurant.Name,
		Total:          total,
		Foods:          orderFoods,
	}

	framework.JSON(w, http.StatusOK, finalRes)

}

func calculateTotalPrice(restaurant *model.Restaurant, foods []*model.Food) int64 {
	var total int64
	for i := range foods {
		total += foods[i].Price
	}
	rawPrice := total + int64(restaurant.Send) + int64(restaurant.Packing)
	totalFloat := float64(rawPrice) * float64(1+(float64(restaurant.Tax)/float64(100)))
	return int64(totalFloat)
}


