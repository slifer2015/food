package controllers

import (
	"context"
	"net/http"

	"strconv"

	"github.com/rs/xmux"
	"test.com/mine/modules/food/model"
	"test.com/mine/services/framework"
)

func getRestaurant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.ParseInt(xmux.Param(ctx, "id"), 10, 0)
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
	res, err := m.FindRestaurant(idInt)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	framework.JSON(w, http.StatusOK, res)

}
