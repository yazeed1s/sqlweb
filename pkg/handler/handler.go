// Package handler provides HTTP request handlers for interacting with SQL databases
// and managing database connections, queries, and operations.
//
// It includes various functionalities for:
//   - Establishing and managing database connections.
//   - Executing SQL queries and handling query results.
//   - Retrieving database schema information.
//   - Performing database operations such as dropping tables, truncating tables,
//     dropping databases, and creating databases.
//   - Exporting table data to JSON and CSV formats.
//   - Handling HTTP request/response for database-related tasks.
//
// This package is designed to work with SQL databases through a RESTful API
package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/yazeed1s/sqlweb/db/connection"
	_sql "github.com/yazeed1s/sqlweb/db/sql"
	_client "github.com/yazeed1s/sqlweb/pkg/client"
	"github.com/yazeed1s/sqlweb/pkg/config"
	"github.com/yazeed1s/sqlweb/pkg/query"
)

type Handler struct {
	client *_client.Client
}

// Response represents a standard response structure for API responses.
type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func NewHandler() *Handler {
	return &Handler{
		client: &_client.Client{},
	}
}

func (h *Handler) GetDB() *sql.DB {
	return h.client.Database
}

// jsonResponse sends a JSON response with the specified HTTP status code.
func jsonResponse(writer http.ResponseWriter, status int, data interface{}) {
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		return
	}
}

// handleBadRequest sends a JSON response for a bad request.
func handleBadRequest(writer http.ResponseWriter, message string, e error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusBadRequest)

	var (
		response Response
		encoder  *json.Encoder
	)

	encoder = json.NewEncoder(writer)
	if e != nil {
		response = Response{
			Message: message,
			Error:   e.Error(),
		}
	}

	if err := encoder.Encode(response); err != nil {
		http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

// handleSuccessRequest sends a JSON response for a successful request.
func handleSuccessRequest(writer http.ResponseWriter, message string, data ...interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	var (
		response Response
		encoder  *json.Encoder
		err      error
	)

	encoder = json.NewEncoder(writer)
	if len(data) == 0 {
		response = createSuccessResponse(message, nil)
		err = encoder.Encode(response)
		if err != nil {
			http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
		}
		return
	}

	if message == "" {
		response = createSuccessResponse("", data...)
		err = encoder.Encode(response)
		if err != nil {
			http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
		}
		return
	}

	response = createSuccessResponse(message, data...)
	err = encoder.Encode(response)
	if err != nil {
		http.Error(writer, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

// createSuccessResponse creates a standard success response structure.
func createSuccessResponse(message string, objects ...interface{}) Response {
	response := Response{
		Message: message,
	}

	if len(objects) > 0 {
		if len(objects) == 1 {
			response.Data = objects[0]
		} else {
			data := mergeMaps(objects...)
			response.Data = data
		}
	}

	return response
}

// mergeMaps merges multiple maps into a single map.
func mergeMaps(maps ...interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	for _, m := range maps {
		if objMap, ok := m.(map[string]interface{}); ok {
			for key, value := range objMap {
				data[key] = value
			}
		}
	}
	return data
}

// parseConnectionRequest parses a JSON request body into a Connection object.
func parseConnectionRequest(request *http.Request) (*connection.Connection, error) {
	var (
		body       []byte
		err        error
		conn       connection.Connection
		bodyReader *bytes.Reader
		decoder    *json.Decoder
	)

	body, err = io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(request.Body)

	bodyReader = bytes.NewReader(body)
	decoder = json.NewDecoder(bodyReader)
	err = decoder.Decode(&conn)
	if err != nil {
		return nil, err
	}

	return &conn, nil
}

// createClient creates a database client from a Connection object.
func createClient(conn *connection.Connection) *_client.Client {
	return &_client.Client{
		Host:     conn.Host,
		Port:     conn.Port,
		User:     conn.User,
		Password: conn.Password,
		Name:     conn.Name,
		Type:     conn.Type,
	}
}

// setSchemaName sets the schema name for the client based on the database type.
func setSchemaName(client *_client.Client) {
	if strings.ToLower(client.Type.String()) == strings.ToLower(_sql.MySQL.String()) {
		client.Schema.Name = client.Name
	} else if strings.ToLower(client.Type.String()) == strings.ToLower(_sql.PostgreSQL.String()) {
		client.Schema.Name = "public"
	}
}

// getColumnsDataForTables retrieves column data for a list of tables.
func getColumnsDataForTables(client *_client.Client, tableNames []string) ([]_client.ColumnData, error) {
	columnsData := make([]_client.ColumnData, 0)
	for _, tableName := range tableNames {
		columns, err := client.GetColumnsData(tableName)
		if err != nil {
			return columnsData, err
		}
		columnsData = append(columnsData, columns)
	}

	return columnsData, nil
}

// checkURLParams checks the number of URL parameters against an expected count.
func checkURLParams(u *url.URL, expectedCount int) error {
	var (
		err       error
		parsedURL *url.URL
		urlParams url.Values
	)
	parsedURL, err = url.Parse(u.String())
	if err != nil {
		return err
	}

	urlParams = parsedURL.Query()
	if len(urlParams) != expectedCount {
		return fmt.Errorf("incorrect number of params: found %d, wanted %d", len(urlParams), expectedCount)
	}

	return nil
}

func (h *Handler) ShowConnectedClient(writer http.ResponseWriter) {
	// writer.Header().Set("Content-Type", "application/json")
	if h.client.Database == nil {
		msg := fmt.Sprintf("Database connection is nil %s", h.client.Name)
		response := Response{
			Message: msg,
			Error:   "Internal Server Error",
		}
		jsonResponse(writer, http.StatusInternalServerError, response)
		return
	}

	response := Response{
		Message: "OK",
		Data:    h.client,
	}
	jsonResponse(writer, http.StatusOK, response)
}

func (h *Handler) SaveConnection() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		conn, err := parseConnectionRequest(request)
		if err != nil {
			msg := fmt.Sprintf("Invalid request body: %v", request.Body)
			handleBadRequest(writer, msg, err)
			return
		}

		savedClient := config.NewConnectionConfig(conn.Name, conn)
		b, err := config.WriteToFile(savedClient)
		if err != nil {
			handleBadRequest(writer, "Error writing connection info to file", err)
			return
		}

		if b == 0 {
			msg := fmt.Sprintf("Error Saving connection info: %s", savedClient)
			handleBadRequest(writer, msg, err)
			return
		}

		handleSuccessRequest(writer, "Success: connection saved", nil)
	}
}

func (h *Handler) SavedConnectionsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			connections []connection.Connection
			err         error
		)

		connections, err = config.GetSavedConnections()
		if err != nil {
			handleBadRequest(writer, "Error retrieving saved connections: ", err)
			return
		}

		handleSuccessRequest(writer, "Success: connection saved", connections)
	}
}

