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

	"golang.org/x/net/html"
)

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
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func FetchPageNode(url string) *html.Node {
	base, err := Request(url) // 9660 9341 6144153
	if err != nil {
		log.Fatalf(err.Error())
	}

	doc, err := html.Parse(strings.NewReader(base))
	if err != nil {
		log.Fatalf(err.Error())
	}

	return doc
}

func StringsToMarkup(s []string, uri string) []string {
	markedUp := make([]string, len(s))

	for i, tag := range s {
		encoded := url.QueryEscape(tag)
		markedUp[i] = fmt.Sprintf("[%s](%sposts?tags=%s&z=1)", tag, uri, encoded)
	}

	return markedUp
}

func searchForElement(node *html.Node, firstAttr, firstVal, secondAttr string, wg *sync.WaitGroup) []string {
	if node == nil {
		log.Println("Node does not exist or is incorrect")
		return nil
	}

	defer wg.Done()

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

func searchForChildElement(node *html.Node, parAttr, parAttrVal, chldTag, chldAttr string, wg *sync.WaitGroup) []string {
	if node == nil {
		log.Println("Node does not exist or is incorrect")
		return nil
	}

	defer wg.Done()

	var content []string

	parentNode := findNode(node, parAttr, parAttrVal)
	val := findValueInChildNode(parentNode, chldTag, chldAttr)
	if val != "" {
		content = append(content, val)
	}

	return content
}