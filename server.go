package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
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
	NumbersA []int
	NumbersB []int
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

	num, _ := strconv.Atoi(req.NewState)
	if req.PlayerName == gameState.PlayerA {
		gameState.NumbersA = append(gameState.NumbersA, num)

	} else if req.PlayerName == gameState.PlayerB {
		gameState.NumbersB = append(gameState.NumbersB, num)
	}

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

	isGameOver, result := checkGameOver(gameState.NumbersA, gameState.NumbersB, gameState.PlayerA, gameState.PlayerB)
	if isGameOver {
		go func() {
			gameState.NotifyA <- result
			gameState.NotifyB <- result
		}()
		fmt.Println(result) // Output: "Player A wins!"
	}
	// Check if the game is over
	// if checkGameOver(gameState.NumbersA, gameState.NumbersB) {
	// 	gameState.GameOver = true
	// 	go func() {
	// 		gameState.NotifyA <- "Game Over!"
	// 		gameState.NotifyB <- "Game Over!"
	// 	}()
	// }

	return &game.UpdateGameStateResponse{Status: "State updated successfully"}, nil
}

func checkGameOver(numbersA []int, numbersB []int, playerA string, playerB string) (bool, string) {

	// Define the winning combinations
	winningCombinations := [][]int{
		{1, 2, 3}, // Top row
		{4, 5, 6}, // Middle row
		{7, 8, 9}, // Bottom row
		{1, 4, 7}, // Left column
		{2, 5, 8}, // Middle column
		{3, 6, 9}, // Right column
		{1, 5, 9}, // Diagonal top-left to bottom-right
		{3, 5, 7}, // Diagonal top-right to bottom-left
	}

	// Helper function to check if a slice contains a combination
	containsCombination := func(playerMoves []int, combination []int) bool {
		count := 0
		for _, pos := range combination {
			for _, move := range playerMoves {
				if pos == move {
					count++
					break
				}
			}
		}
		return count == 3
	}
	// Check if Player A wins
	for _, combo := range winningCombinations {
		if containsCombination(numbersA, combo) {
			return true, playerA + "::wins"
		}
	}

	// Check if Player B wins
	for _, combo := range winningCombinations {
		if containsCombination(numbersB, combo) {
			return true, playerB + "::wins"
		}
	}

	// Check for a draw
	if len(numbersA)+len(numbersB) == 9 {
		return true, "d::draw"
	}

	// Game is not over
	return false, ""
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
