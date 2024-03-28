package main

import (
	"context"
	"flag"
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

	// Getting group information via community VK instance
	group := api_utils.GetGroupInfo(vkCommunity)[0]
	domain := group.ScreenName

	// Setting up user API wallpost handler instance
	userToken := os.Getenv("OTLOZHKA_USER_TOKEN")
	keepAlive, _ := utils.StringToInt(os.Getenv("OTLOZHKA_STORAGE_KEEPALIVE"))
	vkUser := api.NewVK(userToken)
	vkUser.EnableMessagePack()
	vkUser.EnableZstd()
	vkUser.Limit, _ = utils.StringToInt(os.Getenv("OTLOZHKA_USER_RATELIMIT"))
	wallpostStorage := handlers.NewWallpostStorage(int64(keepAlive))

	// Setting up Long Poll
	lp, err := longpoll.NewLongPoll(vkCommunity, group.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Passing NewMessageHandler to a MessageNew event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		handlers.NewMessageHandler(obj, vkCommunity, vkUser, domain, wallpostStorage)
	})

	// Run Bots Long Poll
	log.Println(lng.StartLongPollMsg)
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
