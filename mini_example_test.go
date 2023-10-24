package mini_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/lucasepe/mini"
)

func ExampleNewConf() {
	data := `
width: 1600
subject: Computers,Education,Programming

[background]
art: shapes

[title]
text: Awesome Book!
color: #000000
font-size: 160
offset: 0.3
`
	conf := mini.NewConf()

	width := conf.Int("", "width", 1000)
	height := conf.Int("", "height", 1000)
	subject := conf.StringSlice("", "subject", []string{})
	watermark := conf.String("", "watermark", "")

	art := conf.Enum("background", "art", []string{
		"color-circles",
		"shapes",
		"circles-grid",
		"pixel-hole",
		"dots-wawe",
		"contour-line",
		"dot-line",
		"warp",
	})

	err := conf.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*width)
	fmt.Println(*height)
	for _, x := range *subject {
		fmt.Println(x)
	}
	fmt.Println(*watermark)
	fmt.Println(*art)

	// Output:
	// 1600
	// 1000
	// Computers
	// Education
	// Programming
	//
	// shapes
}
