package client

import (
	"database/sql"
	"fmt"
	"testing"

	_conn "github.com/yazeed1s/sqlweb/db/connection"
	_sql "github.com/yazeed1s/sqlweb/db/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupMySQLConnection() (*Client, error) {
	client := &_conn.Connection{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "11221122",
		Name:     "classicmodels",
		Type:     _sql.MySQL,
	}
	db, err := _conn.ConnectToDatabase(client, client.Type.String())
	if err != nil {
		return nil, err
	}
	return &Client{
		Host:     client.Host,
		Port:     client.Port,
		User:     client.User,
		Password: client.Password,
		Name:     client.Name,
		Type:     client.Type,
		Database: db,
		Schema: Schema{
			Name: client.Name,
		},
	}, nil
}

func TestGetSchemaNamesMySQL(t *testing.T) {
	// Set up the MySQL connection
	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")

	// Test GetSchemaNames
	schemaNames, err := client.GetSchemaNames()
	assert.NoError(t, err, "Failed to retrieve schema names")
	assert.NotNil(t, schemaNames, "Schema names should not be nil")
	require.NotEmpty(t, schemaNames)
	rows, err := client.Database.Query("SHOW DATABASES;")
	require.NoError(t, err)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			t.Fail()
		}
	}(rows)
	var actualSchemaNames []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			require.NoError(t, err)
		}
		actualSchemaNames = append(actualSchemaNames, dbName)
	}
	require.Equal(t, len(actualSchemaNames), len(schemaNames))
	for _, name := range schemaNames {
		assert.Contains(t, actualSchemaNames, name)
	}
	fmt.Printf("Schema Names: %+v\n", schemaNames)
	client.Database.Close()
}

func TestGetSchemaSizeMySQL(t *testing.T) {
	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")
	schemaSize, err := client.GetSchemaSize("classicmodels")
	require.NoError(t, err)
	var actualSize float64
	err = client.Database.QueryRow(`
		SELECT SUM(data_length + index_length) / 1024 / 1024 
		FROM information_schema.TABLES 
		WHERE table_schema = 'classicmodels' 
		GROUP BY table_schema;
	`).Scan(&actualSize)
	require.NoError(t, err)
	assert.Equal(t, "classicmodels", schemaSize.Name)
	assert.Equal(t, actualSize, schemaSize.Size)
	fmt.Printf("Schema Size: %+v\n", schemaSize)
	client.Database.Close()
}

func TestCountTableColumnsMySQL(t *testing.T) {
	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")
	query := fmt.Sprintf(`
		SELECT 
			count(*) AS TotalColumns 
		FROM 
			information_schema.columns 
		WHERE 
			table_schema = '%s' 
		AND 
			table_name = '%s';
	`, "classicmodels", "employees")
	var expectedCount int
	err = client.Database.QueryRow(query).Scan(&expectedCount)
	require.NoError(t, err)
	columnCount, err := client.CountTableColumns("employees")
	require.NoError(t, err)
	assert.Equal(t, expectedCount, columnCount)
	client.Database.Close()
}

func TestCountTableRowsMySQL(t *testing.T) {
	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s.%s", "classicmodels", "employees")
	var expectedCount int
	err = client.Database.QueryRow(query).Scan(&expectedCount)
	require.NoError(t, err)
	rowCount, err := client.CountTableRows("employees")
	require.NoError(t, err)
	assert.Equal(t, expectedCount, rowCount)
	client.Database.Close()
}

func TestGetTableNamesMySQL(t *testing.T) {

	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")

	useQuery := fmt.Sprintf("USE %s;", "classicmodels")
	showTablesQuery := "SHOW TABLES"

	_, err = client.Database.Exec(useQuery)
	require.NoError(t, err)

	res, err := client.Database.Query(showTablesQuery)
	require.NoError(t, err)
	defer func(res *sql.Rows) {
		err = res.Close()
		if err != nil {
			t.Fail()
		}
	}(res)
	var expectedTables []string
	for res.Next() {
		var tableName string
		err = res.Scan(&tableName)
		require.NoError(t, err)
		expectedTables = append(expectedTables, tableName)
	}
	tables, err := client.GetTableNames()
	require.NoError(t, err)
	assert.ElementsMatch(t, expectedTables, tables)
	client.Database.Close()
}

