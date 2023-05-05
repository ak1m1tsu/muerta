package shelflifedetector

import (
	"fmt"
	"regexp"

	"github.com/otiai10/gosseract/v2"
)

type DateDetectorServicer interface {
	Detect(image []byte) ([]string, error)
}

type DateDetectorService struct {
	pattern string
}

func New() *DateDetectorService {
	return &DateDetectorService{
		pattern: `\b(?:0[1-9]|[1-2][0-9]|3[01])[\.\/\-](?:0[1-9]|1[0-2])[\.\/\-](?:\d{4}|\d{2})\b`,
	}
}

func (s *DateDetectorService) Detect(image []byte) ([]string, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("eng", "rus")
	if err := client.SetImageFromBytes(image); err != nil {
		return nil, fmt.Errorf("failed to set image: %w", err)
	}
	text, err := client.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to detect date: %w", err)
	}
	compiler := regexp.MustCompile(s.pattern)
	res := compiler.FindAllString(text, -1)
	return res, nil
}
