package io

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"

	"github.com/lucas-ingemar/clergo/internal/config"
	"github.com/lucas-ingemar/clergo/internal/shared"
)

func addFileInDotManager(filename string) error {
	data := shared.DotManagerData{
		Filename: fullNotesFileName(filename),
	}
	tmpl, err := template.New("dotCmd").Parse(config.CONFIG.DotManagerAddCmd)
	if err != nil {
		return err
	}

	var cmdBuffer bytes.Buffer
	err = tmpl.Execute(&cmdBuffer, data)
	if err != nil {
		return err
	}

	cmdList := strings.Split(cmdBuffer.String(), " ")

	cmd := exec.Command(cmdList[0], cmdList[1:]...)
	if err := cmd.Run(); err != nil {
		return err
	}
}
