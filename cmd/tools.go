package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"

	"github.com/zion8992/textcat/tc"
)

func MakeRequest(Req string, Key string, Value string, Status string, Conn *websocket.Conn) error {
	response := tc.Send{
        Req: Req,
		Key: Key,
		Value: Value,
		Status: Status,
	}
    returnMe, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return Conn.WriteMessage(websocket.TextMessage, returnMe)

}