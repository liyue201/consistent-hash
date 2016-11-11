package chash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type HashFunc func(data []byte) uint32

const (
	DefaultReplicas = 10
)

type Hash struct {
	sync.Mutex
	hash     HashFunc
	replicas int
	keys     []int //  Sorted keys
	hashMap  map[int]string
}

func New(replicas int, fn HashFunc) *Hash {
	h := &Hash{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if h.replicas <= 0 {
		h.replicas = DefaultReplicas
	}
	if h.hash == nil {
		h.hash = crc32.ChecksumIEEE
	}
	return h
}

func (h *Hash) IsEmpty() bool {
	h.Lock()
	defer h.Unlock()

	return len(h.keys) == 0
}

func (h *Hash) Add(nodes ...string) {
	h.Lock()
	defer h.Unlock()

	for _, node := range nodes {
		for i := 0; i < h.replicas; i++ {
			key := int(h.hash([]byte(strconv.Itoa(i) + node)))

			if _, ok := h.hashMap[key]; !ok {
				h.keys = append(h.keys, key)
			}
			h.hashMap[key] = node
		}
	}
	sort.Ints(h.keys)
}

func (h *Hash) Remove(nodes ...string) {
	h.Lock()
	defer h.Unlock()

	deletedKey := make([]int, 0)
	for _, node := range nodes {
		for i := 0; i < h.replicas; i++ {
			key := int(h.hash([]byte(strconv.Itoa(i) + node)))

			if _, ok := h.hashMap[key]; ok {
				deletedKey = append(deletedKey, key)
				delete(h.hashMap, key)
			}
		}
	}
	if len(deletedKey) > 0 {
		h.deleteKeys(deletedKey)
	}
}

func (h *Hash) deleteKeys(deletedKeys []int) {
	sort.Ints(deletedKeys)

	index := 0
	count := 0
	for _, key := range deletedKeys {
		for ; index < len(h.keys); index++ {
			h.keys[index-count] = h.keys[index]

			if key == h.keys[index] {
				count++
				index++
				break
			}
		}
	}

	for ; index < len(h.keys); index++ {
		h.keys[index-count] = h.keys[index]
	}

	h.keys = h.keys[:len(h.keys)-count]
}

func (h *Hash) Get(key string) string {
	if h.IsEmpty() {
		return ""
	}

	hash := int(h.hash([]byte(key)))

	h.Lock()
	defer h.Unlock()

	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= hash
	})

	if idx == len(h.keys) {
		idx = 0
	}

	return h.hashMap[h.keys[idx]]
}
