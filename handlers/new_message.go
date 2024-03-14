package handlers

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/alphatoasterous/otlozhka-bot/api_utils"
	"github.com/alphatoasterous/otlozhka-bot/lang"
	"github.com/alphatoasterous/otlozhka-bot/utils"
	"log"
	"strings"
)

var language = lang.Lang.NewMessageHandler

func getOtlozhka(obj events.MessageNewObject, vk *api.VK, vk_user *api.VK, group object.GroupsGroup) {
	posts, err := api_utils.GetAllPostponedWallposts(vk_user, group.ScreenName)
	if err != nil {
		log.Fatal(err)
	}
	foundPosts := api_utils.FindByFromID(posts, obj.Message.PeerID)
	if len(foundPosts) != 0 { // if foundPosts is not empty (posts found)
		message := api_utils.CreateMessageSendBuilderText(
			utils.GetRandomItemFromStrArray(language.PostponedPostsFound))
		message.PeerID(obj.Message.PeerID)
		_, err := vk.MessagesSend(message.Params)
		if err != nil {
			log.Fatal(err)
		}
		for _, post := range foundPosts {
			msg := api_utils.CreateMessageSendBuilderByPost(post)
			msg.PeerID(obj.Message.PeerID)
			_, err := vk.MessagesSend(msg.Params)
			if err != nil {
				log.Print(err) // TODO: Test out if one failed post will hang all up
				msg = api_utils.CreateMessageSendBuilderText(language.ErrorPostponedPostMessageFailed)
				_, err := vk.MessagesSend(msg.Params)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		message := api_utils.CreateMessageSendBuilderText(
			utils.GetRandomItemFromStrArray(language.NoPostponedPostsFound))
		message.PeerID(obj.Message.PeerID)
		_, err := vk.MessagesSend(message.Params)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewMessageHandler(obj events.MessageNewObject, vk *api.VK, vk_user *api.VK, group object.GroupsGroup) {
	const communityChatID = 2000000004
	if obj.Message.PeerID != communityChatID { // Checks if message didn't come from community messages
		switch {
		case language.PostponedKeywordRegexCompiled.MatchString(strings.ToLower(obj.Message.Text)):
			log.Printf(language.IncomingMessage, obj.Message.PeerID, obj.Message.Text)
			getOtlozhka(obj, vk, vk_user, group)
		}
	}
}
