package utils

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"golang.org/x/net/html"
)

type SafeMap struct {
	sync.RWMutex
	_map map[string]interface{}
}

func (sm *SafeMap) Read(key string) (value interface{}, ok bool) {
	sm.RLock()
	defer sm.RUnlock()

	value, ok = sm._map[key]

	return
}

func (sm *SafeMap) Write(key string, value interface{}) {
	sm.Lock()
	defer sm.Unlock()

	sm._map[key] = value
}

func RandomColor() int {
	rand.Seed(time.Now().UnixNano())

	colors := []int{
		0xFF1567, // razzmatazz 	-> vibrant pink
		0x9B74FF, // medium purple 	-> vibrant purple
		0x0099FF, // azure radiance -> vibrant blue
		0xFFDDE4, // pig pink 		-> light pink
		0xC0A3FF, // biloba flower 	-> light purple
		0x8CA9FF, // portage 		-> light blue
	}

	return colors[rand.Intn(len(colors))]
}

func Request(url string) (string, error) {
	res, err := http.Get(url)
	if res.StatusCode != 200 {
		return res.Status, fmt.Errorf(res.Status)
	}

	if err != nil {
		return err.Error(), err
	}
	defer res.Body.Close()

	// redirect check (debug)
	log.Println(
		color.HiGreenString("URL redirect:"),
		color.HiWhiteString(res.Request.URL.String()),
	)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err.Error(), err
	}

	return string(body), err
}

func FetchPageNode(url string) (*html.Node, error) {
	base, err := Request(url)
	if err != nil {
		log.Printf(
			color.HiRedString("Failed to parse URL: ")+"%s",
			color.HiWhiteString(err.Error()),
		)
		return nil, err
	}

	doc, _ := html.Parse(strings.NewReader(base))
	if err != nil {
		log.Printf(
			color.HiRedString("Failed to parse base: ")+"%s",
			color.HiWhiteString(err.Error()),
		)
		return nil, err
	}

	return doc, nil
}

func EncodeString(tag string) string {
	return url.QueryEscape(tag)
}

func StringsToMarkup(s []string, uri string) []string {
	markedUp := make([]string, len(s))

	for i, tag := range s {
		encoded := url.QueryEscape(tag)
		markedUp[i] = fmt.Sprintf("[%s](%s/posts?tags=%s&z=1)", tag, uri, encoded)
	}

	return markedUp
}

func EvictChars(str string) string {
	if len(str) <= 1024 {
		return str
	}

	for i := 0; i < len(str); i++ {
		if i >= 1024 { // 84 words
			str = strings.Join(
				strings.Split(
					str[:strings.LastIndex(str[:i], ",")], ", "),
				", ",
			)
			break
		}
	}

	return str
}

func searchForTextInElement(node *html.Node, firstAttr, firstVal string) string {
	if node == nil {
		log.Println("Node does not exist or is incorrect")
		return ""
	}

	text := ""
	var findNext func(*html.Node)
	findNext = func(n *html.Node) {
		if n.Type == html.ElementNode {
			found := false
			for _, attr := range n.Attr {
				if attr.Key == firstAttr {
					if strings.Contains(attr.Val, firstVal) {
						found = true
						break
					}
				}
			}

			if found {
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Type == html.TextNode {
						text = c.Data
						break
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNext(c)
		}
	}

	findNext(node)

	return text
}

func searchForElement(node *html.Node, firstAttr, firstVal, secondAttr string) []string {
	if node == nil {
		log.Println("Node does not exist or is incorrect")
		return nil
	}

	var content []string
	var findNext func(*html.Node)
	findNext = func(n *html.Node) {
		if n.Type == html.ElementNode {
			found := false
			for _, attr := range n.Attr {
				if attr.Key == firstAttr {
					if strings.Contains(attr.Val, firstVal) {
						found = true
					}
				}

				if attr.Key == secondAttr {
					if found {
						content = append(content, attr.Val)
						break
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findNext(c)
		}
	}

	findNext(node)

	return content
}

func searchForChildElement(node *html.Node, parAttr, parAttrVal, chldTag, chldAttr string) []string {
	if node == nil {
		log.Println("Node does not exist or is incorrect")
		return nil
	}

	var content []string

	parentNode := findNode(node, parAttr, parAttrVal)
	val := findValueInChildNode(parentNode, chldTag, chldAttr)
	if val != "" {
		content = append(content, val)
	}

	return content
}

func findNode(node *html.Node, nodeAttr, nodeVal string) *html.Node {
	if node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == nodeAttr && attr.Val == nodeVal {
				return node
			}
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		res := findNode(c, nodeAttr, nodeVal)
		if res != nil {
			return res
		}
	}

	return nil
}

func findValueInChildNode(parent *html.Node, chldTag, chldAttr string) string {
	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == chldTag {
			for _, attr := range c.Attr {
				if attr.Key == chldAttr {
					return attr.Val
				}
			}
		}
	}

	return ""
}
