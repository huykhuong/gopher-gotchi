package brain

import (
	"context"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

var ErrorKeywords = []string{
	"error", "exception",
}

func StartClipboardWatcher(ctx context.Context, eventChan chan<- DataEvent) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		lastContent := ""

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				text, err := clipboard.ReadAll()
				if err != nil || text == "" || text == lastContent {
					continue
				}

				lastContent = text

				lowerText := strings.ToLower(text)
				for _, keyword := range ErrorKeywords {
					if strings.Contains(lowerText, keyword) {
						eventChan <- DataEvent{
							Type:    ClipboardErrorDetected,
							Payload: text,
						}
						break
					}
				}
			}
		}
	}()
}
