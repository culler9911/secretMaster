package secret

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/molin0000/secretMaster/rlp"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func contains(s []uint64, e uint64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func canConvert(v1, v2 uint64) bool {
	for _, a := range secretGroup {
		if contains(a, v1) && contains(a, v2) {
			return true
		}
	}
	return false
}

var botMap map[uint64]*Bot

func NewSecretBot(qq, group uint64, groupNick string) *Bot {
	if botMap == nil {
		botMap = make(map[uint64]*Bot)
	}
	bot, ok := botMap[group]
	if ok {
		return bot
	}
	bot = &Bot{QQ: qq, Group: group, Name: groupNick}
	botMap[group] = bot
	return bot
}

func (b *Bot) Run(msg string, fromQQ uint64, nick string) string {
	if !b.talkToMe(msg) {
		return ""
	}

	return b.cmdSwitch(msg, fromQQ)
}

func (b *Bot) cmdSwitch(msg string, fromQQ uint64) string {
	if strings.Contains(msg, "帮助") {
		return `
帮助：回复 帮助 可显示帮助信息。
属性：回复 属性 可查询当前人物的属性信息。
途径：回复 途径 可查询途径列表。
更换：回复 更换+途径序号 可更改当前人物的非凡途径。
排行：回复 排行 可查询当前群内的非凡者排行榜。
探险：回复 探险 可主动触发每日3次的奇遇探险经历。
删除人物：删除当前全部经验和属性，重新创建人物。
购买探险卷轴：花费100金镑购买1次探险机会。
尊名：序列3以后可以自定义尊名显示，方法为@Yami尊名xxxxoooo。
其余详细介绍请见：https://github.com/molin0000/secretMaster/blob/master/README.md`
	}

	if strings.Contains(msg, "购买探险卷轴") {
		return b.adventure(fromQQ, false)
	}

	if strings.Contains(msg, "属性") {
		return b.getProperty(fromQQ)
	}

	if strings.Contains(msg, "途径") {
		return b.getSecretList()
	}

	if strings.Contains(msg, "更换") {
		return b.changeSecretList(msg, fromQQ)
	}

	if strings.Contains(msg, "排行") {
		return b.getRank(fromQQ)
	}

	if strings.Contains(msg, "探险") {
		return b.adventure(fromQQ, true)
	}

	if strings.Contains(msg, "删除人物") {
		return b.deletePerson(fromQQ)
	}

	if strings.Contains(msg, "尊名") {
		return b.setRespectName(msg, fromQQ)
	}

	return ""
}

func (b *Bot) setRespectName(msg string, fromQQ uint64) string {
	p := b.getPersonFromDb(fromQQ)

	if p.SecretLevel < 6 || p.SecretLevel > 10 {
		return "只有序列3以上才可拥有尊名，因为随便的尊名可能会引起高位存在的注意"
	}

	rname := msg[strings.Index(msg, "尊名")+len("尊名"):]

	b.setRNameToDb(fromQQ, &RespectName{Group: b.Group, QQ: fromQQ, Name: rname})

	return "成功设置尊名"
}

func (b *Bot) deletePerson(fromQQ uint64) string {
	p := b.getPersonFromDb(fromQQ)

	if b.getGodFromDb(p.SecretID) == p.QQ {
		getDb().Delete(b.godKey(p.SecretID), nil)
	}

	getDb().Delete(b.keys(fromQQ), nil)
	getDb().Delete(b.ruleKey(fromQQ), nil)
	getDb().Delete(b.moneyKey(fromQQ), nil)
	getDb().Delete(b.advKey(fromQQ), nil)

	return "人物删除成功"
}

func (b *Bot) Update(fromQQ uint64, nick string) string {
	key := b.keys(fromQQ)
	value, err := getDb().Get(key, nil)
	fmt.Println("value:", value)
	fmt.Println("err:", err)
	ret := ""
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

			b.setPersonToDb(fromQQ, p)

			e := b.getExternFromDb(fromQQ)
			b.setExternToDb(fromQQ, e)
		}
	} else {
		v := b.getPersonFromDb(fromQQ)

		ret = b.levelUpdate(v)

		w := b.getWaterRuleFromDb(fromQQ)
		m := b.getMoneyFromDb(fromQQ, v.ChatCount)
		e := b.getExternFromDb(fromQQ)

		if w.DayCnt < 200 {
			v.ChatCount++
			w.DayCnt++
			m.Money++
			e.Magic--
			if v.ChatCount%100 == 0 {
				ret += "\n恭喜！你的战力评价升级了！"
			}
		}

		v.LastChat = uint64(time.Now().Unix())
		v.LevelDown = uint64(time.Now().Unix())
		b.setPersonToDb(fromQQ, v)
		b.setWaterRuleToDb(fromQQ, w)
		b.setMoneyToDb(fromQQ, m)
	}

	return ret
}

