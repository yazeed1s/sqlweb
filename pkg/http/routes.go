package http

import (
	"net/http"
	"net/http/pprof"

	_h "github.com/yazeed1s/sqlweb/pkg/handler"
)

func handleMethod(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func RegisterRoutes(mux *http.ServeMux, handler _h.Handler) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/connect", handleMethod("POST", handler.ConnectHandler()))
	mux.HandleFunc("/save", handleMethod("POST", handler.SaveConnection()))
	mux.HandleFunc("/saved/connections", handleMethod("GET", handler.SavedConnectionsHandler()))
	mux.HandleFunc("/disconnect", handleMethod("POST", handler.DbDisconnect()))
	mux.HandleFunc("/execute", handleMethod("POST", handler.QueryHandler()))
	mux.HandleFunc("/update", handleMethod("POST", handler.UpdateRowHandler()))
	mux.HandleFunc("/export/json", handleMethod("GET", handler.ExportTableToJson()))
	mux.HandleFunc("/export/csv", handleMethod("GET", handler.ExportTableToCSV()))
	mux.HandleFunc("/export/sql", handleMethod("GET", handler.ShowCreateTable()))
	mux.HandleFunc("/schemas", handleMethod("GET", handler.ShowSchemas()))
	mux.HandleFunc("/table", handleMethod("GET", handler.TableDataHandler()))
	mux.HandleFunc("/columns/table", handleMethod("GET", handler.GetColumnData()))
	mux.HandleFunc("/table/size/", handleMethod("GET", handler.TableSizesHandler()))
	// mux.HandleFunc("/client", handleMethod("GET", handler.ShowConnectedClient))
	// mux.HandleFunc("/schema/:name/drop", handleMethod("POST", handler.DropDatabaseHandler))
	// mux.HandleFunc("/schema/create/:name", handleMethod("POST", handler.CreateDatabaseHandler))
	// mux.HandleFunc("/tables", handleMethod("GET", handler.ShowTablesHandler))
	// mux.HandleFunc("/table/:name/columns", handleMethod("GET", handler.CountTableColumnsHandler))
	// mux.HandleFunc("/table/:name/rows", handleMethod("GET", handler.CountTableRowsHandler))
	// mux.HandleFunc("/tables/size", handleMethod("GET", handler.TableSizesHandler))
	// mux.HandleFunc("/table/:name/size", handleMethod("GET", handler.TableSizeHandler))
	// mux.HandleFunc("/table/:name/drop", handleMethod("POST", handler.DropTableHandler))
	// mux.HandleFunc("/table/:name/truncate", handleMethod("POST", handler.TruncateTableHandler))
	// mux.HandleFunc("/schema/size", handler.SchemaSizeHandler)
	// mux.HandleFunc("/schema/:name", handler.HandleFuncSchemaByName)
}
