package prompt_api

import (
	"log"
	"sync"

	"guardian/internal/models/entities"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type ClientManager struct {
	mu      *sync.Mutex
	clients map[primitive.ObjectID]*grpc.ClientConn
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[primitive.ObjectID]*grpc.ClientConn),
	}
}

func (m *ClientManager) GetClient(llm entities.Plugin) (*grpc.ClientConn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, exists := m.clients[llm.ID]; exists {
		return conn, nil
	}

	client, err := grpc.NewClient(llm.Address)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	client.Connect()
	m.clients[llm.ID] = client
	return client, nil
}

func (m *ClientManager) CloseAll() {
    m.mu.Lock()
    defer m.mu.Unlock()

    for id, conn := range m.clients {
        if err := conn.Close(); err != nil {
            log.Printf("Failed to close connection for TargetLLM %s: %v", id, err)
        }
        delete(m.clients, id)
    }
}
