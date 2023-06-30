package config

import (
	"strings"

	"github.com/adrg/xdg"
)

type Config struct {
	LibPath             string
	DotManagerAddCmd    string
	EnableDotManagerAdd bool
}

func initConfig(c Config) Config {
	c.LibPath = strings.ReplaceAll(c.LibPath, "~", xdg.Home)
	return c
}

var (
	CONFIG = initConfig(Config{
		LibPath:             "~/.clergo",
		DotManagerAddCmd:    "yadm add {{.Filename}}",
		EnableDotManagerAdd: true,
	})
)
