// Package client provides a versatile client for interacting with various database types,
// abstracting away the complexities of database connections and queries. It simplifies the
// process of working with different database types, such as MySQL and PostgreSQL.
//
// The central component of this package is the `Client` type, which represents an active
// database client connected to a specific database. It includes essential connection details
// like host, port, username, password, database name, database type (e.g., MySQL or PostgreSQL),
// and schema information.
//
// The package also includes various methods for fetching schema names, schema sizes, table names,
// column details, table data, and more. These methods are designed to work seamlessly with different
// database types, providing a consistent and convenient API for database interactions.
//
// Note that this package assumes an already established database connection, and it focuses on
// simplifying data retrieval and schema exploration tasks.
package client

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_sql "github.com/yazeed1s/sqlweb/db/sql"
)

// Client represent the active client connected to the db
type Client struct {
	Host     string      `json:"host"`
	Port     int         `json:"port"`
	User     string      `json:"user"`
	Password string      `json:"password"`
	Name     string      `json:"database"`
	Type     _sql.DbType `json:"databaseType"`
	Schema   Schema      `json:"schema"`
	Database *sql.DB
}

// Schema represent the db schema connected to
type Schema struct {
	Name      string  `json:"name"`
	NumTables int     `json:"number_of_tables"`
	Size      float64 `json:"size_mb"`
	Tables    []Table `json:"tables"`
}

// Row represents a row within the table, inclusive of its related columns
type Row map[string]interface{}

// Table Represents a table along with its name, data rows, columns, number of columns, number of rows,
// and size in megabytes
type Table struct {
	Name      string   `json:"table_name"`
	Data      []Row    `json:"data"`
	Columns   []Column `json:"columns"`
	N_columns int      `json:"n_columns"`
	N_rows    int      `json:"n_rows"`
	Size      float64  `json:"size_mb"`
}

// Column represents a column within a table, including its field name, data type, key type (e.g., PRI KEY),
// constraint name, and references to other tables and columns.
type Column struct {
	Field            string `json:"field"`
	Type             string `json:"type"`
	Key              string `json:"key"`
	ConstraintName   string `json:"constraint_name"`
	ReferencedTable  string `json:"refrenced_table"`
	ReferencedColumn string `json:"refrenced_column"`
}

// ColumnData represents column-related data for a specific table
type ColumnData struct {
	TableName string   `json:"table_name"`
	Columns   []Column `json:"columns"`
}

// SchemaSize holds information about the size of a schema
type SchemaSize struct {
	Name string  `json:"name"`
	Size float64 `json:"size_mb"`
}

// TableSize holds information about the size of a schema
type TableSize struct {
	Table  string  `json:"table_name"`
	SizeMB float64 `json:"size_mb"`
}

/*
   Functions suffixed with '..Helper' follow a pattern where they utilize a preconstructed query provided by the caller function.

   Caller functions are designed to work with different database types, such as MySQL and PostgreSQL,
   by dynamically selecting the query based on the database type (using switch db.Type()).

   For example, when fetching data or performing operations, the caller prepares a customized query for the task at hand. The '..Helper' function then
   handles the query execution using the database instance, hiding the intenral task-achieving logic from the caller.

   This approach abstracts the database-specific complexities from the caller. If changes or optimizations are needed in the
   way queries are executed, they can be made within the helper function, minimizing the impact on the higher-level code.
*/

func getSchemaNamesHelper(query string, db *sql.DB) ([]string, error) {
	var (
		err         error
		res         *sql.Rows
		schemaNames []string
	)

	res, err = db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			return
		}
	}(res)

	for res.Next() {
		var dbName string
		if err := res.Scan(&dbName); err != nil {
			return nil, err
		}
		schemaNames = append(schemaNames, dbName)
	}
	return schemaNames, nil
}

