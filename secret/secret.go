package secret

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	// "github.com/molin0000/secretMaster/rlp"
)

type Person struct {
	Group       uint64
	QQ          uint64
	Name        string
	JoinTime    uint64
	LastChat    uint64
	SecretID    uint64
	SecretLevel uint64
	ChatCount   uint64
}

type SecretInfo struct {
	SecretName      string
	SecretLevelName [9]string
}

type Bot struct {
	QQ    int64
	Group int64
	Name  string
}

var botMenu = [...]string{
	"帮助",
	"属性",
	"途径",
	"更换途径",
	"查询",
	"排行",
}

var secretInfo = [...]SecretInfo{
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "aa", SecretLevelName: [9]string{"a", "b", "c"}},
}

func NewSecretBot(qq, group int64, groupNick string) *Bot {
	bot := &Bot{QQ: qq, Group: group, Name: groupNick}
	return bot
}

func (b *Bot) Run(msg string, fromQQ int64) string {
	b.update(fromQQ)

	if !b.talkToMe(msg) {
		return ""
	}
	return b.cmdSwitch(msg, fromQQ)
}

func (b *Bot) update(fromQQ int64) {

}

func (b *Bot) talkToMe(msg string) bool {
	if len(msg) == 0 {
		return false
	}

	if strings.Index(msg, "@"+b.Name) != -1 {
		return true
	}

	return false
}

func (b *Bot) cmdSwitch(msg string, fromQQ int64) string {
	if strings.Contains(msg, botMenu[0]) {
		return `
帮助：回复 帮助 可显示帮助信息。
属性：回复 属性 可查询当前人物的属性信息。
途径：回复 途径 可查询全部途径列表。
更换途径：回复 更换途径+途径序号 可更改当前人物的非凡途径。
查询：输入 查询 + QQ号码 可查询指定人员的属性信息。
排行：输入 排行 可查询当前群内的非凡者排行榜。
`
	}

	if strings.Contains(msg, botMenu[1]) {
		return b.getProperty(fromQQ)
	}

	if strings.Contains(msg, botMenu[2]) {
		return b.getSecretList()
	}

	if strings.Contains(msg, botMenu[3]) {
		return b.changeSecretList(msg, fromQQ)
	}

	if strings.Contains(msg, botMenu[4]) {
		return b.checkQQ(msg, fromQQ)
	}

	if strings.Contains(msg, botMenu[5]) {
		return b.getProperty(fromQQ)
	}

	return ""
}

func (b *Bot) getProperty(fromQQ int64) string {
	return "empty"
}

func (b *Bot) getSecretList() string {
	return "empty"
}

func (b *Bot) changeSecretList(msg string, fromQQ int64) string {
	return "empty"
}

func (b *Bot) checkQQ(msg string, fromQQ int64) string {
	return "empty"
}

func (b *Bot) getRank() string {
	return "empty"
}

var db *leveldb.DB

func getDb() *leveldb.DB {
	if db == nil {
		_db, err := leveldb.OpenFile("secret.db", nil)
		if err != nil {
			fmt.Printf("open db error: %+v", err)
			panic(err)
		}
		db = _db
	}
	return db
}

func (b *Bot) keys(fromQQ int64) []byte {
	return []byte(strconv.FormatInt(b.Group, 10) + "_" + strconv.FormatInt(fromQQ, 10))
}
