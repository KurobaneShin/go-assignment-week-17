package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	mo          *MemoryObject
	http        *http.Server
	mongoClient *mongo.Client
}

func NewServer(addr string) *Server {
	http := &http.Server{
		Addr: addr,
	}

	mongoClient := NewMongoClient()

	return &Server{
		http:        http,
		mo:          NewMemoryObject(),
		mongoClient: mongoClient,
	}
}

func (s *Server) Start() error {
	http.HandleFunc("/in-memory", s.inMemoryHandler)

	http.HandleFunc("/fetch-data", s.fetchDataHandler)

	fmt.Printf("Server is running on %s\n", s.http.Addr)

	err := s.http.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	server := NewServer(":3334")

	defer func() {
		server.mongoClient.Disconnect(context.TODO())
	}()

	log.Fatal(server.Start())
}