func (c *Client) GetSchemaNames() ([]string, error) {
	if c.Database == nil {
		return nil, errors.New("database connection is nil")
	}

	var (
		names []string
		err   error
		query string
	)

	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = _sql.MySQLShowDatabases
		names, err = getSchemaNamesHelper(query, c.Database)
		if err != nil {
			return nil, err
		}
		return names, nil

	case strings.ToLower(_sql.PostgreSQL.String()):
		query = _sql.PostgreSQLShowDatabases
		names, err = getSchemaNamesHelper(query, c.Database)
		if err != nil {
			return nil, err
		}
		return names, nil
	}

	return nil, nil
}

func getSchemaSizeHelper(query string, db *sql.DB) (SchemaSize, error) {
	var (
		err        error
		schemaSize SchemaSize
	)
	err = db.QueryRow(query).Scan(&schemaSize.Name, &schemaSize.Size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SchemaSize{}, fmt.Errorf("schema '%s' not found", schemaSize.Name)
		}
		return SchemaSize{}, fmt.Errorf("error executing query: %w", err)
	}
	return schemaSize, nil
}

func (c *Client) GetSchemaSize(name string) (SchemaSize, error) {
	if c.Database == nil {
		return SchemaSize{}, errors.New("database connection is nil")
	}

	var (
		query      string
		err        error
		schemaSize SchemaSize
	)

	schemaSize.Name = name
	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLSchemaSize, name)
		schemaSize, err = getSchemaSizeHelper(query, c.Database)
		if err != nil {
			return SchemaSize{}, nil
		}
		return schemaSize, nil
	case strings.ToLower(_sql.PostgreSQL.String()):
		query = _sql.PostgreSQLSchemaSize
		schemaSize, err = getSchemaSizeHelper(query, c.Database)
		if err != nil {
			return SchemaSize{}, nil
		}
		return schemaSize, nil
	}

	return SchemaSize{}, nil
}

func (c *Client) CountTableColumns(tableName string) (int, error) {
	if c.Database == nil {
		return 0, errors.New("database connection is nil")
	}

	var (
		err   error
		query string
		rows  *sql.Rows
		count int
	)

	query = fmt.Sprintf(_sql.MySQLCountTableColumns, c.Schema.Name, tableName)
	rows, err = c.Database.Query(query)
	if err != nil {
		return 0, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}

func countTableRowsHelper(query string, db *sql.DB) (int, error) {
	var (
		err      error
		rowCount int
	)
	err = db.QueryRow(query).Scan(&rowCount)
	if err != nil {
		return 0, err
	}
	return rowCount, nil
}

func (c *Client) CountTableRows(tableName string) (int, error) {
	if c.Database == nil {
		return 0, errors.New("database connection is nil")
	}

	var (
		query    string
		rowCount int
		err      error
	)

	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLCountTableRows, c.Schema.Name, tableName)
		rowCount, err = countTableRowsHelper(query, c.Database)
		if err != nil {
			return 0, err
		}
		return rowCount, nil

	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(_sql.PostgreSQLCountTableRows, c.Schema.Name, tableName)
		rowCount, err = countTableRowsHelper(query, c.Database)
		if err != nil {
			return 0, err
		}
		return rowCount, nil
	case strings.ToLower(_sql.SQLite.String()):
		query = fmt.Sprintf(_sql.SQLiteCountTableRows, tableName)
		rowCount, err = countTableRowsHelper(query, c.Database)
		if err != nil {
			return 0, err
		}
		return rowCount, nil
	}
	return 0, nil
}

