package processor

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSceneSummaryAddCharacters(t *testing.T) {
	scene := &SceneSummary{}
	
	scene.AddCharacters(100)
	if scene.Characters != 100 {
		t.Errorf("AddCharacters(100) = %d, want 100", scene.Characters)
	}
	
	scene.AddCharacters(50)
	if scene.Characters != 150 {
		t.Errorf("AddCharacters(50) after 100 = %d, want 150", scene.Characters)
	}
	
	scene.AddCharacters(0)
	if scene.Characters != 150 {
		t.Errorf("AddCharacters(0) should not change total = %d, want 150", scene.Characters)
	}
}

func TestSceneSummaryAddWords(t *testing.T) {
	scene := &SceneSummary{}
	
	scene.AddWords(25)
	if scene.Words != 25 {
		t.Errorf("AddWords(25) = %d, want 25", scene.Words)
	}
	
	scene.AddWords(75)
	if scene.Words != 100 {
		t.Errorf("AddWords(75) after 25 = %d, want 100", scene.Words)
	}
	
	scene.AddWords(0)
	if scene.Words != 100 {
		t.Errorf("AddWords(0) should not change total = %d, want 100", scene.Words)
	}
}

func TestSceneSummaryAddFile(t *testing.T) {
	scene := &SceneSummary{}
	
	if scene.Files != 0 {
		t.Errorf("Initial Files count = %d, want 0", scene.Files)
	}
	
	scene.AddFile()
	if scene.Files != 1 {
		t.Errorf("AddFile() = %d, want 1", scene.Files)
	}
	
	scene.AddFile()
	scene.AddFile()
	if scene.Files != 3 {
		t.Errorf("AddFile() called 3 times total = %d, want 3", scene.Files)
	}
}

func TestChapterSummaryAddSceneSummary(t *testing.T) {
	chapter := &ChapterSummary{
		Title: "Test Chapter",
	}
	
	scene1 := SceneSummary{
		Summary: Summary{Characters: 100, Words: 20},
		Files:   2,
	}
	scene2 := SceneSummary{
		Summary: Summary{Characters: 200, Words: 40},
		Files:   3,
	}
	
	chapter.AddSceneSummary(scene1)
	if chapter.Characters != 100 {
		t.Errorf("After adding scene1, Characters = %d, want 100", chapter.Characters)
	}
	if chapter.Words != 20 {
		t.Errorf("After adding scene1, Words = %d, want 20", chapter.Words)
	}
	if len(chapter.SceneSummary) != 1 {
		t.Errorf("After adding scene1, SceneSummary length = %d, want 1", len(chapter.SceneSummary))
	}
	
	chapter.AddSceneSummary(scene2)
	if chapter.Characters != 300 {
		t.Errorf("After adding scene2, Characters = %d, want 300", chapter.Characters)
	}
	if chapter.Words != 60 {
		t.Errorf("After adding scene2, Words = %d, want 60", chapter.Words)
	}
	if len(chapter.SceneSummary) != 2 {
		t.Errorf("After adding scene2, SceneSummary length = %d, want 2", len(chapter.SceneSummary))
	}
}

func TestBookSummaryAddChapterSummary(t *testing.T) {
	book := &BookSummary{}
	
	chapter1 := ChapterSummary{
		Summary: Summary{Characters: 500, Words: 100},
		Title:   "Chapter 1",
	}
	chapter2 := ChapterSummary{
		Summary: Summary{Characters: 300, Words: 60},
		Title:   "Chapter 2",
	}
	
	book.AddChapterSummary(chapter1)
	if book.Characters != 500 {
		t.Errorf("After adding chapter1, Characters = %d, want 500", book.Characters)
	}
	if book.Words != 100 {
		t.Errorf("After adding chapter1, Words = %d, want 100", book.Words)
	}
	if len(book.ChapterSummary) != 1 {
		t.Errorf("After adding chapter1, ChapterSummary length = %d, want 1", len(book.ChapterSummary))
	}
	
	book.AddChapterSummary(chapter2)
	if book.Characters != 800 {
		t.Errorf("After adding chapter2, Characters = %d, want 800", book.Characters)
	}
	if book.Words != 160 {
		t.Errorf("After adding chapter2, Words = %d, want 160", book.Words)
	}
	if len(book.ChapterSummary) != 2 {
		t.Errorf("After adding chapter2, ChapterSummary length = %d, want 2", len(book.ChapterSummary))
	}
}

