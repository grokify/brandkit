package brandkit

import (
	"testing"
)

func TestListIcons(t *testing.T) {
	brands, err := ListIcons()
	if err != nil {
		t.Fatalf("ListIcons() error: %v", err)
	}
	if len(brands) == 0 {
		t.Fatal("ListIcons() returned empty list")
	}
	t.Logf("Found %d brands: %v", len(brands), brands)
}

func TestGetIconWhite(t *testing.T) {
	tests := []string{"aws", "github", "docker", "kubernetes"}
	for _, brand := range tests {
		svg, err := GetIconWhite(brand)
		if err != nil {
			t.Errorf("GetIconWhite(%q) error: %v", brand, err)
			continue
		}
		if len(svg) == 0 {
			t.Errorf("GetIconWhite(%q) returned empty SVG", brand)
		}
		t.Logf("GetIconWhite(%q) = %d bytes", brand, len(svg))
	}
}

func TestIconExists(t *testing.T) {
	if !IconExists("aws") {
		t.Error("IconExists(aws) should be true")
	}
	if IconExists("nonexistent-brand") {
		t.Error("IconExists(nonexistent-brand) should be false")
	}
}
