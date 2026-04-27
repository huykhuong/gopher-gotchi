package ui

import (
	"fmt"
	"strings"

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
	"diana": {
    Happy:   `  (  ˶^ ᴗ ^˶ ) `,
    Neutral: `  (  ˶• ᴗ •˶ ) `,
    Blink:   `  (  ˶- ᴗ -˶ ) `,
    Hungry:  `  (  ˶ó ᴗ ò˶ ) `,
    Dead:    `  (  ˶x ᴗ x˶ ) `,
},
}

func DrawPet(face string, level int, hunger int, mood string, messages []string, cpu int) {
	cyan := color.New(color.FgCyan).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Print("\033[H\033[2J")
	
	fmt.Println(cyan("⚡ [ PRAGMATA PROTOCOL: DIANA ] ⚡"))
	
	// CPU Pulse Meter
	cpuBar := strings.Repeat("█", cpu/10) + strings.Repeat("░", 10-(cpu/10))
	if cpu > 80 {
			cpuBar = red(cpuBar)
	}

	fmt.Printf("%s   %s [CPU: %d%%]\n", white(face), cyan("PULSE:"), cpu)
	fmt.Printf("               %s\n", cpuBar)

	fmt.Println(cyan("----------------------------------"))
	fmt.Printf("Level: %s | Hunger: %d%% | Mood: %s\n", white(level), hunger, cyan(mood))
	fmt.Println(cyan("----------------------------------"))

	// Messages
	for _, msg := range messages {
			fmt.Printf(cyan(">> ") + white(msg) + "\n")
	}
}