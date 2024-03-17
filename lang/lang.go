package lang

import (
	"flag"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"regexp"
)

type (
	ProjectStrings struct {
		Main              mainStrings
		Message           messageStrings
		NewMessageHandler newMessageHandlerStrings
	}
	mainStrings struct {
		StartLongPollMsg string

		ErrorDotenvFailed string
	}
	newMessageHandlerStrings struct {
		IncomingMessage               string
		PostponedKeywordRegex         string
		PostponedKeywordRegexCompiled *regexp.Regexp
		PostponedPostsFound           []string
		NoPostponedPostsFound         []string

		ErrorPostponedPostMessageFailed string
	}
	messageStrings struct {
		MessagePostDate  string
		MessagePostAudio string
		MessagePostText  string
		TimeFormat       string
		Timezone         string

		ErrorLoadingTimeZone string
	}
)

var Lang ProjectStrings

func init() {
	stringsDefault := "strings.toml"
	stringsFilename := flag.String("strings", stringsDefault, "Specify strings filename")
	tomlFile, err := os.ReadFile(*stringsFilename)
	if err != nil {
		log.Fatal("Error reading strings.toml / Ошибка чтения strings.toml")
	}
	err = toml.Unmarshal(tomlFile, &Lang)
	if err != nil {
		log.Fatal("Error parsing strings.toml / Ошибка парсинга strings.toml")
	}
	Lang.NewMessageHandler.PostponedKeywordRegexCompiled =
		regexp.MustCompile(Lang.NewMessageHandler.PostponedKeywordRegex)
}