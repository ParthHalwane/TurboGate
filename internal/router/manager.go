package router

import (
	"net/http"
	"sync/atomic"
)

type RouterManager struct {
	handler atomic.Value
}

func NewRouterManager(initial http.Handler) *RouterManager {
	m := &RouterManager{}
	m.handler.Store(initial)
	return m
}

func (m *RouterManager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := m.handler.Load().(http.Handler)
	handler.ServeHTTP(w, r)
}

func (m *RouterManager) UpdateHandler(newHandler http.Handler) {
	m.handler.Store(newHandler)
}
