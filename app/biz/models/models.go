package models

import "github.com/gorilla/websocket"

// User general type
type User struct {
	ID       uint64 `default:"id"`
	Name     string `default:"name"`
	Email    string `default:"email"`
	Password string `default:"password"`
}

type ListWS struct {
	List []*websocket.Conn
}
