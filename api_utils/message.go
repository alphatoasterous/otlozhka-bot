package api_utils

import (
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/alphatoasterous/otlozhka-bot/config"
	"github.com/alphatoasterous/otlozhka-bot/utils"
	"github.com/rs/zerolog/log"
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
	return utils.RemoveTrailingComma(attachmentString)
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

func CreateMessageSendBuilderByPost(post object.WallWallpost) *params.MessagesSendBuilder {
	msg := params.NewMessagesSendBuilder()
	msg.Message(getMessageText(post))
	if len(post.Attachments) > 0 {
		msg.Attachment(extractFormattedAttachmentsFromWallpost(post.Attachments[0]))
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
