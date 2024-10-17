package domain

import "github.com/gorilla/websocket"

type Client struct {
	Connection *websocket.Conn
	UserID     string
}

type GroupClient struct {
	Connection *websocket.Conn
	UserID     string
	GroupID    string
}
