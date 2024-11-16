package scappers

import (
	"fmt"
	"net/http"

	"github.com/C-L-I-M/chapter-dong-dong/config"
	"github.com/go-viper/mapstructure/v2"
)

type sequentialScrapper struct {
	Url                string `mapstructure:"url"`
	Start              int    `mapstructure:"start"`
	NotFoundStatusCode int    `mapstructure:"not_found_status_code"`
	FoundStatusCode    int    `mapstructure:"found_status_code"`
}

func init() {
	registerScrapper(config.SchedulingModeSequentialPageNotFound, func(parameters map[string]any) (Scrapper, error) {
		s := &sequentialScrapper{}
		if err := mapstructure.Decode(parameters, s); err != nil {
			return nil, fmt.Errorf("invalid parameters: %v", err)
		}

		return s, nil
	})
}

const StateKeyIndex = "index"

func (s *sequentialScrapper) Scrap(ctx *ScrappingContext) ([]NewChapter, error) {
	var chapters []NewChapter
	for {
		i := ctx.GetState(StateKeyIndex).(int)
		url := fmt.Sprintf(s.Url, i)
		found, err := scrapPage(url, s.NotFoundStatusCode, s.FoundStatusCode)
		if err != nil {
			return nil, err
		}

		if !found {
			return chapters, nil
		}

		chapters = append(chapters, NewChapter{
			Name:     fmt.Sprintf("Chapter %d", i),
			Number:   fmt.Sprintf("%d", i),
			Url:      url,
			SagaSlug: ctx.Saga.Slug,
		})
		ctx.SetState(StateKeyIndex, i+1)
	}
}

func scrapPage(url string, notFoundCode int, foundCode int) (bool, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case foundCode:
		return true, nil
	case notFoundCode:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
