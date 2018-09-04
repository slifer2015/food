package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"time"

	"test.com/mine/modules/food/model"
	"test.com/mine/services/framework"
)

type createRestaurantPayload struct {
	Name      string   `json:"name"`
	Locations []string `json:"locations"`
	Tax       int      `json:"tax"`
	Packing   int      `json:"packing"`
	Send      int      `json:"send"`
}

func createRestaurant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}
	var income createRestaurantPayload
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
	var loc = ""
	if len(income.Locations) > 0 {

		loc = strings.Join(income.Locations, ",")
	}

	var restaurant = &model.Restaurant{
		Name:      income.Name,
		Packing:   income.Packing,
		Tax:       income.Tax,
		Send:      income.Send,
		Locations: loc,
		CreatedAt: time.Now(),
	}
	err = m.CreateRestaurant(restaurant)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	framework.JSON(w, http.StatusOK, restaurant)

}
