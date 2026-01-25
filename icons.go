// Package brandkit provides embedded brand icons.
//
// Icons are embedded at compile time and can be retrieved by brand name.
// Each brand has three variants: color, white (for dark backgrounds), and original.
//
// Example:
//
//	svg, err := brandkit.GetIconWhite("aws")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// svg contains the white AWS icon SVG bytes
package brandkit

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

//go:embed brands/*/icon_white.svg brands/*/icon_color.svg brands/*/icon_orig.svg
var brandsFS embed.FS

// IconVariant represents the icon color variant.
type IconVariant string

const (
	// IconVariantWhite is a white foreground icon for dark backgrounds.
	IconVariantWhite IconVariant = "white"
	// IconVariantColor is the full color icon.
	IconVariantColor IconVariant = "color"
	// IconVariantOrig is the original unmodified icon.
	IconVariantOrig IconVariant = "orig"
)

// GetIcon retrieves an icon by brand name and variant.
// Returns the SVG content as bytes.
func GetIcon(brand string, variant IconVariant) ([]byte, error) {
	filename := fmt.Sprintf("icon_%s.svg", variant)
	filepath := path.Join("brands", brand, filename)
	return brandsFS.ReadFile(filepath)
}

// GetIconWhite retrieves the white variant icon for dark backgrounds.
func GetIconWhite(brand string) ([]byte, error) {
	return GetIcon(brand, IconVariantWhite)
}

// GetIconColor retrieves the full color icon.
func GetIconColor(brand string) ([]byte, error) {
	return GetIcon(brand, IconVariantColor)
}

// GetIconOrig retrieves the original unmodified icon.
func GetIconOrig(brand string) ([]byte, error) {
	return GetIcon(brand, IconVariantOrig)
}

// ListIcons returns all available brand names.
func ListIcons() ([]string, error) {
	entries, err := fs.ReadDir(brandsFS, "brands")
	if err != nil {
		return nil, err
	}

	var brands []string
	for _, entry := range entries {
		if entry.IsDir() {
			brands = append(brands, entry.Name())
		}
	}
	sort.Strings(brands)
	return brands, nil
}

// IconExists checks if a brand icon exists.
func IconExists(brand string) bool {
	_, err := GetIconWhite(brand)
	return err == nil
}

// NormalizeIconName converts common aliases to brandkit names.
// For example, "golang" -> "go", "postgresql" -> "postgres".
func NormalizeIconName(name string) string {
	name = strings.ToLower(name)
	aliases := map[string]string{
		"golang":     "go",
		"postgresql": "postgres",
		"k8s":        "kubernetes",
		"gcloud":     "gcp",
	}
	if normalized, ok := aliases[name]; ok {
		return normalized
	}
	return name
}
