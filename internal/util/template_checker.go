package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckTemplates(directory string) ([]string, error) {
	var templates []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			templates = append(templates, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found in directory: %v", directory)
	}
	return templates, nil
}
