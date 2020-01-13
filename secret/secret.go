package secret

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/molin0000/secretMaster/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var debug = false

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

type Persons []Person

func (a Persons) Len() int           { return len(a) }
func (a Persons) Less(i, j int) bool { return a[i].ChatCount > a[j].ChatCount }
func (a Persons) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type SecretInfo struct {
	SecretName      string
	SecretLevelName [10]string
}

type Bot struct {
	QQ    uint64
	Group uint64
	Name  string
}

type WaterRule struct {
	Group  uint64
	QQ     uint64
	DayCnt uint64
	Days   uint64
}

type Money struct {
	Group uint64
	QQ    uint64
	Money uint64
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
	{SecretName: "黑夜", SecretLevelName: [10]string{"不眠者", "午夜诗人", "梦魇", "安魂师", "灵巫", "守夜人", "恐惧主教", "隐秘之仆", "未知", "黑夜"}},
	{SecretName: "死神", SecretLevelName: [10]string{"收尸人", "掘墓人", "通灵者", "死灵导师", "看门人", "不死者", "摆渡人", "死亡执政官", "未知", "死神"}},
	{SecretName: "战神", SecretLevelName: [10]string{"战士", "格斗家", "武器大师", "黎明骑士", "守护者", "猎魔人", "银骑士", "未知", "未知", "战神"}},
	{SecretName: "暗", SecretLevelName: [10]string{"秘祈人", "倾听者", "隐修士", "蔷薇主教", "牧羊人", "黑骑士", "三首圣堂", "未知", "未知", "未知"}},
	{SecretName: "心灵", SecretLevelName: [10]string{"观众", "读心者", "心理医生", "催眠师", "梦境行者", "操纵师", "织梦者", "洞察者", "作家", "空想家"}},
	{SecretName: "风暴", SecretLevelName: [10]string{"水手", "暴怒之民", "航海家", "风眷者", "海洋歌者", "灾难主祭", "海王", "天灾", "雷神", "暴君"}},
	{SecretName: "智慧", SecretLevelName: [10]string{"阅读者", "未知", "侦探", "博学者", "秘术导师", "预言家", "未知", "未知", "未知", "未知"}},
	{SecretName: "太阳", SecretLevelName: [10]string{"歌颂者", "祈光人", "太阳神官", "公证人", "光之祭司", "无暗者", "正义导师", "逐光者", "未知", "太阳"}},
	{SecretName: "异种", SecretLevelName: [10]string{"囚犯", "疯子", "狼人", "活尸", "怨灵", "木偶", "沉默门徒", "古代邪物", "神孽", "未知"}},
	{SecretName: "深渊", SecretLevelName: [10]string{"罪犯", "冷血者", "连环杀手", "恶魔", "欲望使徒", "魔鬼", "呓语者", "未知", "未知", "深渊"}},
	{SecretName: "秘密", SecretLevelName: [10]string{"窥秘人", "格斗教授", "巫师", "卷轴教授", "星象师", "神秘学者", "预言大师", "贤者", "知识皇帝", "未知"}},
	{SecretName: "工匠", SecretLevelName: [10]string{"通识者", "考古学家", "鉴定师", "机械专家", "天文学家", "炼金术士", "奥秘学者", "未知", "未知", "未知"}},
	{SecretName: "愚者", SecretLevelName: [10]string{"占卜家", "小丑", "魔法师", "无面人", "秘偶大师", "诡法师", "古代学者", "奇迹师", "诡秘侍者", "愚者"}},
	{SecretName: "门", SecretLevelName: [10]string{"学徒", "戏法大师", "占星人", "记录官", "旅行家", "秘法师", "漫游者", "旅法师", "星之匙", "未知"}},
	{SecretName: "时间", SecretLevelName: [10]string{"偷盗者", "诈骗师", "解密学者", "盗火人", "窃梦家", "寄生者", "欺瞒导师", "命运木马", "未知", "错误"}},
	{SecretName: "大地", SecretLevelName: [10]string{"耕种者", "医师", "丰收祭司", "生物学家", "德鲁伊", "古代炼金师", "未知", "未知", "未知", "未知"}},
	{SecretName: "月亮", SecretLevelName: [10]string{"药师", "驯兽师", "吸血鬼", "魔药教授", "深红学者", "巫王", "未知", "未知", "未知", "未知"}},
	{SecretName: "命运", SecretLevelName: [10]string{"怪物", "机器", "幸运儿", "灾祸教士", "赢家", "未知", "未知", "先知", "水银之蛇", "命运之轮"}},
	{SecretName: "审判", SecretLevelName: [10]string{"仲裁人", "治安官", "审讯者", "法官", "惩戒骑士", "律令法师", "混乱猎手", "平衡者", "秩序之手", "审判者"}},
	{SecretName: "黑皇帝", SecretLevelName: [10]string{"律师", "野蛮人", "贿赂者", "腐化男爵", "混乱导师", "堕落伯爵", "狂乱法师", "熵之公爵", "弑序亲王", "黑皇帝"}},
	{SecretName: "战争", SecretLevelName: [10]string{"猎人", "挑衅者", "纵火家", "阴谋家", "收割者", "铁血骑士", "战争主教", "天气术士", "征服者", "红祭司"}},
	{SecretName: "魔女", SecretLevelName: [10]string{"刺客", "教唆者", "女巫", "欢愉", "痛苦", "绝望", "不老", "灾难", "末日", "原初"}},
}