func (b *Bot) adventure(fromQQ uint64, limit bool) string {
	a := b.getAdvFromDb(fromQQ)
	if limit && a.DayCnt >= 3 {
		return "对不起，您今日奇遇探险机会已经用完"
	}

	a.DayCnt++
	a.Days = uint64(time.Now().Unix() / (3600 * 24))
	b.setAdvToDb(fromQQ, a)

	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(20)
	info := ""
	m := 0
	e := 0
	if i < 1 { // -钱＆经验5%
		t := rand.Intn(2)
		m = -1 * (10 + rand.Intn(41))
		e = -1 * (10 + rand.Intn(41))
		if t == 0 {
			info = fmt.Sprintf("在探险过程中遇到危险不得不使用保命物品才得以逃脱，损失惨重。经验:%d, 金钱:%d", e, m)
		} else {
			info = fmt.Sprintf("被异教徒盯上并暴力劫掠，损失惨重。经验:%d, 金钱:%d", e, m)
		}
	} else if i < 3 { // -经验10%
		t := rand.Intn(2)
		m = 0
		e = -1 * (10 + rand.Intn(41))
		if t == 0 {
			info = fmt.Sprintf("由于错误的祷告，被不知名的存在盯上，受到打击。经验:%d", e)
		} else {
			info = fmt.Sprintf("在探险过程中用尽了力量不得不返回，灵性枯竭。经验:%d", e)
		}
	} else if i < 5 { // -钱10%
		t := rand.Intn(2)
		m = -1 * (10 + rand.Intn(41))
		e = 0
		if t == 0 {
			info = fmt.Sprintf("在非凡者聚会上买到假货，该死。金钱:%d", m)
		} else {
			info = fmt.Sprintf("交房租的时间到了，老天。金钱:%d", m)
		}
	} else if i < 7 { // +钱＆经验10%
		t := rand.Intn(2)
		m = (20 + rand.Intn(81))
		e = (20 + rand.Intn(81))
		if t == 0 {
			info = fmt.Sprintf("探索了一座遗迹，了解到一些隐秘的知识并发现了一些财宝。经验:%d, 金钱:%d", e, m)
		} else {
			info = fmt.Sprintf("成功打击了一名异教徒，获得了对方的一些物品。经验:%d, 金钱:%d", e, m)
		}
	} else if i < 10 { // +钱15%
		t := rand.Intn(2)
		m = (20 + rand.Intn(81))
		e = 0
		if t == 0 {
			info = fmt.Sprintf("完成一次非凡任务委托，获得报酬。金钱:%d", m)
		} else {
			info = fmt.Sprintf("猎杀了一只凶猛的非凡生物，获得材料。金钱:%d", m)
		}
	} else if i < 13 { // 	+经验15%
		t := rand.Intn(2)
		m = 0
		e = (20 + rand.Intn(81))
		if t == 0 {
			info = fmt.Sprintf("在一出失落的古堡中发现一本记载了神秘学知识的书籍，获得知识。经验:%d", e)
		} else {
			info = fmt.Sprintf("参加非凡者组织的聚会，有幸的到了高阶阅读者的指点，获得知识。经验:%d", e)
		}
	} else { // 无事发生35%
		t := rand.Intn(2)
		m = 0
		e = 0
		if t == 0 {
			info = "美好的一天，何不来点东拜朗甜姜红茶？"
		} else {
			info = "糟糕的一天，并不想出门！"
		}
	}

	person := b.getPersonFromDb(fromQQ)
	_m := b.getMoneyFromDb(fromQQ, person.ChatCount)

	if !limit {
		if _m.Money > 100 {
			_m.Money -= 100
		} else {
			return "钱包空空，买不起了哦"
		}
	}

	if m >= 0 {
		_m.Money += uint64(m)
	} else {
		if _m.Money > uint64(-1*m) {
			_m.Money -= uint64(-1 * m)
		} else {
			_m.Money = 0
		}
	}

	if e >= 0 {
		person.ChatCount += uint64(e)
	} else {
		if person.ChatCount > uint64(-1*e) {
			person.ChatCount -= uint64(-1 * e)
		} else {
			person.ChatCount = 0
		}
	}

	b.setPersonToDb(fromQQ, person)
	b.setMoneyToDb(fromQQ, _m)

	return info
}

