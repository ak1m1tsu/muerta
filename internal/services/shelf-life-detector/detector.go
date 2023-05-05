package shelflifedetector

import (
	"fmt"
	"regexp"

	"github.com/otiai10/gosseract/v2"
)

const datePattern = `\b(?:0[1-9]|[1-2][0-9]|3[01])[\.\/\-](?:0[1-9]|1[0-2])[\.\/\-](?:\d{4}|\d{2})\b`

type DateDetectorServicer interface {
	Detect(path string) (string, error)
}

type DateDetectorService struct {
}

func (s *DateDetectorService) Detect(image []byte) ([]string, error) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("eng", "rus")
	client.SetImageFromBytes(image)
	text, err := client.Text()
	compiler := regexp.MustCompile(datePattern)
	res := compiler.FindAllString(text, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to detect date: %w", err)
	}
	return res, nil
}
