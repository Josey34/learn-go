package usecase

import (
	"fmt"
	"time"
)

type ConcurrentLogger struct {
	counter int
}

func NewConcurrentLogger() *ConcurrentLogger {
	return &ConcurrentLogger{counter: 0}
}

func (cl *ConcurrentLogger) LogAsync(message string) {
	go func(msg string) {
		time.Sleep(100 * time.Millisecond)
		cl.counter++
		fmt.Printf("[Goroutine %d] %s - %s\n",
			cl.counter, msg, time.Now().Format("2006-01-02 15:04:05"))
	}(message)
}

func (cl *ConcurrentLogger) LogMultiple(messages []string) {
	for _, message := range messages {
		cl.LogAsync(message)
	}

	time.Sleep(500 * time.Millisecond)
}
