package handler

import (
	"RevokeBot/util"
	"encoding/json"
	"fmt"
)
import . "RevokeBot/model"

func EventMsgHandler(msg Message) {
	data := msg.CurrentPacket.Data
	msgType := data.EventName
	if data.EventData.UserID == msg.CurrentQQ {
		return
	}
	switch msgType {
	case "ON_EVENT_GROUP_REVOKE":
		{
			handleRevokeContent(data.EventData.GroupID, data.EventData.MsgSeq, 2)
		}
	}

}

func handleRevokeContent(GroupID int, MsgSeq int, SendType int) {
	if EnableRevokePrompt != 1 {
		return
	}
	data := util.QueryRecord(GroupID, MsgSeq)
	if data.MsgType == "" {
		return
	}
	fmt.Println("Query:")
	fmt.Println(data)
	prefixStr := data.NickName + " revoke:\n"
	prefixStrInline := data.NickName + " revoke:"
	switch data.MsgType {
	case "TextMsg":
		{
			Send(GroupID, SendType, prefixStr+data.Content)
		}
	case "PicMsg":
		{
			var picMessage PicMessage
			if err := json.Unmarshal([]byte(data.Content), &picMessage); err == nil {
				SendPicAdvanced(GroupID, SendType, prefixStrInline+
					picMessage.Content, picMessage.FileMd5, picMessage.URL, picMessage.ForwordBuf)
			} else {
				fmt.Println(err)
			}
		}
	case "SmallFaceMsg":
		{
			var smallFaceMessage SmallFaceMessage
			if err := json.Unmarshal([]byte(data.Content), &smallFaceMessage); err == nil {
				SendA(GroupID, SendType, prefixStr+smallFaceMessage.Content, "TextMsg")
			} else {
				fmt.Println(err)
			}
		}
	case "AtMsg":
		{
			var atMessage AtMessage
			if err := json.Unmarshal([]byte(data.Content), &atMessage); err == nil {
				Send(GroupID, SendType, prefixStr+atMessage.Content)
			} else {
				fmt.Println(err)
			}
		}
	case "VoiceMsg":
		{
			var voiceMessage VoiceMessage
			if err := json.Unmarshal([]byte(data.Content), &voiceMessage); err == nil {
				SendVoice(GroupID, SendType, voiceMessage.URL, prefixStr+voiceMessage.Content)
			} else {
				fmt.Println(err)
			}
		}
	case "ReplayMsg":
		{
			var replayMessage ReplayMessage
			if err := json.Unmarshal([]byte(data.Content), &replayMessage); err == nil {
				Send(GroupID, SendType, prefixStr+"replay:"+replayMessage.SrcContent+"\n"+replayMessage.ReplayContent)
			} else {
				fmt.Println(err)
			}
		}
	case "XmlMsg":
		{
			SendA(GroupID, SendType, prefixStr+data.Content, "XmlMsg")
		}
	case "BigFaceMsg":
		{
			var bigFaceMessage BigFaceMessage
			if err := json.Unmarshal([]byte(data.Content), &bigFaceMessage); err == nil {
				SendBigFace(GroupID,SendType,bigFaceMessage.ForwordBuf)
			} else {
				fmt.Println(err)
			}
		}
	}
}
