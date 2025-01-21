package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"hello-world/internal/auth"
)

type Server struct {
	auth *auth.FirebaseAuth
}

func NewServer(ctx context.Context) (*Server, error) {
	firebaseAuth, err := auth.NewFirebaseAuth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase: %v", err)
	}

	return &Server{
		auth: firebaseAuth,
	}, nil
}

func (s *Server) handleHello(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user := r.Context().Value("user")

	response := map[string]interface{}{
		"message": "Hello, World!",
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	ctx := context.Background()

	server, err := NewServer(ctx)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Protected endpoint
	http.HandleFunc("/", server.auth.AuthMiddleware(server.handleHello))

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
