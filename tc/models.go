package tc

import (
	"time"
)

/*
Textcat Websocket Protocol (TWP)
*/


type Recieve struct {
	/* general fields */
	Req string // request type
	Username string
	Token string // also used as password when loggin in or registering
	Value string // a value for storing something, changes depending on request type

	/* request specific fields */
}

type Send struct {
	Req string // request type
	/*
		status: status of client's previous request
	*/
	Key string
	Value string
	Status string
	/*
	ok: went well
	server_error: internal server error
	error: validation error, addon error, any other type of non-server error
	*/
}

type User struct {
	Username string
	Password string
	LastLogin time.Time
	Created time.Time
}


// bridge between the cmd/ and the tc/
type LogicBridge interface {
	HandleReq(msg []byte, conn RequestWriter) error // conn is needed for session manager
	LogMsg(level string, message string, args ...any)
	MakeRequest(Req string, Key string, Value string, Status string, conn RequestWriter) error

	/* Data */
	StoreData(table string, record any) error
	GetDataByID(table string, id int64, out any) error
	GetData(string, func(any) bool, any) error
	CreateTable(name string) error

	/* User */
	UserExists(table string, username string) (bool, error)

	/* Get Data */
	GetMaxCachedMessages() uint16 // from 0 to 65,535, should be enough
	GetMaxUserSessions() uint8 // up to 255 sessions is enough
}

type RequestWriter interface {
	WriteMessage(messageType int, data []byte) error
}
