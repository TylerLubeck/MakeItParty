package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/nlopes/slack"
)

var parrots = []string{
	"aussieparrot",
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
}

func HandleSlackEvents(s *slack.Client) {
	rtm := s.NewRTM()

	auth, err := s.AuthTest()

	if err != nil {
		return
	}

	log.WithFields(log.Fields{
		"slack_user_id": auth.UserID,
	}).Info("Logging in to Slack...")

	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ReactionAddedEvent:
				handleAddReaction(s, ev)
			case *slack.ReactionRemovedEvent:
				handleRemoveReaction(s, ev)
			default:
				//fmt.Printf("Other Event: %v\n", msg.Data)
			}
		}
	}

}

func handleAddReaction(s *slack.Client, event *slack.ReactionAddedEvent) {
	if event.Reaction == "partyparrot" && event.Item.Type == "message" {
		msgRef := slack.NewRefToMessage(event.Item.Channel, event.Item.Timestamp)
		for _, parrot := range parrots {
			s.AddReaction(parrot, msgRef)
		}
	}
}

func handleRemoveReaction(s *slack.Client, event *slack.ReactionRemovedEvent) {
	if event.Reaction == "partyparrot" {
		msgRef := slack.NewRefToMessage(event.Item.Channel, event.Item.Timestamp)
		for _, parrot := range parrots {
			s.RemoveReaction(parrot, msgRef)
		}
	}
}
