package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"strconv"

	"github.com/rs/xmux"
	"test.com/mine/modules/food/model"
	"test.com/mine/services/framework"
)

type editFoodPayload struct {
	Title string `json:"title"`
	Price int64  `json:"price"`
}

func editFood(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idInt, err := strconv.ParseInt(xmux.Param(ctx, "id"), 10, 0)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	// find restaurant
	m := model.NewFoodManager()
	// find restaurant by id
	res, err := m.FindFood(idInt)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}
	var income editFoodPayload
	err = json.Unmarshal(data, &income)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	res.Price = income.Price
	res.Title = income.Title

	err = m.UpdateFood(res)
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