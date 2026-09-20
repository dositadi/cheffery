package authapp

type Bcrypt interface {
	Compare(hashedPassword, password []byte) error
}
