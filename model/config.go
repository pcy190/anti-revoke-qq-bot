package iotqq

import (
	"gopkg.in/fatih/set.v0"
)

// login qq id
var LoginQQ int64 = 0
// master qq id
var Master int64 = 0
var EnableRevokePrompt = 1
var RevokeGroupList = set.New(set.ThreadSafe)
var GrantList = set.New(set.ThreadSafe)

var DEBUG bool = false

var GrantAll = false // allow all users to set revoke bot

// replace the following config
var MysqlIP = "127.0.0.1"
var MysqlPort=3306
var MysqlUsername=""
var MysqlPassword =""
var MysqlDbName="qq"
var ApiIp = "0.0.0.0"
var ApiPort = 0
var ConfigName="bot.conf"

var PublicHelpPrompt = `!revoke on/off/status
!status
!grant add/remove AT
!grant status/all/off
!zuan on/off/level(0-4)
!zuan add/remove AT
!p add/remove AT`
