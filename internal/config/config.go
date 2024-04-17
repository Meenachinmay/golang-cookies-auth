package config

import (
	"github.com/redis/go-redis/v9"
	"golang-cookies/internal/database"
)

type ApiConfig struct {
	DB          *database.Queries
	RedisClient *redis.Client
}
