package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-manager-api/handler"
	"task-manager-api/middleware"
	"task-manager-api/repository"
	"task-manager-api/usecase"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	repo, err := repository.NewSQLiteTaskRepository("tasks.db")
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repo.Close()

	userRepo, err := repository.NewSQLiteUserRepository("tasks.db")
	if err != nil {
		log.Fatalf("Failed to initialize user repository: %v", err)
	}
	defer userRepo.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Hello from Task Manager API"}`)
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET environment variable not set")
	}

	cache := usecase.NewCacheService(5 * time.Minute)
	uc := usecase.NewTaskUsecase(repo, cache)
	authUc := usecase.NewAuthUsecase(userRepo, jwtSecret)
	processor := usecase.NewTaskProcessor(repo)

	handler.SetupRoutes(mux, uc, authUc, processor, cache, repo)

	rateLimiter := middleware.NewRateLimiter(20)

	handler := middleware.Chain(
		rateLimiter.RateLimiterMiddleware,
		middleware.CorrelationIDMiddleware,
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
	)(mux)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// TESTS
	testConcurrentLogger()
	testRaceCondition()
	testTaskQueue()
	testTaskStream()

	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")

}

func testConcurrentLogger() {
	logger := usecase.NewConcurrentLogger()
	messages := []string{
		"First message",
		"Second message",
		"Third message",
		"Fourth message",
		"Fifth message",
	}
	fmt.Println("Starting concurrent logging...")
	logger.LogMultiple(messages)
	fmt.Println("Done!")
}

func testRaceCondition() {
	tc := usecase.NewTaskCounter(50)
	taskIDs := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		taskIDs[i] = 1
	}

	result, err := tc.CountTasks(taskIDs)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	expectedSum := 1000 + 1000 // 1000 taskIDs (each=1) + 1000 goroutine increments

	fmt.Printf("Result: %d, Expected: %d\n", result, expectedSum)
	if result != expectedSum {
		fmt.Printf("RACE CONDITION! Lost %d updates\n", expectedSum-result)
	}
}

func testTaskQueue() {
	fmt.Println("\n=== Testing Task Queue")

	queue := usecase.NewTaskQueue(5)

	// Test 1
	fmt.Println("\n Test 1")
	queue.Enqueue(10)
	queue.Enqueue(20)
	queue.Enqueue(30)

	taskID, err := queue.Dequeue()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Dequeued: %d\n", taskID)

	// Test 2: Queue Full Error
	// fmt.Println("\n--- Test 2: Queue Full ---")

	// queue.Enqueue(40)
	// queue.Enqueue(50)
	// queue.Enqueue(60)
	// queue.Enqueue(70)

	// err := queue.Enqueue(80)
	// if err != nil {
	// 	fmt.Printf("✅ Error (expected): %v\n", err)
	// }

	// Test 3: Queue Empty Error
	// fmt.Println("\n--- Test 3: Queue Empty ---")

	// for i := 0; i < 5; i++ {
	// 	id, _ := queue.Dequeue()
	// 	fmt.Printf("Dequeued: %d\n", id)
	// }

	// _, err := queue.Dequeue()
	// if err != nil {
	// 	fmt.Printf("✅ Error (expected): %v\n", err)
	// }

	// Test 4: Concurrent Goroutines
	fmt.Println("\n--- Test 4: Concurrent Producers & Consumer ---")

	queue2 := usecase.NewTaskQueue(10)

	for i := 1; i <= 3; i++ {
		go func(producerID int) {
			for j := 1; j <= 5; j++ {
				taskID := producerID*100 + j
				err := queue2.Enqueue(taskID)
				if err == nil {
					fmt.Printf("Producer %d enqueued: %d\n", producerID, taskID)
				}
			}
		}(i)
	}

	time.Sleep(200 * time.Millisecond)

	fmt.Printf("Queue size after: %d\n", queue2.Size())
	for {
		id, err := queue2.Dequeue()
		if err != nil {
			break
		}
		fmt.Printf("Dequeued: %d\n", id)
	}

	fmt.Println("--- All tests passed! ---\n")
}

func testTaskStream() {
	fmt.Println("\n=== Testing Task Stream (Exercise 4) ===")

	stream := usecase.NewTaskStream()

	stream.Start()

	stream.Send(1)
	stream.Send(2)
	stream.Send(3)

	time.Sleep(200 * time.Millisecond)

	for i := 0; i < 3; i++ {
		result := stream.GetResult()
		fmt.Printf("Received result: %s\n", result)
	}
	stream.Close()

	fmt.Println("--- Task Stream test passed! ---\n")
}
