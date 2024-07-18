package cache

import (
	"fmt"
	"testing"
)

func TestNewRedisCache(t *testing.T) {
	_, err := NewRedisCache("localhost:6379", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Redis connection was successful.")
}
