package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// OutputFilename is a type that represents the name of the output file
type OutputFilename string

// InkwellConfig is a struct that represents the configuration of the book
type InkwellConfig struct {
	Title   string   `yaml:"title"`
	Summary string   `yaml:"summary"`
	Authors []string `yaml:"authors"`

	DedicationFilename string          `yaml:"dedication"`
	SceneSeparator     string          `yaml:"scene_separator"`
	Sections           []SectionConfig `yaml:"sections"`
	Chapters           []ChapterConfig `yaml:"chapters"`
	OutputFilename     OutputFilename  `yaml:"output_filename,omitempty"`
	OutputNumbers      bool            `yaml:"number_paragraphs,omitempty"`
	SummaryFilename    OutputFilename  `yaml:"summary_filename,omitempty"`
}

// SectionConfig is a struct that represents the configuration of a section
type SectionConfig struct {
	Title          string         `yaml:"title"`
	Files          []string       `yaml:"files"`
	OutputFilename OutputFilename `yaml:"output_filename,omitempty"`
	OutputNumbers  bool           `yaml:"number_paragraphs,omitempty"`
}

// ChapterConfig is a struct that represents the configuration of a chapter
type ChapterConfig struct {
	Title          string `yaml:"title"`
	Scenes         []SceneConfig
	OutputFilename OutputFilename `yaml:"output_filename,omitempty"`
	OutputNumbers  bool           `yaml:"number_paragraphs,omitempty"`
}

// SceneConfig is a struct that represents the configuration of a scene
type SceneConfig struct {
	Files          []string       `yaml:"files"`
	OutputFilename OutputFilename `yaml:"output_filename,omitempty"`
	OutputNumbers  bool           `yaml:"number_paragraphs,omitempty"`
}

// NewInkwellConfig reads the configuration from a file and returns an InkwellConfig
func NewInkwellConfig(filename string) (*InkwellConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config InkwellConfig
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	if config.SceneSeparator == "" {
		config.SceneSeparator = "*&#9;*&#9;*"
	}

	return &config, nil
}
