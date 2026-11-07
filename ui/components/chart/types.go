package chart

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

type Config []Series

type Series struct {
	Key   string
	Label string
	Color string

	Theme *SeriesTheme

	Icon templ.Component
}

type SeriesTheme struct {
	Light string
	Dark  string
}

func (c Config) styleBlock(id string) string {
	hasColor := false
	for _, s := range c {
		if s.Color != "" || s.Theme != nil {
			hasColor = true
			break
		}
	}
	if !hasColor {
		return ""
	}
	var sb strings.Builder
	for _, theme := range []struct{ Name, Prefix string }{{"light", ""}, {"dark", ".dark "}} {
		fmt.Fprintf(&sb, "%s[data-chart=%s] {\n", theme.Prefix, id)
		for _, s := range c {
			color := s.Color
			if s.Theme != nil {
				themed := s.Theme.Light
				if theme.Name == "dark" {
					themed = s.Theme.Dark
				}
				if themed != "" {
					color = themed
				}
			}
			if color != "" {
				fmt.Fprintf(&sb, "  --color-%s: %s;\n", s.Key, color)
			}
		}
		sb.WriteString("}\n")
	}
	return sb.String()
}

type ContainerProps struct {
	ID         string
	Config     Config
	Class      string
	Attributes templ.Attributes
}

type StyleProps struct {
	ID     string
	Config Config
}
