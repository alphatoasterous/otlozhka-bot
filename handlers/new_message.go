package handlers

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/alphatoasterous/otlozhka-bot/api_utils"
	"github.com/alphatoasterous/otlozhka-bot/config"
	"github.com/alphatoasterous/otlozhka-bot/utils"
	"github.com/rs/zerolog/log"
	"slices"
	"strings"
)

var messages = config.BotConfig.MessageHandler
var regexes = config.BotConfig.CompiledRegexes

func messagePosts(peerID int, vkCommunity *api.VK, foundPosts []object.WallWallpost) {
	if len(messages.PostponedPostsFoundMsgs) != 0 { // if post found messages are defined
		message := api_utils.CreateMessageSendBuilderText(
			utils.GetRandomItemFromStrArray(messages.PostponedPostsFoundMsgs)) // send random message to user
		message.PeerID(peerID)
		_, err := vkCommunity.MessagesSend(message.Params)
		if err != nil {
			log.Fatal().Err(err)
		}
	}
	for _, post := range foundPosts {
		msg := api_utils.CreateMessageSendBuilderByPost(post)
		msg.PeerID(peerID)
		_, err := vkCommunity.MessagesSend(msg.Params)
		if err != nil {
			log.Fatal().Err(err)
		}
	}
}

func NewMessageHandler(obj events.MessageNewObject, vkCommunity *api.VK,
	vkUser *api.VK, domain string, groupManagerIDs []int, storage *WallpostStorage) {
	const communityChatID = 2000000004 // Community group chat
	incomingMessageText := strings.ToLower(obj.Message.Text)
	if obj.Message.PeerID != communityChatID { // Checks if message camen't from community group chat
		if slices.Contains(groupManagerIDs, obj.Message.PeerID) { // If message came from community management
			switch {
			case regexes.UpdateStorage.MatchString(incomingMessageText):
				log.Debug().Msgf("Update storage message[id%d]: %s", obj.Message.PeerID, obj.Message.Text)
				previousWallpostCount := storage.GetWallpostCount()
				storage.UpdateWallpostStorage(vkUser, domain)
				message := api_utils.CreateMessageSendBuilderText("")
				if storage.GetWallpostCount()-previousWallpostCount >= 10 {
					message.Message(utils.GetRandomItemFromStrArray(messages.StorageUpdatedCommendMsgs))
				} else {
					message.Message(utils.GetRandomItemFromStrArray(messages.StorageUpdatedMsgs))
				}
				message.PeerID(obj.Message.PeerID)
				_, err := vkCommunity.MessagesSend(message.Params)
				if err != nil {
					log.Fatal().Err(err)
				}
			case regexes.PrintStorage.MatchString(incomingMessageText):
				log.Debug().Msgf("Update storage message[id%d]: %s", obj.Message.PeerID, obj.Message.Text)
				formattedCalendar, err := api_utils.GetFormattedCalendar(storage.GetWallposts(), "Europe/Moscow")
				if err != nil {
					log.Fatal().Err(err)
				}
				message := api_utils.CreateMessageSendBuilderText(formattedCalendar)
				message.PeerID(obj.Message.PeerID)
				_, err = vkCommunity.MessagesSend(message.Params)
				if err != nil {
					log.Fatal().Err(err)
				}
			}
		}

	}

	switch {
	case regexes.Otlozhka.MatchString(incomingMessageText):
		log.Printf("Incoming message[id%d]: %s", obj.Message.PeerID, obj.Message.Text)
		if storage.CheckWallpostStorageNeedsUpdate() {
			storage.UpdateWallpostStorage(vkUser, domain)
		}
		posts := storage.GetWallposts()
		foundPosts := GetWallpostsByPeerID(obj.Message.PeerID, posts)
		if len(foundPosts) != 0 {
			messagePosts(obj.Message.PeerID, vkCommunity, foundPosts)
		} else {
			message := api_utils.CreateMessageSendBuilderText(
				utils.GetRandomItemFromStrArray(messages.NoPostponedPostsFoundMsgs))
			message.PeerID(obj.Message.PeerID)
			_, err := vkCommunity.MessagesSend(message.Params)
			if err != nil {
				log.Fatal().Err(err)
			}
		}
	}
}
