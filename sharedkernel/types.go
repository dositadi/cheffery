package sharedkernel

type JWT string

func (j JWT) String() string {
	return string(j)
}

type HeaderName string

var (
	Authorization HeaderName = "Authorization"
	RefreshCustom HeaderName = "X-Refresh"
)

func (h HeaderName) String() string {
	return string(h)
}
