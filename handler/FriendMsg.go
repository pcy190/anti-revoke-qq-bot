package handler

import (
	. "RevokeBot/model"
	"RevokeBot/util"
	"fmt"
	"strconv"
	"strings"
)

func FriendMsgHandler(msg Message) {
	data := msg.CurrentPacket.Data
	fmt.Println(data.Content)
	if data.Content == "!grant clean" {
		GrantList.Clear()
		Send(int(data.FromUin), 1, "Grant: "+GrantList.String())
	} else if data.Content == "!status" {
		var FlashStatus string
		if RevokeGroupList.IsEmpty() {
			Send(int(data.FromUin), 1, "Revoke is "+strconv.Itoa(EnableRevokePrompt)+"\nRevoke list is empty"+FlashStatus+
				"\nGrant: "+GrantList.String()+"\nisGrantAll:"+strconv.FormatBool(GrantAll))
		} else {
			Send(int(data.FromUin), 1, "Revoke is "+strconv.Itoa(EnableRevokePrompt)+"\nRevoke list is "+RevokeGroupList.String()+FlashStatus+
				"\nGrant: "+GrantList.String()+"\nisGrantAll:"+strconv.FormatBool(GrantAll))
		}
	} else if data.Content == "!revoke on" {
		EnableRevokePrompt = 1
		Send(int(data.FromUin), 1, "RevokePrompt is on")
	} else if data.Content == "!revoke off" {
		EnableRevokePrompt = 0
		Send(int(data.FromUin), 1, "Revoke is off")
	}  else {
		arg := strings.Split(data.Content, " ")
		if len(arg) < 3 {
			return
		}
		if arg[0] == "!revoke" && arg[1] == "add" {
			RevokeGroupList.Add(strconv.Atoi(arg[2]))
			Send(int(data.FromUin), 1, "Add successfully\nUse !status to check status")

		}
		if arg[0] == "!revoke" && arg[1] == "sub" {
			RevokeGroupList.Remove(strconv.Atoi(arg[2]))
			Send(int(data.FromUin), 1, "Remove successfully\nUse !status to check status")
		}
	}
	util.WriteConfig()
}
