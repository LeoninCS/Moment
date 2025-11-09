package gateway

import (
	"moment/internal/config"
	"strings"
)

// Router 路由管理器
type Router struct {
	routes []*config.RouteConfig
}

// NewRouter 创建新的路由器
func NewRouter(routes []*config.RouteConfig) *Router {
	return &Router{
		routes: routes,
	}
}

// FindRoute 根据路径查找匹配的路由配置
func (r *Router) FindRoute(path string) *config.RouteDetail {
	// 遍历查找匹配的路由
	for _, route := range r.routes {
		// 检查路径是否匹配
		for _, detail := range route.Route {
			if strings.HasPrefix(path, detail.Path) {
				return detail
			}
		}
	}
	return nil
}

// GetRouteDetailByPath 根据路径获取具体的路由详情
func (r *Router) GetRouteDetailByPath(path string) (*config.RouteDetail, *config.RouteConfig) {
	for _, route := range r.routes {
		for _, detail := range route.Route {
			if strings.HasPrefix(path, detail.Path) {
				return detail, route
			}
		}
	}
	return nil, nil
}
