package config

import (
	"flag"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
)

type (
	BotConfiguration struct {
		Main            mainConfig
		MessageBuilder  messageBuilderConfig
		MessageHandler  messageHandlerConfig
		CompiledRegexes compiledRegexes
	}

	mainConfig struct {
		UserToken             string
		CommunityToken        string
		CommunityAPIRateLimit int
		UserAPIRateLimit      int
		StorageKeepAlive      int
	}

	messageBuilderConfig struct {
		MessageFormat string
		TimeFormat    string
		Timezone      string
	}

	messageHandlerConfig struct {
		OtlozhkaRegex      string
		UpdateStorageRegex string

		StorageUpdatedMsgs        []string
		StorageUpdatedCommendMsgs []string

		PostponedPostsFoundMsgs   []string
		NoPostponedPostsFoundMsgs []string
	}

	compiledRegexes struct {
		PostponedKeyword *regexp.Regexp
		UpdateStorage    *regexp.Regexp
	}
)

func DefaultBotConfiguration() BotConfiguration {
	const OtlozhkaRegex = "отложк[ауе]"
	const UpdateStorageRegex = "обнови"
	return BotConfiguration{

		Main: mainConfig{
			UserToken:             "",
			CommunityToken:        "",
			CommunityAPIRateLimit: 5,
			UserAPIRateLimit:      1,
			StorageKeepAlive:      900,
		},
		MessageBuilder: messageBuilderConfig{
			MessageFormat: "📅 : %s\n📝: %s",
			TimeFormat:    "02.01.2006 15:04:05",
			Timezone:      "Europe/Moscow",
		},
		MessageHandler: messageHandlerConfig{
			OtlozhkaRegex:             OtlozhkaRegex,
			UpdateStorageRegex:        UpdateStorageRegex,
			StorageUpdatedMsgs:        []string{"Хранилище синхронизировано. Следующее обновление через 15 минут."},
			StorageUpdatedCommendMsgs: []string{"Хранилище синхронизировано. Спасибо за Ваш труд!"},
			PostponedPostsFoundMsgs:   []string{""},
			NoPostponedPostsFoundMsgs: []string{"Отложенных постов не найдено."},
		},
		CompiledRegexes: compiledRegexes{
			PostponedKeyword: regexp.MustCompile(OtlozhkaRegex),
			UpdateStorage:    regexp.MustCompile(UpdateStorageRegex),
		},
	}
}

var BotConfig BotConfiguration

func init() {
	configFilename := flag.String("config", "config.toml", "Specify config filename")

	BotConfig = DefaultBotConfiguration()

	// Check if config.toml exists
	_, err := os.Stat(*configFilename)
	if os.IsNotExist(err) {
		log.Warn().Msgf("%s does not exist, creating it with default parameters", *configFilename)
		tomlBotConfig, err := toml.Marshal(BotConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Error marshalling BotConfig")
		}
		err = os.WriteFile(*configFilename, tomlBotConfig, 0644)
		if err != nil {
			log.Fatal().Err(err).Msgf("Error writing marshalled BotConfig to a %s", *configFilename)
		}
	} else {
		log.Warn().Msgf("%s does exist, unmarshalling it", *configFilename)
		tomlFile, err := os.ReadFile(*configFilename)
		if err != nil {
			log.Fatal().Err(err).Msgf("Error reading %s", *configFilename)
		}
		err = toml.Unmarshal(tomlFile, &BotConfig)
		if err != nil {
			log.Fatal().Err(err).Msgf("Error unmarshalling %s", *configFilename)
		}
	}
}
