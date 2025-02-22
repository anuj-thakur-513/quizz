package services

import (
	"log"
	"net/http"
	"sync"

	"encoding/json"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var wsMutex *sync.Mutex = &sync.Mutex{}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := map[string]bool{
			"http://localhost:5173":            true,
			"https://quizz.anujthakur.dev":     true,
			"https://www.quizz.anujthakur.dev": true,
		}

		origin := r.Header.Get("Origin")
		return allowedOrigins[origin]
	},
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

// send question detail over WS
func SendQuestion(conn *websocket.Conn, question primitive.M, wg *sync.WaitGroup, mu *sync.RWMutex) {
	questionId := question["_id"].(primitive.ObjectID).Hex()
	questionText := question["question_text"].(string)
	isMultipleCorrect := question["is_multiple_correct"].(bool)
	options := question["options"].(primitive.A)
	finalOptions := []string{}
	for _, option := range options {
		option := option.(primitive.M)
		finalOptions = append(finalOptions, option["option"].(string))
	}

	data := map[string]interface{}{
		"questionText":      questionText,
		"questionId":        questionId,
		"isMultipleCorrect": isMultipleCorrect,
		"options":           finalOptions,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal data:", err)
		return
	}
	mu.Lock()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonData)); err != nil {
		log.Println("Failed to write message:", err)
	}
	mu.Unlock()
	wg.Done()
}

func SendLeaderboard(conn *websocket.Conn, key string, wg *sync.WaitGroup, mu *sync.RWMutex) {
	data := GetZSet(key)
	var leaderboard []map[string]interface{}
	for _, d := range data {
		curr := map[string]interface{}{}
		if err := json.Unmarshal([]byte(d), &curr); err != nil {
			log.Println("Failed to unmarshal data:", err)
			continue
		}
		s := GetZScore(key, &LeaderboardSetMember{
			UserId:   curr["user_id"].(string),
			Username: curr["username"].(string),
		})
		curr["score"] = s
		leaderboard = append(leaderboard, curr)
	}

	jsonData, err := json.Marshal(leaderboard)
	if err != nil {
		log.Println("Failed to marshal data:", err)
		return
	}

	mu.Lock()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonData)); err != nil {
		log.Println("Failed to write message:", err)
	}
	mu.Unlock()
	wg.Done()
}

func SendQuizEnded(conn *websocket.Conn, wg *sync.WaitGroup, mu *sync.RWMutex) {
	mu.Lock()
	data := map[string]interface{}{
		"type": "quiz_ended",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal data:", err)
		return
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonData)); err != nil {
		log.Println("Failed to write message:", err)
	}
	mu.Unlock()
	wg.Done()
}
