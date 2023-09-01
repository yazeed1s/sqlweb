package bin

import (
	"embed"
	"io/fs"
	"log"

	"net/http"
)

//go:embed all:build
var staticFiles embed.FS

func buildHTTPFS() http.FileSystem {
	build, err := fs.Sub(staticFiles, "build")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(build)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	fileSystem := buildHTTPFS()
	filePath := r.URL.Path
	if _, err := fileSystem.Open(filePath); err != nil {
		filePath = "index.html" // TODO: 404.hml
	}
	http.FileServer(fileSystem).ServeHTTP(w, r)
}
