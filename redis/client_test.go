package redis

import (
	"context"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewClient("47.237.141.189:17503", "dvmIO8dJn%", 10)
	if err != nil {
		t.Fatalf("连接 Redis 服务器失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Test "SET"
	t.Run("Set", func(t *testing.T) {
		resp, err := client.Set(ctx, "key_set_1", "value_set_1")
		if err != nil {
			t.Fatalf("Error set key: %v", err)
		}
		if resp != "OK" {
			t.Fatalf("Expected OK, got %v", resp)
		}
	})

	// Test "SET with expiration"
	t.Run("Set with expiration", func(t *testing.T) {
		resp, err := client.Set(ctx, "key_set_2", "value_set_2", 2*time.Second)
		if err != nil {
			t.Fatalf("Error set key: %v", err)
		}
		if resp != "OK" {
			t.Fatalf("Expected OK, got %v", resp)
		}
	})

	// Test "GET"
	t.Run("Get", func(t *testing.T) {
		resp, err := client.Get(ctx, "key_set_1")
		if err != nil {
			t.Fatalf("Error get key: %v", err)
		}
		if resp != "value_set_1" {
			t.Fatalf("Expected value_set_1, got %v", resp)
		}
	})

	// Test "GET with expiration"
	t.Run("Get with expiration", func(t *testing.T) {
		resp, err := client.Get(ctx, "key_set_2")
		if err != nil {
			t.Fatalf("Error get key: %v", err)
		}
		if resp != "value_set_2" {
			t.Fatalf("Expected value_set_2, got %v", resp)
		}
		time.Sleep(2 * time.Second)
		resp, err = client.Get(ctx, "key_set_2")
		if err != nil {
			t.Fatalf("Error get key: %v", err)
		}
		if resp != "" {
			t.Fatalf("Expected empty, got %v", resp)
		}
	})

	// Test "DEL"
	t.Run("Del", func(t *testing.T) {
		_, err := client.Set(ctx, "key_del_1", "value_del")
		if err != nil {
			t.Fatalf("Error set key: %v", err)
		}
		_, err = client.Set(ctx, "key_del_2", "value_del")
		if err != nil {
			t.Fatalf("Error set key: %v", err)
		}

		resp, err := client.Del(ctx, "key_del_1", "key_del_2")
		if err != nil {
			t.Fatalf("Error del key: %v", err)
		}
		if resp != "2" {
			t.Fatalf("Expected 2, got %v", resp)
		}
	})

	// Test "EXPIRE"
	t.Run("Expire", func(t *testing.T) {
		_, err := client.Set(ctx, "key_expire", "value_expire")
		if err != nil {
			t.Fatalf("Error set key: %v", err)
		}

		resp, err := client.Expire(ctx, "key_expire", 2*time.Second)
		if err != nil {
			t.Fatalf("Error expire key: %v", err)
		}
		if resp != "1" {
			t.Fatalf("Expected 1, got %v", resp)
		}

		time.Sleep(2 * time.Second)

		resp, err = client.Get(ctx, "key_expire")
		if err != nil {
			t.Fatalf("Error get key: %v", err)
		}
		if resp != "" {
			t.Fatalf("Expected empty, got %v", resp)
		}
	})

	// Test "INCR"
	t.Run("Incr", func(t *testing.T) {
		_, err := client.Set(ctx, "key_incr", "1")
		if err != nil {
			t.Fatalf("Error set key: %v", err)
		}

		resp, err := client.Incr(ctx, "key_incr")
		if err != nil {
			t.Fatalf("Error incr key: %v", err)
		}
		if resp != "2" {
			t.Fatalf("Expected 2, got %v", resp)
		}
	})

	// Test "Wrong Password"
	t.Run("Wrong Password", func(t *testing.T) {
		_, err := NewClient("x.x.x.x:xxxx", "wrong password", 10)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}
