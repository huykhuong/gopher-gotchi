package ui

type Theme struct {
	Happy   string
	Neutral string
	Blink   string
	Hungry  string
	Dead    string
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
