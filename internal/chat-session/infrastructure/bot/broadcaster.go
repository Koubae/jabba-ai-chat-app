package bot

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		connections: make(map[string]map[string][]*Connection),
	}
}

type Connection struct {
	Conn          *websocket.Conn
	ApplicationID string
	SessionID     string
	UserID        int64
	Username      string
	MemberID      string
	Channel       string
	writeMutex    sync.Mutex
}

type Broadcaster struct {
	// Map structure: ApplicationID -> SessionID -> Connection
	connections map[string]map[string][]*Connection
	mutex       sync.RWMutex
}

func (cm *Broadcaster) Connect(
	conn *websocket.Conn,
	applicationID string,
	sessionID string,
	userID int64,
	username string,
	memberID string,
	channel string,
) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if cm.connections[applicationID] == nil {
		cm.connections[applicationID] = make(map[string][]*Connection)
	}

	connections, _ := cm.connections[applicationID][sessionID]
	for _, connection := range connections {
		if connection.MemberID == memberID {
			return errors.New(fmt.Sprintf("Member %s %s already connected to session %s", username, memberID, sessionID))
		}
	}

	newConnection := &Connection{
		Conn:          conn,
		ApplicationID: applicationID,
		SessionID:     sessionID,
		UserID:        userID,
		Username:      username,
		MemberID:      memberID,
		Channel:       channel,
	}
	cm.connections[applicationID][sessionID] = append(cm.connections[applicationID][sessionID], newConnection)
	return nil
}

func (cm *Broadcaster) Disconnect(conn *websocket.Conn, applicationID string, sessionID string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	appConnections, exists := cm.connections[applicationID]
	if !exists {
		return
	}

	sessionConnections, existsSession := appConnections[sessionID]
	if existsSession {
		for i, connection := range sessionConnections {
			if connection.Conn == conn {
				cm.connections[applicationID][sessionID] = append(sessionConnections[:i], sessionConnections[i+1:]...)
				break
			}
		}

		if len(cm.connections[applicationID][sessionID]) == 0 {
			delete(cm.connections[applicationID], sessionID)
		}
	}

	if len(appConnections) == 0 {
		delete(cm.connections, applicationID)
	}

	// Forcing client to disconnect (in case is still connected)
	// TODO: check if we should close connection here or not.
	// I say that we should attempt to do it since it may fail during broadcasting (and we remove the connection)
	// but the actual WebSocket may be connected still.
	err := conn.Close()
	if err != nil {
		log.Printf("Failed to close connection: %s", err)
	}

}

func (cm *Broadcaster) GetSessionConnections(applicationID string, sessionID string) []*Connection {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	if appConnections, exists := cm.connections[applicationID]; exists {
		if sessionConnections, exists := appConnections[sessionID]; exists {
			// Return a copy to avoid concurrent access issues
			result := make([]*Connection, len(sessionConnections))
			copy(result, sessionConnections)
			return result
		}
	}
	return []*Connection{}
}

func (cm *Broadcaster) Broadcast(applicationID string, sessionID string, messageType int, message []byte) {
	cm.mutex.RLock()
	connections := cm.GetSessionConnections(applicationID, sessionID)
	cm.mutex.RUnlock()

	var failedConnections []*websocket.Conn
	var failedMutex sync.Mutex

	var wg sync.WaitGroup

	for _, conn := range connections {
		wg.Add(1)
		go func(connection *Connection) {
			defer wg.Done()

			connection.writeMutex.Lock()
			defer connection.writeMutex.Unlock()

			if err := connection.Conn.WriteMessage(messageType, message); err != nil {
				failedMutex.Lock()
				failedConnections = append(failedConnections, connection.Conn)
				failedMutex.Unlock()
			}
		}(conn)
	}

	wg.Wait()

	for _, failedConn := range failedConnections {
		cm.Disconnect(failedConn, applicationID, sessionID)
	}
}
