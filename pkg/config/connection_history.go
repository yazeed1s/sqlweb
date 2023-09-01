package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/yazeed1s/sqlweb/db/connection"
)

// ConnectionHistory represents an object for storing connection information in a file.
// The 'Schema' field holds the schema name, which serves as the key for retrieving objects from the file later.
// The 'Connection' field contains the actual connection data.
type ConnectionHistory struct {
	Schema     string                `json:"key"`
	Connection connection.Connection `json:"connection"`
}

const (
	//configDirName  = ".config"
	appDirName     = "sqlweb"
	configFileName = "connection_history.json"
)

// NewConnectionConfig creates a new ConnectionHistory object with the provided key and connection data.
func NewConnectionConfig(key string, connection *connection.Connection) *ConnectionHistory {
	return &ConnectionHistory{
		Schema:     key,
		Connection: *connection,
	}
}

// WriteToFile appends a ConnectionHistory object to a JSON file for persistent storage,
// ensuring that the file maintains an array of JSON objects.
//
// If the file doesn't exist, it creates the file and initializes it with a JSON array containing the provided object.
// If the file already exists, it appends the JSON object to the existing array.
func WriteToFile(conf *ConnectionHistory) (int, error) {
	// os.UserHomeDir():
	// - On Unix, including macOS, it returns the $HOME environment variable
	// - On Windows, it returns %USERPROFILE%
	// - On Plan 9, it returns the $home environment variable
	// os.UserConfigDir():
	//   - On Unix systems, it returns $XDG_CONFIG_HOME as specified by
	//     https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
	//     if non-empty, else $HOME/.config.
	//   - On Darwin, it returns $HOME/Library/Application Support
	//   - On Windows, it returns %AppData%
	//   - On Plan 9, it returns $home/lib.
	var (
		err        error
		file       *os.File
		appDirPath string
		fileName   string
		data       []byte
		configDir  string
		bits       int
	)
	configDir, err = os.UserConfigDir()
	if err != nil {
		return 0, err
	}
	appDirPath = filepath.Join(configDir, appDirName)
	if _, err = os.Stat(appDirPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(appDirPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	fileName = filepath.Join(appDirPath, configFileName)

	data, err = json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return 0, err
	}

	if _, err = os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(fileName)
		if err != nil {
			return 0, err
		}
		// Insert the JSON object between [ ] to represent a JSON array of objects.
		bits, err = file.WriteString("[\n" + string(data) + "\n]")
		if err != nil {
			return 0, err
		}
		return bits, nil
	}

	file, err = os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return 0, err
	}

	// Move the file pointer to the end, just before the closing square bracket.
	_, err = file.Seek(-2, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	// Append the new ConnectionHistory object to the existing file.
	bits, err = file.WriteString("\n," + string(data) + "\n]")
	if err != nil {
		return 0, err
	}

	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			return
		}
	}(file)

	return bits, nil
}

// ReadFromFile reads a ConnectionHistory object from the configuration file based on the provided key.
func ReadFromFile(key string) (connection.Connection, error) {
	var (
		err          error
		file         *os.File
		fullFilePath string
		bytes        []byte
		configDir    string
	)
	configDir, err = os.UserConfigDir()
	fullFilePath = filepath.Join(configDir, appDirName, configFileName)
	file, err = os.Open(fullFilePath)
	if err != nil {
		return connection.Connection{}, err
	}

	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			return
		}
	}(file)

	bytes, err = io.ReadAll(file)
	if err != nil {
		return connection.Connection{}, err
	}
	var connections []ConnectionHistory
	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return connection.Connection{}, err
	}
	for _, conn := range connections {
		if conn.Schema == key {
			return conn.Connection, nil
		}
	}
	return connection.Connection{}, fmt.Errorf("connection not found for key: %s", key)
}

// GetSavedConnections retrieves all saved connection configurations from the configuration file.
func GetSavedConnections() ([]connection.Connection, error) {
	var (
		err              error
		file             *os.File
		fullFilePath     string
		bytes            []byte
		configDir        string
		connections      []ConnectionHistory
		savedConnections []connection.Connection
	)
	configDir, err = os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	fullFilePath = filepath.Join(configDir, appDirName, configFileName)
	file, err = os.Open(fullFilePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	bytes, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return nil, err
	}

	for _, conn := range connections {
		savedConnections = append(savedConnections, conn.Connection)
	}

	return savedConnections, nil
}
