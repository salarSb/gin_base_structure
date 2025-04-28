package services

import (
	"base_structure/src/config"
	"base_structure/src/data/cache"
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-redis/redis/v7"
	"time"
)

type BlacklistService struct {
	redis *redis.Client
}

func NewBlacklistService(cfg *config.Config) *BlacklistService {
	return &BlacklistService{redis: cache.GetRedis(cfg)}
}

func (b *BlacklistService) Blacklist(token string, ttl time.Duration) error {
	return b.redis.Set(hashToken(token), 1, ttl).Err()
}

func (b *BlacklistService) IsBlacklisted(token string) (bool, error) {
	n, err := b.redis.Exists(hashToken(token)).Result()
	return n == 1, err
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return hex.EncodeToString(sum[:])
}
