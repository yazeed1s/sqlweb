package main

import (
	"database/sql"
	"fmt"
	_ "net/http/pprof"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_a "sqlweb/pkg/app"
)

func main() {
	// profiler.StartProfiling()
	// defer profiler.StopProfiling()
	app := _a.NewApp()
	err := app.ParseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	app.SetupRouter()
	app.StartServer()
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			return
		}
	}(app.Handler.GetDB())
}
