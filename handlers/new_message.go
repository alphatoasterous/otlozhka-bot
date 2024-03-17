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

func getPostponedPostsByPeerID(peerID int, vkCommunity *api.VK, vkUser *api.VK, domain string) {
	// getPostponedPostsByPeerID retrieves postponed posts that match given peer ID.
	//
	// This function fetches postponed posts by the specified peer ID using the VK API.
	// It requires a valid VK API community token, a VK user token, and the short link of the searchable group.
	//
	// Parameters:
	//   - peerID: An integer representing the ID of the peer.
	//   - vkCommunity: A pointer to an initialized VK instance with Community access token.
	//   - vkUser: A pointer to an initialized VK instance with User access token.
	//   - domain: Short address of the searchable community. Obtained via object.GroupsGroup.ScreenName.
	posts, err := api_utils.GetAllPostponedWallposts(vkUser, domain)
	if err != nil {
		log.Fatal(err)
	}
	foundPosts := api_utils.FindByFromID(posts, peerID)
	if len(foundPosts) != 0 { // if posts found
		message := api_utils.CreateMessageSendBuilderText(
			utils.GetRandomItemFromStrArray(lng.PostponedPostsFound))
		message.PeerID(peerID)
		_, err := vkCommunity.MessagesSend(message.Params)
		if err != nil {
			log.Fatal(err)
		}
		for _, post := range foundPosts {
			msg := api_utils.CreateMessageSendBuilderByPost(post)
			msg.PeerID(peerID)
			_, err := vkCommunity.MessagesSend(msg.Params)
			if err != nil {
				log.Print(err)
				msg = api_utils.CreateMessageSendBuilderText(lng.ErrorPostponedPostMessageFailed)
				_, err := vkCommunity.MessagesSend(msg.Params)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	} else {
		message := api_utils.CreateMessageSendBuilderText(
			utils.GetRandomItemFromStrArray(lng.NoPostponedPostsFound))
		message.PeerID(peerID)
		_, err := vkCommunity.MessagesSend(message.Params)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewMessageHandler(obj events.MessageNewObject, vkCommunity *api.VK, vkUser *api.VK, group object.GroupsGroup) {
	const communityChatID = 2000000004         // Community group chat
	if obj.Message.PeerID != communityChatID { // Checks if message didn't come from community group chat
		switch {
		case lng.PostponedKeywordRegexCompiled.MatchString(strings.ToLower(obj.Message.Text)):
			log.Printf(lng.IncomingMessage, obj.Message.PeerID, obj.Message.Text)
			getPostponedPostsByPeerID(obj.Message.PeerID, vkCommunity, vkUser, group.ScreenName)
		}
	}
}
