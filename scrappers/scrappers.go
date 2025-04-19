package scrappers

import (
	"fmt"

	"github.com/C-L-I-M/chapter-dong-dong/config"
)

type Scrapper interface {
	Scrap(ctx *ScrappingContext) ([]NewChapter, error)
}

type ScrapperFactory func(parameters map[string]any) (Scrapper, error)

var scrappers = make(map[config.ScrappingMode]ScrapperFactory)

func registerScrapper(name config.ScrappingMode, new ScrapperFactory) {
	if _, ok := scrappers[name]; ok {
		panic("scrapper already registered")
	}

	scrappers[name] = new
}

func NewScrapper(name config.ScrappingMode, parameters map[string]any) (Scrapper, error) {
	if factory, ok := scrappers[name]; ok {
		return factory(parameters)
	}

	return nil, fmt.Errorf("scrapper not found: %s", name)
}
