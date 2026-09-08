package userapp

type Bcrypt interface {
	Compare(hashedPassword, password []byte) error
	GenerateHash(password []byte) (hashed []byte, err error)
}
