package utils

import (
	"fmt"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Tags struct {
	Data SafeMap
	Node *html.Node
}

func NewTag(node *html.Node) *Tags {
	return &Tags{
		Data: NewSafeMap(),
		Node: node,
	}
}

func (t *Tags) FindArtistName(swg *sync.WaitGroup) {
	defer swg.Done()

	artistName := searchForElement(
		t.Node,
		"class",
		"tag-type-1",
		"data-tag-name",
	)
	if artistName == nil {
		t.Data.Write("artist", []string{"no artist found"})
		return
	}

	t.Data.Write("artist", artistName)
}

func (t *Tags) FindImageUrl(swg *sync.WaitGroup) {
	defer swg.Done()

	// sometimes, this will fail because we don't have permissions to view the image
	imageContent := searchForElement(
		t.Node,
		"class",
		"image-container note-container",
		"data-file-url",
	)
	if imageContent == nil {
		// image doesn't exist so throw this instead
		t.Data.Write("image", []string{
			"https://pbs.twimg.com/media/FZ8WhlkXkAAug7p?format=png&name=large",
		})
		return
	}

	t.Data.Write("image", imageContent)
}

func (t *Tags) FindImageSource(swg *sync.WaitGroup) {
	defer swg.Done()

	// find the source link associated with the image
	check := searchForChildElement(
		t.Node,
		"id",
		"post-info-source",
		"a",
		"href",
	)
	if check == nil {
		// try just getting the content inside of post-info-source.
		// if the text after (Source: ) is empty, go to next, otherwise, extract it
		// if it doesn't exist, return "no source"
		sourceContent := searchForTextInElement(
			t.Node,
			"id",
			"post-info-source",
		)

		// in post-info-source, the inner text content is formatted like so: "Source: ..."
		// where "..." is the actual content
		// the first match is "Source:", the second match is anything found after it
		re := regexp.MustCompile(`^Source:\s+(.*)$`)
		match := re.FindStringSubmatch(sourceContent)
		if match[1] == "" {
			t.Data.Write("imgsrc", []string{"no source found"})
			return
		}

		// found a non-link source, still ok!
		t.Data.Write("imgsrc", []string{match[1]})
		return
	}

	t.Data.Write("imgsrc", check)
}

func (t *Tags) FindCharacters(swg *sync.WaitGroup) {
	defer swg.Done()

	characterContents := searchForElement(
		t.Node,
		"class",
		"tag-type-4",
		"data-tag-name",
	)
	if characterContents == nil {
		t.Data.Write("characters", []string{"original"})
		return
	}

	t.Data.Write("characters", characterContents)
}

func (t *Tags) FindCopyright(swg *sync.WaitGroup) {
	defer swg.Done()

	copyrightContent := searchForElement(
		t.Node,
		"class",
		"tag-type-3",
		"data-tag-name",
	)
	if copyrightContent == nil {
		// highly unlikely that this occurs, but we should have it just in case
		t.Data.Write("copyright", []string{"no copyright found"})
		return
	}

	t.Data.Write("copyright", copyrightContent)
}

func (t *Tags) FindGeneralTags(swg *sync.WaitGroup) {
	defer swg.Done()

	generalContents := searchForElement(
		t.Node,
		"class",
		"tag-type-0",
		"data-tag-name")
	if generalContents == nil {
		t.Data.Write("general", []string{"no tags found"})
		return
	}

	t.Data.Write("general", generalContents)
}

func MakeTagsTable(s []string, length *int) string {
	var tags strings.Builder

	sliceLen := 10
	if len(s) < 10 {
		sliceLen = len(s)
	}

	for _, tag := range s[:sliceLen] {
		if len(tag) > *length {
			*length = len(tag)
		}
	}

	for _, tag := range s[:sliceLen] {
		tags.WriteString(
			fmt.Sprintf(
				"%[2]*[1]s\n", tag, (*length+len(tag))/2,
			),
		)
	}

	return tags.String()
}
