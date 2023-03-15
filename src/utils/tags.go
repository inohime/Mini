package utils

import (
	"sync"

	"golang.org/x/net/html"
)

type Tags struct {
	sync.Mutex
	Data map[string][]string
	Node *html.Node
	Sync sync.WaitGroup
}

func NewTag(node *html.Node) *Tags {
	return &Tags{
		Data: make(map[string][]string, 6),
		Node: node,
	}
}

func (t *Tags) FindArtistName() {
	t.Lock()
	defer t.Unlock()
	t.Data["artist"] = searchForElement(t.Node, "class", "tag-type-1", "data-tag-name", &t.Sync)
}

func (t *Tags) FindImageUrl() {
	t.Lock()
	defer t.Unlock()
	t.Data["image"] = searchForElement(t.Node, "class", "image-container note-container", "data-file-url", &t.Sync)
}

func (t *Tags) FindImageSource() {
	t.Lock()
	defer t.Unlock()
	t.Data["imgsrc"] = searchForChildElement(t.Node, "id", "post-info-source", "a", "href", &t.Sync)
}

func (t *Tags) FindCharacters() {
	t.Lock()
	defer t.Unlock()
	t.Data["characters"] = searchForElement(t.Node, "class", "tag-type-4", "data-tag-name", &t.Sync)
}

func (t *Tags) FindCopyright() {
	t.Lock()
	defer t.Unlock()
	t.Data["copyright"] = searchForElement(t.Node, "class", "tag-type-3", "data-tag-name", &t.Sync)
}

func (t *Tags) FindGeneralTags() {
	t.Lock()
	defer t.Unlock()
	t.Data["general"] = searchForElement(t.Node, "class", "tag-type-0", "data-tag-name", &t.Sync)
}
