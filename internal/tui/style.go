package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/goccy/go-yaml"
)

type Theme string

const (
	ThemeAscii      Theme = styles.AsciiStyle
	ThemeAuto       Theme = styles.AutoStyle
	ThemeDark       Theme = styles.DarkStyle
	ThemeDracula    Theme = styles.DraculaStyle
	ThemeTokyoNight Theme = styles.TokyoNightStyle
	ThemeLight      Theme = styles.LightStyle
	ThemeNoTTY      Theme = styles.NoTTYStyle
	ThemePink       Theme = styles.PinkStyle
)

type SlideStyle struct {
	LipGlossStyle lipgloss.Style
	Theme         Theme
}

type StyleConfig struct {
	Layout      lipgloss.Style  `yaml:"layout"`
	Border      lipgloss.Border `yaml:"border"`
	BorderColor string          `yaml:"border_color"`
	Theme       Theme           `yaml:"theme"`
}

func (s *StyleConfig) UnmarshalYAML(bytes []byte) error {
	aux := struct {
		Layout      string `yaml:"layout"`
		Border      string `yaml:"border"`
		BorderColor string `yaml:"border_color"`
		Theme       string `yaml:"theme"`
	}{}

	var err error

	if err = yaml.Unmarshal(bytes, &aux); err != nil {
		return err
	}

	s.Layout, err = getLayout(aux.Layout)
	if err != nil {
		return err
	}

	s.Border = getBorder(aux.Border)
	s.BorderColor = aux.BorderColor
	s.Theme = getTheme(aux.Theme)

	return nil
}

func (s StyleConfig) ApplyStyle(width, height int) SlideStyle {
	borderColor := "#9999CC" // Blueish
	if s.BorderColor != "" {
		borderColor = s.BorderColor
	}

	style := s.Layout.
		Border(s.Border).
		BorderForeground(lipgloss.Color(borderColor)).
		Width(width - 4).
		Height(height - 2)

	return SlideStyle{
		LipGlossStyle: style,
		Theme:         s.Theme,
	}
}

func getBorder(border string) lipgloss.Border {
	switch border {
	case "rounded":
		return lipgloss.RoundedBorder()
	case "double":
		return lipgloss.DoubleBorder()
	case "thick":
		return lipgloss.ThickBorder()
	case "hidden":
		return lipgloss.HiddenBorder()
	case "block":
		return lipgloss.BlockBorder()
	case "innerHalfBlock":
		return lipgloss.InnerHalfBlockBorder()
	case "outerHalfBlock":
		return lipgloss.OuterHalfBlockBorder()
	case "normal":
		fallthrough
	default:
		return lipgloss.NormalBorder()
	}
}

func getLayout(layout string) (lipgloss.Style, error) {
	style := lipgloss.NewStyle()

	layout = strings.TrimSpace(layout)
	if layout == "" {
		return style, nil
	}

	positions := strings.Split(layout, ",")
	if len(positions) > 2 {
		return style, fmt.Errorf("invalid layout configuration: %s", layout)
	}

	p1, err := getLayoutPosition(positions[0])
	if err != nil {
		return style, err
	}

	if len(positions) == 1 {
		return style.Align(p1, p1), nil
	}

	p2, err := getLayoutPosition(positions[1])
	if err != nil {
		return style, err
	}

	return style.Align(p1, p2), nil
}

func getLayoutPosition(p string) (lipgloss.Position, error) {
	switch strings.TrimSpace(p) {
	case "center":
		return lipgloss.Center, nil
	case "left":
		return lipgloss.Left, nil
	case "right":
		return lipgloss.Right, nil
	case "top":
		return lipgloss.Top, nil
	case "bottom":
		return lipgloss.Bottom, nil
	default:
		return 0, fmt.Errorf("invalid position: %s", strings.TrimSpace(p))
	}
}

func getTheme(theme string) Theme {
	switch theme {
	case "ascii":
		return ThemeAscii
	case "auto":
		return ThemeAuto
	case "dark":
		return ThemeDark
	case "dracula":
		return ThemeDracula
	case "tokyo-night":
		return ThemeTokyoNight
	case "light":
		return ThemeLight
	case "notty":
		return ThemeNoTTY
	case "pink":
		return ThemePink
	default:
		return ThemeDark
	}
}
