package utils

import (
	"regexp"
	"sync"

	"golang.org/x/net/html"
)

type Tags struct {
	Data SafeMap
	Node *html.Node
}

func NewTag(node *html.Node) *Tags {
	return &Tags{
		Data: SafeMap{_map: make(map[string]interface{}, 6)},
		Node: node,
	}
}

func (t *Tags) FindArtistName(swg *sync.WaitGroup) {
	defer swg.Done()

	check := searchForElement(t.Node, "class", "tag-type-1", "data-tag-name")
	if check == nil {
		t.Data.Write("artist", []string{"no artist found"})
		return
	}

	t.Data.Write("artist", check)
}

func (t *Tags) FindImageUrl(swg *sync.WaitGroup) {
	defer swg.Done()
	// sometimes, this will fail because we don't have permissions to view the image
	check := searchForElement(t.Node, "class", "image-container note-container", "data-file-url")
	if check == nil {
		// image doesn't exist so throw this instead
		t.Data.Write("image", []string{
			"https://pbs.twimg.com/media/FZ8WhlkXkAAug7p?format=png&name=large",
		})
		return
	}

	t.Data.Write("image", check)
}

func (t *Tags) FindImageSource(swg *sync.WaitGroup) {
	defer swg.Done()

	check := searchForChildElement(t.Node, "id", "post-info-source", "a", "href")

	if check == nil {
		// try just getting the content inside of post-info-source.
		// if the text after (Source: ) is empty, go to next, otherwise, extract it
		// if it doesn't exist, return "no source"
		sourceCheck := searchForTextInElement(t.Node, "id", "post-info-source")
		re := regexp.MustCompile(`^Source:\s+(.*)$`)
		match := re.FindStringSubmatch(sourceCheck)

		if match[1] == "" {
			t.Data.Write("imgsrc", []string{"no source found"})
			return
		}

		t.Data.Write("imgsrc", []string{match[1]})
		return
	}

	t.Data.Write("imgsrc", check)
}

func (t *Tags) FindCharacters(swg *sync.WaitGroup) {
	defer swg.Done()

	check := searchForElement(t.Node, "class", "tag-type-4", "data-tag-name")
	if check == nil {
		t.Data.Write("characters", []string{"original"})
		return
	}

	t.Data.Write("characters", check)
}

func (t *Tags) FindCopyright(swg *sync.WaitGroup) {
	defer swg.Done()

	t.Data.Write("copyright", searchForElement(t.Node, "class", "tag-type-3", "data-tag-name"))
}

func (t *Tags) FindGeneralTags(swg *sync.WaitGroup) {
	defer swg.Done()

	check := searchForElement(t.Node, "class", "tag-type-0", "data-tag-name")
	if check == nil {
		t.Data.Write("general", []string{"no tags found"})
		return
	}

	t.Data.Write("general", check)
}
