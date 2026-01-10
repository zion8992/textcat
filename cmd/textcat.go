package main

import (
	/* std */
	"log/slog"
	"encoding/json"
	"fmt"
	
	/* database */
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	/* stuff */
	"errors"
	"bytes"
	"encoding/gob"
	"reflect"

	/* textcat */
	"github.com/zion8992/textcat/tc"
)

type Application struct {
	Log *slog.Logger
	Database *sql.DB
	Textcat *tc.Textcat // session manager, auth manager, channel manager
}

func (app *Application) StoreData(table string, record any) error {
	// Create table if it doesn't exist
	createQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB NOT NULL
		)
	`, table)
	if _, err := app.Database.Exec(createQuery); err != nil {
		return err
	}

	// Encode the record
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(record); err != nil {
		return err
	}

	// Insert the data
	insertQuery := fmt.Sprintf("INSERT INTO %s (data) VALUES (?)", table)
	_, err := app.Database.Exec(insertQuery, buf.Bytes())
	return err
}


func (app *Application) GetDataByID(table string, id int64, out any) error {
	query := fmt.Sprintf("SELECT data FROM %s WHERE id = ?", table)

	var blob []byte
	if err := app.Database.QueryRow(query, id).Scan(&blob); err != nil {
		return err
	}

	return gob.NewDecoder(bytes.NewReader(blob)).Decode(out)
}

func (app *Application) GetData(table string, match func(any) bool, out any) error {
	query := fmt.Sprintf("SELECT data FROM %s", table)

	rows, err := app.Database.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var blob []byte
		if err := rows.Scan(&blob); err != nil {
			return err
		}

		tmp := reflect.New(reflect.TypeOf(out).Elem()).Interface()
		if err := gob.NewDecoder(bytes.NewReader(blob)).Decode(tmp); err != nil {
			return err
		}

		if match(tmp) {
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(tmp).Elem())
			return nil
		}
	}

	return sql.ErrNoRows
}


func (app *Application) UserExists(table string, username string) (bool, error) {
	query := fmt.Sprintf("SELECT data FROM %s", table)

	rows, err := app.Database.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var blob []byte
		if err := rows.Scan(&blob); err != nil {
			return false, err
		}

		var u tc.User
		if err := gob.NewDecoder(bytes.NewReader(blob)).Decode(&u); err != nil {
			return false, err
		}

		if u.Username == username {
			return true, nil
		}
	}

	return false, nil
}


func (app *Application) CreateTable(name string) error {
	// Create table if it doesn't exist
	createQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			data BLOB NOT NULL
		)
	`, name)
	if _, err := app.Database.Exec(createQuery); err != nil {
		return err
	}
	return nil
}



func (app *Application) LogMsg(level string, message string, args ...any) {
	switch level {
		case "info":
			app.Log.Info(message, args...)
		case "warn":
			app.Log.Warn(message, args...)
		case "error":
			app.Log.Error(message, args...)
		default:
			app.Log.Info(message, args...)
	}
}


func (app *Application) HandleReq(msg []byte, conn any) error {
    var data tc.Recieve

    if err := json.Unmarshal(msg, &data); err != nil {
        return tc.MakeError("error", err)
    }

    fmt.Println(data)
	switch data.Req {
		case "register":
			if data.Username == "" || data.Token == "" {
				return errors.New("error: empty username or password")
			}
			err := app.Textcat.CreateUser(data.Username, data.Token)
			if err != nil {
				return err
			}
		case "login":
			if data.Username == "" || data.Token == "" {
				return errors.New("error: empty username or password")
			}
			err := app.Textcat.LoginUser(data.Username, data.Token, conn)
			if err != nil {
				return err
			}
	}

	
    return nil
}

func (app *Application) GetMaxCachedMessages() uint16 {
	// TODO: add this feature
	return 10
}

func (app *Application) GetMaxUserSessions() uint8 {
	// TODO: add a config
	return 32
}