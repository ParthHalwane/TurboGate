package router

import (
	"net/http"
	"sync/atomic"
)

type RouterManager struct {
	current atomic.Value // holds http.Handler
}

func NewRouterManager(handler http.Handler) *RouterManager {
	rm := &RouterManager{}
	rm.current.Store(handler)
	return rm
}

func (rm *RouterManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := rm.current.Load().(http.Handler)
	handler.ServeHTTP(w, r)
}

func (rm *RouterManager) UpdateHandler(newHandler http.Handler) {
	rm.current.Store(newHandler)
}