func TestBookSummaryString(t *testing.T) {
	// Create a book with nested structure
	book := &BookSummary{}
	
	// Create scenes
	scene1 := SceneSummary{
		Summary: Summary{Characters: 100, Words: 20},
		Files:   2,
	}
	scene2 := SceneSummary{
		Summary: Summary{Characters: 200, Words: 40},
		Files:   1,
	}
	scene3 := SceneSummary{
		Summary: Summary{Characters: 150, Words: 30},
		Files:   3,
	}
	
	// Create chapters
	chapter1 := ChapterSummary{
		Title: "Chapter 1",
	}
	chapter1.AddSceneSummary(scene1)
	chapter1.AddSceneSummary(scene2)
	
	chapter2 := ChapterSummary{
		Title: "Chapter 2",
	}
	chapter2.AddSceneSummary(scene3)
	
	// Add chapters to book
	book.AddChapterSummary(chapter1)
	book.AddChapterSummary(chapter2)
	
	// Test String() method
	yamlStr, err := book.String()
	if err != nil {
		t.Fatalf("BookSummary.String() error = %v", err)
	}
	
	// Parse back to verify structure
	var parsedBook BookSummary
	err = yaml.Unmarshal([]byte(yamlStr), &parsedBook)
	if err != nil {
		t.Fatalf("Failed to parse YAML output: %v", err)
	}
	
	// Verify totals
	expectedChars := 450 // 100 + 200 + 150
	expectedWords := 90  // 20 + 40 + 30
	if parsedBook.Characters != expectedChars {
		t.Errorf("BookSummary total Characters = %d, want %d", parsedBook.Characters, expectedChars)
	}
	if parsedBook.Words != expectedWords {
		t.Errorf("BookSummary total Words = %d, want %d", parsedBook.Words, expectedWords)
	}
	
	// Verify averages were calculated
	expectedBookAverage := expectedWords / 2 // 2 chapters
	if parsedBook.Average != expectedBookAverage {
		t.Errorf("BookSummary Average = %d, want %d", parsedBook.Average, expectedBookAverage)
	}
	
	// Verify chapter averages
	if len(parsedBook.ChapterSummary) != 2 {
		t.Fatalf("Expected 2 chapters, got %d", len(parsedBook.ChapterSummary))
	}
	
	// Chapter 1 average: (20 + 40) / 2 scenes = 30
	expectedChapter1Avg := 30
	if parsedBook.ChapterSummary[0].Average != expectedChapter1Avg {
		t.Errorf("Chapter 1 Average = %d, want %d", parsedBook.ChapterSummary[0].Average, expectedChapter1Avg)
	}
	
	// Chapter 2 average: 30 / 1 scene = 30
	expectedChapter2Avg := 30
	if parsedBook.ChapterSummary[1].Average != expectedChapter2Avg {
		t.Errorf("Chapter 2 Average = %d, want %d", parsedBook.ChapterSummary[1].Average, expectedChapter2Avg)
	}
	
	// Verify scene averages per file
	// Scene 1: 20 words / 2 files = 10
	expectedScene1Avg := 10
	if parsedBook.ChapterSummary[0].SceneSummary[0].AverageWordsPerFile != expectedScene1Avg {
		t.Errorf("Scene 1 AverageWordsPerFile = %d, want %d", 
			parsedBook.ChapterSummary[0].SceneSummary[0].AverageWordsPerFile, expectedScene1Avg)
	}
	
	// Scene 2: 40 words / 1 file = 40
	expectedScene2Avg := 40
	if parsedBook.ChapterSummary[0].SceneSummary[1].AverageWordsPerFile != expectedScene2Avg {
		t.Errorf("Scene 2 AverageWordsPerFile = %d, want %d", 
			parsedBook.ChapterSummary[0].SceneSummary[1].AverageWordsPerFile, expectedScene2Avg)
	}
	
	// Scene 3: 30 words / 3 files = 10
	expectedScene3Avg := 10
	if parsedBook.ChapterSummary[1].SceneSummary[0].AverageWordsPerFile != expectedScene3Avg {
		t.Errorf("Scene 3 AverageWordsPerFile = %d, want %d", 
			parsedBook.ChapterSummary[1].SceneSummary[0].AverageWordsPerFile, expectedScene3Avg)
	}
	
	// Verify YAML structure contains expected keys
	expectedKeys := []string{"characters", "words", "average", "chapters"}
	for _, key := range expectedKeys {
		if !strings.Contains(yamlStr, key+":") {
			t.Errorf("YAML output missing expected key: %s", key)
		}
	}
}

