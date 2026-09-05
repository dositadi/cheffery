package userservice

import (
	"github.com/dositadi/cheffery/services/shared/logger"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	logger logger.Logger
}

func (b Bcrypt) Compare(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (b Bcrypt) GenerateHash(password []byte) (hashed []byte, err error) {
	hashed, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	return hashed, err
}
