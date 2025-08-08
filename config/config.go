package config

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/pelletier/go-toml/v2"
)

type (
	BotConfiguration struct {
		Main            mainConfig
		ZerologConfig   ZerologConfiguration
		MessageBuilder  messageBuilderConfig
		MessageHandler  messageHandlerConfig
		CompiledRegexes compiledRegexes
	}

	ZerologConfiguration struct {
		// Enable console logging
		ConsoleLoggingEnabled bool

		// EncodeLogsAsJson makes the log framework log JSON
		EncodeLogsAsJson bool
		// FileLoggingEnabled makes the framework log to a file
		// the fields below can be skipped if this value is false!
		FileLoggingEnabled bool
		// Directory to log to to when filelogging is enabled
		Directory string
		// Filename is the name of the logfile which will be placed inside the directory
		Filename string
		// MaxSize the max size in MB of the logfile before it's rolled
		MaxSize int
		// MaxBackups the max number of rolled files to keep
		MaxBackups int
		// MaxAge the max age in days to keep a logfile
		MaxAge int
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
		PrintStorageRegex  string

		StorageUpdatedMsgs        []string
		StorageUpdatedCommendMsgs []string
		StorageEmptyMsgs          []string

		PostponedPostsFoundMsgs   []string
		NoPostponedPostsFoundMsgs []string
	}

	compiledRegexes struct {
		Otlozhka      *regexp.Regexp
		UpdateStorage *regexp.Regexp
		PrintStorage  *regexp.Regexp
	}
)

func DefaultBotConfiguration() BotConfiguration {
	return BotConfiguration{

		Main: mainConfig{
			UserToken:             "",
			CommunityToken:        "",
			CommunityAPIRateLimit: 5,
			UserAPIRateLimit:      1,
			StorageKeepAlive:      900,
		},
		ZerologConfig: ZerologConfiguration{
			ConsoleLoggingEnabled: true,
			EncodeLogsAsJson:      true,
			FileLoggingEnabled:    true,
			Directory:             "logs",
			Filename:              "otlozhka-bot.log",
			MaxSize:               5,
			MaxBackups:            5,
			MaxAge:                30,
		},
		MessageBuilder: messageBuilderConfig{
			MessageFormat: "📅 : %s\n📝: %s",
			TimeFormat:    "02.01.2006 15:04:05",
			Timezone:      "Europe/Moscow",
		},
		MessageHandler: messageHandlerConfig{
			OtlozhkaRegex:             "отложк[ауе]",
			UpdateStorageRegex:        "обнови",
			PrintStorageRegex:         "календарь",
			StorageUpdatedMsgs:        []string{"Хранилище синхронизировано. Следующее обновление через 15 минут."},
			StorageUpdatedCommendMsgs: []string{"Хранилище синхронизировано. Спасибо за Ваш труд!"},
			StorageEmptyMsgs:          []string{"В хранилище пусто. Вероятно, в сообществе нет отложенных постов."},
			PostponedPostsFoundMsgs:   []string{""},
			NoPostponedPostsFoundMsgs: []string{"Отложенных постов не найдено."},
		},
	}
}

var BotConfig BotConfiguration

func init() {
	configFilename := flag.String("config", "config.toml", "Specify config filename")

	BotConfig = DefaultBotConfiguration()

	fmt.Printf("DEBUG: Loading configuration from a file")
	// Check if config.toml exists
	_, err := os.Stat(*configFilename)
	if os.IsNotExist(err) {
		fmt.Printf("WARNING: %s does not exist, creating it with default parameters\n", *configFilename)
		tomlBotConfig, err := toml.Marshal(BotConfig)
		if err != nil {
			fmt.Printf("ERROR: Cannot marshal BotConfig: %v, path: %s\n", err, *configFilename)
			return
		}
		err = os.WriteFile(*configFilename, tomlBotConfig, 0644)
		if err != nil {
			fmt.Printf("ERROR: Error writing marshalled BotConfig to a %s", *configFilename)
			return
		}
	} else {
		fmt.Printf("DEBUG: %s does exist, unmarshalling it\n", *configFilename)
		tomlFile, err := os.ReadFile(*configFilename)
		if err != nil {
			fmt.Printf("ERROR: Error reading %s: %v\n", *configFilename, err)
			return
		}
		err = toml.Unmarshal(tomlFile, &BotConfig)
		if err != nil {
			fmt.Printf("ERROR: Error unmarshalling %s: %v\n", *configFilename, err)
			return
		}
	}

	if BotConfig.Main.UserToken == "" || BotConfig.Main.CommunityToken == "" {
		fmt.Println("ERROR: No UserToken or CommunityToken provided")
		return
	}

	BotConfig.CompiledRegexes.Otlozhka = regexp.MustCompile(BotConfig.MessageHandler.OtlozhkaRegex)
	BotConfig.CompiledRegexes.UpdateStorage = regexp.MustCompile(BotConfig.MessageHandler.UpdateStorageRegex)
	BotConfig.CompiledRegexes.PrintStorage = regexp.MustCompile(BotConfig.MessageHandler.PrintStorageRegex)

}
