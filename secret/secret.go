package secret

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/molin0000/secretMaster/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

type Person struct {
	Group       uint64
	QQ          uint64
	Name        string
	JoinTime    uint64
	LastChat    uint64
	LevelDown   uint64
	SecretID    uint64
	SecretLevel uint64
	ChatCount   uint64
}

type SecretInfo struct {
	SecretName      string
	SecretLevelName [9]string
}

type Bot struct {
	QQ    uint64
	Group uint64
	Name  string
}

var botMenu = [...]string{
	"帮助",
	"属性",
	"途径",
	"更换",
	"查询",
	"排行",
}

var secretInfo = [...]SecretInfo{
	{SecretName: "黑夜", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "死神", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "战神", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "暗", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "心灵", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "风暴", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "智慧", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "太阳", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "异种", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "深渊", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "秘密", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "工匠", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "愚者", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "门", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "时间", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "大地", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "月亮", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "命运", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "审判", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "黑皇帝", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "战争", SecretLevelName: [9]string{"a", "b", "c"}},
	{SecretName: "魔女", SecretLevelName: [9]string{"a", "b", "c"}},
}

func NewSecretBot(qq, group uint64, groupNick string) *Bot {
	bot := &Bot{QQ: qq, Group: group, Name: groupNick}
	return bot
}

func (b *Bot) Run(msg string, fromQQ uint64, nick string) string {
	b.update(fromQQ, nick)
	if !b.talkToMe(msg) {
		return ""
	}
	return b.cmdSwitch(msg, fromQQ)
}

func (b *Bot) update(fromQQ uint64, nick string) {
	key := b.keys(fromQQ)
	value, err := getDb().Get(key, nil)
	fmt.Println("value:", value)
	fmt.Println("err:", err)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			fmt.Println("a new man.")
			p := &Person{
				Group:       b.Group,
				QQ:          fromQQ,
				Name:        nick,
				JoinTime:    uint64(time.Now().Unix()),
				LastChat:    uint64(time.Now().Unix()),
				LevelDown:   uint64(time.Now().Unix()),
				SecretID:    99,
				SecretLevel: 99,
				ChatCount:   1,
			}

			buf, _ := rlp.EncodeToBytes(p)
			getDb().Put(b.keys(fromQQ), buf, nil)
			verify, _ := getDb().Get(b.keys(fromQQ), nil)
			var v Person
			rlp.DecodeBytes(verify, &v)
			fmt.Printf("%+v", v)
		}
	} else {
		verify, _ := getDb().Get(b.keys(fromQQ), nil)
		var v Person
		rlp.DecodeBytes(verify, &v)
		fmt.Printf("before level:%+v", v)
		b.levelUpdate(&v)
		fmt.Printf("after level:%+v", v)
		v.ChatCount++
		v.LastChat = uint64(time.Now().Unix())
		v.LevelDown = uint64(time.Now().Unix())
		buf, _ := rlp.EncodeToBytes(v)
		getDb().Put(b.keys(fromQQ), buf, nil)
		fmt.Println("update finish")
	}
}

func (b *Bot) levelUpdate(p *Person) {
	if p.SecretID > 22 {
		return
	}

	if p.ChatCount > 10000 {
		p.SecretLevel = 9
	} else if p.ChatCount > 8000 {
		p.SecretLevel = 8
	} else if p.ChatCount > 5500 {
		p.SecretLevel = 7
	} else if p.ChatCount > 3000 {
		p.SecretLevel = 6
	} else if p.ChatCount > 1500 {
		p.SecretLevel = 5
	} else if p.ChatCount > 1000 {
		p.SecretLevel = 4
	} else if p.ChatCount > 800 {
		p.SecretLevel = 3
	} else if p.ChatCount > 400 {
		p.SecretLevel = 2
	} else if p.ChatCount > 200 {
		p.SecretLevel = 1
	} else if p.ChatCount > 50 {
		p.SecretLevel = 0
	} else {
		p.SecretLevel = 99
	}

	fmt.Println("level:", p.SecretLevel)

	if (uint64(time.Now().Unix()) - p.LevelDown) > 3600*24*7 {
		if p.SecretLevel >= (uint64(time.Now().Unix())-p.LevelDown)/(3600*24*7) {
			p.SecretLevel -= (uint64(time.Now().Unix()) - p.LevelDown) / (3600 * 24 * 7)
		}
		p.LevelDown = uint64(time.Now().Unix())
	}
}

