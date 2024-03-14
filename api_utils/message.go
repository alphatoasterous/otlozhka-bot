package api_utils

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
	"github.com/alphatoasterous/otlozhka-bot/lang"
	"github.com/alphatoasterous/otlozhka-bot/utils"
	"log"
	"time"
)

const RandomId = 0

var strings = lang.Lang.Message

func formatAttachments(attachment object.WallWallpostAttachment) string {
	// TODO: Ugly.
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
	return utils.RemoveTrailingComma(attachmentString)
}

func getPublicationDate(postDate int) string {
	t := time.Unix(int64(postDate), 0)
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Fatal("Error loading timezone:", err)
	}
	t = t.In(loc)
	formattedTime := t.Format("02.01.2006 15:04:05 MSK")
	return formattedTime
}

func getMessageText(post object.WallWallpost) string {
	// TODO: Ugly. Put magic strings somewhere else for ease of translation. i18n maybe?
	// Get beautiful-ish formatted string for message text
	// 1. Get publication date
	const newline = "\n"
	var msgText string
	msgText += strings.MessagePostDate + getPublicationDate(post.Date) + newline
	// 2. Get audio attachment info in case of it not attaching to the message properly(thanks, VK!)
	audioObject := post.Attachments[0].Audio
	if audioObject.ToAttachment() != "audio0_0" {
		msgText += strings.MessagePostAudio + audioObject.Artist + "—" + audioObject.Title + newline
	}
	// 3. Get post text
	if post.Text != "" {
		msgText += strings.MessagePostText + post.Text + newline
	}

	return msgText
}

func CreateMessageSendBuilderByPost(post object.WallWallpost) *params.MessagesSendBuilder {
	msg := params.NewMessagesSendBuilder()
	msg.Message(getMessageText(post))
	msg.Attachment(formatAttachments(post.Attachments[0]))
	msg.RandomID(RandomId)
	return msg
}

func CreateMessageSendBuilderText(text string) *params.MessagesSendBuilder {
	msg := params.NewMessagesSendBuilder()
	msg.Message(text)
	msg.RandomID(RandomId)
	return msg
}
