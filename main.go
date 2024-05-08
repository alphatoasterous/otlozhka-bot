package main

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/alphatoasterous/otlozhka-bot/api_utils"
	"github.com/alphatoasterous/otlozhka-bot/config"
	"github.com/alphatoasterous/otlozhka-bot/handlers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setting up zerolog logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting up otlozhka-bot...")
	botConfig := config.BotConfig.Main

	// Setting up community API instance
	vkCommunity := api.NewVK(botConfig.CommunityToken)
	vkCommunity.Limit = botConfig.CommunityAPIRateLimit
	vkCommunity.EnableMessagePack()
	vkCommunity.EnableZstd()
	log.Debug().Msg("Community API instance set up")

	// Setting up user API instance
	vkUser := api.NewVK(botConfig.UserToken)
	vkUser.EnableMessagePack()
	vkUser.EnableZstd()
	vkUser.Limit = botConfig.UserAPIRateLimit
	log.Debug().Msg("User API instance set up")

	// Getting group information via community VK instance
	group := api_utils.GetGroupInfo(vkCommunity)[0]
	domain := group.ScreenName
	groupManagerIDs := api_utils.GetGroupManagerIDs(vkUser, domain)

	// Setting up wallpost storage
	keepAlive := botConfig.StorageKeepAlive
	wallpostStorage := handlers.NewWallpostStorage(int64(keepAlive))
	wallpostStorage.UpdateWallpostStorage(vkUser, domain)
	log.Debug().Msg("Wallpost Storage instance set up")

	// Setting up Long Poll
	lp, err := longpoll.NewLongPoll(vkCommunity, group.ID)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Debug().Msg("Long Poll set up")

	// Passing NewMessageHandler to a MessageNew event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		handlers.NewMessageHandler(obj, vkCommunity, vkUser, domain, groupManagerIDs, wallpostStorage)
	})

	// Run Bots Long Poll
	log.Info().Msg("otlozhka-bot set, running Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal().Err(err)
	}
}
