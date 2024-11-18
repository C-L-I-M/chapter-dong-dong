package cmd

import (
	"maps"
	"slices"
	"sync"
	"time"

	"github.com/C-L-I-M/chapter-dong-dong/config"
	"github.com/C-L-I-M/chapter-dong-dong/discord"
	"github.com/C-L-I-M/chapter-dong-dong/scappers"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, _ []string) {
	cfg, err := config.Load()
	cobra.CheckErr(err)

	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	cobra.CheckErr(err)

	channels, err := discord.LoadChannels(session, cfg.ServerId, cfg.StaticChannels)
	cobra.CheckErr(err)

	cobra.CheckErr(channels.ResolveConfigDiff(slices.Sorted(maps.Keys(cfg.Sagas))))

	var wg sync.WaitGroup
	wg.Add(len(cfg.Sagas))

	for sagaSlug, saga := range cfg.Sagas {
		scrapper, err := scappers.NewScrapper(saga.SchedulingMode, saga.Parameters)
		cobra.CheckErr(err)

		go func(saga *config.Saga, scrapper scappers.Scrapper) {
			defer wg.Done()

			ticker := time.NewTicker(saga.Interval)
			defer ticker.Stop()
			for range ticker.C {
				log.Info(sagaSlug + ": tick start")
				ctx := scappers.FromSaga(sagaSlug, saga)
				chapters, err := scrapper.Scrap(ctx)
				if err != nil {
					log.Errorf("%s: failed to scrap: %v", sagaSlug, err)
					continue
				}

				for _, chapter := range chapters {
					log.Infof("%s: new chapter: %q - %s", sagaSlug, chapter.Name, chapter.Number)
					if err := channels.Send(chapter); err != nil {
						log.Errorf("Failed to send chapter: %v", err)
					}
				}

				if ctx.HasChanged() {
					if err := config.Save(); err != nil {
						log.Errorf("%s: failed to save config: %v", sagaSlug, err)
					}
				}
				log.Info(sagaSlug + ": tick end")
			}
		}(saga, scrapper)
	}

	wg.Wait()
}
