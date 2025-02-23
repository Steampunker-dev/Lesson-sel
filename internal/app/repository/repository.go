package repository

import (
	"awesomeProject/internal/app/config"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)

type Repository struct {
	db     *gorm.DB
	logger *logrus.Logger
	rd     *redis.Client
	cfg    config.RedisConfig
}

func NewRepository(dsn string, l *logrus.Logger, cfg config.RedisConfig) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Password:    cfg.Password,
		Addr:        cfg.Host + ":" + strconv.Itoa(cfg.Port),
		DB:          0,
		DialTimeout: cfg.DialTimeout,
		ReadTimeout: cfg.ReadTimeout,
	})

	return &Repository{
		db:     db,
		logger: l,
		rd:     redisClient,
		cfg:    cfg,
	}, nil
}
func (r *Repository) CloseRedis() error {
	return r.rd.Close()
}
