package processor

import (
	"github.com/nivthefox/inkwell/config"
	"io"
	"os"
	"strings"
)

// ProcessBook iterates over each of the files in every scene in the config
// and builds the appropriate output files by concatenating the contents
// of the files in each scene.
func ProcessBook(config config.InkwellConfig) error {
	builder := &strings.Builder{}
	tperr := createTitlePage(config.Title, config.Author, builder)
	if tperr != nil {
		return tperr
	}

	derr := createDedication(config.DedicationFilename, builder)
	if derr != nil {
		return derr
	}

	for _, chapter := range config.Chapters {
		chapter, err := ProcessChapter(chapter, config.SceneSeparator)
		if err != nil {
			return err
		}
		builder.WriteString(chapter.String())
	}

	if config.OutputFilename != "" {
		err := writeToFile(builder, config.OutputFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

// ProcessChapter iterates over each of the scenes in the chapter in the config
// and builds the appropriate output files by concatenating the contents
// of the files in each scene.
func ProcessChapter(config config.ChapterConfig, separator string) (*strings.Builder, error) {
	builder := &strings.Builder{}
	builder.WriteString("\n## " + config.Title + "\n")

	for idx, scene := range config.Scenes {
		if idx > 0 {
			builder.WriteString("\n" + separator + "\n")
		}

		sceneBuilder, err := ProcessScene(scene)
		if err != nil {
			return nil, err
		}

		builder.WriteString(sceneBuilder.String())
	}

	if config.OutputFilename != "" {
		err := writeToFile(builder, config.OutputFilename)
		if err != nil {
			return nil, err
		}
	}

	return builder, nil
}

// ProcessScene concatenates the contents of the files in the scene in the config
// and writes the output to the appropriate output file.
func ProcessScene(config config.SceneConfig) (*strings.Builder, error) {
	builder := &strings.Builder{}
	for int, path := range config.Files {
		if int > 0 {
			builder.WriteString("\n")
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Read the contents of the file
		_, err = io.Copy(builder, file)
		if err != nil {
			return nil, err
		}
	}

	if config.OutputFilename != "" {
		err := writeToFile(builder, config.OutputFilename)
		if err != nil {
			return nil, err
		}
	}

	return builder, nil
}

// createDedication reads the contents of the dedication file and writes it to the builder.
func createDedication(filename string, builder *strings.Builder) error {
	if filename == "" {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	builder.WriteString("## Dedication\n")

	// Read the contents of the file
	_, err = io.Copy(builder, file)
	if err != nil {
		return err
	}

	builder.WriteString("\n")
	return nil
}

// createTitlePage writes the title and author of the book to the builder.
func createTitlePage(title, author string, builder *strings.Builder) error {
	builder.WriteString("# " + title + "\n")
	builder.WriteString("- By " + author + "\n")
	return nil
}

// writeToFile writes the contents of the strings.Builder to a file with the given filename.
func writeToFile(builder *strings.Builder, filename config.OutputFilename) error {
	file, err := os.Create(string(filename))
	if err != nil {
		return err
	}

	_, err = file.WriteString(builder.String())
	if err != nil {
		return err
	}

	return nil
}
