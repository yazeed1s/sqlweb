package query

import (
	"fmt"
	"testing"

	_conn "sqlweb/db/connection"
	_sql "sqlweb/db/sql"
	_cl "sqlweb/pkg/client"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func SetupMySQLConnection() (*_cl.Client, error) {
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
	return &_cl.Client{
		Host:     client.Host,
		Port:     client.Port,
		User:     client.User,
		Password: client.Password,
		Name:     client.Name,
		Type:     client.Type,
		Database: db,
		Schema: _cl.Schema{
			Name: client.Name,
		},
	}, nil
}

func TestDropTable(t *testing.T) {
	// Set up the MySQL connection
	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")
	table := `CREATE TABLE test_p (
		PersonID int,
		LastName varchar(255),
		FirstName varchar(255),
		Address varchar(255),
		City varchar(255)
	);`
	addedTable := "test_p"
	rows, err := client.Database.Exec(table)
	assert.NoError(t, err)
	affected, err := rows.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), affected, "Expected rows affected to be 0")

	tables, err := client.GetTableNames()
	assert.NoError(t, err)
	assert.Contains(t, tables, addedTable)

	// Perform the test for DropTable
	result, err := DropTable(addedTable, client.Name, client.Database)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.AffectedRows, "Expected affected rows to be 0")
	assert.Contains(t, result.Msg, "dropped successfully")

	// Verify that the table no longer exists
	tables, err = client.GetTableNames()
	assert.NoError(t, err)
	assert.NotContains(t, tables, addedTable)
	err = client.Database.Close()
	if err != nil {
		return
	}
}

func TestTruncateTable(t *testing.T) {
	client, err := SetupMySQLConnection()
	assert.NoError(t, err, "Failed to set up MySQL connection")
	table := `CREATE TABLE test_p (
		PersonID int,
		LastName varchar(255),
		FirstName varchar(255),
		Address varchar(255),
		City varchar(255)
	);`
	addedTable := "test_p"
	rows, err := client.Database.Exec(table)
	assert.NoError(t, err)
	affected, err := rows.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), affected, "Expected rows affected to be 0")

	// Insert data into the table
	insertQuery := fmt.Sprintf(`
		INSERT INTO %s (
			PersonID, 
			LastName, 
			FirstName, 
			Address, 
			City
		) VALUES (?, ?, ?, ?, ?)`, addedTable,
	)
	_, err = client.Database.Exec(insertQuery, 1, "Doe", "John", "123 Main St", "New York")
	assert.NoError(t, err)
	r, err := client.CountTableRows(addedTable)
	assert.NoError(t, err)
	assert.Equal(t, 1, r, "Expected rows affected to be 1")
	// Perform the test for TruncateTable
	result, err := TruncateTable(addedTable, client.Name, client.Database)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// // Verify that the table is empty
	emptyQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", addedTable)
	var count int
	err = client.Database.QueryRow(emptyQuery).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Expected the table to be empty")
	// Perform the test for DropTable
	result, err = DropTable(addedTable, client.Name, client.Database)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.AffectedRows, "Expected affected rows to be 0")
	assert.Contains(t, result.Msg, "dropped successfully")
	// Verify that the table no longer exists
	tables, err := client.GetTableNames()
	assert.NoError(t, err)
	assert.NotContains(t, tables, addedTable)
	err = client.Database.Close()
	if err != nil {
		return
	}
}