func (b *Bot) levelUpdate(p *Person) string {
	if p.SecretID > 22 {
		return ""
	}

	ret := ""
	levelOld := p.SecretLevel
	god := b.getGodFromDb(p.SecretID)
	if p.ChatCount > 10000 {
		if god == 0 || god == p.QQ {
			b.setGodToDb(p.SecretID, &p.QQ)
			p.SecretLevel = 9
		}
	} else if p.ChatCount > 8000 {
		if god == 0 {
			p.SecretLevel = 8
		}
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

	if p.SecretLevel != levelOld {
		ret = "恭喜！你的序列晋升了！"
	}

	fmt.Println("level:", p.SecretLevel)

	if (uint64(time.Now().Unix()) - p.LevelDown) > 3600*24*7 {
		if p.SecretLevel >= (uint64(time.Now().Unix())-p.LevelDown)/(3600*24*7) {
			p.SecretLevel -= (uint64(time.Now().Unix()) - p.LevelDown) / (3600 * 24 * 7)
		}
		p.LevelDown = uint64(time.Now().Unix())
	}

	return ret
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

func (b *Bot) getProperty(fromQQ uint64) string {
	v := b.getPersonFromDb(fromQQ)
	var secretName string
	var secretLevelName string
	var startTime string

	if v.SecretID > 22 {
		secretName = "普通人"
	} else {
		secretName = secretInfo[v.SecretID].SecretName
	}

	if v.SecretLevel > 10 {
		secretLevelName = "普通人"
	} else {
		secretLevelName = fmt.Sprintf("序列%d：%s", 9-v.SecretLevel, secretInfo[v.SecretID].SecretLevelName[v.SecretLevel])
		// secretLevelName = fmt.Sprintf("序列%d", 9-v.SecretLevel)
	}

	startTime = fmt.Sprintf("%d小时", (time.Now().Unix()-time.Unix(int64(v.JoinTime), 0).Unix())/3600)
	// startTime = time.Unix(int64(v.JoinTime), 0).Format("2006-01-02 15:04:05")

	myFightIndex := v.ChatCount / 100
	reLive := uint64(0)
	sReLive := ""
	if myFightIndex > 99 {
		myFightIndex = v.ChatCount / 100 % 100
		reLive = v.ChatCount / uint64(10000)
		if reLive > 0 {
			sReLive = fmt.Sprintf("(转生+%d)", reLive)
		}
	}

	money := b.getMoneyFromDb(fromQQ, v.ChatCount)
	if money.Money > 2 {
		money.Money -= 2
	}
	b.setMoneyToDb(fromQQ, money)

	info := fmt.Sprintf("\n昵称：%s\n途径：%s\n序列：%s\n经验：%d\n金镑：%d\n修炼时间：%s\n战力评价：%s%s\n尊名：%s",
		v.Name, secretName, secretLevelName, v.ChatCount, money.Money,
		startTime, fight[myFightIndex], sReLive,
		b.getRNameFromDb(fromQQ),
	)

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

	v := b.getPersonFromDb(fromQQ)
	if v.SecretLevel < 5 {
		return "很抱歉，在您到达序列4之前，无法再次转换途径"
	}

	if v.SecretID != 99 && !canConvert(v.SecretID, uint64(value-1)) {
		return "很抱歉，非凡者只能在相近途径互相转换"
	}

	v.SecretID = uint64(value - 1)
	b.setPersonToDb(fromQQ, v)

	money := b.getMoneyFromDb(fromQQ, v.ChatCount)
	if money.Money > 100 {
		money.Money -= 100
	}
	b.setMoneyToDb(fromQQ, money)

	return fmt.Sprintf("成功更换到途径：%d", value)
}

func (b *Bot) checkQQ(msg string, fromQQ uint64) string {
	return "暂不支持"
}

func (b *Bot) getRank(fromQQ uint64) string {
	iter := getDb().NewIterator(util.BytesPrefix(b.getKeyPrefix()), nil)
	persons := make([]Person, 0)
	cnt := 0
	for iter.Next() {
		fmt.Printf("key:%+v", string(iter.Key()))
		fmt.Printf("value:%+v", iter.Value())

		if strings.Index(string(iter.Key()), "money") != -1 ||
			strings.Index(string(iter.Key()), "adv") != -1 ||
			strings.Index(string(iter.Key()), "water") != -1 {
			continue
		}

		verify := iter.Value()
		var v Person
		rlp.DecodeBytes(verify, &v)
		fmt.Printf("%+v", v)
		persons = append(persons, v)
	}
	iter.Release()
	err := iter.Error()
	fmt.Println(err)

	retValue := ""

	sort.Sort(Persons(persons))

	for i := 0; i < len(persons); i++ {
		v := persons[i]
		retValue = fmt.Sprintf("%s\n第%d名：%s，经验：%d", retValue, i+1, v.Name, v.ChatCount)
		cnt++
		if cnt > 30 {
			break
		}
	}

	v := b.getPersonFromDb(fromQQ)
	money := b.getMoneyFromDb(fromQQ, v.ChatCount)
	if money.Money > 2 {
		money.Money -= 2
	}
	b.setMoneyToDb(fromQQ, money)

	return retValue
}

var eventChan chan string

func TickerInit() {
	fmt.Println("Ticker init")
	eventChan = make(chan string, 100)
	t := time.NewTicker(time.Minute)
	ticker := func() {
		fmt.Println("Ticker reached", time.Now())
		timeNow := time.Now()
		hour, min, _ := timeNow.Clock()
		fmt.Println(hour, min)
		if hour == 8 && min == 30 {
			eventChan <- "GoodMorning"
		}

		if hour == 0 && min == 16 {
			eventChan <- "GoodEvening"
		}
	}
	go func() {
		for {
			select {
			case <-t.C:
				ticker()
			}
		}
	}()
}

func TickerProc() string {
	select {
	case e := <-eventChan:
		fmt.Println(e)
		if e == "GoodMorning" {
			return "宝贝们早上好呀~已经8点30分了哦~早安~喵"
		}

		if e == "GoodEvening" {
			return "宝贝们很晚啦~已经23点整了哦~早点休息呀~晚安~喵"
		}
	}
	return ""
}
