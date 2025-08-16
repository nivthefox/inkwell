package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/nivthefox/inkwell/config"
)

func TestProcessWikiLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple wiki link",
			input:    "This is a [[wiki link]] in text",
			expected: "This is a wiki link in text",
		},
		{
			name:     "multiple wiki links",
			input:    "[[first link]] and [[second link]]",
			expected: "first link and second link",
		},
		{
			name:     "nested brackets",
			input:    "[[link with [inner] brackets]]",
			expected: "[[link with [inner] brackets]]", // regex stops at first ], so nested brackets don't work
		},
		{
			name:     "no wiki links",
			input:    "Regular text without links",
			expected: "Regular text without links",
		},
		{
			name:     "empty wiki link",
			input:    "Text with [[]] empty link",
			expected: "Text with [[]] empty link", // regex requires at least one non-] character
		},
		{
			name:     "wiki link with special characters",
			input:    "[[link with @#$%^&*()]]",
			expected: "link with @#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processWikiLinks(tt.input)
			if result != tt.expected {
				t.Errorf("processWikiLinks(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateMetadata(t *testing.T) {
	config := config.InkwellConfig{
		Title:   "Test Book",
		Summary: "A test summary",
		Authors: []string{"Author One", "Author Two"},
	}

	builder := &strings.Builder{}
	createMetadata(config, builder)

	result := builder.String()

	// Check that metadata contains expected content
	expectedParts := []string{
		"---",
		"Title: Test Book",
		"Summary: A test summary",
		"Date: " + time.Now().Format("2006-01-02"), // Check date prefix
		"Authors: Author One",
		"        Author Two",
		"---",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("createMetadata() result missing expected part: %q", part)
		}
	}

	// Should start and end with ---
	if !strings.HasPrefix(result, "---\n") {
		t.Error("createMetadata() should start with ---")
	}
	if !strings.HasSuffix(result, "---\n") {
		t.Error("createMetadata() should end with ---")
	}
}

func TestCreateTitlePage(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		authors  []string
		expected string
	}{
		{
			name:     "single author",
			title:    "My Book",
			authors:  []string{"John Doe"},
			expected: "# My Book\nBy John Doe\n",
		},
		{
			name:     "multiple authors",
			title:    "Collaborative Work",
			authors:  []string{"Jane Smith", "Bob Johnson"},
			expected: "# Collaborative Work\nBy Jane Smith, Bob Johnson\n",
		},
		{
			name:     "empty title",
			title:    "",
			authors:  []string{"Author"},
			expected: "# \nBy Author\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &strings.Builder{}
			err := createTitlePage(tt.title, tt.authors, builder)
			if err != nil {
				t.Errorf("createTitlePage() error = %v", err)
				return
			}

			result := builder.String()
			if result != tt.expected {
				t.Errorf("createTitlePage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCreateDedication(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	dedicationFile := filepath.Join(tempDir, "dedication.txt")
	dedicationContent := "To my family and friends"

	err := os.WriteFile(dedicationFile, []byte(dedicationContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		wantErr  bool
		expected string
	}{
		{
			name:     "valid dedication file",
			filename: dedicationFile,
			wantErr:  false,
			expected: "## Dedication\n" + dedicationContent + "\n",
		},
		{
			name:     "empty filename",
			filename: "",
			wantErr:  false,
			expected: "",
		},
		{
			name:     "non-existent file",
			filename: "non-existent.txt",
			wantErr:  true,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := &strings.Builder{}
			err := createDedication(tt.filename, builder)

			if tt.wantErr && err == nil {
				t.Error("createDedication() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("createDedication() unexpected error = %v", err)
				return
			}

			result := builder.String()
			if result != tt.expected {
				t.Errorf("createDedication() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestWriteToFile(t *testing.T) {
	tempDir := t.TempDir()
	
	tests := []struct {
		name     string
		output   string
		numbers  bool
		filename string
		wantErr  bool
	}{
		{
			name:     "simple text without numbers",
			output:   "Hello\nWorld",
			numbers:  false,
			filename: "test1.txt",
			wantErr:  false,
		},
		{
			name:     "text with line numbering",
			output:   "---\ntitle: test\n---\n\n# Header\n\nFirst paragraph.\n\nSecond paragraph.",
			numbers:  true,
			filename: "test2.txt",
			wantErr:  false,
		},
		{
			name:     "text with windows line endings",
			output:   "Line 1\r\nLine 2\r\n",
			numbers:  false,
			filename: "test3.txt",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath := filepath.Join(tempDir, tt.filename)
			err := writeToFile(tt.output, tt.numbers, config.OutputFilename(fullPath))

			if tt.wantErr && err == nil {
				t.Error("writeToFile() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("writeToFile() unexpected error = %v", err)
				return
			}

			if !tt.wantErr {
				// Verify file was created and contains expected content
				content, err := os.ReadFile(fullPath)
				if err != nil {
					t.Errorf("Failed to read written file: %v", err)
					return
				}

				// Check that Windows line endings were converted
				if strings.Contains(string(content), "\r\n") {
					t.Error("writeToFile() should convert Windows line endings")
				}

				// For numbered output, check that line numbers were added
				if tt.numbers {
					contentStr := string(content)
					if strings.Contains(tt.output, "First paragraph.") && !strings.Contains(contentStr, "<1>") {
						t.Error("writeToFile() with numbers=true should add line numbers to paragraphs")
					}
				}
			}
		})
	}
}

func TestProcessScene(t *testing.T) {
	// Create temporary files for testing
	tempDir := t.TempDir()
	file1 := filepath.Join(tempDir, "scene1.txt")
	file2 := filepath.Join(tempDir, "scene2.txt")

	content1 := "This is the first part of the scene."
	content2 := "This is the [[wiki link]] second part."

	err := os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name           string
		config         config.SceneConfig
		stripWikiLinks bool
		wantErr        bool
		expectedWords  int
	}{
		{
			name: "scene with multiple files, no wiki link stripping",
			config: config.SceneConfig{
				Files: []string{file1, file2},
			},
			stripWikiLinks: false,
			wantErr:        false,
			expectedWords:  15, // Total words in both files
		},
		{
			name: "scene with wiki link stripping",
			config: config.SceneConfig{
				Files: []string{file2},
			},
			stripWikiLinks: true,
			wantErr:        false,
			expectedWords:  7, // Words after removing wiki link markup
		},
		{
			name: "scene with non-existent file",
			config: config.SceneConfig{
				Files: []string{"non-existent.txt"},
			},
			stripWikiLinks: false,
			wantErr:        true,
			expectedWords:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chapter := &ChapterSummary{}
			result, err := ProcessScene(tt.config, tt.stripWikiLinks, chapter)

			if tt.wantErr && err == nil {
				t.Error("ProcessScene() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ProcessScene() unexpected error = %v", err)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("ProcessScene() returned nil result")
					return
				}

				// Check that chapter summary was updated
				if len(chapter.SceneSummary) != 1 {
					t.Errorf("ProcessScene() should add one scene summary, got %d", len(chapter.SceneSummary))
				}

				scene := chapter.SceneSummary[0]
				if scene.Words != tt.expectedWords {
					t.Errorf("ProcessScene() scene words = %d, want %d", scene.Words, tt.expectedWords)
				}

				if scene.Files != len(tt.config.Files) {
					t.Errorf("ProcessScene() scene files = %d, want %d", scene.Files, len(tt.config.Files))
				}
			}
		})
	}
}

func TestProcessSection(t *testing.T) {
	// Create temporary files for testing
	tempDir := t.TempDir()
	file1 := filepath.Join(tempDir, "section1.txt")

	content1 := "This is section content with [[wiki links]]."

	err := os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name           string
		config         config.SectionConfig
		stripWikiLinks bool
		wantErr        bool
		expectedTitle  string
	}{
		{
			name: "valid section with wiki link stripping",
			config: config.SectionConfig{
				Title: "Chapter One",
				Files: []string{file1},
			},
			stripWikiLinks: true,
			wantErr:        false,
			expectedTitle:  "# Chapter One\n",
		},
		{
			name: "section with non-existent file",
			config: config.SectionConfig{
				Title: "Bad Chapter",
				Files: []string{"non-existent.txt"},
			},
			stripWikiLinks: false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ProcessSection(tt.config, tt.stripWikiLinks)

			if tt.wantErr && err == nil {
				t.Error("ProcessSection() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ProcessSection() unexpected error = %v", err)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("ProcessSection() returned nil result")
					return
				}

				resultStr := result.String()
				if !strings.HasPrefix(resultStr, tt.expectedTitle) {
					t.Errorf("ProcessSection() result should start with %q", tt.expectedTitle)
				}

				// Check wiki link processing
				if tt.stripWikiLinks && strings.Contains(resultStr, "[[") {
					t.Error("ProcessSection() should strip wiki links when stripWikiLinks=true")
				}
			}
		})
	}
}

func TestProcessChapter(t *testing.T) {
	// Create temporary files for testing
	tempDir := t.TempDir()
	file1 := filepath.Join(tempDir, "scene1.txt")
	file2 := filepath.Join(tempDir, "scene2.txt")

	content1 := "First scene content."
	content2 := "Second scene content."

	err := os.WriteFile(file1, []byte(content1), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	err = os.WriteFile(file2, []byte(content2), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name      string
		config    config.ChapterConfig
		separator string
		wantErr   bool
	}{
		{
			name: "chapter with multiple scenes",
			config: config.ChapterConfig{
				Title: "Test Chapter",
				Scenes: []config.SceneConfig{
					{Files: []string{file1}},
					{Files: []string{file2}},
				},
			},
			separator: "\\* \\* \\*",
			wantErr:   false,
		},
		{
			name: "chapter with invalid scene file",
			config: config.ChapterConfig{
				Title: "Bad Chapter",
				Scenes: []config.SceneConfig{
					{Files: []string{"non-existent.txt"}},
				},
			},
			separator: "\\* \\* \\*",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			book := &BookSummary{}
			result, err := ProcessChapter(tt.config, tt.separator, false, book)

			if tt.wantErr && err == nil {
				t.Error("ProcessChapter() expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("ProcessChapter() unexpected error = %v", err)
				return
			}

			if !tt.wantErr {
				if result == nil {
					t.Error("ProcessChapter() returned nil result")
					return
				}

				resultStr := result.String()
				expectedTitle := "## " + tt.config.Title + "\n"
				if !strings.HasPrefix(resultStr, expectedTitle) {
					t.Errorf("ProcessChapter() should start with %q", expectedTitle)
				}

				// Check that separator is used between scenes
				if len(tt.config.Scenes) > 1 && !strings.Contains(resultStr, tt.separator) {
					t.Errorf("ProcessChapter() should include separator %q between scenes", tt.separator)
				}

				// Check that book summary was updated
				if len(book.ChapterSummary) != 1 {
					t.Errorf("ProcessChapter() should add one chapter summary, got %d", len(book.ChapterSummary))
				}

				chapter := book.ChapterSummary[0]
				if chapter.Title != tt.config.Title {
					t.Errorf("ProcessChapter() chapter title = %q, want %q", chapter.Title, tt.config.Title)
				}

				if len(chapter.SceneSummary) != len(tt.config.Scenes) {
					t.Errorf("ProcessChapter() should have %d scene summaries, got %d", len(tt.config.Scenes), len(chapter.SceneSummary))
				}
			}
		})
	}
}