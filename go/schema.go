package jfrogappsconfig

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	version        = "1.0"
	dotJFrogDir    = ".jfrog"
	configFileName = "jfrog-apps-config.yml"
)

type JFrogAppsConfig struct {
	Version string   `yaml:"version,omitempty"`
	Modules []Module `yaml:"modules,omitempty"`
}

type Module struct {
	Name            string   `yaml:"name,omitempty"`
	SourceRoot      string   `yaml:"source_root,omitempty"`
	ExcludePatterns []string `yaml:"exclude_patterns,omitempty"`
	ExcludeScanners []string `yaml:"exclude_scanners,omitempty"`
	Scanners        Scanners `yaml:"scanners,omitempty"`
}

type Scanners struct {
	Secrets *Scanner     `yaml:"secrets,omitempty"`
	Iac     *Scanner     `yaml:"iac,omitempty"`
	Sast    *SastScanner `yaml:"sast,omitempty"`
}

type Scanner struct {
	WorkingDirs     []string `yaml:"working_dirs,omitempty"`
	ExcludePatterns []string `yaml:"exclude_patterns,omitempty"`
}

type SastScanner struct {
	Scanner       `yaml:",inline"`
	Language      string   `yaml:"language,omitempty"`
	ExcludedRules []string `yaml:"excluded_rules,omitempty"`
}

func LoadConfigIfExist() (*JFrogAppsConfig, error) {
	jfrogAppsConfigBytes, err := os.ReadFile(filepath.Join(dotJFrogDir, configFileName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// File does not exist
			return nil, nil
		}
		return nil, err
	}

	jfrogAppsConfig := &JFrogAppsConfig{}
	err = yaml.Unmarshal(jfrogAppsConfigBytes, jfrogAppsConfig)
	return jfrogAppsConfig, err
}
