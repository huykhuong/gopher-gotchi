package ui

import (
	"fmt"

	"github.com/fatih/color"
)

type Theme struct {
	Happy string
	Neutral string
	Blink string
	Hungry string
	Dead string
}

var Themes = map[string]Theme{
	"gopher": {
		Happy:   `  ( ^ ▽ ^ ) `,
		Neutral: `  ( ・ ▽ ・ ) `,
		Blink:   `  ( - ▽ - ) `,
		Hungry:  `  ( º﹃ º ) `,
		Dead:    `  ( x _ x ) `,
	},
	"robot": {
		Happy:   `  [ ^ _ ^ ] `,
		Neutral: `  [ o _ o ] `,
		Blink:   `  [ - _ - ] `,
		Hungry:  `  [ ﹃ _ ﹃ ] `,
		Dead:    `  [ # _ # ] `,
	},
	"cat": {
		Happy:   ` (= ^ ⩊ ^ =) `,
		Neutral: ` (= ・ ⩊ ・ =) `,
		Blink:   ` (= - ⩊ - =) `,
		Hungry:  ` (= º ⩊ º =) `,
		Dead:    ` (= x ⩊ x =) `,
	},
}

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