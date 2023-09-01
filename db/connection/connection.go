package connection

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	_sql "sqlweb/db/sql"
)

// var connectionPools map[string]*sqlx.DB

// const defaultPostgresSqlSchema string = "public" //don't yell, this is just for now

type Connection struct {
	Host     string      `json:"host"`
	Port     int         `json:"port"`
	User     string      `json:"user"`
	Password string      `json:"password"`
	Name     string      `json:"database"`
	Type     _sql.DbType `json:"databaseType"`
}

func (c *Connection) UnmarshalJSON(data []byte) error {
	type clientAlias Connection
	aux := &struct {
		*clientAlias
		Type string `json:"databaseType"`
	}{clientAlias: (*clientAlias)(c)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.Type = parseDbType(aux.Type)
	return nil
}

func (c *Connection) MarshalJSON() ([]byte, error) {
	type clientAlias Connection
	aux := &struct {
		*clientAlias
		Type string `json:"databaseType"`
	}{clientAlias: (*clientAlias)(c)}
	aux.Type = c.Type.String()
	return json.Marshal(aux)
}

func parseDbType(dbType string) _sql.DbType {
	switch strings.ToLower(dbType) {
	case "mysql":
		return _sql.MySQL
	case "postgresql":
		return _sql.PostgreSQL
	case "sqlite":
		return _sql.SQLite
	default:
		return _sql.Unsupported
	}
}
func (c *Connection) optionalDBurl() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/",
		c.User,
		c.Password,
		c.Host,
		c.Port,
	)
}

func (c *Connection) mySqlUrl() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}

func (c *Connection) postgresUrl() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
	)
}

func ConnectToDatabase(c *Connection, dbType string) (*sql.DB, error) {
	log.Println("here: ", dbType)
	log.Println(c)
	log.Println(c.mySqlUrl())
	if len(dbType) == 0 {
		return nil, fmt.Errorf("database type cannot be empty")
	}
	var (
		db  *sql.DB
		err error
	)
	switch strings.ToLower(dbType) {
	case strings.ToLower(_sql.MySQL.String()):
		db, err = sql.Open("mysql", c.mySqlUrl())
	case strings.ToLower(_sql.PostgreSQL.String()):
		db, err = sql.Open("postgres", c.postgresUrl())
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	err = testQuery(db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func testQuery(db *sql.DB) error {
	_, err := db.Exec("SELECT 1")
	return err
}

func OptionalConnectToDatabase(config *Connection, dbType string) (*sql.DB, error) {
	log.Println(dbType)
	switch strings.ToLower(dbType) {
	case strings.ToLower(_sql.MySQL.String()):
		db, err := sql.Open("mysql",
			fmt.Sprintf(
				"%s:%s@tcp(%s:%d)/",
				config.User,
				config.Password,
				config.Host,
				config.Port,
			),
		)
		if err != nil {
			return nil, err
		}
		if err = db.Ping(); err != nil {
			err := db.Close()
			if err != nil {
				return nil, err
			}
			return nil, err
		}
		return db, nil

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

func Disconnect(db *sql.DB) error {
	return db.Close()
}

// func init() {
// 	connectionPools = make(map[string]*sqlx.DB)
// }

// func GetConnectionPool(dbType string, config *client.Client) (*sqlx.DB, error) {
// 	poolKey := dbType + ":" + config.Host + ":" + strconv.Itoa(int(config.Port))
// 	pool, ok := connectionPools[poolKey]
// 	if ok {
// 		return pool, nil
// 	}
// 	pool, err := createConnectionPool(dbType, config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	connectionPools[poolKey] = pool
// 	return pool, nil
// }

// func createConnectionPool(dbType string, config *client.Client) (*sqlx.DB, error) {
// 	switch strings.ToLower(dbType) {
// 	case strings.ToLower(_sql.MySQL.String()):
// 		db, err := sqlx.Open("mysql", mySqlUrl(*config))
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Set pool configuration
// 		db.SetMaxOpenConns(10)
// 		db.SetMaxIdleConns(5)
// 		db.SetConnMaxLifetime(5 * time.Minute)
// 		if err = db.Ping(); err != nil {
// 			err := db.Close()
// 			if err != nil {
// 				return nil, err
// 			}
// 			return nil, err
// 		}
// 		return db, nil
// 	default:
// 		return nil, fmt.Errorf("unsupported database type: %s", dbType)
// 	}
// }

//func ensureConnectionAlive(db *sqlx.DB) error {
//	if err := db.Ping(); err != nil {
//		newDB := db.DB// Get the underlying *sql.DB from sqlx.DB
//		if err := newDB.Close(); err != nil {
//			return err
//		}
//		newConn, err := newDB.Conn(context.Background())
//		if err != nil {
//			return err
//		}
//		db.DB = newConn // Update sqlx.DB with the new connection
//	}
//	return nil
//}
