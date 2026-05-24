package brain

import "fmt"

type Task func() error

// StartWorkerPool spawns a worker that drains taskChan until it is closed.
// The returned channel is closed when the worker has finished draining,
// so callers can wait for in-flight work before final shutdown.
func StartWorkerPool(taskChan <-chan Task) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		for task := range taskChan {
			if err := task(); err != nil {
				fmt.Printf("⚠️ Background Worker Error: %v\n", err)
			}
		}
	}()

	return done
}
