package tray

import (
	"fmt"

	"github.com/getlantern/systray"
)

// This function tells macOS about the pet's status
func Update(level int, hunger int, mood string) {
	title := fmt.Sprintf("🐹 Lvl %d (%d%%)", level, hunger)
	systray.SetTitle(title)

	tooltip := fmt.Sprintf("Mood: %s", mood)
	systray.SetTooltip(tooltip)
}

func Init(onQuit func()) {
	systray.Run(func() {
		systray.SetTitle("🐹")
		systray.SetTooltip("Gopher-Gotchi is waking up...")

		mQuit := systray.AddMenuItem("Quit Gopher-Gotchi", "Stop the companion")

		go func ()  {
			<-mQuit.ClickedCh
			onQuit()
			systray.Quit()
		}()
	}, func() {})
}