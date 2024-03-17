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

var lng = lang.Lang.Message

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
	loc, err := time.LoadLocation(lng.Timezone)
	if err != nil {
		log.Fatal(lng.ErrorLoadingTimeZone, err)
	}
	t = t.In(loc)
	formattedTime := t.Format(lng.TimeFormat)
	return formattedTime
}

func getMessageText(post object.WallWallpost) string {
	// Get beautiful-ish formatted string for message text
	// 1. Get publication date
	const newline = "\n"
	var msgText string
	msgText += lng.MessagePostDate + getPublicationDate(post.Date) + newline
	// 2. Get audio attachment info in case of it not attaching to the message properly(thanks, VK!)
	// weather update: after live testing it is deemed unnecessary. Remove it, future me.
	/*
		audioObject := post.Attachments[0].Audio
		if audioObject.ToAttachment() != "audio0_0" {
			msgText += lng.MessagePostAudio + audioObject.Artist + "—" + audioObject.Title + newline
		}
	*/
	// 3. Get post text
	if post.Text != "" {
		msgText += lng.MessagePostText + post.Text + newline
	}

	return msgText
}

func CreateMessageSendBuilderByPost(post object.WallWallpost) *params.MessagesSendBuilder {
	msg := params.NewMessagesSendBuilder()
	msg.Message(getMessageText(post))
	msg.Attachment(extractFormattedAttachmentsFromWallpost(post.Attachments[0]))
	msg.RandomID(RandomId)
	return msg
}

func CreateMessageSendBuilderText(text string) *params.MessagesSendBuilder {
	msg := params.NewMessagesSendBuilder()
	msg.Message(text)
	msg.RandomID(RandomId)
	return msg
}
