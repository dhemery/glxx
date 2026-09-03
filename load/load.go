// Package load loads a GLX archive
package load

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	glx "github.com/genealogix/glx/go-glx"
)

func Load(root string) (*glx.GLXFile, error) {
	files, err := readFiles(root)
	if err != nil {
		return nil, err
	}

	s := glx.NewSerializer(nil)
	archive, duplicates, err := s.DeserializeMultiFileFromMap(files)

	if err != nil {
		return nil, err
	}

	if len(duplicates) > 0 {
		return nil, fmt.Errorf("duplicates: %s", duplicates)
	}

	return archive, nil
}

func readFiles(rootDir string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !isGLXFile(d.Name()) {
			return nil
		}

		path = filepath.Clean(path)
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		files[relPath] = data

		return nil
	})

	return files, err
}

func isGLXFile(filename string) bool {
	return filepath.Ext(filename) == glx.FileExtGLX
}
