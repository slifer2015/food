package controllers

import (
	"context"
	"net/http"
	"strconv"

	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/rs/xmux"
	"test.com/mine/modules/food/model"
	"test.com/mine/services/framework"
)

type editRestaurantPayload struct {
	Name      string   `json:"name"`
	Locations []string `json:"locations"`
	Tax       int      `json:"tax"`
	Packing   int      `json:"packing"`
	Send      int      `json:"send"`
}

func editRestaurant(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}
	var income editRestaurantPayload
	err = json.Unmarshal(data, &income)
	if err != nil {
		framework.JSON(w, http.StatusBadRequest, struct {
			Err string `json:"err"`
		}{
			Err: err.Error(),
		})
		return
	}

	var loc = ""
	if len(income.Locations) > 0 {

		loc = strings.Join(income.Locations, ",")
	}

	res.Name = income.Name
	res.Tax = income.Tax
	res.Send = income.Send
	res.Packing = income.Packing
	res.Locations = loc

	err = m.UpdateRestaurant(res)
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
