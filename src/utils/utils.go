package utils

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"golang.org/x/net/html"
)

// RandomColor is a helper function that generates a random colour for interaction responses
//
// Return:
//   - Random int containing the value of a pre-defined colour
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

// Request is a helper function that sends HTTP requests to a given url with redirect debugging enabled
//
// Params:
//   - url: the requested url to get
//
// Return:
//   - The body formatted into a string and nil if no errors occur,
//     otherwise the error will be returned for both
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

// FetchPageNode is a helper function that attempts to get the document node for a page
//
// Params:
//   - url: the requested url to get
//
// Return:
//   - The document's root node and nil if no errors occur,
//     otherwise a nil and an error
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

// EncodeString is a helper function that escapes the tag given
//
// Params:
//   - tag: the tag given from generatecmd which is used to find images
//
// Return:
//   - A query escaped string (( :D -> %3AD ))
func EncodeString(tag string) string {
	return url.QueryEscape(tag)
}

// StringsToMarkup is a helper function that formats a slice of strings into a booru clickable link
// Note:
//   - URI can be either safe/danbooru
//
// Params:
//   - s: the slice of strings to be formatted into a clickable link
//   - uri: the booru uri to be used in the formatting process
//
// Return:
//   - An array of strings with url formatting applied
//
// Example:
//
//	> StringsToMarkup(["apples", "carrots", "potatoes"], "https://safebooru.donmai.us")
func StringsToMarkup(s []string, uri string) []string {
	markedUp := make([]string, len(s))

	for i, tag := range s {
		encoded := url.QueryEscape(tag)
		markedUp[i] = fmt.Sprintf("[%s](%s/posts?tags=%s&z=1)", tag, uri, encoded)
	}

	return markedUp
}

// EvictChars is a helper function that removes any characters until the string is 1024 characters
// (or less, depends on where the last comma is)
//
// Params:
//   - str: the string to evict characters from
//
// Return:
//   - A string with 1024 characters or less
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

// contentSearch is a helper function for element searches in the document's root node
//
// Params:
//   - node: the root node of the HTML tree to search
//   - firstAttr: the name of the first attribute of the element to match
//   - secondAttr: the name of the second attribute of the element to match
//   - findInner: allows searching through the inner content of an element (text)
//
// Return:
//   - A slice of strings found from the desired elements
func contentSearch(node *html.Node, firstAttr, firstVal, secondAttr string, findInner bool) (content []string) {
	if node == nil {
		fmt.Println("Node does not exist or is incorrect")
		return nil
	}

	doc := goquery.NewDocumentFromNode(node)
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		if val, exists := s.Attr(firstAttr); exists && strings.Contains(val, firstVal) {
			if secondAttr == "" && findInner {
				content = append(content, s.Text())
				return

			} else if val, exists := s.Attr(secondAttr); exists {
				content = append(content, val)
			}
		}
	})

	return
}

// searchForElement finds an element that has a two specific attributes and an attribute value
//
// Params:
//   - node: the root node of the HTML tree to search
//   - firstAttr: the name of the element attribute to match
//   - firstVal: the value of the first element attribute to match
//   - secondAttr: the name of the second attribute of the element to match
//
// Return:
//   - A string slice containing the contents of the element that matches the specified attributes and values
//
// Example:
//
//	> searchForElement(root, "class", "tag-type-1", "data-tag-name")
func searchForElement(node *html.Node, firstAttr, firstVal, secondAttr string) []string {
	return contentSearch(node, firstAttr, firstVal, secondAttr, false)
}

func searchForTextInElement(node *html.Node, firstAttr, firstVal string) string {
	contents := contentSearch(node, firstAttr, firstVal, "", true)
	if contents == nil {
		log.Println(
			color.HiRedString(
				"Failed to find text in element that met these requirements: %s",
				color.New(color.FgHiWhite).Sprintf(
					"%s %s",
					firstAttr,
					firstVal,
				),
			),
		)
		return ""
	}
	return contents[0]
}

// searchForChildElement finds a child element within a parent element
// that contains a specific attribute
//
// Params:
//   - node: the root node of the HTML tree to search
//   - parAttr: the name of the parent element attribute to match
//   - parAttrVal: the value of the parent element attribute to match
//   - childTag: the tag name of the child element to search for
//   - childAttr: the name of the child element attribute to retrieve the value from
//
// Return:
//   - A string slice containing the inner content (values) of a child element
//
// Example:
//
//	> searchForChildElement(node, "id", "my-fav-colour", "a", "href")
func searchForChildElement(node *html.Node, parAttr, parAttrVal, childTag, childAttr string) []string {
	if node == nil {
		log.Println("Node does not exist or is incorrect")
		return nil
	}

	var content []string

	parentNode := findNode(node, parAttr, parAttrVal)
	val := findValueInChildNode(parentNode, childTag, childAttr)
	if val != "" {
		content = append(content, val)
	}

	return content
}

// findNode is a helper function that searches for the first node that has a specific attribute and value
//
// Params:
//   - node: the root node of the HTML tree to start the search from
//   - nodeAttr: the name of the attribute to match
//   - nodeVal: the value of the attribute to match
//
// Return:
//   - The first node that has the specified attribute and value
//
// Example:
//
//	> findNode(node, "class", "search-tag")
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

// findValueInChildNode is a helper function that searches for a value
// given a tag and an attribute of a child node
//
// Params:
//   - node: the root node of the HTML tree to start the search from
//   - childTag: the tag of the element to match
//   - childAttr: the name of the attribute to match
//
// Return:
//   - A string that contains the inner content of a node
//     or empty if the node has no value
//
// Example:
//
//	> findValueInChildNode(parentNode, "a", "href")
func findValueInChildNode(parent *html.Node, childTag, childAttr string) string {
	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Data == childTag {
			for _, attr := range c.Attr {
				if attr.Key == childAttr {
					return attr.Val
				}
			}
		}
	}
	return ""
}
