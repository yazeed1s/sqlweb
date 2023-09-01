package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	_con "sqlweb/db/connection"
)

type ConnectionHistory struct {
	Schema     string          `json:"key"`
	Connection _con.Connection `json:"connection"`
}

const (
	//configDirName  = ".config"
	appDirName     = "sqlweb"
	configFileName = "connection_history.json"
)

func NewConnectionConfig(key string, connection *_con.Connection) *ConnectionHistory {
	return &ConnectionHistory{
		Schema:     key,
		Connection: *connection,
	}
}

// WriteToFile
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
func WriteToFile(conf *ConnectionHistory) (int, error) {
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
	_, err = file.Seek(-2, io.SeekEnd)
	if err != nil {
		return 0, err
	}
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

func ReadFromFile(key string) (_con.Connection, error) {
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
		return _con.Connection{}, err
	}

	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			return
		}
	}(file)

	bytes, err = io.ReadAll(file)
	if err != nil {
		return _con.Connection{}, err
	}
	var connections []ConnectionHistory
	err = json.Unmarshal(bytes, &connections)
	if err != nil {
		return _con.Connection{}, err
	}
	for _, conn := range connections {
		if conn.Schema == key {
			return conn.Connection, nil
		}
	}
	return _con.Connection{}, fmt.Errorf("connection not found for key: %s", key)
}

func GetSavedConnections() ([]_con.Connection, error) {
	var (
		err              error
		file             *os.File
		fullFilePath     string
		bytes            []byte
		configDir        string
		connections      []ConnectionHistory
		savedConnections []_con.Connection
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