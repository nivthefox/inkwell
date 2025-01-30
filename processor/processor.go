package processor

import (
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nivthefox/inkwell/config"
)

var spaces = regexp.MustCompile(`[ \t]+`)

// ProcessBook iterates over each of the files in every scene in the config
// and builds the appropriate output files by concatenating the contents
// of the files in each scene.
func ProcessBook(config config.InkwellConfig) error {
	builder := &strings.Builder{}
	summary := BookSummary{}

	createMetadata(config, builder)

	if config.Title != "" {
		tperr := createTitlePage(config.Title, config.Authors, builder)
		if tperr != nil {
			return tperr
		}
	}

	derr := createDedication(config.DedicationFilename, builder)
	if derr != nil {
		return derr
	}

	for _, section := range config.Sections {
		_, secerr := ProcessSection(section)
		if secerr != nil {
			return secerr
		}
	}

	for _, chapter := range config.Chapters {
		text, err := ProcessChapter(chapter, config.SceneSeparator, &summary)
		if err != nil {
			return err
		}
		builder.WriteString("\n" + text.String())
	}

	if config.OutputFilename != "" {
		ferr := writeToFile(builder.String(), config.OutputNumbers, config.OutputFilename)
		if ferr != nil {
			return ferr
		}
	}

	if config.SummaryFilename != "" {
		sum, serr := summary.String()
		if serr != nil {
			return serr
		}
		ferr := writeToFile(sum, false, config.SummaryFilename)
		if ferr != nil {
			return ferr
		}
	}

	return nil
}

// ProcessChapter iterates over each of the scenes in the chapter in the config
// and builds the appropriate output files by concatenating the contents
// of the files in each scene.
func ProcessChapter(config config.ChapterConfig, separator string, book *BookSummary) (*strings.Builder, error) {
	builder := &strings.Builder{}
	builder.WriteString("## " + config.Title + "\n")
	summary := ChapterSummary{
		Title: config.Title,
	}

	for idx, scene := range config.Scenes {
		if idx > 0 {
			builder.WriteString("\n\\* \\* \\*\n\n")
		}

		sceneBuilder, err := ProcessScene(scene, &summary)
		if err != nil {
			return nil, err
		}

		builder.WriteString(sceneBuilder.String())
	}

	if config.OutputFilename != "" {
		err := writeToFile(builder.String(), config.OutputNumbers, config.OutputFilename)
		if err != nil {
			return nil, err
		}
	}

	book.AddChapterSummary(summary)
	return builder, nil
}

// ProcessScene concatenates the contents of the files in the scene in the config
// and writes the output to the appropriate output file.
func ProcessScene(config config.SceneConfig, chapter *ChapterSummary) (*strings.Builder, error) {
	scene := &strings.Builder{}
	summary := SceneSummary{}

	for idx, path := range config.Files {
		builder := &strings.Builder{}
		if idx > 0 {
			scene.WriteString("\n")
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		// Read the contents of the file
		_, err = io.Copy(builder, file)
		if err != nil {
			return nil, err
		}
		_ = file.Close()

		// Trim extra newlines
		content := strings.TrimSpace(builder.String())
		content = spaces.ReplaceAllString(content, " ")

		summary.AddCharacters(len(content))
		summary.AddWords(len(strings.Fields(content)))
		summary.AddFile()

		scene.WriteString(content + "\n")
	}

	if config.OutputFilename != "" {
		err := writeToFile(scene.String(), config.OutputNumbers, config.OutputFilename)
		if err != nil {
			return nil, err
		}
	}

	chapter.AddSceneSummary(summary)
	return scene, nil
}

// ProcessSection concatenates the contents of the files in the section in the config
// and writes the output to the appropriate output file.
func ProcessSection(config config.SectionConfig) (*strings.Builder, error) {
	section := &strings.Builder{}
	section.WriteString("# " + config.Title + "\n")

	for idx, path := range config.Files {
		builder := &strings.Builder{}
		if idx > 0 {
			section.WriteString("\n")
		}

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		// Read the contents of the file
		_, err = io.Copy(builder, file)
		if err != nil {
			return nil, err
		}
		_ = file.Close()

		// Trim extra newlines
		content := strings.TrimSpace(builder.String())
		content = spaces.ReplaceAllString(content, " ")

		section.WriteString(content + "\n")
	}

	if config.OutputFilename != "" {
		err := writeToFile(section.String(), config.OutputNumbers, config.OutputFilename)
		if err != nil {
			return nil, err
		}
	}

	return section, nil
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

// createMetadata writes the title and summary of the book to the builder.
func createMetadata(config config.InkwellConfig, builder *strings.Builder) {
	builder.WriteString("---\n")
	builder.WriteString("Title: " + config.Title + "\n")
	builder.WriteString("Summary: " + config.Summary + "\n")
	builder.WriteString("Date: " + time.Now().Format(time.RFC3339) + "\n")
	builder.WriteString("Authors:")
	for idx, author := range config.Authors {
		if idx > 0 {
			builder.WriteString("        ")
		}
		builder.WriteString(" " + author + "\n")
	}
	builder.WriteString("---\n")

	return
}

// createTitlePage writes the title and author of the book to the builder.
func createTitlePage(title string, authors []string, builder *strings.Builder) error {
	builder.WriteString("# " + title + "\n")
	builder.WriteString("By " + strings.Join(authors, ", ") + "\n")
	return nil
}

// writeToFile writes the contents of the strings.Builder to a file with the given filename.
func writeToFile(output string, numbers bool, filename config.OutputFilename) error {
	// trim windows line endings
	output = strings.ReplaceAll(output, "\r\n", "\n")

	file, err := os.Create(string(filename))
	if err != nil {
		return err
	}

	if numbers {
		o := strings.Split(output, "\n")
		no := 1
		inMeta := false
		inBlockquote := false

		for idx, line := range o {
			trimmedLine := strings.TrimSpace(line)

			// Handle metadata section
			if trimmedLine == "---" {
				inMeta = !inMeta
				continue
			}
			if inMeta {
				continue
			}

			// Skip empty lines, headers, and scene separators
			if trimmedLine == "" ||
				strings.HasPrefix(trimmedLine, "#") ||
				trimmedLine == "\\* \\* \\*" ||
				strings.HasPrefix(trimmedLine, "-") {
				continue
			}

			// Handle blockquotes
			if strings.HasPrefix(trimmedLine, ">") {
				inBlockquote = true
				continue
			}

			// If we were in a blockquote and hit a non-blockquote line,
			// mark the end of blockquote
			if inBlockquote && !strings.HasPrefix(trimmedLine, ">") {
				inBlockquote = false
			}

			// Only number actual paragraph lines (not in blockquote or other special elements)
			if !inBlockquote && trimmedLine != "" {
				o[idx] = line + " <" + strconv.Itoa(no) + ">"
				no += 1
			}
		}
		output = strings.TrimSpace(strings.Join(o, "\n"))
	}

	_, err = file.WriteString(output)
	if err != nil {
		return err
	}

	return nil
}
