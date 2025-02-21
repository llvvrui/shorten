package storage

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestSaveAndGet(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}

	defer mr.Close()
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	store := NewRedisStore(client)

	// 测试保存和获取
	_, err = store.SaveShortUrl("", "http://example.com", 0)
	if err != nil {
		t.Fatal(err)
	}

	// 获取生成的shortCode，假设counter初始为0，调用INCR后得到1，编码为"1"
	longUrl, err := store.GetLongUrl("1")
	if err != nil {
		t.Fatal(err)
	}

	if longUrl != "http://example.com" {
		t.Errorf("expected http://example.com, got %s", longUrl)
	}

}

// 自定义后缀的测试
func TestCustomSuffixConflict(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}

	defer mr.Close()
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	store := NewRedisStore(client)

	// 第一次保存customSuffix
	_, err = store.SaveShortUrl("mylink", "http://example.com", 0)
	if err != nil {
		t.Fatal(err)
	}

	// 再次保存相同的customSuffix
	// _, err = store.SaveShortUrl("mylink", "http://example.com", 0)
	// if err != ErrShortCodeExists {
	// 	t.Errorf("expected ErrShortCodeExists, got %v", err)
	// }
}

func TestExpiration(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}

	defer mr.Close()
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	store := NewRedisStore(client)
	expiration := time.Second * 10
	_, err = store.SaveShortUrl("", "http://example.com", expiration)
	if err != nil {
		t.Fatal(err)
	}

	key := "short:1"
	if mr.TTL(key) != 10*time.Second {
		t.Errorf("expected TTL ~10s, got %v", mr.TTL(key))
	}
}