func getTableNamesHelper(query string, db *sql.DB) ([]string, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	var (
		rows   *sql.Rows
		err    error
		tables []string
	)

	rows, err = db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var tableName string
		if err = rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (c *Client) GetTableNames() ([]string, error) {
	if c.Database == nil {
		return nil, errors.New("database connection is nil")
	}

	var (
		tables []string
		err    error
		query  string
	)

	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLUse, c.Schema.Name)
		_, err = c.Database.Exec(query)
		if err != nil {
			return nil, err
		}
		query = _sql.MySQLShowTables
		tables, err = getTableNamesHelper(query, c.Database)
		if err != nil {
			return nil, err
		}

	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(_sql.PostgreSQLShowTables, c.Schema.Name)
		tables, err = getTableNamesHelper(query, c.Database)
		if err != nil {
			return nil, err
		}
	case strings.ToLower(_sql.SQLite.String()):
		query = _sql.SQLiteShowTables
		tables, err = getTableNamesHelper(query, c.Database)
		if err != nil {
			return nil, err
		}
	}

	return tables, nil
}

func getColumnsHelper(query string, db *sql.DB) ([]Column, error) {

	var (
		rows    *sql.Rows
		err     error
		columns []Column
	)

	rows, err = db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	for rows.Next() {
		var column Column
		err = rows.Scan(
			&column.Field,
			&column.Type,
			&column.Key,
			&column.ConstraintName,
			&column.ReferencedTable,
			&column.ReferencedColumn,
		)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

func (c *Client) GetColumns(tableName string) ([]Column, error) {
	if c.Database == nil {
		return nil, errors.New("database connection is nil")
	}

	var (
		err   error
		query string
		cols  []Column
	)

	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLColumnsInfo, c.Schema.Name, tableName)
		cols, err = getColumnsHelper(query, c.Database)
		if err != nil {
			return nil, err
		}
		return cols, nil
	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(_sql.PostgreSQLColumnsInfo, c.Schema.Name, tableName)
		cols, err = getColumnsHelper(query, c.Database)
		if err != nil {
			return nil, err
		}
		return cols, nil
	case strings.ToLower(_sql.SQLite.String()):
		query = fmt.Sprintf(_sql.SQLiteColumnsInfo, tableName)
		cols, err = getColumnsHelper(query, c.Database)
		if err != nil {
			return nil, err
		}
		return cols, nil
	}

	return nil, nil
}

func (c *Client) GetColumnsData(tableName string) (ColumnData, error) {
	if c.Database == nil {
		return ColumnData{}, errors.New("database connection is nil")
	}

	var (
		err   error
		query string
		cols  []Column
		data  ColumnData
	)

	data.TableName = tableName
	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLColumnsInfo, c.Schema.Name, tableName)
		cols, err = getColumnsHelper(query, c.Database)
		data.Columns = cols
		if err != nil {
			return ColumnData{}, err
		}
		return data, nil

	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(_sql.PostgreSQLColumnsInfo, c.Schema.Name, tableName)
		cols, err = getColumnsHelper(query, c.Database)
		data.Columns = cols
		if err != nil {
			return ColumnData{}, err
		}
		return data, nil
	case strings.ToLower(_sql.SQLite.String()):
		query = fmt.Sprintf(_sql.SQLiteColumnsInfo, tableName)
		cols, err = getColumnsHelper(query, c.Database)
		data.Columns = cols
		if err != nil {
			return ColumnData{}, err
		}
		return data, nil
	}

	return ColumnData{}, nil
}

/*
- buildSelectAll constructs the SQL query to select all columns from a table.
- Based on the database type, it formats the query string with the appropriate placeholders and values.
- It returns the formatted query string.
*/
func buildSelectAll(cols []Column, DbType, schema, table string, perPage, offset int) string {
	var (
		columnList string
		query      string
	)
	for i, columnName := range cols {
		if i > 0 {
			columnList += ", "
		}
		// handle column names with spaces
		switch strings.ToLower(DbType) {
		case strings.ToLower(_sql.MySQL.String()):
			columnList += fmt.Sprintf("`%s`", columnName.Field)
		case strings.ToLower(_sql.PostgreSQL.String()):
			columnList += fmt.Sprintf("\"%s\"", columnName.Field)
		default:
			columnList += columnName.Field
		}
	}

	switch strings.ToLower(DbType) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLSelectAllWithLimit, columnList, schema, table, perPage, offset)
	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(_sql.PostgreSQLSelectAllWithLimit, columnList, schema, table, perPage, offset)
	case strings.ToLower(_sql.SQLite.String()):
		query = fmt.Sprintf(_sql.SQLiteSelectAllWithLimit, columnList, table, perPage, offset)
	}

	return query
}

