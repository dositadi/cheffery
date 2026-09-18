package cmd

import "github.com/dositadi/cheffery/services/repository/internal/platform/app"

func main() {
	app.New().StartServer()
}
