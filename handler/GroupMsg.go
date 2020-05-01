package handler

import (
	. "RevokeBot/model"
	"RevokeBot/util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func isAuthorized(msg Message) bool {
	if GrantAll {
		return true
	}
	// NOTE : convert user id to int (not int64)
	uid := int(msg.CurrentPacket.Data.FromUserID)
	list := GrantList.List()
	for i := 0; i < len(list); i++ {
		if list[i] == uid {
			return true
		}
	}
	if msg.CurrentPacket.Data.FromUserID == Master || GrantList.Has(uid) {
		return true
	} else {
		Send(int(msg.CurrentPacket.Data.FromGroupID), 2, "Unauthorized")
		return false
	}
}

func GroupMsgHandler(msg Message) {
	if msg.CurrentQQ == msg.CurrentPacket.Data.FromUserID {
		return
	}

	var data = msg.CurrentPacket.Data
	inputs := strings.ReplaceAll(data.Content, "！", "!")
	if inputs == "!help" {
		Send(int(data.FromGroupID), 2, PublicHelpPrompt)
		return
	}
	if inputs == "!grant status" {
		if data.FromUserID == Master {
			Send(int(data.FromGroupID), 2, "root")
			return
		}
		if GrantList.Has(data.FromUserID) || GrantAll {
			Send(int(data.FromGroupID), 2, "Authorized")
		} else {
			Send(int(data.FromGroupID), 2, "Unauthorized")
		}
		return
	}
	if inputs == "!grant all" && msg.CurrentPacket.Data.FromUserID == Master {
		GrantAll = true
		Send(int(data.FromGroupID), 2, "Grant All")
	} else if inputs == "!grant off" && msg.CurrentPacket.Data.FromUserID == Master {
		GrantAll = false
		Send(int(data.FromGroupID), 2, "Grant Off")
	}

	changed := false
	if inputs == "!revoke on" {
		if !isAuthorized(msg) {
			return
		}
		changed = true
		RevokeGroupList.Add(data.FromGroupID)
		Send(data.FromGroupID, 2, "Revoke is on")
	} else if inputs == "!revoke off" {
		if !isAuthorized(msg) {
			return
		}
		changed = true
		RevokeGroupList.Remove(data.FromGroupID)
		Send(data.FromGroupID, 2, "Revoke is off")
	} else if inputs == "!revoke status" {
		if RevokeGroupList.Has(data.FromGroupID) {
			Send(data.FromGroupID, 2, "Revoke is on")
		} else {
			Send(data.FromGroupID, 2, "Revoke is off")
		}
	}
	if changed {
		util.WriteConfig()
	}

	if msg.CurrentPacket.Data.FromUserID == Master || GrantList.Has(msg.CurrentPacket.Data.FromUserID) {
		if strings.Index(inputs, "!grant add") != -1 {
			var atMessage AtMessage
			if err := json.Unmarshal([]byte(inputs), &atMessage); err == nil {
				id := atMessage.UserID
				if id != 0 {
					GrantList.Add(id)
					Send(data.FromGroupID, 2, "Grant "+strconv.Itoa(int(id)))
				}
			} else {
				fmt.Println(err)
			}
			util.WriteConfig()
		} else if strings.Index(inputs, "!grant remove") != -1 {
			var atMessage AtMessage
			if err := json.Unmarshal([]byte(inputs), &atMessage); err == nil {
				id := atMessage.UserID
				if id != 0 {
					GrantList.Remove(id)
					Send(data.FromGroupID, 2, "Grant remove:"+strconv.Itoa(int(id)))
				}
			} else {
				fmt.Println(err)
			}
			util.WriteConfig()
		}
	}
	if RevokeGroupList.Has(data.FromGroupID) {
		util.AddRecord(data.FromGroupID, data.MsgSeq, data.MsgRandom, data.MsgType, data.Content, data.FromNickName, data.MsgTime)
	}
}

type ResultData struct {
	Reason    string   `json:"reason"`
	Result    DataList `json:"result"`
	ErrorCode int      `json:"error_code"`
}

type DataList struct {
	Data []DataInfo `json:"data"`
}

type DataInfo struct {
	Content string
}
