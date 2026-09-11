package admin

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"golf-booking-go/global"
	"golf-booking-go/internal/utils"
	"golf-booking-go/pkg/setting"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func TestTokenFlow(t *testing.T) {
	oldConfig, oldRedis := global.Config, global.RDB
	t.Cleanup(func() { global.Config, global.RDB = oldConfig, oldRedis })
	global.Config = &setting.Config{Secret: setting.SecretSetting{JwtSecret: "test-secret-only"}}
	access, refresh, err := utils.GenerateTokens(1)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		raw, kind string
		valid     bool
	}{
		{access, "access", true}, {refresh, "refresh", true}, {access, "refresh", false}, {refresh, "access", false}, {"garbage", "refresh", false},
	} {
		_, err := utils.ParseToken(tc.raw, tc.kind)
		if (err == nil) != tc.valid {
			t.Fatalf("type %s: valid=%v, err=%v", tc.kind, tc.valid, err)
		}
	}
	for _, tc := range []struct {
		name   string
		expiry *jwt.NumericDate
		method jwt.SigningMethod
		key    string
	}{
		{"expired", jwt.NewNumericDate(time.Now().Add(-time.Minute)), jwt.SigningMethodHS256, "test-secret-only"},
		{"missing expiry", nil, jwt.SigningMethodHS256, "test-secret-only"},
		{"wrong algorithm", jwt.NewNumericDate(time.Now().Add(time.Hour)), jwt.SigningMethodHS384, "test-secret-only"},
		{"wrong signature", jwt.NewNumericDate(time.Now().Add(time.Hour)), jwt.SigningMethodHS256, "other-secret"},
	} {
		raw, err := jwt.NewWithClaims(tc.method, utils.TokenClaims{AdminID: 1, TokenType: "refresh", RegisteredClaims: jwt.RegisteredClaims{ID: "id", IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: tc.expiry}}).SignedString([]byte(tc.key))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := utils.ParseToken(raw, "refresh"); !errors.Is(err, utils.ErrInvalidToken) {
			t.Fatalf("%s accepted: %v", tc.name, err)
		}
	}
	t.Run("rotation", func(t *testing.T) {
		binary, err := exec.LookPath("redis-server")
		if err != nil {
			t.Skip("redis-server required for rotation integration check")
		}
		socket := filepath.Join(t.TempDir(), "redis.sock")
		cmd := exec.Command(binary, "--port", "0", "--unixsocket", socket, "--save", "", "--appendonly", "no")
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = cmd.Process.Kill(); _ = cmd.Wait() })
		client := redis.NewClient(&redis.Options{Network: "unix", Addr: socket})
		t.Cleanup(func() { _ = client.Close() })
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for client.Ping(ctx).Err() != nil {
			if ctx.Err() != nil {
				t.Fatal(ctx.Err())
			}
			time.Sleep(10 * time.Millisecond)
		}
		global.RDB = client
		service := &AdminService{}
		_, first, err := service.issueTokens(ctx, 1, "")
		if err != nil {
			t.Fatal(err)
		}
		key := fmt.Sprintf("admin:refresh:%x", sha256.Sum256([]byte(first)))
		ttl, err := client.TTL(ctx, key).Result()
		if err != nil || ttl <= 0 || ttl > utils.RefreshTokenTTL {
			t.Fatalf("invalid TTL: %v, %v", ttl, err)
		}
		var wg sync.WaitGroup
		results := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); _, _, err := service.issueTokens(ctx, 1, first); results <- err }()
		}
		wg.Wait()
		close(results)
		success, rejected := 0, 0
		for err := range results {
			if err == nil {
				success++
			} else if errors.Is(err, ErrUnauthorized) {
				rejected++
			} else {
				t.Fatal(err)
			}
		}
		if success != 1 || rejected != 1 {
			t.Fatalf("success=%d rejected=%d", success, rejected)
		}
		if client.Exists(ctx, key).Val() != 0 {
			t.Fatal("old token remains usable")
		}
		if _, _, err := service.issueTokens(ctx, 1, "unregistered"); !errors.Is(err, ErrUnauthorized) {
			t.Fatalf("unknown token: %v", err)
		}
		global.RDB = nil
		if _, _, err := service.issueTokens(ctx, 1, ""); err == nil {
			t.Fatal("issued tokens without store")
		}
	})
}
