package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	// "time"

	game "game_backend/proto/game_backend" // Replace with the actual path to your game package

	"google.golang.org/grpc"
)

type GameServer struct {
	game.UnimplementedGameServiceServer
	mu           sync.Mutex
	createdGames map[string]*GameState
}

type GameState struct {
	PlayerA  string
	PlayerB  string
	State    string // Example: "XOXOOX___"
	GameOver bool
	NotifyA  chan string
	NotifyB  chan string
}

func (s *GameServer) CreateGame(ctx context.Context, req *game.CreateGameRequest) (*game.CreateGameResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gameID := fmt.Sprintf("game-%d", len(s.createdGames)+1) // Simple game ID generation
	gameState := &GameState{
		PlayerA: req.PlayerName,
		State:   "", // Initial empty game state
		NotifyA: make(chan string, 10),
	}

	s.createdGames[gameID] = gameState

	// Simulate waiting for Player B
	// go func() {
	// 	for {
	// 		if gameState.PlayerB != "" {
	// 			break
	// 		}
	// 		time.Sleep(30 * time.Second)
	// 	}
	// 	// gameState.NotifyA <- fmt.Sprintf("%s has joined the game", gameState.PlayerB)
	// }()

	return &game.CreateGameResponse{GameId: gameID}, nil
}

func (s *GameServer) JoinGame(ctx context.Context, req *game.JoinGameRequest) (*game.JoinGameResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gameState, exists := s.createdGames[req.GameId]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if gameState.PlayerB != "" {
		return nil, fmt.Errorf("game is already full")
	}

	gameState.PlayerB = req.PlayerName
	gameState.NotifyB = make(chan string, 10)

	// Notify Player A that Player B has joined
	// gameState.NotifyA <- fmt.Sprintf("%s has joined the game", gameState.PlayerB)

	return &game.JoinGameResponse{Status: "Successfully joined"}, nil
}

func (s *GameServer) UpdateGameState(ctx context.Context, req *game.UpdateGameStateRequest) (*game.UpdateGameStateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	gameState, exists := s.createdGames[req.GameId]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if gameState.GameOver {
		return nil, fmt.Errorf("game is already over")
	}

	// Update the game state
	gameState.State = req.NewState

	log.Printf("Game ID: %s updated by Player: %s, New State: %s", req.GameId, req.PlayerName, req.NewState)

	// Notify both players automatically
	go func() {
		if gameState.PlayerA != "" {
			gameState.NotifyA <- fmt.Sprintf("%s::%s", req.PlayerName, req.NewState)
			// gameState.NotifyA <- req.NewState
		}
	}()
	go func() {
		if gameState.PlayerB != "" {
			gameState.NotifyB <- fmt.Sprintf("%s::%s", req.PlayerName, req.NewState)
			// gameState.NotifyB <- req.NewState
		}
	}()

	// Check if the game is over
	if checkGameOver(gameState.State) {
		gameState.GameOver = true
		go func() {
			gameState.NotifyA <- "Game Over!"
			gameState.NotifyB <- "Game Over!"
		}()
	}

	return &game.UpdateGameStateResponse{Status: "State updated successfully"}, nil
}

func checkGameOver(state string) bool {
	// Placeholder logic to check if the game is over
	return false
}

func (s *GameServer) StreamGameUpdates(req *game.StreamGameUpdatesRequest, stream game.GameService_StreamGameUpdatesServer) error {
	s.mu.Lock()
	gameState, exists := s.createdGames[req.GameId]
	s.mu.Unlock()

	if !exists {
		return fmt.Errorf("game not found")
	}

	var notifyChannel chan string
	if req.PlayerName == gameState.PlayerA {
		notifyChannel = gameState.NotifyA
		fmt.Println("Notification channel created for Player A")
	} else if req.PlayerName == gameState.PlayerB {
		notifyChannel = gameState.NotifyB
		fmt.Println("Notification channel created for Player B")
	} else {
		return fmt.Errorf("player not found in game")
	}

	for msg := range notifyChannel {
		if err := stream.Send(&game.StreamGameUpdatesResponse{Update: msg}); err != nil {
			log.Printf("Error sending message to client: %v", err)
			return err
			// break
		}
	}
	return nil

	// for {
	// 	fmt.Print(notifyChannel)
	// 	fmt.Println()
	// 	select {
	// 	case msg := <-notifyChannel:
	// 		err := stream.Send(&game.StreamGameUpdatesResponse{Update: msg})
	// 		if err != nil {
	// 			log.Printf("Error sending update: %v", err)
	// 			return err
	// 		}
	// 	case <-time.After(20 * time.Second):
	// 		// Keep connection alive
	// 		err := stream.Send(&game.StreamGameUpdatesResponse{Update: "Waiting for updates..."})
	// 		if err != nil {
	// 			log.Printf("Error sending keep-alive: %v", err)
	// 			return err
	// 		}
	// 	}
	// }
}

func main() {
	server := grpc.NewServer()
	gameServer := &GameServer{
		createdGames: make(map[string]*GameState),
	}
	game.RegisterGameServiceServer(server, gameServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Server started on port :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
