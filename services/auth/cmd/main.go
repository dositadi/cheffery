package cmd

import "github.com/dositadi/cheffery/services/auth/internal/platform/platformapp"

func main() {
	platformapp.New().Run()
}
