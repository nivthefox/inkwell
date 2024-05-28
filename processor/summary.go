package processor

import "gopkg.in/yaml.v3"

type Summary struct {
	Characters int `yaml:"characters"`
	Words      int `yaml:"words"`
}

type BookSummary struct {
	Summary        `yaml:",inline"`
	Average        int              `yaml:"average"`
	ChapterSummary []ChapterSummary `yaml:"chapters"`
}

type ChapterSummary struct {
	Summary      `yaml:",inline"`
	Title        string         `yaml:"title"`
	Average      int            `yaml:"average"`
	SceneSummary []SceneSummary `yaml:"scenes"`
}

type SceneSummary struct {
	Summary             `yaml:",inline"`
	Files               int `yaml:"files"`
	AverageWordsPerFile int `yaml:"average"`
}

func (s *BookSummary) AddChapterSummary(c ChapterSummary) {
	s.Characters += c.Characters
	s.Words += c.Words
	s.ChapterSummary = append(s.ChapterSummary, c)
}

func (s *BookSummary) String() (string, error) {
	// compute averages
	s.Average = s.Words / len(s.ChapterSummary)
	for i := range s.ChapterSummary {
		s.ChapterSummary[i].Average = s.ChapterSummary[i].Words / len(s.ChapterSummary[i].SceneSummary)
		for j := range s.ChapterSummary[i].SceneSummary {
			s.ChapterSummary[i].SceneSummary[j].AverageWordsPerFile = s.ChapterSummary[i].SceneSummary[j].Words / s.ChapterSummary[i].SceneSummary[j].Files
		}
	}

	// write to yaml
	out, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

func (c *ChapterSummary) AddSceneSummary(s SceneSummary) {
	c.Characters += s.Characters
	c.Words += s.Words
	c.SceneSummary = append(c.SceneSummary, s)
}

func (s *SceneSummary) AddCharacters(n int) {
	s.Characters += n
}

func (s *SceneSummary) AddWords(n int) {
	s.Words += n
}

func (s *SceneSummary) AddFile() {
	s.Files += 1
}
