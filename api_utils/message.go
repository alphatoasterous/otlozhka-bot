package api_utils

import (
	"errors"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/alphatoasterous/otlozhka-bot/config"
	"github.com/alphatoasterous/otlozhka-bot/utils"
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
	"time"
)

const RandomId = 0

var messageBuilderConfig = config.BotConfig.MessageBuilder

func extractFormattedAttachmentsFromWallpost(attachment object.WallWallpostAttachment) string {
	var attachmentString string
	if api.FmtValue(attachment.Photo, 1) != "photo0_0" {
		attachmentString += attachment.Photo.ToAttachment() + ","
	}
	if api.FmtValue(attachment.Video, 1) != "video0_0" {
		attachmentString += attachment.Video.ToAttachment() + ","
	}
	if api.FmtValue(attachment.Audio, 1) != "audio0_0" {
		attachmentString += attachment.Audio.ToAttachment() + ","
	}
	if api.FmtValue(attachment.Doc, 1) != "doc0_0" {
		attachmentString += attachment.Doc.ToAttachment() + ","
	}
	return attachmentString
}

func getPublicationDate(postDate int) string {
	t := time.Unix(int64(postDate), 0)
	loc, err := time.LoadLocation(messageBuilderConfig.Timezone)
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading timezone")
	}
	t = t.In(loc)
	formattedTime := t.Format(messageBuilderConfig.TimeFormat)
	return formattedTime
}

func getMessageText(post object.WallWallpost) string {
	return fmt.Sprintf(messageBuilderConfig.MessageFormat, getPublicationDate(post.Date), post.Text)
}

func getPostAudios(post object.WallWallpost) ([]object.AudioAudio, error) {
	if len(post.Attachments) > 0 {
		var postAudios []object.AudioAudio
		for _, attachment := range post.Attachments {
			if attachment.Audio.ID != 0 {
				postAudios = append(postAudios, attachment.Audio)
			}
		}
		if postAudios == nil {
			return nil, errors.New("getPostAudios: Post doesn't contain any audio")
		}
		return postAudios, nil
	}
	return nil, errors.New("getPostAudios: Post doesn't contain any audio")
}

func getAudioArtistTitle(audio object.AudioAudio) string {
	return fmt.Sprintf("%s - %s", audio.Artist, audio.Title)
}

func CreateMessageSendBuilderByPost(post object.WallWallpost) *params.MessagesSendBuilder {
	msg := params.NewMessagesSendBuilder()
	msg.Message(getMessageText(post))
	if len(post.Attachments) > 0 {
		var collectedAttachments string
		for _, attachment := range post.Attachments {
			collectedAttachments += extractFormattedAttachmentsFromWallpost(attachment)
		}
		msg.Attachment(collectedAttachments)
	}
	msg.RandomID(RandomId)
	return msg
}

func CreateMessageSendBuilderText(text string) *params.MessagesSendBuilder {
	msg := params.NewMessagesSendBuilder()
	msg.Message(text)
	msg.RandomID(RandomId)
	return msg
}

func GetFormattedCalendar(posts []object.WallWallpost, timeZone string) (string, error) {
	loc, err := time.LoadLocation(timeZone)
	if err != nil {
		return "", err // Return an error if the timezone is invalid
	}

	// Group posts by date
	groupedPosts := make(map[time.Time][]object.WallWallpost)
	for _, post := range posts {
		dateTime := utils.UnixToTime(int64(post.Date), loc)
		date := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, loc) // Round to midnight
		groupedPosts[date] = append(groupedPosts[date], post)
	}

	// Sort keys of groupedPosts map
	var dates []time.Time
	for date := range groupedPosts {
		dates = append(dates, date)
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })

	// Create the formatted output
	var result string
	for _, date := range dates {
		dailyPosts := groupedPosts[date]
		dateStr := date.Format("02.01.2006")
		result += fmt.Sprintf("\n📅 %s:\n", dateStr)
		for _, post := range dailyPosts {
			timeStr := utils.UnixToTime(int64(post.Date), loc).Format("15:04")
			link := fmt.Sprintf("vk.com/wall%d_%d", post.OwnerID, post.ID)
			postAudios, err := getPostAudios(post)
			if err != nil {
				if err.Error() != "getPostAudios: Post doesn't contain any audio" {
					return "", err
				}
			} else {
				var audioTexts []string
				for _, audio := range postAudios {
					audioTexts = append(audioTexts, getAudioArtistTitle(audio))
				}
				link += " | 🎧: " + strings.Join(audioTexts, "; ")
			}
			result += fmt.Sprintf("%s: %s\n", timeStr, link)
		}
	}
	return result, nil
}