func getTableHelper(query string, db *sql.DB) (*Table, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	var (
		rows      *sql.Rows
		tableData *Table
		err       error
		columns   []string
		results   []Row
		values    []interface{}
		valuePtrs []interface{}
		numRows   int
		numCols   int
	)

	rows, err = db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	// TODO: (Optimize memory allocation) preallocating 'results' with the exact number of rows
	// results := make([]Row, 0, rowCount)
	values = make([]interface{}, len(columns))
	valuePtrs = make([]interface{}, len(columns))
	for rows.Next() {
		row := make(Row, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err = rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}
		results = append(results, row)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	numRows, numCols = len(results), len(columns)
	if err != nil {
		return nil, err
	}

	tableData = &Table{
		Data:      results,
		N_columns: numCols,
		N_rows:    numRows,
	}

	return tableData, nil
}

func (c *Client) GetTable(tableName string, page, perPage int) (*Table, error) {
	if c.Database == nil {
		return nil, errors.New("database connection is nil")
	}

	var (
		cols      []Column
		tableData *Table
		table     *Table
		size      TableSize
		err       error
		offset    int
		query     string
	)

	offset = (page - 1) * perPage
	cols, err = c.GetColumns(tableName)
	if err != nil {
		return nil, err
	}

	query = buildSelectAll(cols, c.Type.String(), c.Schema.Name, tableName, perPage, offset)
	tableData, err = getTableHelper(query, c.Database)
	if err != nil {
		return nil, err
	}

	// sqlite3 driver does not set SQLITE_ENABLE_DBSTAT_VTAB,
	// dbstat is needed to get table size in sqlite
	// for now, just skip the size funcion
	if !strings.EqualFold(c.Type.String(), _sql.SQLite.String()) {
		size, err = c.GetTableSize(tableName)
		if err != nil {
			return nil, err
		}
	}

	table = &Table{
		Name:      tableName,
		Data:      tableData.Data,
		Columns:   cols,
		N_columns: len(cols),
		N_rows:    len(tableData.Data),
		Size:      size.SizeMB,
	}

	return table, nil
}

func getTableSizes(query string, db *sql.DB) ([]TableSize, error) {

	var (
		rows   *sql.Rows
		tables []TableSize
		err    error
	)

	rows, err = db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	tables = make([]TableSize, 0)
	for rows.Next() {
		var tableSize TableSize
		err = rows.Scan(&tableSize.Table, &tableSize.SizeMB)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		tables = append(tables, tableSize)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return tables, nil
}

func (c *Client) GetTablesSize() ([]TableSize, error) {
	if c.Database == nil {
		return nil, errors.New("database connection is nil")
	}
	// tableSizes := make([]TableSize, 0)
	var (
		tableSizes []TableSize
		err        error
		query      string
	)

	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLGetTablesSize, c.Schema.Name)
		tableSizes, err = getTableSizes(query, c.Database)
		if err != nil {
			return nil, err
		}
		return tableSizes, nil

	case strings.ToLower(_sql.PostgreSQL.String()):
		query = _sql.PostgreSQLTableSizes
		tableSizes, err = getTableSizes(query, c.Database)
		if err != nil {
			return nil, err
		}
		return tableSizes, nil
	case strings.ToLower(_sql.SQLite.String()):
		query = _sql.SQLiteTablesSize
		tableSizes, err = getTableSizes(query, c.Database)
		if err != nil {
			return nil, err
		}
		return tableSizes, nil
	}

	return nil, nil
}

