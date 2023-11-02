package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dieklingel/core/internal/core"
	"github.com/redis/go-redis/v9"
)

type AppService struct {
	storageService core.StorageService
	redisClient    *redis.Client
}

func NewAppService(storageService core.StorageService) *AppService {
	return &AppService{
		storageService: storageService,
		redisClient: redis.NewClient(&redis.Options{
			Addr: storageService.Read().Redis.Host,
		}),
	}
}

func (appService *AppService) Register(name string, token string) {
	_, err := appService.redisClient.Set(context.Background(), name, token, 24*30*time.Hour).Result()
	if err != nil {
		fmt.Printf("could not register app '%s'; error: %v", name, err)
	}
}

func (appService *AppService) Unregister(name string) {
	_, err := appService.redisClient.Del(context.Background(), name).Result()
	if err != nil {
		fmt.Printf("could not unregister app '%s'; error: %v", name, err.Error())
	}
}

func (appService *AppService) IsRegisterd(name string) bool {
	result := appService.redisClient.Exists(context.Background(), name).Val()
	return result == 1
}
