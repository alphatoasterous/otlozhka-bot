package main

import (
	"context"
	"flag"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/alphatoasterous/otlozhka-bot/api_utils"
	"github.com/alphatoasterous/otlozhka-bot/configs"
	"github.com/alphatoasterous/otlozhka-bot/handlers"
	"github.com/alphatoasterous/otlozhka-bot/utils"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
)

var lng = configs.Lang.Main

func init() {
	// Loading environment variables from a filename provided by arguments / default .env file
	dotenvDefault := ".env"
	dotenvFilename := flag.String("dotenv", dotenvDefault, "Specify dotenv filename")
	flag.Parse()
	err := godotenv.Load(*dotenvFilename)
	if err != nil {
		log.Print(lng.ErrorDotenvFailed)
	}
}

func main() {
	// Setting up community API instance
	communityToken := os.Getenv("OTLOZHKA_COMMUNITY_TOKEN")
	vkCommunity := api.NewVK(communityToken)
	vkCommunity.Limit, _ = utils.StringToInt(os.Getenv("OTLOZHKA_COMMUNITY_RATELIMIT"))
	vkCommunity.EnableMessagePack()
	vkCommunity.EnableZstd()

	// Setting up user API instance
	userToken := os.Getenv("OTLOZHKA_USER_TOKEN")
	vkUser := api.NewVK(userToken)
	vkUser.EnableMessagePack()
	vkUser.EnableZstd()
	vkUser.Limit, _ = utils.StringToInt(os.Getenv("OTLOZHKA_USER_RATELIMIT"))

	// Setting up wallpost storage
	keepAlive, _ := utils.StringToInt(os.Getenv("OTLOZHKA_STORAGE_KEEPALIVE"))
	wallpostStorage := handlers.NewWallpostStorage(int64(keepAlive))

	// Getting group information via community VK instance
	group := api_utils.GetGroupInfo(vkCommunity)[0]
	domain := group.ScreenName
	groupManagerIDs := api_utils.GetGroupManagerIDs(vkUser, domain)

	// Setting up Long Poll
	lp, err := longpoll.NewLongPoll(vkCommunity, group.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Passing NewMessageHandler to a MessageNew event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		handlers.NewMessageHandler(obj, vkCommunity, vkUser, domain, groupManagerIDs, wallpostStorage)
	})

	// Run Bots Long Poll
	log.Println(lng.StartLongPollMsg)
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
