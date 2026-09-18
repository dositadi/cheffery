package bcryptapp

import "golang.org/x/crypto/bcrypt"

func (b *Bcrypt) GenerateHash(password []byte) (hashed []byte, err error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}
