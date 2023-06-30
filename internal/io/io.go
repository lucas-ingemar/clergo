package io

import (
	"errors"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/lucas-ingemar/clergo/internal/config"
	"github.com/lucas-ingemar/clergo/internal/markdown"
	"github.com/lucas-ingemar/clergo/internal/shared"
)

func saveToFile(item shared.Item, NewFilename, OldFilename string) error {
	if NewFilename == "" {
		return errors.New("filename cannot be empty")
	}

	if NewFilename == OldFilename || OldFilename == "" {
		return markdown.WriteFile(item, NewFilename)
	}

	// FIXME: Move old file to "oldname.mdYYYYMMDDHHMM"
	// And remove the deletion
	err := os.Rename(OldFilename, OldFilename+"YYYYMMMDDHHMM")
	if err != nil {
		return err
	}

	return markdown.WriteFile(item, NewFilename)
}

func WriteFile(item *shared.Item) error {
	newFilename := strings.TrimSpace(strings.ToLower(strings.ReplaceAll(item.TitleVar, " ", "_")))
	if newFilename == "" {
		return errors.New("title cannot be empty")
	}
	newFilename += ".md"

	err := saveToFile(*item, newFilename, item.Filename)
	if err != nil {
		return err
	}

	item.Filename = newFilename
	return nil
}

func ReadFiles() (items []list.Item, err error) {
	notesPath := path.Join(config.CONFIG.LibPath, "notes")
	files, err := os.ReadDir(notesPath)
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		infoI, errI := files[i].Info()
		infoJ, errJ := files[j].Info()
		if errors.Join(errI, errJ) != nil {
			return false
		}
		return infoI.ModTime().After(infoJ.ModTime())
	})

	errorlist := []error{}
	for _, file := range files {
		if path.Ext(file.Name()) == ".md" {
			item, err := markdown.ParseFile(path.Join(notesPath, file.Name()))
			if err != nil {
				errorlist = append(errorlist, err)
				continue
			}
			item.Filename = file.Name()
			items = append(items, item)
		}
	}
	return items, errors.Join(errorlist...)
}

func DeleteFile(item shared.Item) error {
	// FIXME: Should probably add a trashcan or something
	notesPath := path.Join(config.CONFIG.LibPath, "notes")
	return os.Remove(path.Join(notesPath, item.Filename))
}
