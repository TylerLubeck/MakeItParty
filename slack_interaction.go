package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/deckarep/golang-set"
	"github.com/nlopes/slack"
	"math/rand"
)

var parrots = []string{
	"aussieparrot",
	"aussiecongaparrot",
	"boredparrot",
	"chillparrot",
	"congaparrot",
	"dealwithitparrot",
	"explodyparrot",
	"fastparrot",
	"fiestaparrot",
	"gothparrot",
	"ice-cream-parrot",
	"middleparrot",
	"oldtimeyparrot",
	"parrot",
	"parrotcop",
	"parrotdad",
	"reversecongaparrot",
	"rightparrot",
	"sadparrot",
	"shufflefurtherparrot",
	"shuffleparrot",
	"shufflepartyparrot",
	"slowparrot",
	"harrypotterparrot",
	"parrotmustache",
	"creepyparrot",
	"parrotbeer",
	"bluescluesparrot",
	"blondesassyparrot",
	"burgerparrot",
	"witnessprotectionparrot",
	"bradford",
}

func HandleSlackEvents(s *slack.Client) {
	rtm := s.NewRTM()

	auth, err := s.AuthTest()

	if err != nil {
		log.Panic("Failed auth check. Aborting")
		return
	}

	log.WithFields(log.Fields{
		"slack_user_id": auth.UserID,
	}).Info("Logging in to Slack...")

	channelIDs := getPublicChannelIDs(s)
	channelIDs = channelIDs.Union(getPrivateChannelIDs(s))

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ReactionAddedEvent:
				handleAddReaction(s, ev, channelIDs)
			case *slack.ReactionRemovedEvent:
				handleRemoveReaction(s, ev, channelIDs)
			case *slack.ChannelJoinedEvent:
				channelIDs.Add(ev.Channel.ID)

				log.WithFields(log.Fields{
					"channel_id": ev.Channel.ID,
					"name":       ev.Channel.Name,
					"event":      ev.Type,
				}).Info("Channel List Modified")

			/*
			 * In theory, these next three cases could be combined,
			 * but I can't figure out the golang ducktyping so ¯\_(ツ)_/¯
			 */
			case *slack.ChannelLeftEvent:
				channelIDs.Remove(ev.Channel)

				log.WithFields(log.Fields{
					"channel_id": ev.Channel,
					"event":      ev.Type,
				}).Info("Channel List Modified")

			case *slack.GroupJoinedEvent:
				channelIDs.Add(ev.Channel)

				log.WithFields(log.Fields{
					"channel_id": ev.Channel,
					"event":      ev.Type,
				}).Info("Channel List Modified")

			case *slack.GroupLeftEvent:
				channelIDs.Remove(ev.Channel)

				log.WithFields(log.Fields{
					"channel_id": ev.Channel,
					"event":      ev.Type,
				}).Info("Channel List Modified")

			default:
				//fmt.Printf("Other Event: %v\n", msg.Data)
			}
		}
	}

}

func getPublicChannelIDs(s *slack.Client) mapset.Set {
	channels, err := s.GetChannels(true)
	if err != nil {
		log.Panic("Failed to get public slack channels. Aborting")
	}

	channelIDs := mapset.NewSet()

	/*
	 * This is a list of all channels in our slack instance, so we need to
	 * filter for only the ones we're a member of
	 */
	for _, channel := range channels {
		if !channel.IsMember {
			continue
		}
		channelIDs.Add(channel.ID)
	}

	return channelIDs
}

func getPrivateChannelIDs(s *slack.Client) mapset.Set {
	groups, err := s.GetGroups(true)
	if err != nil {
		log.Panic("Failed to get private slack channels. Aborting")
	}

	channelIDs := mapset.NewSet()

	/*
	 * This is a list of private channels the bot belongs to,
	 * so all IDs are valid
	 */
	for _, group := range groups {
		channelIDs.Add(group.ID)
	}

	return channelIDs
}

func handleAddReaction(s *slack.Client, event *slack.ReactionAddedEvent, validChannels mapset.Set) {
	if !validChannels.Contains(event.Item.Channel) {
		log.WithFields(log.Fields{
			"channel_id": event.Item.Channel,
		}).Info("Skipping channel")

		return
	}
	if event.Reaction == "partyparrot" && event.Item.Type == "message" {

		log.WithFields(log.Fields{
			"channel_id": event.Item.Channel,
		}).Info("Making it party")

		msgRef := slack.NewRefToMessage(event.Item.Channel, event.Item.Timestamp)
		indexes := rand.Perm(len(parrots))
		for _, parrotIndex := range indexes[:23] {
			s.AddReaction(parrots[parrotIndex], msgRef)
		}
	}
}

func handleRemoveReaction(s *slack.Client, event *slack.ReactionRemovedEvent, validChannels mapset.Set) {
	if !validChannels.Contains(event.Item.Channel) {
		log.WithFields(log.Fields{
			"channel_id": event.Item.Channel,
		}).Info("Skipping channel")
		return
	}
	if event.Reaction == "partyparrot" {

		log.WithFields(log.Fields{
			"channel_id": event.Item.Channel,
		}).Info("Removing the party")

		msgRef := slack.NewRefToMessage(event.Item.Channel, event.Item.Timestamp)
		for _, parrot := range parrots {
			s.RemoveReaction(parrot, msgRef)
		}
	}
}
