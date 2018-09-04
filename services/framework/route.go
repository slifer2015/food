package framework

import (
	"context"
	"net/http"
	"sync"

	"log"

	"encoding/json"

	"github.com/rs/xhandler"
	"github.com/rs/xmux"
	"test.com/mine/services/initializer"
)

var (
	all    []Router
	once   = sync.Once{}
	engine *xmux.Mux
)

// Base controller struct
type Base struct {
}

// Router type for route registration
type Router interface {
	Routes(*xmux.Mux)
}

// Register routes
func Register(r Router) {
	all = append(all, r)
}

type router struct {
}

func (router) Initial(ctx context.Context) {
	once.Do(func() {
		engine = xmux.New()
		// add registered route
		for i := range all {
			all[i].Routes(engine)
		}
		log.Fatal(http.ListenAndServe(":8080", xhandler.New(context.Background(), engine)))
	})
}

func init() {
	initializer.Register(router{}, 1000)
}

func JSON(w http.ResponseWriter, code int, i interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	_ = enc.Encode(i)
}
