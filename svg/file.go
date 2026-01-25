package svg

import (
	"os"
	"path/filepath"
	"strings"
)

// FileInfo contains information about a file or directory path.
type FileInfo struct {
	Path  string
	IsDir bool
}

// GetPathInfo returns information about a path.
func GetPathInfo(path string) (*FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Path:  path,
		IsDir: info.IsDir(),
	}, nil
}

// ListSVGFiles returns all SVG files in a directory (non-recursive).
func ListSVGFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".svg") {
			continue
		}
		files = append(files, filepath.Join(dirPath, entry.Name()))
	}

	return files, nil
}

// IsSVGFile returns true if the path is an SVG file.
func IsSVGFile(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".svg")
}

// ListSVGFilesRecursive returns all SVG files in a directory tree.
func ListSVGFilesRecursive(dirPath string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".svg") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
