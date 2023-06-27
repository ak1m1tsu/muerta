package shelflifedetector

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/otiai10/gosseract/v2"
)

type DateDetectorServicer interface {
	Detect(image []byte) ([]time.Time, error)
}

type DateDetectorService struct {
	client *gosseract.Client
}

var reDate = regexp.MustCompile(`\b(?:0[1-9]|[1-2][0-9]|3[01])[\.\/\-](?:0[1-9]|1[0-2])[\.\/\-](?:\d{4}|\d{2})\b`)

func New(cl chan struct{}) *DateDetectorService {
	client := gosseract.NewClient()
	client.SetLanguage("eng", "rus")
	go func() {
		<-cl
		client.Close()
	}()
	return &DateDetectorService{client: client}
}

func (s *DateDetectorService) Detect(image []byte) ([]time.Time, error) {
	if err := s.client.SetImageFromBytes(image); err != nil {
		return nil, fmt.Errorf("failed to set image: %w", err)
	}
	text, err := s.client.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to detect date: %w", err)
	}
	cleanText := strings.Replace(text, "\n", " ", -1)
	matches := reDate.FindAllString(cleanText, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no matches found")
	}
	dates := make([]time.Time, 2)
	replacer := strings.NewReplacer("/", "-", ".", "-")
	for i, match := range matches {
		fixedMatch := replacer.Replace(match)
		pieces := strings.Split(fixedMatch, "-")
		if len(pieces) != 3 {
			continue
		}
		if len(pieces[2]) == 2 {
			pieces[2] = strconv.Itoa(time.Now().Year())[:2] + pieces[2]
		}
		converted := make([]int, len(pieces))
		var err error
		for i, p := range pieces {
			converted[i], err = strconv.Atoi(p)
			if err != nil {
				return nil, err
			}
		}
		date := time.Date(converted[2], time.Month(converted[1]), converted[0], 0, 0, 0, 0, time.UTC)
		dates[i] = date
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })
	return dates, nil
}

func (s *DateDetectorService) Close() error {
	return s.client.Close()
}
