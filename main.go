package main

import (
	"RevokeBot/handler"
	"RevokeBot/model"
	"RevokeBot/util"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"
	"unsafe"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var url1, qq string
var qqNum int64
var conf iotqq.Conf
var zanOkData, qd []int64

func periodCall(d time.Duration, f func()) {
	for x := range time.Tick(d) {
		f()
		log.Println(x)
	}
}

func resetPause() {

	m1 := len(zanOkData)
	for m := 0; m < m1; m++ {
		i := 0
		zanOkData = append(zanOkData[:i], zanOkData[i+1:]...)
	}
	m2 := len(qd)
	for m := 0; m < m2; m++ {
		i := 0
		qd = append(qd[:i], qd[i+1:]...)
	}
}
func SendJoin(c *gosocketio.Client) {
	log.Println("Get Connection Successfully!")
	result, err := c.Ack("GetWebConn", qq, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("emit", result)
	}
}

/*

	** set DEBUG=false at config.go in the product environment. **

	Cross Compile
	SET CGO_ENABLED=0
	SET GOOS=linux
	SET GOARCH=amd64
	go build main.go

*/
func main() {
	var site string
	var port int
	port = iotqq.ApiPort

	// *******  CREATE THE DATABASE TABLE ONLY AT THE FIRST TIME   ******
	util.CreateTable()

	fmt.Println("[+] Revoke Bot Online.")
	site = iotqq.ApiIp
	qq = strconv.Itoa(int(iotqq.LoginQQ))
	runtime.GOMAXPROCS(runtime.NumCPU())
	url1 = site + ":" + strconv.Itoa(port)
	iotqq.Set(url1, qq)

	c, err := gosocketio.Dial(
		gosocketio.GetUrl(site, port, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}
	err = c.On("OnEvents", func(h *gosocketio.Channel, msg iotqq.Message) {
		fmt.Println(msg)
		handler.EventMsgHandler(msg)
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected Successfully")
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On("OnGroupMsgs", func(h *gosocketio.Channel, msg iotqq.Message) {
		fmt.Println(msg)
		handler.GroupMsgHandler(msg)
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On("OnFriendMsgs", func(h *gosocketio.Channel, msg iotqq.Message) {
		fmt.Println(msg)
		if int64(msg.CurrentPacket.Data.FromUin) != iotqq.Master {
			iotqq.Send(int(msg.CurrentPacket.Data.FromUserID), 1, "Auth Fail")
			fmt.Println(msg.CurrentPacket.Data.FromUin)
			return
		}
		handler.FriendMsgHandler(msg)
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	go SendJoin(c)
	periodCall(24*time.Hour, resetPause)
home:
	time.Sleep(600 * time.Second)
	SendJoin(c)
	goto home
}

func init() {
	file, err := os.Open(iotqq.ConfigName)
	conf = iotqq.Conf{true, make(map[string]int), []int{}, []int{}}
	if err != nil {
		log.Println(err)
		_, _ = os.Create(iotqq.ConfigName)
		f, _ := os.OpenFile(iotqq.ConfigName, os.O_WRONLY, 0666)
		defer f.Close()
		enc := json.NewEncoder(f)
		conf.Enable = true
		conf.GData = make(map[string]int)
		tmp := []int{}
		for _, i := range iotqq.RevokeGroupList.List() {
			t := reflect.ValueOf(i).Int()
			tmp = append(tmp, *(*int)(unsafe.Pointer(&t)))
		}
		conf.RevokeList = tmp
		tmpGrant := []int{}
		for _, i := range iotqq.GrantList.List() {
			t := reflect.ValueOf(i).Int()
			tmp = append(tmp, *(*int)(unsafe.Pointer(&t)))
		}
		conf.GrantList = tmpGrant
		_ = enc.Encode(conf)
	}
	defer file.Close()
	tmp := json.NewDecoder(file)
	//log.Println(tmp)
	for tmp.More() {
		err := tmp.Decode(&conf)
		if err != nil {
			fmt.Println("Error:", err)
		}
		tmp := conf.RevokeList
		for _, i := range tmp {
			iotqq.RevokeGroupList.Add(i)
		}
		tmp = conf.GrantList
		for _, i := range tmp {
			iotqq.GrantList.Add(i)
		}
		//fmt.Println(conf)
	}
	fmt.Println("Init config successfully.")
}
