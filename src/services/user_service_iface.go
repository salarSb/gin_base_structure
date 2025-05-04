package services

import "base_structure/src/api/dto"

type UserServiceIface interface {
	SendOtp(req *dto.GetOtpRequest) error
	LoginByUsername(req *dto.LoginByUsernameRequest) (*dto.TokenDetail, error)
	RegisterByUsername(req *dto.RegisterByUsernameRequest) error
	RegisterLoginByMobileNumber(req *dto.RegisterLoginByMobileRequest) (*dto.TokenDetail, error)
}
