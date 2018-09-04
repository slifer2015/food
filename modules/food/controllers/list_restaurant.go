package controllers

import (
	"context"
	"net/http"

	"test.com/mine/modules/food/model"
	"test.com/mine/services/framework"
)

func listRestaurant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	m := model.NewFoodManager()
	// find restaurant by id
	res := m.ListRestaurant()
	framework.JSON(w, http.StatusOK, res)

}
