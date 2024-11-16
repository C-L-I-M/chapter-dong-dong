package main

import (
	"fmt"

	"github.com/C-L-I-M/chapter-dong-dong/scappers"
	"github.com/bwmarrin/discordgo"
)

type Channels struct {
	session  *discordgo.Session
	serverId string
	slugToId map[string]string
}

func (c Channels) Send(chapter scappers.NewChapter) error {
	if _, ok := c.slugToId[chapter.SagaSlug]; !ok {
		channel, err := c.session.GuildChannelCreate(c.serverId, chapter.SagaSlug, discordgo.ChannelTypeGuildText)
		if err != nil {
			return fmt.Errorf("%s: failed to create channel: %v", chapter.SagaSlug, err)
		}

		c.slugToId[chapter.SagaSlug] = channel.ID
	}

	if _, err := c.session.ChannelMessageSend(c.slugToId[chapter.SagaSlug], chapter.String()); err != nil {
		return fmt.Errorf("%s: failed to send message: %v", chapter.SagaSlug, err)
	}

	return nil
}

func LoadChannels(session *discordgo.Session, serverId string) (*Channels, error) {
	channelsAPI, err := session.GuildChannels(serverId)
	if err != nil {
		return nil, err
	}

	channels := make(map[string]string, len(channelsAPI))
	for _, channel := range channelsAPI {
		channels[channel.Name] = channel.ID
	}

	return &Channels{
		session:  session,
		serverId: serverId,
		slugToId: channels,
	}, nil
}