func (b *Bot) talkToMe(msg string) bool {
	if len(msg) == 0 {
		return false
	}

	cp := fmt.Sprintf("CQ:at,qq=%d", b.QQ)

	fmt.Println("talkToMe?", msg, cp)

	if strings.Index(msg, cp) != -1 {
		return true
	}

	return false
}

func (b *Bot) cmdSwitch(msg string, fromQQ uint64) string {
	if strings.Contains(msg, botMenu[0]) {
		return `
帮助：回复 帮助 可显示帮助信息。
属性：回复 属性 可查询当前人物的属性信息。
途径：回复 途径 可查询途径列表。
更换：回复 更换+途径序号 可更改当前人物的非凡途径。
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

func (b *Bot) getProperty(fromQQ uint64) string {
	verify, _ := getDb().Get(b.keys(fromQQ), nil)
	var v Person
	rlp.DecodeBytes(verify, &v)
	fmt.Printf("%+v", v)
	var secretName string
	var secretLevelName string
	var startTime string
	var lastTime string

	if v.SecretID > 22 {
		secretName = "普通人"
	} else {
		secretName = secretInfo[v.SecretID].SecretName
	}

	if v.SecretLevel > 10 {
		secretLevelName = "普通人"
	} else {
		// secretLevelName = fmt.Sprintf("序列%d：%s", 9-v.SecretLevel, secretInfo[v.SecretID].SecretLevelName[v.SecretLevel])
		secretLevelName = fmt.Sprintf("序列%d", 9-v.SecretLevel)
	}

	startTime = time.Unix(int64(v.JoinTime), 0).Format("2006-01-02 15:04:05")
	lastTime = time.Unix(int64(v.LastChat), 0).Format("2006-01-02 15:04:05")

	info := fmt.Sprintf("\n尊名：%s\n途径：%s\n序列：%s\n经验：%d\n修炼开始时间：%s\n最近修炼时间：%s\n技能：%s\n",
		v.Name, secretName, secretLevelName, v.ChatCount, startTime, lastTime, "无")

	fmt.Print(info)
	return info
}

func (b *Bot) getSecretList() string {
	secretList := "途径列表：\n"

	for i := 0; i < len(secretInfo); i++ {
		secretList += fmt.Sprintf("%d: %s\n", i+1, secretInfo[i].SecretName)
	}

	return secretList
}

func (b *Bot) changeSecretList(msgRaw string, fromQQ uint64) string {
	index := strings.Index(msgRaw, "更换")
	msg := msgRaw[index:]
	fmt.Println("changeSecretList:", msg, msgRaw)
	defer fmt.Println("finish changed")

	var valid = regexp.MustCompile("[0-9]")
	r := valid.FindAllStringSubmatch(msg, -1)
	if len(r) == 0 || len(r) > 2 {
		return "数值无效"
	}

	var value int
	if len(r) == 1 {
		value, _ = strconv.Atoi(r[0][0])
	} else {
		a1, _ := strconv.Atoi(r[0][0])
		a2, _ := strconv.Atoi(r[1][0])
		value = a1*10 + a2
	}

	if value < 1 || value > 22 {
		return "数值无效"
	}

	verify, _ := getDb().Get(b.keys(fromQQ), nil)
	var v Person
	rlp.DecodeBytes(verify, &v)
	v.SecretID = uint64(value - 1)
	buf, _ := rlp.EncodeToBytes(v)
	fmt.Printf("%+v", v)

	getDb().Put(b.keys(fromQQ), buf, nil)

	return fmt.Sprintf("成功更换到途径：%d", value)
}

func (b *Bot) checkQQ(msg string, fromQQ uint64) string {
	return "暂不支持"
}

func (b *Bot) getRank() string {
	return "来不及做了，等下一版吧"
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

func (b *Bot) keys(fromQQ uint64) []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10))
}
