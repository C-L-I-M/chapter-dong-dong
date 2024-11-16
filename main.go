package main

import (
	"flag"
	"time"

	"github.com/C-L-I-M/chapter-dong-dong/scappers"
	log "github.com/sirupsen/logrus"

	"github.com/C-L-I-M/chapter-dong-dong/config"
	"github.com/bwmarrin/discordgo"
)

var (
	ConfigPath = flag.String("config", "config.json", "Path to the config file")
	BotToken   = flag.String("token", "", "Bot access token")
	ServerId   = flag.String("server", "", "Server id")
)

func init() {
	log.SetReportCaller(true)
}

func main() {
	flag.Parse()
	var err error
	session, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	cfg, err := config.Load(*ConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	channels, err := LoadChannels(session, *ServerId)
	if err != nil {
		log.Fatalf("Failed to load channels: %v", err)
	}

	scrappers := make([]scappers.Scrapper, len(cfg.Sagas))
	for i, saga := range cfg.Sagas {
		scrappers[i], err = scappers.NewScrapper(saga.SchedulingMode, saga.Parameters)
		if err != nil {
			log.Fatalf("%s: failed to create scrapper: %v", saga.Slug, err)
		}

		go func() {
			ticker := time.NewTicker(saga.Interval)
			defer ticker.Stop()
			for range ticker.C {
				ctx := scappers.FromSaga(saga)
				chapters, err := scrappers[i].Scrap(ctx)
				if err != nil {
					log.Errorf("%s: failed to scrap: %v", saga.Slug, err)
					continue
				}

				for _, chapter := range chapters {
					log.Infof("%s: new chapter: %q - %d", saga.Slug, chapter.Name, chapter.Number)
					if err := channels.Send(chapter); err != nil {
						log.Errorf("Failed to send chapter: %v", err)
					}
				}

				if ctx.HasChanged() {
					if err := config.Save(*ConfigPath, cfg); err != nil {
						log.Errorf("%s: failed to save config: %v", saga.Slug, err)
					}
				}
			}
		}()
	}
}
