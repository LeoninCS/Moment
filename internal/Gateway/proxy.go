package gateway

import (
	"log"
	"moment/internal/Gateway/loadbalancer"
	"moment/internal/config"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Proxy 代理管理器
type Proxy struct {
	loadbalancer *loadbalancer.LoadBalancer
}

// NewProxy 创建新的代理管理器
func NewProxy(configs []*config.RouteConfig) *Proxy {
	lb := loadbalancer.NewLoadBalancer()
	for _, routeConfig := range configs {
		for _, routeDetail := range routeConfig.Route {
			bk, err := loadbalancer.NewBackend(routeDetail.Path, routeDetail.Destination, routeDetail.Weight)
			if err != nil {
				log.Printf("创建后端服务失败: %v", err)
				continue
			}
			lb.AddBackend(bk)
		}
	}
	return &Proxy{
		loadbalancer: lb,
	}
}

// ServeHTTP 处理HTTP请求并转发到后端服务
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request, route *config.RouteDetail) {
	lb := p.loadbalancer
	// 选择后端服务
	backend := lb.GetNextTarget()
	if backend == nil {
		http.Error(w, "无可用后端服务", http.StatusServiceUnavailable)
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)

	// 保存原始路径用于日志记录
	originalPath := r.URL.Path

	// 设置反向代理的Director
	proxy.Director = func(req *http.Request) {
		// 设置目标服务的协议和主机
		req.URL.Scheme = backend.URL.Scheme
		req.URL.Host = backend.URL.Host
		// 处理路径映射
		req.URL.Path = strings.TrimPrefix(req.URL.Path, route.Path)
		// 保留原始请求头
		req.Header = r.Header
	}

	// 设置错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("代理错误: %v", err)
		http.Error(w, "后端服务错误", http.StatusBadGateway)
	}

	// 记录转发日志
	log.Printf("转发请求: %s %s -> %s://%s%s", r.Method, originalPath, backend.URL.Scheme, backend.URL.Host, r.URL.Path)

	// 转发请求
	proxy.ServeHTTP(w, r)
}