func TestBookSummaryStringWithDivisionByZero(t *testing.T) {
	// Test edge case where there are no chapters (should cause division by zero)
	book := &BookSummary{
		Summary: Summary{Characters: 100, Words: 50},
	}
	
	// This should panic due to division by zero in the average calculation
	defer func() {
		if r := recover(); r == nil {
			t.Error("BookSummary.String() should panic when there are no chapters")
		}
	}()
	
	_, _ = book.String()
}

func TestBookSummaryStringWithZeroScenes(t *testing.T) {
	// Test edge case where chapter has no scenes
	book := &BookSummary{}
	chapter := ChapterSummary{
		Title:   "Empty Chapter",
		Summary: Summary{Characters: 0, Words: 0},
	}
	book.AddChapterSummary(chapter)
	
	// This should panic due to division by zero in scene average calculation
	defer func() {
		if r := recover(); r == nil {
			t.Error("BookSummary.String() should panic when chapter has no scenes")
		}
	}()
	
	_, _ = book.String()
}

func TestBookSummaryStringWithZeroFiles(t *testing.T) {
	// Test edge case where scene has no files
	book := &BookSummary{}
	scene := SceneSummary{
		Summary: Summary{Characters: 100, Words: 20},
		Files:   0, // No files
	}
	chapter := ChapterSummary{
		Title: "Chapter with zero-file scene",
	}
	chapter.AddSceneSummary(scene)
	book.AddChapterSummary(chapter)
	
	// This should panic due to division by zero in file average calculation
	defer func() {
		if r := recover(); r == nil {
			t.Error("BookSummary.String() should panic when scene has no files")
		}
	}()
	
	_, _ = book.String()
}

func TestSummaryStructIntegration(t *testing.T) {
	// Integration test that verifies the complete flow
	book := &BookSummary{}
	
	// Simulate realistic usage
	scene := &SceneSummary{}
	scene.AddCharacters(250)
	scene.AddWords(50)
	scene.AddFile()
	scene.AddFile() // 2 files total
	
	chapter := &ChapterSummary{Title: "Integration Test Chapter"}
	chapter.AddSceneSummary(*scene)
	
	book.AddChapterSummary(*chapter)
	
	// Verify all data is properly aggregated
	if book.Characters != 250 {
		t.Errorf("Integration test: book Characters = %d, want 250", book.Characters)
	}
	if book.Words != 50 {
		t.Errorf("Integration test: book Words = %d, want 50", book.Words)
	}
	if chapter.Characters != 250 {
		t.Errorf("Integration test: chapter Characters = %d, want 250", chapter.Characters)
	}
	if chapter.Words != 50 {
		t.Errorf("Integration test: chapter Words = %d, want 50", chapter.Words)
	}
	if scene.Files != 2 {
		t.Errorf("Integration test: scene Files = %d, want 2", scene.Files)
	}
	
	// Test YAML generation
	yamlOutput, err := book.String()
	if err != nil {
		t.Errorf("Integration test: YAML generation failed: %v", err)
	}
	
	if yamlOutput == "" {
		t.Error("Integration test: YAML output should not be empty")
	}
	
	// Verify YAML contains expected structure
	expectedPatterns := []string{
		"characters: 250",
		"words: 50",
		"average: 50", // 50 words / 1 chapter
		"title: Integration Test Chapter",
	}
	
	for _, pattern := range expectedPatterns {
		if !strings.Contains(yamlOutput, pattern) {
			t.Errorf("Integration test: YAML missing pattern %q", pattern)
		}
	}
}