package connection

import (
	"database/sql"
	"encoding/json"
	_sql "github.com/yazeed1s/sqlweb/db/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestClientJSONMarshaling(t *testing.T) {
	// Create a sample Client instance
	conn := &Connection{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "password",
		Name:     "mydb",
		Type:     _sql.MySQL,
	}
	data, err := json.Marshal(conn)
	assert.NoError(t, err)
	var parsedConnection Connection
	err = json.Unmarshal(data, &parsedConnection)
	assert.NoError(t, err)
	assert.Equal(t, conn, &parsedConnection)
}

func TestConnectToDatabase(t *testing.T) {
	client := &Connection{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "11221122",
		Name:     "classicmodels",
		Type:     _sql.MySQL,
	}
	db, err := ConnectToDatabase(client, client.Type.String())
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer func(db *sql.DB) {
		err := Disconnect(db)
		if err != nil {
			return
		}
	}(db)
}

func TestNoDatabaseType(t *testing.T) {
	client := &Connection{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "11221122",
		Name:     "classicmodels",
		Type:     _sql.Unsupported, // Empty database type
	}

	db, err := ConnectToDatabase(client, client.Type.String())
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestOptionalConnectToDatabase(t *testing.T) {
	// Test case 1: MySQL
	config := &Connection{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "11221122",
		Name:     "",
		Type:     _sql.MySQL,
	}
	db, err := OptionalConnectToDatabase(config, _sql.MySQL.String())
	assert.NoError(t, err)
	assert.NotNil(t, db)
	defer func(db *sql.DB) {
		err := Disconnect(db)
		if err != nil {

		}
	}(db)
}

func TestUnsupportedDatabaseType(t *testing.T) {
	client := &Connection{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "11221122",
		Name:     "classicmodels",
		Type:     _sql.SQLite, // Unsupported database type
	}

	db, err := ConnectToDatabase(client, client.Type.String())
	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestProperDisconnection(t *testing.T) {
	client := &Connection{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "11221122",
		Name:     "classicmodels",
		Type:     _sql.MySQL,
	}

	db, err := ConnectToDatabase(client, client.Type.String())
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = Disconnect(db)
	assert.NoError(t, err)
	_, err = db.Exec("SELECT 1")
	assert.EqualError(t, err, "sql: database is closed")
}