func (h *Handler) ConnectHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			conn        *connection.Connection
			client      *_client.Client
			db          *sql.DB
			data        map[string]interface{}
			err         error
			msg         string
			tableNames  []string
			schema      string
			columnsData []_client.ColumnData
		)

		conn, err = parseConnectionRequest(request)
		if err != nil {
			msg = fmt.Sprintf("Invalid request body: %v", request.Body)
			handleBadRequest(writer, msg, err)
			return
		}

		client = createClient(conn)
		h.client = client
		db, err = connection.ConnectToDatabase(conn, conn.Type.String())
		if err != nil {
			handleBadRequest(writer, "Failed to connect to the database", err)
			return
		}

		h.client.Database = db
		if !strings.EqualFold(h.client.Type.String(), _sql.SQLite.String()) {
			setSchemaName(h.client)
		}

		tableNames, err = h.client.GetTableNames()
		if err != nil {
			msg = fmt.Sprintf("Failed to get available tables from %s", h.client.Name)
			handleBadRequest(writer, msg, err)
			return
		}

		columnsData, err = getColumnsDataForTables(h.client, tableNames)
		if err != nil {
			msg = fmt.Sprintf("Failed to get columns data for tables from %s", h.client.Name)
			handleBadRequest(writer, msg, err)
			return
		}

		h.client.Schema.NumTables = len(tableNames)
		msg = fmt.Sprintf("Successfully connected to %s", h.client.Name)
		// for PostgreSQL, avoid sending 'public' as schema name to the frontend
		if strings.EqualFold(h.client.Type.String(), _sql.PostgreSQL.String()) {
			schema = h.client.Name
		} else {
			schema = h.client.Schema.Name
		}
		data = map[string]interface{}{"schema": schema, "tables": columnsData}
		// log.Println("hey", h.client.Schema.Name)
		handleSuccessRequest(writer, msg, data)
	}
}

