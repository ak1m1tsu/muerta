package shelflifedetector

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

type DateDetectorServicer interface {
	Detect(image []byte) ([]string, error)
}

type DateDetectorService struct {
	pattern string
	client  *gosseract.Client
}

const DefaultPattern = `\b(?:0[1-9]|[1-2][0-9]|3[01])[\.\/\-](?:0[1-9]|1[0-2])[\.\/\-](?:\d{4}|\d{2})\b`

func New(cl chan struct{}) *DateDetectorService {
	client := gosseract.NewClient()
	client.SetLanguage("eng", "rus")
	go func() {
		<-cl
		client.Close()
	}()
	return &DateDetectorService{
		pattern: DefaultPattern,
		client:  client,
	}
}

func (s *DateDetectorService) Detect(image []byte) ([]string, error) {
	if err := s.client.SetImageFromBytes(image); err != nil {
		return nil, fmt.Errorf("failed to set image: %w", err)
	}
	text, err := s.client.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to detect date: %w", err)
	}
	text = strings.TrimSpace(strings.Join(strings.Split(text, "\n"), " "))
	result := regexp.MustCompile(s.pattern).FindAllString(text, 2)
	return result, nil
}

func (s *DateDetectorService) Close() error {
	return s.client.Close()
}
