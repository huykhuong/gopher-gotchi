package brain

import "fmt"

type Task func() error

func StartWorkerPool(taskChan <-chan Task) {
	go func() {
		for task := range taskChan {
			err := task()
			if err != nil {
				fmt.Printf("⚠️ Background Worker Error: %v\n", err)
			}
		}
	}()
}