func NewSecretBot(qq, group uint64, groupNick string) *Bot {
	bot := &Bot{QQ: qq, Group: group, Name: groupNick}
	return bot
}

func (b *Bot) Run(msg string, fromQQ uint64, nick string) string {
	if len(msg) > 9 {
		fmt.Println(msg, "大于3", len(msg))
		b.update(fromQQ, nick)
	} else {
		fmt.Println(msg, "小于3", len(msg))
	}

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

			b.setPersonToDb(fromQQ, p)
		}
	} else {
		v := b.getPersonFromDb(fromQQ)
		b.levelUpdate(v)

		w := b.getWaterRuleFromDb(fromQQ)

		if debug {
			v.ChatCount += 1000
		} else {
			if w.DayCnt < 200 {
				v.ChatCount++
				w.DayCnt++
			}
		}

		v.LastChat = uint64(time.Now().Unix())
		v.LevelDown = uint64(time.Now().Unix())
		b.setPersonToDb(fromQQ, v)
		b.setWaterRuleToDb(fromQQ, w)
	}
}

func (b *Bot) getPersonFromDb(fromQQ uint64) *Person {
	verify, _ := getDb().Get(b.keys(fromQQ), nil)
	var v Person
	rlp.DecodeBytes(verify, &v)
	return &v
}

func (b *Bot) setPersonToDb(fromQQ uint64, p *Person) {
	buf, _ := rlp.EncodeToBytes(p)
	getDb().Put(b.keys(fromQQ), buf, nil)
}

func (b *Bot) getWaterRuleFromDb(fromQQ uint64) *WaterRule {
	ret, err := getDb().Get(b.ruleKey(fromQQ), nil)
	if err != nil {
		return &WaterRule{Group: b.Group, QQ: fromQQ, DayCnt: 1, Days: uint64(time.Now().Unix() / (3600 * 24))}
	}
	var w WaterRule
	rlp.DecodeBytes(ret, &w)
	if w.Days != uint64(time.Now().Unix()/(3600*24)) {
		w.DayCnt = 0
	}
	return &w
}

func (b *Bot) setWaterRuleToDb(fromQQ uint64, w *WaterRule) {
	buf, _ := rlp.EncodeToBytes(w)
	getDb().Put(b.ruleKey(fromQQ), buf, nil)
}

func (b *Bot) getMoneyFromDb(fromQQ uint64, chatCnt uint64) *Money {
	ret, err := getDb().Get(b.moneyKey(fromQQ), nil)
	if err != nil {
		return &Money{Group: b.Group, QQ: fromQQ, Money: chatCnt}
	}
	var m Money
	rlp.DecodeBytes(ret, &m)
	return &m
}

