package services

import (
	"base_structure/src/api/dto"
	"base_structure/src/common"
	"base_structure/src/config"
	"base_structure/src/constants"
	"base_structure/src/data/db"
	"base_structure/src/data/models"
	"base_structure/src/pkg/logging"
	"base_structure/src/pkg/service_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	logger       logging.Logger
	cfg          *config.Config
	otpService   *OtpService
	tokenService *TokenService
	database     *gorm.DB
}

func NewUserService(cfg *config.Config) *UserService {
	database := db.GetDb(cfg)
	logger := logging.NewLogger(cfg)
	return &UserService{
		cfg:          cfg,
		database:     database,
		logger:       logger,
		otpService:   NewOtpService(cfg),
		tokenService: NewTokenService(cfg),
	}
}

func (s *UserService) RegisterByUsername(req *dto.RegisterByUsernameRequest) error {
	u := models.User{Username: req.Username, FirstName: req.FirstName, LastName: req.LastName, Email: req.Email}
	exists, err := s.existsByEmail(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
	}
	exists, err = s.existsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessage: service_errors.UsernameExists}
	}
	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	u.Password = string(hp)
	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return err
	}
	tx := s.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	err = tx.Create(&models.RoleUser{RoleId: roleId, UserId: u.ID}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	tx.Commit()
	return nil
}

func (s *UserService) RegisterLoginByMobileNumber(req *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error) {
	err := s.otpService.ValidateOtp(req.MobileNumber, req.Otp)
	if err != nil {
		return nil, err
	}
	exists, err := s.existsByMobileNumber(req.MobileNumber)
	if err != nil {
		return nil, err
	}
	u := models.User{MobileNumber: req.MobileNumber, Username: req.MobileNumber}
	if exists {
		var user models.User
		err = s.database.
			Model(&models.User{}).
			Where("mobile_number = ?", u.MobileNumber).
			Preload("RoleUsers", func(tx *gorm.DB) *gorm.DB {
				return tx.Preload("Role")
			}).
			First(&user).Error
		if err != nil {
			return nil, err
		}
		tokenDto := tokenDto{
			UserId:       user.ID,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
			Username:     user.Username,
			MobileNumber: user.MobileNumber,
			Email:        user.Email,
		}
		if len(*user.RoleUsers) > 0 {
			for _, ru := range *user.RoleUsers {
				tokenDto.Roles = append(tokenDto.Roles, ru.Role.Name)
			}
		}
		token, err := s.tokenService.GenerateToken(&tokenDto)
		if err != nil {
			return nil, err
		}
		return token, nil
	}
	bp := []byte(common.GeneratePassword())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	u.Password = string(hp)
	roleId, err := s.getDefaultRole()
	if err != nil {
		s.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return nil, err
	}
	tx := s.database.Begin()
	err = tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	err = tx.Create(&models.RoleUser{RoleId: roleId, UserId: u.ID}).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	var user models.User
	err = s.database.
		Model(&models.User{}).
		Where("mobile_number = ?", u.MobileNumber).
		Preload("RoleUsers", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	tokenDto := tokenDto{
		UserId:       user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		Email:        user.Email,
	}
	if len(*user.RoleUsers) > 0 {
		for _, ru := range *user.RoleUsers {
			tokenDto.Roles = append(tokenDto.Roles, ru.Role.Name)
		}
	}
	token, err := s.tokenService.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *UserService) LoginByUsername(req *dto.LoginByUsernameRequest) (*dto.TokenDetail, error) {
	var user models.User
	err := s.database.
		Model(&models.User{}).
		Where("username = ?", req.Username).
		Preload("RoleUsers", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("Role")
		}).First(&user).Error
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.InvalidCredentials, Err: err}
	}
	tokenDto := tokenDto{
		UserId:       user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		MobileNumber: user.MobileNumber,
		Email:        user.Email,
	}
	if len(*user.RoleUsers) > 0 {
		for _, ru := range *user.RoleUsers {
			tokenDto.Roles = append(tokenDto.Roles, ru.Role.Name)
		}
	}
	token, err := s.tokenService.GenerateToken(&tokenDto)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *UserService) SendOtp(req *dto.GetOtpRequest) error {
	otp := common.GenerateOtp()
	err := s.otpService.SetOtp(req.MobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) existsByEmail(email string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) existsByUsername(username string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) existsByMobileNumber(mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.Model(&models.User{}).
		Select("count(*) > 0").
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) getDefaultRole() (roleId uint, err error) {
	if err = s.database.Model(&models.Role{}).
		Select("id").
		Where("name = ?", constants.DefaultRoleName).
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}
