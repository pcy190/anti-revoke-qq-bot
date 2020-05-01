package util

import (
	"RevokeBot/model"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"reflect"
	"unsafe"
)

var (
	userName  string = iotqq.MysqlUsername
	password  string = iotqq.MysqlPassword
	ipAddrees string = iotqq.MysqlIP
	port      int    = iotqq.MysqlPort
	dbName    string = iotqq.MysqlDbName
	charset   string = "utf8"
)
var Db *sqlx.DB

func connectMysql() (*sqlx.DB) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddrees, port, dbName, charset)
	Db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("mysql connect failed, detail is [%v]", err.Error())
	}
	return Db
}

func AddRecord(GroupID int, MsgSeq int, MsgRandom int, MsgType string, Content string, NickName string, MsgTime int) {
	result, err := Db.Exec("insert into msg_cache(GroupID, MsgSeq, MsgRandom,MsgType,Content,NickName,MsgTime)  values(?,?,?,?,?,?,?)", GroupID, MsgSeq, MsgRandom, MsgType, Content, NickName, MsgTime)
	if err != nil {
		fmt.Printf("data insert faied, error:[%v]", err.Error())
		return
	}
	id, _ := result.LastInsertId()
	fmt.Printf("insert success, last id:[%d]\n", id)
}

func QueryRecord(GroupID int, MsgSeq int) (iotqq.RevokedMessage) {
	result, err := Db.Query("select * from msg_cache where `GroupID`= ? and `MsgSeq` = ?", GroupID, MsgSeq)

	if err != nil {
		fmt.Printf("data query faied, error:[%v]", err.Error())
		return iotqq.RevokedMessage{}
	}
	for result.Next() {
		var GroupID int
		var MsgSeq int
		var MsgRandom int
		var MsgType string
		var Content string
		var NickName string
		var MsgTime int
		err = result.Scan(&GroupID, &MsgSeq, &MsgRandom, &MsgType, &Content, &NickName, &MsgTime)
		if err != nil {
			fmt.Printf("data query success but scan failed, error:[%v]", err.Error())
			return iotqq.RevokedMessage{}
		}
		return iotqq.RevokedMessage{GroupID: GroupID, MsgSeq: MsgSeq, MsgRandom: MsgRandom, MsgType: MsgType, Content: Content, NickName: NickName, MsgTime: MsgTime}
	}
	return iotqq.RevokedMessage{}
}

func updateRecord() {
	result, err := Db.Exec("update userinfo set username = 'anson' where uid = 1")
	if err != nil {
		fmt.Printf("update faied, error:[%v]", err.Error())
		return
	}
	num, _ := result.RowsAffected()
	fmt.Printf("update success, affected rows:[%d]\n", num)
}

func deleteRecord() {
	result, err := Db.Exec("delete from userinfo where uid = 2")
	if err != nil {
		fmt.Printf("delete faied, error:[%v]", err.Error())
		return
	}
	num, _ := result.RowsAffected()
	fmt.Printf("delete success, affected rows:[%d]\n", num)
}

func CreateTable() {
	_, err := Db.Exec("CREATE TABLE IF NOT EXISTS `qqmsg`.`msg_cache` ( \n" +
		"`GroupID` BIGINT NOT NULL,\n" +
		"`MsgSeq` BIGINT NOT NULL,\n" +
		"`MsgRandom` BIGINT NOT NULL,\n" +
		"`MsgType` VARCHAR(40),\n" +
		"`Content` TEXT, \n" +
		"`NickName` VARCHAR(100), \n" +
		"`MsgTime` BIGINT ) \n" +
		" ENGINE = InnoDB \n" +
		"DEFAULT CHARACTER SET = utf8;")
	if err != nil {
		fmt.Printf("delete faied, error:[%v]", err.Error())
		return
	} else {
		fmt.Println("Successfully create table")
		return
	}

}

func init() {
	Db = connectMysql()
	//defer Db.Close()
}

func WriteConfig() {
	file, err := os.OpenFile(iotqq.ConfigName, os.O_WRONLY|os.O_TRUNC, 666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	tmp := []int{}
	for _, i := range iotqq.RevokeGroupList.List() {
		t := reflect.ValueOf(i).Int()
		tmp = append(tmp, *(*int)(unsafe.Pointer(&t)))
	}
	tmpGrant := []int{}
	for _, i := range iotqq.GrantList.List() {
		t := reflect.ValueOf(i).Int()
		tmpGrant = append(tmpGrant, *(*int)(unsafe.Pointer(&t)))
	}
	conf := iotqq.Conf{Enable: true, GData: make(map[string]int), RevokeList: tmp, GrantList: tmpGrant}
	_ = enc.Encode(conf)
	fmt.Println("Write config successfully.")
}

func ReadFromConfig() bool {
	normal := true
	file, err := os.Open(iotqq.ConfigName)
	conf := iotqq.Conf{true, make(map[string]int), []int{}, []int{}}
	if err != nil {
		log.Println(err)
		_, _ = os.Create(iotqq.ConfigName)
		f, _ := os.OpenFile(iotqq.ConfigName, os.O_WRONLY, 0666)
		defer f.Close()
		enc := json.NewEncoder(f)
		conf.Enable = true
		conf.GData = make(map[string]int)
		var tmp []int
		for _, i := range iotqq.RevokeGroupList.List() {
			t := reflect.ValueOf(i).Int()
			tmp = append(tmp, *(*int)(unsafe.Pointer(&t)))
		}
		conf.RevokeList = tmp
		var tmpGrant []int
		for _, i := range iotqq.GrantList.List() {
			t := reflect.ValueOf(i).Int()
			tmp = append(tmp, *(*int)(unsafe.Pointer(&t)))
		}
		conf.GrantList = tmpGrant
		_ = enc.Encode(conf)
		normal = false
	}
	defer file.Close()
	tmp := json.NewDecoder(file)
	for tmp.More() {
		err := tmp.Decode(&conf)
		if err != nil {
			fmt.Println("Error:", err)
			normal = false
		}
		tmp := conf.RevokeList
		for _, i := range tmp {
			iotqq.RevokeGroupList.Add(i)
		}
		tmp = conf.GrantList
		for _, i := range tmp {
			iotqq.GrantList.Add(i)
			fmt.Println(i)
		}
	}
	return normal
}

// run it for the first time to create the database table
func main() {
	//sqlx.DB = connectMysql()
	//defer Db.Close()
	//createTable()
	//addRecord(Db)
	//updateRecord(Db)
	//deleteRecord(Db)
}
