package chash

import (
	"testing"
)

func TestHash(t *testing.T) {
	hash := New(10, nil)
	if hash == nil {
		t.Errorf("expected obj")
	}
	hash.Add("127.0.0.1:545", "143.43.4.4:124", "176.26.72.1:43")

	node := hash.Get("key1")
	t.Log("key1 node =", node)

	node = hash.Get("key2")
	t.Log("key2 node =", node)

	hash.Remove("127.0.0.1:545", "143.43.4.4:124")

	node = hash.Get("key1")
	t.Log("key1 node =", node)

	hash.Add("127.555.0.1:545")

	node = hash.Get("key1")
	t.Log("key1 node =", node)

	hash.Remove("176.26.672.1:43")

	node = hash.Get("key1")
	t.Log("key1 node =", node)
}
