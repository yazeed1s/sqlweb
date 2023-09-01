package main

import (
	"database/sql"
	"fmt"
	_ "net/http/pprof"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yazeed1s/sqlweb/pkg/app"
)

func main() {
	// profiler.StartProfiling()
	// defer profiler.StopProfiling()
	a := app.NewApp()
	err := a.ParseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	a.SetupRouter()
	a.StartServer()
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			return
		}
	}(a.Handler.GetDB())
}
