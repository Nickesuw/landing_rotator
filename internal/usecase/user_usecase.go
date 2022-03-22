package usecase

import (
	"github.com/Ferluci/ip2loc"
	_ "gitlab.tubecorporate.com/platform-go/core/pkg/chlog"
	"landing_rotator/internal/repository"
)

type UserUseCase struct {
	UserRepository repository.UserRepository
}

func NewUserUseCase(geoDB *ip2loc.DB) UserUseCase {
	userRepository := repository.NewUserRepository(geoDB)
	return UserUseCase{UserRepository: userRepository}
}

func (c UserUseCase) ParseIp(ip string) (ip2loc.IP2LocationRecord, error) {
	return c.UserRepository.ParseIp(ip)
}
