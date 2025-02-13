package util

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func InitApp() error {
	projectRoot, err := FindProjectRoot()
	if err != nil {
		logrus.WithError(err).Error("Failed to find project root")
		return err
	}

	templatesDir := filepath.Join(projectRoot, "web", "templates")
	templates, err := CheckTemplates(templatesDir)
	if err != nil {
		logrus.WithError(err).Error("Template check failed")
		return err
	}

	if err := LoadTemplates(templatesDir, templates); err != nil {
		logrus.WithError(err).Error("Failed to load templates")
		return err
	}

	logrus.Info("All templates loaded successfully")
	return nil
}