func (b *Bot) setMoneyToDb(fromQQ uint64, m *Money) {
	buf, _ := rlp.EncodeToBytes(m)
	getDb().Put(b.moneyKey(fromQQ), buf, nil)
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
		return b.getRank()
	}

	return ""
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
	if money.Money > 1 {
		money.Money--
	}
	b.setMoneyToDb(fromQQ, money)

	info := fmt.Sprintf("\n尊名：%s\n途径：%s\n序列：%s\n经验：%d\n金镑：%d\n修炼时间：%s\n战力评价：%s%s\n",
		v.Name, secretName, secretLevelName, v.ChatCount, money.Money, startTime, fight[myFightIndex], sReLive)

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
	v.SecretID = uint64(value - 1)
	b.setPersonToDb(fromQQ, v)

	return fmt.Sprintf("成功更换到途径：%d", value)
}

func (b *Bot) checkQQ(msg string, fromQQ uint64) string {
	return "暂不支持"
}

func (b *Bot) getRank() string {
	iter := getDb().NewIterator(util.BytesPrefix(b.getKeyPrefix()), nil)
	persons := make([]Person, 0)
	cnt := 0
	for iter.Next() {
		fmt.Printf("key:%+v", iter.Key())
		fmt.Printf("value:%+v", iter.Value())
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
	return retValue
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

func (b *Bot) ruleKey(fromQQ uint64) []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10) + "_water")
}

func (b *Bot) moneyKey(fromQQ uint64) []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10) + "_money")
}

func (b *Bot) getKeyPrefix() []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_")
}

var fight = [...]string{
	"不堪一击",
	"毫不足虑",
	"不足挂齿",
	"初学乍练",
	"勉勉强强",
	"初窥门径",
	"初出茅庐",
	"略知一二",
	"普普通通",
	"平平常常",
	"平淡无奇",
	"粗懂皮毛",
	"半生不熟",
	"登堂入室",
	"略有小成",
	"已有小成",
	"鹤立鸡群",
	"驾轻就熟",
	"青出於蓝",
	"融会贯通",
	"心领神会",
	"炉火纯青",
	"了然於胸",
	"波澜不惊",
	"略有大成",
	"已有大成",
	"豁然贯通",
	"不同凡响",
	"非比寻常",
	"出类拔萃",
	"罕有敌手",
	"波澜老成",
	"冠绝一时",
	"锐不可当",
	"断蛟伏虎",
	"百战不殆",
	"技冠群雄",
	"超群绝伦",
	"神乎其技",
	"出神入化",
	"傲视群雄",
	"超尘拔俗",
	"无出其右",
	"登峰造极",
	"无与伦比",
	"靡坚不摧",
	"所向披靡",
	"一代宗师",
	"世外高人",
	"精深奥妙",
	"神功盖世",
	"望而生遁",
	"勇冠三军",
	"横扫千军",
	"万敌莫开",
	"举世无双",
	"独步天下",
	"指点江山",
	"惊世骇俗",
	"撼天动地",
	"震古铄今",
	"亘古一人",
	"席卷八荒",
	"超凡入圣",
	"威镇寰宇",
	"空前绝后",
	"天人合一",
	"深藏不露",
	"深不可测",
	"高深莫测",
	"岿然如峰",
	"势如破竹",
	"势涌彪发",
	"叱咤喑呜",
	"跌宕昭彰",
	"拿云攫石",
	"摧朽拉枯",
	"鹰撮霆击",
	"鳌掷鲸吞",
	"潮鸣电掣",
	"风云变色",
	"风行雷厉",
	"拔山举鼎",
	"山呼海啸",
	"回山倒海",
	"倒峡泻河",
	"熏天赫地",
	"拔地参天",
	"万仞之巅",
	"气贯长虹",
	"气吞山河",
	"欺霜傲雪",
	"立海垂云",
	"移山竭海",
	"凤翥龙翔",
	"摘星逐月",
	"气凌霄汉",
	"扭转乾坤",
	"洗尽铅华",
	"返璞归真",
}