func getTableSize(query string, db *sql.DB) (TableSize, error) {
	var (
		tableSize TableSize
		err       error
	)
	err = db.QueryRow(query).Scan(&tableSize.Table, &tableSize.SizeMB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TableSize{}, fmt.Errorf("table '%s' not found", tableSize.Table)
		}
		return TableSize{}, fmt.Errorf("error executing query: %w", err)
	}
	return tableSize, nil
}

func (c *Client) GetTableSize(table string) (TableSize, error) {
	if c.Database == nil {
		return TableSize{}, errors.New("database connection is nil")
	}

	var (
		err   error
		t     TableSize
		query string
	)

	log.Println("get table sizes for ", table)
	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLGetTableSize, c.Schema.Name, table)
		t, err = getTableSize(query, c.Database)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return TableSize{}, fmt.Errorf("table '%s' not found", table)
			}
			return TableSize{}, fmt.Errorf("error executing query: %w", err)
		}
		return t, nil

	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(_sql.PostgreSQLTableSize, c.Schema.Name, table)
		t, err = getTableSize(query, c.Database)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return TableSize{}, fmt.Errorf("table '%s' not found", table)
			}
			return TableSize{}, fmt.Errorf("error executing query: %w", err)
		}
		return t, nil
	case strings.ToLower(_sql.SQLite.String()):
		query = fmt.Sprintf(_sql.SQLiteTableSize, table)
		log.Println("query size = ", query)
		t, err = getTableSize(query, c.Database)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return TableSize{}, fmt.Errorf("table '%s' not found", table)
			}
			return TableSize{}, fmt.Errorf("error executing query: %w", err)
		}
		return t, nil
	}

	return TableSize{}, nil
}

