package query

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_sql "sqlweb/db/sql"
	_cl "sqlweb/pkg/client"
)

type Query struct {
	SQLQuery string `json:"query"`
}

type Result struct {
	AffectedRows int64                    `json:"affected_rows"`
	Time         string                   `json:"time_taken"`
	Data         []map[string]interface{} `json:"data"`
	Msg          string                   `json:"message"`
}

// stringDataTypes contains substrings of data types
// that require quoting in SQL update statement
var stringDataTypes = []string{"char", "text", "date", "time", "year"}

func checkDatabaseConnection(db *sql.DB) error {
	if db == nil {
		return errors.New("database connection is nil")
	}
	return nil
}

// getColumnDataType returns the data type of a given column
func getColumnDataType(table, schema, column, dbType string, db *sql.DB) (string, error) {
	if db == nil {
		return "", errors.New("database connection is nil")
	}

	var (
		query    string
		err      error
		dataType string
	)

	switch strings.ToLower(dbType) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLGetColumnDataType, schema, table, column)
	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(_sql.PostgreSQLGetColumnDataType, schema, table, column)
	}

	err = db.QueryRow(query).Scan(&dataType)
	if err != nil {
		return "", err
	}

	return dataType, nil
}

// wrapValue: Wraps a value in single quotes if it requires quoting in SQL update statements.
// This applies to data types containing any of the substrings "char", "text", "date", "time", and "year".
// it returns the wrapped value if its column's data type matches any of the above substrings,
// or return the original value otherwise.
func wrapValue(dataType, value string) string {
	lowerCase := strings.ToLower(dataType)
	for _, substr := range stringDataTypes {
		if strings.Contains(lowerCase, substr) {
			return fmt.Sprintf("'%s'", value)
		}
	}
	return value
}

// wrapPrimaryKey: Wraps the primary key in single quotes if it requires quoting in SQL update statements.
// This applies to data types containing any of the substrings "char", "text", "date", "time", and "year".
// it returns the wrapped primary key if its column's data type matches any of the above substrings,
// or return the original primary key otherwise.
func wrapPrimaryKey(dataType, priKey string) string {
	lowerCase := strings.ToLower(dataType)
	for _, substr := range stringDataTypes {
		if strings.Contains(lowerCase, substr) {
			return fmt.Sprintf("'%s'", priKey)
		}
	}
	return priKey
}

// UpdateRow constructs and executes an SQL UPDATE statement to modify a row in the specified table.
// The function handles checking the column data type, and wraps its value in single quotes if necessary.
// Returns the result of the update operation or any encountered errors.
func UpdateRow(table, parentCol, newVal, priKeyVal, priKeyCol string, client *_cl.Client) (*Result, error) {
	if err := checkDatabaseConnection(client.Database); err != nil {
		return nil, err
	}
	var (
		err               error
		query             string
		msg               string
		sqlResult         sql.Result
		result            *Result
		startTime         time.Time
		rows              int64
		elapsedTime       time.Duration
		wrappedValue      string
		wrappedPrimaryKey string
		columnDataType    string
	)

	columnDataType, err = getColumnDataType(
		table, client.Schema.Name, parentCol,
		client.Type.String(), client.Database,
	)
	if err != nil {
		return nil, err
	}
	wrappedValue = wrapValue(columnDataType, newVal)

	columnDataType, err = getColumnDataType(
		table, client.Schema.Name, priKeyCol,
		client.Type.String(), client.Database,
	)
	if err != nil {
		return nil, err
	}
	wrappedPrimaryKey = wrapPrimaryKey(columnDataType, priKeyVal)

	query = fmt.Sprintf(_sql.SQLUpdateRow, table, parentCol, wrappedValue, priKeyCol, wrappedPrimaryKey)
	log.Println("query is: ", query)
	startTime = time.Now()
	sqlResult, err = client.Database.Exec(query)
	if err != nil {
		return nil, err
	}
	elapsedTime = time.Since(startTime)

	rows, err = sqlResult.RowsAffected()
	if err != nil {
		return nil, err
	}

	msg = fmt.Sprintf(
		"Row update successfully (%d rows affected, time taken %.3f)",
		rows, elapsedTime.Seconds(),
	)
	result = &Result{
		AffectedRows: rows,
		Time:         fmt.Sprintf("%.3f", elapsedTime.Seconds()),
		Msg:          msg,
	}
	return result, nil
}

