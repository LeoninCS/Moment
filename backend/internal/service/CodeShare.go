// 并发安全 + TTL + 基于 FIFO slice 实现的LRU
package service

import (
	"sync"
	"time"
)

const (
	MaxEntries     = 1000
	MaxContentSize = 100000
)

type CodeShare struct {
	mu    sync.RWMutex
	store map[string]*Code
	order []string
}

type Code struct {
	Author      string `json:"author"`
	Language    string `json:"language"`
	Content     string `json:"content"`
	Hash        string `json:"hash"`
	DestroyTime int64  `json:"destroy_time"`
}

func NewCodeShareService() *CodeShare {

	cs := &CodeShare{
		store: make(map[string]*Code),
		order: make([]string, 0, MaxEntries),
	}

	go func() {
		for range time.Tick(3 * time.Minute) {
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
	/*
		目前用户较少，暂时不使用复杂的哈希算法，因为后缀过长
		data := fmt.Sprintf("%s|%s|%s|%d", author, lang, content, destroy)
		hashBytes := sha256.Sum256([]byte(data))
		hash := hex.EncodeToString(hashBytes[:])
	*/
	// 简单哈希

	hash := GetHash(10)
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
	cs.mu.Lock()
	code, ok := cs.store[hash]

	if ok {
		// 1. 从 order 中移除旧位置
		for i := range cs.order {
			if cs.order[i] == hash {
				cs.order = append(cs.order[:i], cs.order[i+1:]...)
				break
			}
		}
		// 2. 放到末尾（表示最近使用）
		cs.order = append(cs.order, hash)
	}

	cs.mu.Unlock()

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
