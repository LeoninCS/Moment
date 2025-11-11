package loadbalancer

import (
	"net/url"
	"sync"
)

type Backend struct {
	Name   string
	URL    *url.URL
	mu     sync.Mutex
	Alive  bool
	Weight int
}

func NewBackend(name string, destination string, weight int) (*Backend, error) {
	url, err := url.Parse(destination)
	if err != nil {
		return nil, err
	}
	return &Backend{
		Name:   name,
		URL:    url,
		Weight: weight,
		Alive:  true,
		mu:     sync.Mutex{},
	}, nil
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Alive
}
