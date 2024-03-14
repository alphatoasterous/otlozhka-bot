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

var lng = lang.Lang.NewMessageHandler

func getOtlozhka(obj events.MessageNewObject, vk *api.VK, vk_user *api.VK, group object.GroupsGroup) {
	posts, err := api_utils.GetAllPostponedWallposts(vk_user, group.ScreenName)
	if err != nil {
		log.Fatal(err)
	}
	foundPosts := api_utils.FindByFromID(posts, obj.Message.PeerID)
	if len(foundPosts) != 0 { // if foundPosts is not empty (posts found)
		message := api_utils.CreateMessageSendBuilderText(
			utils.GetRandomItemFromStrArray(lng.PostponedPostsFound))
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
				log.Print(err)
				msg = api_utils.CreateMessageSendBuilderText(lng.ErrorPostponedPostMessageFailed)
				_, err := vk.MessagesSend(msg.Params)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		message := api_utils.CreateMessageSendBuilderText(
			utils.GetRandomItemFromStrArray(lng.NoPostponedPostsFound))
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
		case lng.PostponedKeywordRegexCompiled.MatchString(strings.ToLower(obj.Message.Text)):
			log.Printf(lng.IncomingMessage, obj.Message.PeerID, obj.Message.Text)
			getOtlozhka(obj, vk, vk_user, group)
		}
	}
}
