package main

import "main/src/synthetic"

func main() {
	synthetic.Boot()
}

// package main

// import (
// 	"fmt"
// 	"net/url"
// 	"strings"
// )

// func stringsToMarkup(s []string, uri string) []string {
// 	markedUp := make([]string, len(s))
// 	for i, tag := range s {
// 		encoded := url.QueryEscape(tag)
// 		markedUp[i] = fmt.Sprintf("[%s](%sposts?tags=%s&z=1)", tag, uri, encoded)
// 	}
// 	return markedUp
// }

// func main() {
// 	data := make(map[string][]string, 1)
// 	data["characters"] = []string{"metako_(machikado_mazoku)", "yoshida_ryouko"}

// 	beep := stringsToMarkup(data["characters"], "https://danbooru.donmai.us/")

// 	x := strings.Join(beep[:], ", ")
// 	fmt.Println(x)

// 	// for _, x := range beep {
// 	// 	fmt.Println(x)
// 	// }
// }
