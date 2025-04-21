package services

import (
	"base_structure/src/config"
	"base_structure/src/constants"
	"base_structure/src/data/cache"
	"base_structure/src/pkg/logging"
	"base_structure/src/pkg/service_errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type OtpService struct {
	Logger      logging.Logger
	Cfg         *config.Config
	RedisClient *redis.Client
}

type OtpDto struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *config.Config) *OtpService {
	logger := logging.NewLogger(cfg)
	redisClient := cache.GetRedis(cfg)
	return &OtpService{Logger: logger, Cfg: cfg, RedisClient: redisClient}
}

func (s *OtpService) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	val := &OtpDto{Value: otp, Used: false}
	res, err := cache.Get[OtpDto](s.RedisClient, key)
	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	}
	err = cache.Set(s.RedisClient, key, val, s.Cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (s *OtpService) ValidateOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	res, err := cache.Get[OtpDto](s.RedisClient, key)
	if err != nil {
		return err
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpUsed}
	} else if err == nil && !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: service_errors.OtpNotValid}
	} else if err == nil && !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(s.RedisClient, key, res, s.Cfg.Otp.ExpireTime*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
