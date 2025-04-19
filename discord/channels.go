package discord

import (
	"fmt"
	"maps"
	"slices"

	"github.com/C-L-I-M/chapter-dong-dong/scrappers"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type Channels struct {
	session  *discordgo.Session
	serverId string
	slugToId map[string]string
}

func LoadChannels(session *discordgo.Session, serverId string, staticChans []string) (*Channels, error) {
	channelsAPI, err := session.GuildChannels(serverId)
	if err != nil {
		return nil, err
	}

	channels := make(map[string]string, len(channelsAPI))
	for _, channel := range channelsAPI {
		if slices.Contains(staticChans, channel.Name) {
			continue
		}
		channels[channel.Name] = channel.ID
	}

	return &Channels{
		session:  session,
		serverId: serverId,
		slugToId: channels,
	}, nil
}

func partitionOldAndNewChannels(currentChannels, maybeNewChannels []string) ([]string, []string) {
	var toDelete, toCreate []string
	for _, channel := range currentChannels {
		if !slices.Contains(maybeNewChannels, channel) {
			toDelete = append(toDelete, channel)
		}
	}

	for _, channel := range maybeNewChannels {
		if !slices.Contains(currentChannels, channel) {
			toCreate = append(toCreate, channel)
		}
	}

	return toDelete, toCreate
}

func (c *Channels) Send(chapter scrappers.NewChapter) error {
	if _, ok := c.slugToId[chapter.SagaSlug]; !ok {
		channel, err := c.session.GuildChannelCreate(c.serverId, chapter.SagaSlug, discordgo.ChannelTypeGuildText)
		if err != nil {
			return fmt.Errorf("%s: failed to create channel: %v", chapter.SagaSlug, err)
		}

		log.Info("Channel created: ", channel.Name)
		c.slugToId[chapter.SagaSlug] = channel.ID
	}

	if _, err := c.session.ChannelMessageSend(c.slugToId[chapter.SagaSlug], chapter.String()); err != nil {
		return fmt.Errorf("%s: failed to send message: %v", chapter.SagaSlug, err)
	}

	return nil
}

func (c *Channels) ResolveConfigDiff(cfgChannels []string) error {
	currentChans := slices.Sorted(maps.Keys(c.slugToId))
	toDelete, toCreate := partitionOldAndNewChannels(currentChans, cfgChannels)
	log.Info("Channels that should be deleted: ", toDelete)

	for _, channelName := range toCreate {
		channel, err := c.session.GuildChannelCreate(c.serverId, channelName, discordgo.ChannelTypeGuildText)
		if err != nil {
			log.Errorf("Failed to create channel %s: %v", channelName, err)
			return err
		}
		log.Info("Channel created: ", channelName)
		c.slugToId[channelName] = channel.ID
	}

	return nil
}
