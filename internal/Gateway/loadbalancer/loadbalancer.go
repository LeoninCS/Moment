package loadbalancer

import (
	"math"
	"sync"
)

type LoadBalancer struct {
	baclend   []*Backend
	mutex     sync.Mutex
	curweight map[*Backend]int
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		baclend:   make([]*Backend, 0),
		curweight: make(map[*Backend]int),
	}
}
func (lb *LoadBalancer) AddBackend(backend *Backend) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	lb.baclend = append(lb.baclend, backend)
}

func (lb *LoadBalancer) RemoveBackend(backend *Backend) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	for i, b := range lb.baclend {
		if b.Name == backend.Name {
			lb.baclend = append(lb.baclend[:i], lb.baclend[i+1:]...)
			return
		}
	}
}
func (lb *LoadBalancer) GetBackends() []*Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	return lb.baclend
}

func (lb *LoadBalancer) GetNextTarget() *Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	if len(lb.baclend) == 0 {
		return nil
	}
	idx, maxWeight, sumWeight := -1, math.MinInt32, 0
	for i, backend := range lb.baclend {
		if !backend.IsAlive() {
			continue
		}
		lb.curweight[backend] += backend.Weight
		sumWeight += backend.Weight
		if lb.curweight[backend] > maxWeight {
			maxWeight = lb.curweight[backend]
			idx = i
		}
	}
	if idx == -1 {
		return nil
	}
	lb.curweight[lb.baclend[idx]] -= sumWeight
	return lb.baclend[idx]
}
