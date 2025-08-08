package api_utils

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/alphatoasterous/otlozhka-bot/config"
	"github.com/alphatoasterous/otlozhka-bot/logging"
	"github.com/alphatoasterous/otlozhka-bot/utils"
)

const RandomId = 0

var messageBuilderConfig = config.BotConfig.MessageBuilder

// extractFormattedAttachmentsFromWallpost extracts and formats attachments from a WallWallpostAttachment.
// It compiles a string of attachment identifiers for photos, videos, audio, and documents.
// Each type of attachment is checked for existence before appending its identifier to the result string.
// The returned string can be directly used in API calls that require attachment string.
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

// getReadableDate formats a UNIX timestamp into a readable date and time based on a specified timezone.
// If an error occurs while loading the timezone, it logs the error and exits fatally.
// Returns the formatted time as a string.
func getReadableDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	loc, err := time.LoadLocation(messageBuilderConfig.Timezone)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("Error loading timezone")
	}
	t = t.In(loc)
	formattedTime := t.Format(messageBuilderConfig.TimeFormat)
	return formattedTime
}

// getMessageText constructs the message text for a given post using a configurable format.
// It integrates the post's publication date and text content, formatting the date using getReadableDate.
func getMessageText(post object.WallWallpost) string {
	return fmt.Sprintf(messageBuilderConfig.MessageFormat, getReadableDate(int64(post.Date)), post.Text)
}

// getPostAudios extracts audio attachments from a WallWallpost.
// Returns a slice of AudioAudio objects or an error if no audio attachments are found.
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

// getAudioArtistTitle formats the artist and title of an audio into a single string.
func getAudioArtistTitle(audio object.AudioAudio) string {
	return fmt.Sprintf("%s - %s", audio.Artist, audio.Title)
}

// CreateMessageSendBuilderByPost prepares a message builder for sending messages,
// incorporating text and attachments based on a provided WallWallpost.
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

// CreateMessageSendBuilderText creates a simple message send builder with text content.
func CreateMessageSendBuilderText(text string) *params.MessagesSendBuilder {

	// TODO: Rework this and entire messaging routine altogether by segmenting text into array of messages if needed
	VK_MESSAGE_LIMIT := 4096
	if len(text) > VK_MESSAGE_LIMIT {
		text = text[:VK_MESSAGE_LIMIT-3] + "..."
	}
	msg := params.NewMessagesSendBuilder()
	msg.Message(text)
	msg.RandomID(RandomId)
	return msg
}

// GetFormattedCalendar groups wall posts by date and formats them into a readable calendar view.
// The formatting takes into account the timezone, sorting posts by date.
// Returns a formatted string representing the post calendar or an error if an issue occurs during formatting.
func GetFormattedCalendar(posts []object.WallWallpost, timezone string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		// If loading timezone errors out: default to timezone provided via config. timezone variable may be used in
		// future, if I could somehow get user's timezone from VK API.
		loc, err = time.LoadLocation(messageBuilderConfig.Timezone)
		if err != nil {
			return "", err // Return an error if the timezone is invalid
		}
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