func (h *Handler) DbDisconnect() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		err := connection.Disconnect(h.client.Database)
		if err != nil {
			handleBadRequest(writer, "Failed to disconnect from database", err)
			return
		}
		handleSuccessRequest(writer, "Disconnected successfully")
	}
}

func (h *Handler) ShowSchemas() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err     error
			schemas []string
		)

		schemas, err = h.client.GetSchemaNames()
		if err != nil {
			handleBadRequest(writer, "Failed to get schemas from database", err)
			return
		}

		handleSuccessRequest(writer, "", schemas)
	}
}

func (h *Handler) ShowTablesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err        error
			tableNames []string
			msg        string
		)

		tableNames, err = h.client.GetTableNames()
		if err != nil {
			msg = fmt.Sprintf("Failed to get available tables from %s", h.client.Schema.Name)
			handleBadRequest(writer, msg, err)
			return
		}

		h.client.Schema.NumTables = len(tableNames)
		handleSuccessRequest(writer, "", tableNames)
	}
}

func (h *Handler) CountTableColumnsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			res       map[string]interface{}
			msg       string
			tableName string
			cols      int
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		cols, err = h.client.CountTableColumns(tableName)
		if err != nil {
			msg = fmt.Sprintf("Failed to count columns for table %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{
			"table":   tableName,
			"columns": cols,
		}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) CountTableRowsHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			res       map[string]interface{}
			msg       string
			tableName string
			rows      int
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		rows, err = h.client.CountTableRows(tableName)
		if err != nil {
			msg = fmt.Sprintf("Failed to count rows for table %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{
			"table": tableName,
			"rows":  rows,
		}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) GetColumnData() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			cols      _client.ColumnData
			msg       string
			tableName string
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		cols, err = h.client.GetColumnsData(tableName)
		if err != nil {
			msg = fmt.Sprintf("Failed to get columns data for table %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}
		handleSuccessRequest(writer, "", cols)
	}
}

func (h *Handler) ShowCreateTable() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err  error
			data string
			msg  string
		)

		data, err = h.client.ShowCreateTable()
		if err != nil {
			msg = "Failed to get table statement for tables"
			handleBadRequest(writer, msg, err)
			return
		}
		handleSuccessDownloadRequest(writer, data)
	}
}

func (h *Handler) TableDataHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err        error
			tableData  *_client.Table
			res        map[string]interface{}
			msg        string
			tableName  string
			page       string
			perPage    string
			rows       int
			pageInt    int
			perPageInt int
			totalPages float64
		)

		err = checkURLParams(request.URL, 3)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		page = request.URL.Query().Get("page")
		perPage = request.URL.Query().Get("perPage")
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			msg = fmt.Sprintf("invalid 'page' parameter: %s", page)
			handleBadRequest(writer, msg, err)
			return
		}

		perPageInt, err = strconv.Atoi(perPage)
		if err != nil {
			msg = fmt.Sprintf("invalid 'perPage' parameter: %s", perPage)
			handleBadRequest(writer, msg, err)
			return
		}

		rows, err = h.client.CountTableRows(tableName)
		if err != nil {
			msg = fmt.Sprintf("Failed to count table rows: %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}
		totalPages = float64(rows) / float64(perPageInt)
		if totalPages < 1 {
			totalPages = 1
		} else {
			totalPages = math.Round(totalPages)
		}

		tableData, err = h.client.GetTable(tableName, pageInt, perPageInt)
		if err != nil {
			msg = fmt.Sprintf("Failed to get table data: %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{
			"table":       tableData,
			"total_rows":  rows,
			"total_pages": totalPages,
		}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) TableSizeHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err          error
			tableSize    _client.TableSize
			responseData map[string]interface{}
			tableName    string
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, "", err)
			return
		}

		tableName = request.URL.Query().Get("name")
		if tableName == "" {
			handleBadRequest(writer, "Table name is missing or empty", nil)
			return
		}

		tableSize, err = h.client.GetTableSize(tableName)
		if err != nil {
			handleBadRequest(writer, fmt.Sprintf("Failed to get table size for %s", tableName), err)
			return
		}

		responseData = map[string]interface{}{
			"table": map[string]interface{}{
				"name": tableName,
				"size": tableSize,
			},
		}
		handleSuccessRequest(writer, "", responseData)
	}
}

