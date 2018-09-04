package controllers

import (
	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"test.com/mine/services/framework"
)

type ctrl struct {
	framework.Base
}

func (*ctrl) Routes(m *xmux.Mux) {
	m.POST("/api/restaurant", xhandler.HandlerFuncC(createRestaurant))
	m.GET("/api/restaurant/:id", xhandler.HandlerFuncC(getRestaurant))
	m.GET("/api/restaurant", xhandler.HandlerFuncC(listRestaurant))
	m.PUT("/api/restaurant/:id", xhandler.HandlerFuncC(editRestaurant))
	m.DELETE("/api/restaurant/:id", xhandler.HandlerFuncC(deleteRestaurant))

	m.POST("/api/food/:id", xhandler.HandlerFuncC(addFood))
	m.PUT("/api/food/:id", xhandler.HandlerFuncC(editFood))


	m.POST("/api/order", xhandler.HandlerFuncC(setOrder))
}

func init() {
	framework.Register(&ctrl{})
}
