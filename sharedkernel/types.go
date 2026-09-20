package sharedkernel

type JWT string

func (j JWT) String() string {
	return string(j)
}
