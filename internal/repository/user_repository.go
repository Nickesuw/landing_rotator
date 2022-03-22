package repository

import (
	"github.com/Ferluci/ip2loc"
)

type UserRepository struct {
	geoDB *ip2loc.DB
}

func NewUserRepository(geoDB *ip2loc.DB) UserRepository {
	return UserRepository{geoDB: geoDB}
}

func (r UserRepository) ParseIp(ip string) (ip2loc.IP2LocationRecord, error) {
	return r.geoDB.GetAll(ip)
}