func createFile(fileName string) (*os.File, error) {
	var (
		err      error
		file     *os.File
		filePath string
		homeDir  string
		appPath  string
	)

	homeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	appPath = filepath.Join(homeDir, "sqlweb")
	if _, err = os.Stat(appPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(appPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	filePath = filepath.Join(appPath, fileName)
	if _, err = os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *Client) ExportToJsonFile(tableName string) (int, error) {
	if c.Database == nil {
		return 0, errors.New("database connection is nil")
	}

	var (
		err          error
		file         *os.File
		table        *Table
		jsonFileName string
		query        string
		data         []byte
		bytes        int
	)

	jsonFileName = fmt.Sprintf("%s.json", tableName)
	file, err = createFile(jsonFileName)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	query = fmt.Sprintf(_sql.SQLSelectAll, c.Schema.Name, tableName)
	table, err = getTableHelper(query, c.Database)
	if err != nil {
		return 0, err
	}

	data, err = json.MarshalIndent(table.Data, "", "\t")
	if err != nil {
		return 0, err
	}

	bytes, err = file.WriteString(string(data))
	if err != nil {
		return 0, err
	}

	return bytes, nil
}

func (c *Client) ExportToCSVFile(tableName string) (int, error) {
	if c.Database == nil {
		return 0, errors.New("database connection is nil")
	}

	var (
		err         error
		file        *os.File
		table       *Table
		writer      *csv.Writer
		firstRow    Row
		csvFileName string
		query       string
		header      []string
		bits        int
	)

	csvFileName = fmt.Sprintf("%s.csv", tableName)
	file, err = createFile(csvFileName)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	query = fmt.Sprintf(_sql.SQLSelectAll, c.Schema.Name, tableName)
	table, err = getTableHelper(query, c.Database)
	if err != nil {
		return 0, err
	}

	writer = csv.NewWriter(file)
	defer writer.Flush()

	if len(table.Data) > 0 {
		firstRow = table.Data[0]
		header = make([]string, 0, len(firstRow))
		for key := range firstRow {
			header = append(header, key)
		}
		if err = writer.Write(header); err != nil {
			return 0, err
		}
	}

	for _, row := range table.Data {
		var values []string
		for _, v := range row {
			values = append(values, fmt.Sprintf("%v", v))
		}
		if err = writer.Write(values); err != nil {
			return 0, err
		}
		bits += len([]byte(strings.Join(values, ","))) + len("\n")
	}
	return bits, nil
}

func (c *Client) ShowCreateTableFile() (int, error) {
	if c.Database == nil {
		return 0, fmt.Errorf("database connection is nil")
	}
	var (
		file         *os.File
		writer       *bufio.Writer
		err          error
		query        string
		tableName    string
		sqlStatement string
		tables       []string
		sqlFileName  string
		header       string
		totalBytes   int
		b            int
	)
	tables, err = c.GetTableNames()
	if err != nil {
		return 0, nil
	}

	sqlFileName = fmt.Sprintf("%s.sql", c.Schema.Name)
	file, err = createFile(sqlFileName)
	writer = bufio.NewWriter(file)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err = file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()
	header = `
========================================================================
========================================================================
`
	for _, t := range tables {
		query = fmt.Sprintf(_sql.MySQLShowCreateTable, c.Schema.Name, t)
		err = c.Database.QueryRow(query).Scan(&tableName, &sqlStatement)
		if err != nil {
			return 0, err
		}
		b, err = writer.WriteString(header + "\n")
		totalBytes = totalBytes + b
		if err != nil {
			return totalBytes, err
		}
		b, err = writer.WriteString("==== TABLE:" + tableName + "\n")
		totalBytes = totalBytes + b
		if err != nil {
			return totalBytes, err
		}
		b, err = writer.WriteString(sqlStatement + "\n")
		totalBytes = totalBytes + b
		if err != nil {
			return totalBytes, err
		}
	}

	if err = writer.Flush(); err != nil {
		return 0, err
	}

	return totalBytes, nil
}

func (c *Client) ExportToJson(tableName string) ([]byte, error) {
	if c.Database == nil {
		return nil, errors.New("database connection is nil")
	}

	var (
		err   error
		table *Table
		query string
		data  []byte
	)

	query = fmt.Sprintf(_sql.SQLSelectAll, c.Schema.Name, tableName)
	table, err = getTableHelper(query, c.Database)
	if err != nil {
		return nil, err
	}

	data, err = json.MarshalIndent(table.Data, "", "\t")
	if err != nil {
		return nil, err
	}

	return data, nil
}

// TODO: fix bug where NULL SQL values are preventing the export
// invistgate why tables with lots of null values aren't exported
func sqlToCsv(rows *sql.Rows) (string, error) {

	var (
		//err         error
		builder strings.Builder
		writer  *csv.Writer
	)
	writer = csv.NewWriter(&builder)
	defer writer.Flush()
	writer.Comma = ','
	columnNames, err := rows.Columns()
	if err != nil {
		return "", nil
	}
	headers := columnNames
	err = writer.Write(headers)
	if err != nil {
		return "", fmt.Errorf("failed to write headers: %w", err)
	}
	values := make([]interface{}, len(columnNames))
	valuePtrs := make([]interface{}, len(columnNames))
	for rows.Next() {
		row := make([]string, len(columnNames))
		for i := range columnNames {
			valuePtrs[i] = &values[i]
		}

		if err = rows.Scan(valuePtrs...); err != nil {
			return "", err
		}
		for i := range columnNames {
			var value interface{}
			rawValue := values[i]

			byteArray, ok := rawValue.([]byte)
			if ok {
				value = string(byteArray)
			} else {
				value = rawValue
			}
			float64Value, ok := value.(float64)
			if ok {
				value = fmt.Sprintf("%v", float64Value)
			} else {
				float32Value, ok := value.(float32)
				if ok {
					value = fmt.Sprintf("%v", float32Value)
				}
			}
			timeValue, ok := value.(time.Time)
			if ok {
				value = timeValue.Format(time.RFC822)
			}
			row[i] = fmt.Sprintf("%v", value)
		}
		err = writer.Write(row)
		if err != nil {
			return "", fmt.Errorf("failed to write data row to csv %w", err)
		}
	}
	if err = rows.Err(); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (c *Client) ExportToCSV(tableName string) (string, error) {
	if c.Database == nil {
		return "", errors.New("database connection is nil")
	}

	query := fmt.Sprintf(_sql.SQLSelectAll, c.Schema.Name, tableName)
	rows, err := c.Database.Query(query)
	if err != nil {
		return "", err
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	csvStr, err := sqlToCsv(rows)
	if err != nil {
		return "", err
	}

	return csvStr, nil
}

func (c *Client) ShowCreateTable() (string, error) {
	tables, err := c.GetTableNames()
	if err != nil {
		return "", err
	}

	seperator := `
========================================================================
========================================================================
`
	switch strings.ToLower(c.Type.String()) {
	case strings.ToLower(_sql.MySQL.String()):
		result, err := c.ShowCreateTableMySQL(tables, seperator)
		if err != nil {
			return "", nil
		}
		return result, nil
	case strings.ToLower(_sql.PostgreSQL.String()):
		result, err := c.ShowCreateTablePostgreSQL(tables, seperator)
		if err != nil {
			return "", nil
		}
		return result, nil
	}
	return "", nil
}

func (c *Client) ShowCreateTablePostgreSQL(tables []string, seperator string) (string, error) {
	if c.Database == nil {
		return "", fmt.Errorf("database connection is nil")
	}

	var (
		err          error
		query        string
		sqlStatement string
		builder      strings.Builder
	)

	_, err = c.Database.Exec(_sql.PostgreSQLShowCreateFunction)
	if err != nil {
		return builder.String(), err
	}

	defer func() {
		_, dropErr := c.Database.Exec(_sql.PostgreSQLDropShowCreateFunction)
		if dropErr != nil {
			return
		}
	}()

	for _, t := range tables {
		query = fmt.Sprintf(_sql.PostgreSQLShowCreate, c.Schema.Name, t)
		err = c.Database.QueryRow(query).Scan(&sqlStatement)
		if err != nil {
			return builder.String(), err
		}
		builder.WriteString(seperator + "\n")
		builder.WriteString("===== TABLE: " + t + " =====" + "\n")
		builder.WriteString(sqlStatement + "\n")
	}

	return builder.String(), nil
}

func (c *Client) ShowCreateTableMySQL(tables []string, seperator string) (string, error) {
	if c.Database == nil {
		return "", fmt.Errorf("database connection is nil")
	}

	var (
		err          error
		tableName    string
		sqlStatement string
		builder      strings.Builder
		query        string
	)

	for _, t := range tables {
		query = fmt.Sprintf(_sql.MySQLShowCreateTable, c.Schema.Name, t)
		err = c.Database.QueryRow(query).Scan(&tableName, &sqlStatement)
		if err != nil {
			return builder.String(), err
		}

		builder.WriteString(seperator + "\n")
		builder.WriteString("===== TABLE: " + tableName + " =====" + "\n")
		builder.WriteString(sqlStatement + "\n")
	}

	return builder.String(), nil
}

func (c *Client) ShowCreateTableSQLite(tables []string, seperator string) (string, error) {
	if c.Database == nil {
		return "", fmt.Errorf("database connection is nil")
	}

	var (
		err          error
		tableName    string
		sqlStatement string
		builder      strings.Builder
		query        string
	)

	for _, t := range tables {
		query = fmt.Sprintf(_sql.SQLiteShowCreateTable, t)
		err = c.Database.QueryRow(query).Scan(&tableName, &sqlStatement)
		if err != nil {
			return builder.String(), err
		}

		builder.WriteString(seperator + "\n")
		builder.WriteString("===== TABLE: " + tableName + " =====" + "\n")
		builder.WriteString(sqlStatement + "\n")
	}

	return builder.String(), nil
}
