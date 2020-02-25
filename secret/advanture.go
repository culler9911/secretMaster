package secret

import (
	"fmt"
	"math/rand"
	"time"
)

func (b *Bot) adventure(fromQQ uint64, limit bool) string {
	a := b.getAdvFromDb(fromQQ)
	if limit && a.DayCnt >= 3 {
		return "对不起，您今日奇遇探险机会已经用完"
	}

	a.DayCnt++
	a.Days = uint64(time.Now().Unix() / (3600 * 24))
	b.setAdvToDb(fromQQ, a)
	externProperty := b.getExternFromDb(fromQQ)

	rand.Seed(time.Now().UnixNano())

	i := int(externProperty.Luck*5) + rand.Intn(100-int(externProperty.Luck*5))
	info := ""
	m := 0
	e := 0
	fmt.Println("i:", i)
	for p := 0; p < len(adventureEvents); p++ {
		i = i - adventureEvents[p].Probability
		fmt.Println("i:", i)
		if i < 0 {
			switch adventureEvents[p].Type {
			case 1:
				m = -1 * (10 + rand.Intn(41))
				e = -1 * (10 + rand.Intn(41))
				info = fmt.Sprintf("%s经验:%d, 金镑:%d", adventureEvents[p].Messages[rand.Intn(len(adventureEvents[p].Messages))], e, m)
			case 2:
				m = 0
				e = -1 * (10 + rand.Intn(41))
				info = fmt.Sprintf("%s经验:%d", adventureEvents[p].Messages[rand.Intn(len(adventureEvents[p].Messages))], e)
			case 3:
				m = -1 * (10 + rand.Intn(41))
				e = 0
				info = fmt.Sprintf("%s金镑:%d", adventureEvents[p].Messages[rand.Intn(len(adventureEvents[p].Messages))], m)
			case 4:
				m = (20 + rand.Intn(81))
				e = (20 + rand.Intn(81))
				info = fmt.Sprintf("%s经验:%d, 金镑:%d", adventureEvents[p].Messages[rand.Intn(len(adventureEvents[p].Messages))], e, m)
			case 5:
				m = (20 + rand.Intn(81))
				e = 0
				info = fmt.Sprintf("%s金镑:%d", adventureEvents[p].Messages[rand.Intn(len(adventureEvents[p].Messages))], m)
			case 6:
				m = 0
				e = (20 + rand.Intn(81))
				info = fmt.Sprintf("%s经验:%d", adventureEvents[p].Messages[rand.Intn(len(adventureEvents[p].Messages))], e)
			default:
				m = 0
				e = 0
				info = fmt.Sprintf("%s", adventureEvents[p].Messages[rand.Intn(len(adventureEvents[p].Messages))])
			}
			break
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
