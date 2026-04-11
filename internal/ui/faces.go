package ui

import (
	"fmt"

	"github.com/fatih/color"
)

const (
	FaceHappy   = `  ( ^ ▽ ^ ) `
	FaceNeutral = `  ( ・ ▽ ・ ) `
	FaceBlink   = `  ( - ▽ - ) `
	FaceHungry  = `  ( º﹃ º ) `
	FaceLevelUp = `  ( ✧ ▽ ✧ ) `
	FaceDead    = `  ( x _ x ) `
)

func DrawPet(face string, level int, hunger int, mood string, messages[]string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Print("\033[H\033[2J") // Clear screen
	fmt.Println("\n" + bold("--- GOPHER-GOTCHI ---"))
	fmt.Println(face)
	fmt.Println("---------------------")
	fmt.Printf("Level:  %s\n", yellow(level))
	fmt.Printf("Hunger: %d%%\n", hunger)
	fmt.Printf("Mood:   %s\n", cyan(mood))
	fmt.Println("---------------------")

	fmt.Println(bold("\n[ Activity Log ]"))
	if len(messages) == 0 {
		fmt.Println("No recent activity...")
	}
	for _, msg := range messages {
		fmt.Printf("> %s\n", msg)
	}
}