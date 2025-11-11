package gateway

import (
	"log"
	"moment/internal/config"
	"net/http"
)

// Gateway API网关
type Gateway struct {
	config *config.Config
	router *Router
	proxy  *Proxy
}

// NewGateway 创建新的API网关
func NewGateway(config *config.Config) *Gateway {
	return &Gateway{
		config: config,
		router: NewRouter(config.Routes),
		proxy:  NewProxy(config.Routes),
	}
}

// ServeHTTP 实现HTTP处理器接口
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 查找匹配的路由详情和配置
	routeDetail, routeConfig := g.router.GetRouteDetailByPath(r.URL.Path)
	if routeDetail == nil || routeConfig == nil {
		log.Printf("未找到匹配的路由: %s", r.URL.Path)
		http.Error(w, "路由未找到", http.StatusNotFound)
		return
	}

	// 转发请求
	g.proxy.ServeHTTP(w, r, routeDetail)
}

// Start 启动网关服务
func (g *Gateway) Start() error {
	port := g.config.Server.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("API网关启动在端口 %s", port)
	return http.ListenAndServe(":"+port, g)
}
