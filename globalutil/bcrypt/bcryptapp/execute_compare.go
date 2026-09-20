package bcryptapp

import "golang.org/x/crypto/bcrypt"

func (b *Bcrypt) Compare(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