func (h *Handler) TableSizesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			tableSize []_client.TableSize
			res       map[string]interface{}
		)

		tableSize, err = h.client.GetTablesSize()
		if err != nil {
			handleBadRequest(writer, "Failed to get table size", err)
			return
		}

		res = map[string]interface{}{"table_size": tableSize}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) UpdateRowHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		type JsonRequest struct {
			CellValue       string `json:"cellValue"`
			EditedCellValue string `json:"editedCellValue"`
			HeaderValue     string `json:"headerValue"`
			ParentColumn    string `json:"parentColumn"`
			TableName       string `json:"tableName"`
		}

		var (
			err    error
			result *query.Result
			res    map[string]interface{}
			msg    string
			req    JsonRequest
		)

		err = json.NewDecoder(request.Body).Decode(&req)
		if err != nil {
			handleBadRequest(writer, "Invalid JSON", err)
			return
		}

		result, err = query.UpdateRow(
			req.TableName, req.ParentColumn,
			req.EditedCellValue, req.CellValue,
			req.HeaderValue, h.client,
		)

		if err != nil {
			msg = "Failed to update row table data"
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{"result": result}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) QueryHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err    error
			q      *query.Query
			result *query.Result
			res    map[string]interface{}
			msg    string
		)

		if err = json.NewDecoder(request.Body).Decode(&q); err != nil {
			msg = fmt.Sprintf("invalid query: %s", q)
			handleBadRequest(writer, msg, err)
			return
		}

		result, err = query.ExecuteQuery(q, h.client)
		if err != nil {
			handleBadRequest(writer, "Failed to execute query", err)
			return
		}

		res = map[string]interface{}{"result": result}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) DropTableHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			result    *query.Result
			res       map[string]interface{}
			tableName string
			msg       string
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		result, err = query.DropTable(tableName, h.client.Schema.Name, h.client.Database)
		if err != nil {
			msg = fmt.Sprintf("Failed to drop table: %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{"result": result}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) TruncateTableHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			result    *query.Result
			res       map[string]interface{}
			tableName string
			msg       string
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		result, err = query.TruncateTable(tableName, h.client.Schema.Name, h.client.Database)
		if err != nil {
			msg = fmt.Sprintf("Failed to truncate table: %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{"result": result}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) DropDatabaseHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err    error
			result *query.Result
			res    map[string]interface{}
			dbName string
			msg    string
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)

			return
		}

		dbName = request.URL.Query().Get("name")
		result, err = query.DropDatabase(dbName, h.client.Database)
		if err != nil {
			msg = fmt.Sprintf("Failed to drop database: %s", dbName)
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{"result": result}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) CreateDatabaseHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err    error
			result *query.Result
			res    map[string]interface{}
			dbName string
			msg    string
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		dbName = request.URL.Query().Get("name")
		result, err = query.CreateDatabase(dbName, h.client.Database)
		if err != nil {
			msg = fmt.Sprintf("Failed to create database: %s", dbName)
			handleBadRequest(writer, msg, err)
			return
		}

		res = map[string]interface{}{"result": result}
		handleSuccessRequest(writer, "", res)
	}
}

func (h *Handler) ExportTableToJson() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			tableName string
			msg       string
			data      []byte
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		data, err = h.client.ExportToJson(tableName)
		if err != nil {
			msg = fmt.Sprintf("Failed to export table data: %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}

		handleSuccessDownloadRequest(writer, string(data))
	}
}

func (h *Handler) ExportTableToCSV() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(request.Body)

		var (
			err       error
			tableName string
			msg       string
			data      string
		)

		err = checkURLParams(request.URL, 1)
		if err != nil {
			handleBadRequest(writer, msg, err)
			return
		}

		tableName = request.URL.Query().Get("name")
		data, err = h.client.ExportToCSV(tableName)
		if err != nil {
			msg = fmt.Sprintf("Failed to export table data: %s", tableName)
			handleBadRequest(writer, msg, err)
			return
		}
		handleSuccessDownloadRequest(writer, data)
	}
}

func handleSuccessDownloadRequest(writer http.ResponseWriter, data string) {
	writer.Header().Set("Content-Type", "application/octet-stream")
	// writer.Header().Set("Filename", fileName)
	writer.WriteHeader(http.StatusAccepted)
	_, err := writer.Write([]byte(data))
	if err != nil {
		http.Error(writer, "Error writing response", http.StatusInternalServerError)
		return
	}
}
