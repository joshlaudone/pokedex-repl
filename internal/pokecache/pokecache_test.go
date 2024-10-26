package pokecache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/joshlaudone/pokedex-repl/internal/pokecache"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.New(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestEmptyGet(t *testing.T) {
	const interval = 5 * time.Second

	t.Run("Test Empty Get", func(t *testing.T) {
		cache := pokecache.New(interval)
		val, ok := cache.Get("nonexistent key")
		if ok {
			t.Errorf("expected to not find key")
			return
		}
		if len(val) != 0 {
			t.Errorf("expected value to be empty")
		}
	})
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond

	key := "https://example.com"

	cache := pokecache.New(baseTime)
	cache.Add(key, []byte("testdata"))

	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
