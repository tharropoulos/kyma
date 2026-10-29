package cmd

import (
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	_ "golang.org/x/image/webp"
	"golang.org/x/term"

	"github.com/museslabs/kyma/internal/config"
	"github.com/museslabs/kyma/internal/export"
	"github.com/museslabs/kyma/internal/img"
	"github.com/museslabs/kyma/internal/img/placeholder"
	"github.com/museslabs/kyma/internal/logger"
	"github.com/museslabs/kyma/internal/textsize"
	"github.com/museslabs/kyma/internal/tui"
)

var (
	exportOutput    string
	exportSize      string
	exportScale     float64
	exportFonts     []string
	exportWallpaper string
	exportOpacity   float64
)

func init() {
	exportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Path to write the PDF to (default: the presentation's name, ending in .pdf)")
	exportCmd.Flags().StringVar(&exportSize, "size", "", "Size of the screen to draw slides on, in cells, such as 160x45 (default: this terminal's)")
	exportCmd.Flags().Float64Var(&exportScale, "scale", 1, "How many times larger than this terminal's cells to draw them")
	exportCmd.Flags().StringArrayVar(&exportFonts, "font", nil, "Font file to draw text with, such as the terminal's own; repeat it for fonts to fall back on")
	exportCmd.Flags().StringVar(&exportWallpaper, "wallpaper", "", "Image to show through the background, as behind a transparent terminal")
	exportCmd.Flags().Float64Var(&exportOpacity, "opacity", 0.9, "Opacity of the background over --wallpaper, from 0 to 1")
	exportCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to config file")
}

var exportCmd = &cobra.Command{
	Use:   "export <filename>",
	Short: "Export every slide of a presentation to a PDF",
	Long: "Export every slide of a presentation to a PDF, a page a slide, drawn the way this terminal shows them:\n" +
		"at its size, in its colors and with its scaled headings, with images at their full resolution.\n" +
		"Run it in the terminal you present from, and pass the font that terminal uses with --font.",
	Args:         markdownArg,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := logger.Load(logPath); err != nil {
			return fmt.Errorf("failed to initialize slog: %w", err)
		}

		if err := config.Load(configPath); err != nil {
			return err
		}

		// Ask the terminal before parsing, since slides decide how to draw
		// their headings when they are created.
		textsize.Detect()
		colors := export.QueryColors()

		// Images are drawn from their own pixels, not uploaded to a
		// terminal, whatever this one can show.
		placeholder.Keep()

		filename := args[0]
		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		root, err := parseSlides(string(data))
		if err != nil {
			return err
		}

		t, err := exportTerminal(colors)
		if err != nil {
			return err
		}

		output := exportOutput
		if output == "" {
			output = strings.TrimSuffix(filename, filepath.Ext(filename)) + ".pdf"
		}

		n, err := exportDeck(output, root, t)
		if err != nil {
			_ = os.Remove(output)
			return err
		}

		slog.Info("Exported presentation", "slides", n, "output", output)
		fmt.Fprintf(os.Stderr, "Exported %d slides to %s\n", n, output)
		return nil
	},
}

// exportTerminal works out what the slides are drawn as: this terminal, or
// what the flags say instead.
func exportTerminal(colors export.Colors) (export.Terminal, error) {
	t := export.Terminal{
		Fg:      colors.Fg,
		Bg:      colors.Bg,
		Palette: colors.Palette,
		Opacity: exportOpacity,
	}

	if exportSize != "" {
		if _, err := fmt.Sscanf(exportSize, "%dx%d", &t.Cols, &t.Rows); err != nil || t.Cols <= 0 || t.Rows <= 0 {
			return t, fmt.Errorf("expected a size in cells such as 160x45, got: %v", exportSize)
		}
	} else {
		cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			return t, errors.New("not running in a terminal, so give the size to draw slides at with --size")
		}
		t.Cols, t.Rows = cols, rows
	}

	if exportScale <= 0 {
		return t, fmt.Errorf("expected a scale above 0, got: %v", exportScale)
	}
	cw, ch := img.CellSize()
	t.CellWidth = max(1, int(math.Round(float64(cw)*exportScale)))
	t.CellHeight = max(1, int(math.Round(float64(ch)*exportScale)))

	for _, path := range exportFonts {
		family, err := export.LoadFamily(path)
		if err != nil {
			return t, fmt.Errorf("failed to load font: %w", err)
		}
		t.Fonts = append(t.Fonts, family)
	}

	if exportWallpaper != "" {
		f, err := os.Open(exportWallpaper)
		if err != nil {
			return t, err
		}
		defer f.Close()
		wallpaper, _, err := image.Decode(f)
		if err != nil {
			return t, fmt.Errorf("failed to read wallpaper %s, which has to be a PNG, JPEG, GIF or WebP: %w", exportWallpaper, err)
		}
		t.Wallpaper = wallpaper
	}

	return t, nil
}

// exportDeck draws every slide from root on into a PDF at path, and returns
// how many there were.
func exportDeck(path string, root *tui.Slide, t export.Terminal) (int, error) {
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	deck, err := export.NewDeck(f, t)
	if err != nil {
		return 0, err
	}

	n := 0
	for slide := root; slide != nil; slide = slide.Next {
		n++
		fmt.Fprintf(os.Stderr, "\rDrawing slide %d", n)

		view := slide.Render(t.Cols, t.Rows)
		if strings.Count(view, "\n") >= t.Rows {
			fmt.Fprintf(os.Stderr, "\rSlide %d doesn't fit on a %dx%d screen and is cut off\n", n, t.Cols, t.Rows)
		}
		if err := deck.Add(textsize.Render(view, t.Cols, t.Rows)); err != nil {
			return n, err
		}
	}
	fmt.Fprint(os.Stderr, "\r\x1b[K")

	if err := deck.Close(); err != nil {
		return n, err
	}
	return n, f.Close()
}
