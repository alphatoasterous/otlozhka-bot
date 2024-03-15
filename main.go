package main

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/alphatoasterous/otlozhka-bot/api_utils"
	"github.com/alphatoasterous/otlozhka-bot/handlers"
	"github.com/alphatoasterous/otlozhka-bot/lang"
	"github.com/alphatoasterous/otlozhka-bot/utils"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
)

var lng = lang.Lang.Main

func init() {
	// Loading environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Print(lng.ErrorDotenvFailed)
	}

}

func main() {
	// Setting up user and community API instances
	userToken := os.Getenv("OTLOZHKA_USER_TOKEN")
	communityToken := os.Getenv("OTLOZHKA_COMMUNITY_TOKEN")
	vk := api.NewVK(communityToken)
	vk_user := api.NewVK(userToken)
	vk.Limit, _ = utils.StringToInt(os.Getenv("OTLOZHKA_COMMUNITY_RATELIMIT"))
	vk.EnableMessagePack()
	vk.EnableZstd()
	vk_user.Limit, _ = utils.StringToInt(os.Getenv("OTLOZHKA_USER_RATELIMIT"))
	vk_user.EnableMessagePack()
	vk_user.EnableZstd()

	// Getting group information via community token
	group := api_utils.GetGroupInfo(vk)[0]

	// Setting up Long Poll
	lp, err := longpoll.NewLongPoll(vk, group.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Passing NewMessageHandler to a MessageNew event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		handlers.NewMessageHandler(obj, vk, vk_user, group)
	})

	// Run Bots Long Poll
	log.Println(lng.StartLongPollMsg)
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
