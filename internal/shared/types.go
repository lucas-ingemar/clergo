package shared

import (
	"fmt"
	"strings"

	"github.com/muesli/reflow/truncate"
)

type Item struct {
	TitleVar string
	// DescriptionVar string
	TagsVar  []string
	BodyText string
	Filename string
}

func (i Item) Title() string {
	return truncate.StringWithTail(i.TitleVar, 25, "...")
}

func (i Item) Description() string {
	return truncate.StringWithTail(strings.Join(i.TagsVar, ", "), 25, "...")
}

func (i Item) FilterValue() string {
	return fmt.Sprintf("%s: %s", strings.Join(i.TagsVar, ", "), i.TitleVar)
}

func (i Item) Body() string             { return i.BodyText }
func (i *Item) SetBodyText(text string) { i.BodyText = text }

type DotManagerData struct {
	Filename string
}
