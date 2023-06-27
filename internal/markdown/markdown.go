package markdown

import (
	"os"
	"path"
	"strings"

	"git2.borje.zone/public/clergo/internal/config"
	"git2.borje.zone/public/clergo/internal/shared"
)

func Generate(item shared.Item) string {
	outstr := "# " + item.TitleVar
	for _, tag := range item.TagsVar {
		outstr += "\n- " + tag
	}
	outstr += "\n\n"
	outstr += item.Body()
	return outstr
}

func Parse(mdData string) shared.Item {
	item := shared.Item{
		TitleVar: "",
		TagsVar:  []string{},
		BodyText: "",
	}

	mdData = strings.TrimRight(mdData, "\n")

	lines := strings.Split(mdData, "\n")
	bodylines := []string{}

	bodyFound := false

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			item.TitleVar = stripNTrim(line, "#")
			continue
		}
		if strings.HasPrefix(line, "-") {
			item.TagsVar = append(item.TagsVar, stripNTrim(line, "-"))
			continue
		}
		if strings.TrimSpace(line) == "" && !bodyFound {
			continue
		}
		bodyFound = true
		bodylines = append(bodylines, line)
	}
	item.BodyText = strings.Join(bodylines, "\n")
	return item
}

func ParseFile(filename string) (shared.Item, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return shared.Item{}, err
		// log.Err(err).Msgf("could not read file %s", filename)
	}

	return Parse(string(fileContent)), nil
}

func WriteFile(item shared.Item, filename string) error {
	mdData := Generate(item)
	return os.WriteFile(path.Join(config.CONFIG.LibPath, "notes", filename), []byte(mdData), os.ModePerm)
}

func stripNTrim(line, separator string) string {
	return strings.TrimSpace(strings.ReplaceAll(line, separator, ""))
}