func ExecuteQuery(q *Query, client *_cl.Client) (*Result, error) {
	if err := checkDatabaseConnection(client.Database); err != nil {
		return nil, err
	}

	var (
		err   error
		query string
		res   *Result
	)

	switch strings.ToLower(strings.ToLower(client.Type.String())) {
	case strings.ToLower(_sql.MySQL.String()):
		query = fmt.Sprintf(_sql.MySQLUse, client.Schema.Name)
		_, err = client.Database.Exec(query)
		if err != nil {
			return nil, err
		}
		query = fmt.Sprintf(q.SQLQuery)
		res, err = execQueryHelper(client.Database, query)
		if err != nil {
			return nil, err
		}
		return res, nil

	case strings.ToLower(_sql.PostgreSQL.String()):
		query = fmt.Sprintf(q.SQLQuery)
		res, err = execQueryHelper(client.Database, query)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, nil
}

func execQueryHelper(db *sql.DB, query string) (*Result, error) {
	var (
		err       error
		columns   []string
		msg       string
		rows      *sql.Rows
		startTime time.Time
		result    *Result
		row       map[string]interface{}
		values    []interface{}
		pointers  []interface{}
	)

	startTime = time.Now()
	result = &Result{
		AffectedRows: 0,
		Data:         make([]map[string]interface{}, 0),
	}

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

	for rows.Next() {
		row = make(map[string]interface{})
		values = make([]interface{}, len(columns))
		pointers = make([]interface{}, len(columns))
		for i := range columns {
			pointers[i] = &values[i]
		}
		err = rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		for i, column := range columns {
			val := values[i]
			if byteVal, ok := val.([]byte); ok {
				row[column] = string(byteVal)
			} else {
				row[column] = val
			}
		}
		result.Data = append(result.Data, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	result.Time = fmt.Sprintf("%.5f", time.Since(startTime).Seconds())
	result.AffectedRows = int64(len(result.Data))
	msg = fmt.Sprintf("Query executed successfully (%d rows affected, time taken %s)", result.AffectedRows, result.Time)
	result.Msg = msg
	return result, nil
}

func DropTable(table, dbname string, db *sql.DB) (*Result, error) {
	if err := checkDatabaseConnection(db); err != nil {
		return nil, err
	}

	var (
		err         error
		query       string
		res         sql.Result
		result      *Result
		startTime   time.Time
		elapsedTime time.Duration
		rows        int64
	)

	query = fmt.Sprintf(_sql.MySQLUse, dbname)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf(_sql.MySQLDropTable, table)
	startTime = time.Now()
	res, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	rows, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}

	elapsedTime = time.Since(startTime)
	result = &Result{
		AffectedRows: rows,
		Time:         fmt.Sprintf("%.3f", elapsedTime.Seconds()),
		Msg:          fmt.Sprintf("Table '%s' dropped successfully (%s)", table, elapsedTime.String()),
	}
	return result, nil
}

func TruncateTable(table, dbname string, db *sql.DB) (*Result, error) {
	if err := checkDatabaseConnection(db); err != nil {
		return nil, err
	}
	var (
		err         error
		query       string
		res         sql.Result
		result      *Result
		startTime   time.Time
		elapsedTime time.Duration
		rows        int64
	)

	query = fmt.Sprintf(_sql.MySQLUse, dbname)
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	query = fmt.Sprintf(_sql.MySQLTruncateTable, table)
	startTime = time.Now()
	res, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	elapsedTime = time.Since(startTime)
	result = &Result{
		AffectedRows: rows,
		Time:         fmt.Sprintf("%.3f", elapsedTime.Seconds()),
		Msg:          fmt.Sprintf("Table '%s' truncateed successfully (%s)", table, elapsedTime.String()),
	}
	return result, nil
}

func DropDatabase(dbname string, db *sql.DB) (*Result, error) {
	if err := checkDatabaseConnection(db); err != nil {
		return nil, err
	}
	var (
		err         error
		query       string
		res         sql.Result
		result      *Result
		startTime   time.Time
		elapsedTime time.Duration
		rows        int64
	)

	query = fmt.Sprintf(_sql.MySQLDropDatabase, dbname)
	startTime = time.Now()
	res, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	elapsedTime = time.Since(startTime)
	result = &Result{
		AffectedRows: rows,
		Time:         fmt.Sprintf("%.3f", elapsedTime.Seconds()),
		Msg:          fmt.Sprintf("Database '%s' dropped successfully (%s)", dbname, elapsedTime.String()),
	}
	return result, nil
}

func CreateDatabase(dbname string, db *sql.DB) (*Result, error) {
	if err := checkDatabaseConnection(db); err != nil {
		return nil, err
	}

	var (
		err         error
		query       string
		res         sql.Result
		result      *Result
		startTime   time.Time
		elapsedTime time.Duration
		rows        int64
	)

	query = fmt.Sprintf(_sql.MySQLCreateDatabase, dbname)
	startTime = time.Now()
	res, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	rows, err = res.RowsAffected()
	if err != nil {
		return nil, err
	}
	elapsedTime = time.Since(startTime)
	result = &Result{
		AffectedRows: rows,
		Time:         fmt.Sprintf("%.3f", elapsedTime.Seconds()),
		Msg:          fmt.Sprintf("Database '%s' dropped successfully (%s)", dbname, elapsedTime.String()),
	}
	return result, nil
}