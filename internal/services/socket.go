package services

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var wsMutex *sync.Mutex = &sync.Mutex{}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func UpgradeWsConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return nil, err
	}
	return conn, nil
}

// Add a connection to the activeConnections
func AddConnection(userId string, conn *websocket.Conn) {
	wsMutex.Lock()
	defer wsMutex.Unlock()
	SetCache(userId, conn.RemoteAddr().String())

	log.Printf("Connection added for user: %s", userId)
}

// Remove a connection from the activeConnections
func RemoveConnection(userId string) {
	wsMutex.Lock()
	defer wsMutex.Unlock()
	DeleteCache(userId)
	log.Printf("Connection removed for user: %s", userId)
}