// func TestGetTableMySQL(t *testing.T) {

// 	client, err := SetupMySQLConnection()
// 	assert.NoError(t, err, "Failed to set up MySQL connection")

// 	dbName := "classicmodels"
// 	tableName := "employees"

// 	query := fmt.Sprintf("SELECT * FROM %s.%s LIMIT 3000", dbName, tableName)

// 	rows, err := client.Database.Query(query)
// 	require.NoError(t, err)
// 	defer func(rows *sql.Rows) {
// 		err := rows.Close()
// 		if err != nil {
// 			t.Fail()
// 		}
// 	}(rows)

// 	columns, err := rows.Columns()
// 	require.NoError(t, err)

// 	var expectedData []map[string]interface{}
// 	for rows.Next() {
// 		values := make([]interface{}, len(columns))
// 		valuePtrs := make([]interface{}, len(columns))
// 		for i := range columns {
// 			valuePtrs[i] = &values[i]
// 		}
// 		err := rows.Scan(valuePtrs...)
// 		require.NoError(t, err)
// 		row := make(map[string]interface{})
// 		for i, col := range columns {
// 			var v interface{}
// 			val := values[i]
// 			b, ok := val.([]byte)
// 			if ok {
// 				v = string(b)
// 			} else {
// 				v = val
// 			}
// 			row[col] = v
// 		}
// 		expectedData = append(expectedData, row)
// 	}
// 	require.NoError(t, rows.Err())

// 	table, err := client.GetTable(tableName)
// 	require.NoError(t, err)

// 	assert.Equal(t, tableName, table.Name)
// 	assert.ElementsMatch(t, expectedData, table.Data)
// 	assert.Equal(t, len(columns), table.N_columns)
// 	assert.Equal(t, len(expectedData), table.N_rows)
// 	client.Database.Close()
// }

func TestGetTableSizeMySQL(t *testing.T) {

	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")

	schema := "classicmodels"
	table := "employees"

	query := fmt.Sprintf(`
		SELECT
			table_name AS "Table",
			ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS "Size (MB)"
		FROM
			information_schema.TABLES
		WHERE
			TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s';
	`, schema, table)

	var expectedSize TableSize
	err = client.Database.QueryRow(query).Scan(&expectedSize.Table, &expectedSize.SizeMB)
	require.NoError(t, err)

	tableSize, err := client.GetTableSize(table)
	require.NoError(t, err)

	assert.Equal(t, expectedSize, tableSize)
	client.Database.Close()
}

func TestGetTablesSize(t *testing.T) {

	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")

	schema := "classicmodels"

	query := fmt.Sprintf(`
		SELECT
			TABLE_NAME AS "Table",
			ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS "Size (MB)"
		FROM
			information_schema.TABLES
		WHERE
			TABLE_SCHEMA = '%s'
		ORDER BY
			(DATA_LENGTH + INDEX_LENGTH) DESC;
	`, schema)

	expectedSizes := make([]TableSize, 0)
	rows, err := client.Database.Query(query)
	require.NoError(t, err)
	defer rows.Close()
	for rows.Next() {
		var tableSize TableSize
		err = rows.Scan(&tableSize.Table, &tableSize.SizeMB)
		require.NoError(t, err)
		expectedSizes = append(expectedSizes, tableSize)
	}
	tableSizes, err := client.GetTablesSize()
	require.NoError(t, err)
	assert.Equal(t, expectedSizes, tableSizes)
	client.Database.Close()
}
