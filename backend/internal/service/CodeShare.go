package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

const (
	MaxEntries     = 1000
	MaxContentSize = 100000
)

type Code struct {
	Author      string `json:"author"`
	Language    string `json:"language"`
	Content     string `json:"content"`
	Hash        string `json:"hash"`
	DestroyTime int64  `json:"destroy_time"`
}

type CodeShare struct {
	mu    sync.RWMutex
	store map[string]*Code
	order []string
}

func NewCodeShareService() *CodeShare {

	cs := &CodeShare{
		store: make(map[string]*Code),
		order: make([]string, 0, MaxEntries),
	}

	go func() {
		for range time.Tick(3 * time.Hour) {
			cs.cleanExpired()
		}
	}()

	return cs
}

// 上传
func (cs *CodeShare) Upload(author, lang, content string, ttl int64) *Code {
	if len(content) > MaxContentSize {
		return nil
	}

	destroy := time.Now().Unix() + ttl
	data := fmt.Sprintf("%s|%s|%s|%d", author, lang, content, destroy)
	hashBytes := sha256.Sum256([]byte(data))
	hash := hex.EncodeToString(hashBytes[:])

	code := &Code{
		Author:      author,
		Language:    lang,
		Content:     content,
		Hash:        hash,
		DestroyTime: destroy,
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// 内存存储
	if len(cs.store) >= MaxEntries {
		oldest := cs.order[0]
		delete(cs.store, oldest)
		cs.order = cs.order[1:]
	}
	cs.store[hash] = code
	cs.order = append(cs.order, hash)

	return code
}

// 获取
func (cs *CodeShare) Get(hash string) (*Code, bool) {
	cs.mu.RLock()
	code, ok := cs.store[hash]
	cs.mu.RUnlock()

	if ok && time.Now().Unix() <= code.DestroyTime {
		return code, true
	}

	return nil, false
}

// 清理过期的
func (cs *CodeShare) cleanExpired() {
	now := time.Now().Unix()
	cs.mu.Lock()
	defer cs.mu.Unlock()

	newOrder := cs.order[:0]
	for _, h := range cs.order {
		c := cs.store[h]
		if c == nil || c.DestroyTime < now {
			delete(cs.store, h)
			continue
		}
		newOrder = append(newOrder, h)
	}
	cs.order = newOrder
}
