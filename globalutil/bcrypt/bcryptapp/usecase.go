package bcryptapp

import "github.com/dositadi/cheffery/services/shared/logger"

type Bcrypt struct {
	logger logger.Logger
}

func New(logger logger.Logger) *Bcrypt {
	return &Bcrypt{
		logger: logger,
	}
}
