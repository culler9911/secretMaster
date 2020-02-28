package secret

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var botMap map[uint64]*Bot

func NewSecretBot(qq, group uint64, groupNick string, private bool, api interface{}) *Bot {
	setCqpCall(api.(CqpCall))
	if botMap == nil {
		botMap = make(map[uint64]*Bot)
	}
	bot, ok := botMap[group]
	if ok {
		bot.Private = private
		return bot
	}
	bot = &Bot{QQ: qq, Group: group, Name: groupNick, Private: private}
	bot.information = &Information{bot}
	bot.career = &Career{bot}
	bot.organization = &Organization{bot}
	bot.store = &Store{bot}

	botMap[group] = bot
	return bot
}

func (b *Bot) Run(msg string, fromQQ uint64, nick string) string {
	if !b.talkToMe(msg) {
		return ""
	}

	b.CurrentNick = nick

	return b.searchMenu(msg, fromQQ, &menus)
}

func (b *Bot) RunPrivate(msg string, fromQQ uint64, nick string) string {
	return b.searchMenu(msg, fromQQ, &menus)
}

func (b *Bot) Update(fromQQ uint64, nick string) string {
	if !b.getSwitch() {
		return ""
	}

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

		// ret = b.levelUpdate(v)

		w := b.getWaterRuleFromDb(fromQQ)
		m := b.getMoneyFromDb(fromQQ, 0)
		fmt.Printf("money:%+v", m)
		e := b.getExternFromDb(fromQQ)
		fmt.Printf("extern:%+v", e)

		if e.Magic > 0 {
			e.Magic--
		}

		normalHumanStop := false
		if v.SecretID > 22 && v.ChatCount > 400 {
			normalHumanStop = true
		}

		if w.DayCnt < 200 && e.Magic > 0 && !normalHumanStop {
			v.ChatCount++
			w.DayCnt++
			m.Money++
			if v.ChatCount%100 == 0 {
				ret += "\n恭喜！你的战力评价升级了！"
			}
		}

		v.LastChat = uint64(time.Now().Unix())
		v.LevelDown = uint64(time.Now().Unix())
		b.setPersonToDb(fromQQ, v)
		b.setWaterRuleToDb(fromQQ, w)
		b.setMoneyToDb(fromQQ, m)
		b.setExternToDb(fromQQ, e)
	}

	return ret
}

func (b *Bot) printMenu(menu *Menu) string {
	if menu == nil {
		return ""
	}

	if menu.ID == 7 && !b.Private {
		return ""
	}

	if !b.getSwitch() {
		return ""
	}

	info := fmt.Sprintf("\n%s: %s \n", menu.Title, menu.Info)
	if menu.SubMenu != nil && len(menu.SubMenu) > 0 {
		for _, mu := range menu.SubMenu {
			if mu.ID == 7 && !b.Private {
				continue
			}

			info += fmt.Sprintf("%s: %s \n", mu.Title, mu.Info)
		}
	}
	if len(menu.Commit) > 0 {
		info += menu.Commit
	}
	return info
}

func (b *Bot) searchMenu(msg string, fromQQ uint64, menu *Menu) string {
	if strings.Contains(msg, menu.Title) {
		if menu.SubMenu != nil && len(menu.SubMenu) > 0 {
			return b.printMenu(menu)
		}
		return b.cmdRun(msg, fromQQ)
	}

	if menu.SubMenu != nil && len(menu.SubMenu) > 0 {
		for _, mu := range menu.SubMenu {
			info := b.searchMenu(msg, fromQQ, &mu)
			if len(info) > 0 {
				return info
			}
		}
	}
	return ""
}

func (b *Bot) cmdRun(msg string, fromQQ uint64) string {
	if strings.Contains(msg, "序列战争关") {
		return b.botSwitch(fromQQ, false)
	}

	if strings.Contains(msg, "序列战争开") {
		return b.botSwitch(fromQQ, true)
	}

	if !b.getSwitch() {
		return ""
	}

	if strings.Contains(msg, "探险卷轴") {
		return b.adventure(fromQQ, false)
	}

	if strings.Contains(msg, "红剧场门票") {
		return b.redTheater(fromQQ)
	}

	if strings.Contains(msg, "属性") {
		return b.getProperty(fromQQ)
	}

	if strings.Contains(msg, "查询") {
		rankStr := msg[strings.Index(msg, "查询")+len("查询"):]
		fmt.Println("查询排名：", rankStr)
		rank, err := strconv.Atoi(rankStr)
		if err != nil {
			return "是找我吗？查询排名要查询加数字哦。如果不是，请不要艾特我。"
		}
		if len(b.Rank) > rank-1 {
			return b.getProperty(b.Rank[rank-1])
		}
		return "请先查看最新排行"
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

	if strings.Contains(msg, "自杀") {
		return b.deletePerson(fromQQ)
	}

	if strings.Contains(msg, "尊名") {
		return b.setRespectName(msg, fromQQ)
	}

	if strings.Contains(msg, "GM") {
		return b.gmCmd(fromQQ, msg)
	}

	if strings.Contains(msg, "货币升级") {
		return b.moneyUpdate(fromQQ, true)
	}

	if strings.Contains(msg, "货币降级") {
		return b.moneyUpdate(fromQQ, false)
	}

	if strings.Contains(msg, "货币映射") {
		return b.moneyMap(fromQQ, msg)
	}

	if strings.Contains(msg, "查看映射") {
		bind := b.getMoneyBind()
		return fmt.Sprintf("%+v\n", bind)
	}

	if strings.Contains(msg, ".master") {
		return b.setMaster(fromQQ, msg)
	}

	if strings.Contains(msg, "道具") {
		return b.getItems(fromQQ)
	}

	if strings.Contains(msg, "神位") {
		return b.getGod()
	}

	if strings.Contains(msg, "灵性药剂") {
		return b.buyMagicPotion(fromQQ)
	}

	if strings.Contains(msg, "灵性材料") {
		return b.buyMagicItem(fromQQ)
	}

	if strings.Contains(msg, "至高权杖") {
		return b.buyMace(fromQQ)
	}

	if strings.Contains(msg, "技能") {
		return b.getSkill(fromQQ)
	}

	if strings.Contains(msg, "晋升") {
		return b.promotion(fromQQ)
	}

	if strings.Contains(msg, "创建") {
		return b.createChurch(fromQQ, msg)
	}

	if strings.Contains(msg, "解散") {
		return b.deleteChurch(fromQQ, msg)
	}

	if strings.Contains(msg, "寻访") {
		return b.listChurch()
	}

	return ""
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
