package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"sqlweb/pkg/cli"
	"sqlweb/pkg/handler"
	_http "sqlweb/pkg/http"
	_static "sqlweb/static"
)

type App struct {
	Args    *cli.Args
	Router  *http.ServeMux
	Handler *handler.Handler
}

func NewApp() *App {
	return &App{
		Args:    cli.NewArgs(),
		Router:  http.NewServeMux(),
		Handler: handler.NewHandler(),
	}
}

func (app *App) ParseFlags() error {
	var (
		showVersion *bool
		showHelp    *bool
		err         error
	)
	flag.IntVar(&app.Args.Port, "p", app.Args.Port, "Set the port number (default: 3000)")
	flag.BoolVar(&app.Args.Log, "l", app.Args.Log, "Enable logging")
	flag.StringVar(&app.Args.Connection, "c", app.Args.Connection, "Use saved connection")
	showVersion = flag.Bool("v", false, "Display version")
	showHelp = flag.Bool("h", false, "Show help")
	flag.Parse()
	if *showVersion {
		fmt.Println(app.Args.Version)
		os.Exit(0)
	}
	if *showHelp {
		fmt.Println(app.Args.Help)
		os.Exit(0)
	}
	if err = app.Args.ValidatePortRange(); err != nil {
		return err
	}
	return nil
}

func (app *App) SetupRouter() {
	app.Router.HandleFunc("/", _static.ServeStaticFiles)
	_http.RegisterRoutes(app.Router, *app.Handler)
}

func (app *App) StartServer() {
	serveMux := _http.CorsMiddleware(app.Router)
	log.Print("Listening...", app.Args.Port)
	// log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.Args.Port), app.Router))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", app.Args.Port), serveMux))
}
