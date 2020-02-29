package secret

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func calcInterest(bank *Bank) *Bank {
	today := uint64(time.Now().Unix() / (3600 * 24))
	if bank.Date != 0 && bank.Date != today {
		bank.Amount = uint64(float64(bank.Amount) * math.Pow(1.03, float64(today-bank.Date)))
		bank.Date = today
	}

	return bank
}

func (b *Bot) bank(fromQQ uint64, msg string) string {
	strs := strings.Split(msg, ";")
	if len(strs) < 2 {
		return "指令错误" + fmt.Sprintf("%+v", strs)
	}

	switch strs[1] {
	case "存款":
		if len(strs) < 3 {
			return "指令错误" + fmt.Sprintf("%+v", strs)
		}
		money, err := strconv.Atoi(strs[2])
		if err != nil {
			return "指令错误" + fmt.Sprintf("%+v", err)
		}
		if money < 100 {
			return "每次存款最少100金镑。"
		}
		if uint64(money) > b.getMoney(fromQQ) {
			return "你并没有那么多钱。"
		}

		bank := b.getPersonValue("Bank", fromQQ, &Bank{}).(*Bank)
		bank = calcInterest(bank)
		bank.Amount += uint64(money)
		b.setPersonValue("Bank", fromQQ, bank)
		b.setMoney(fromQQ, -1*money)
		return fmt.Sprintf("成功存入%d金镑，当前存款%d金镑。", money, bank.Amount)
	case "取款":
		if len(strs) < 3 {
			return "指令错误" + fmt.Sprintf("%+v", strs)
		}
		money, err := strconv.Atoi(strs[2])
		if err != nil {
			return "指令错误" + fmt.Sprintf("%+v", err)
		}
		if money <= 0 {
			return "你想取个西瓜？"
		}
		bank := b.getPersonValue("Bank", fromQQ, &Bank{}).(*Bank)
		bank = calcInterest(bank)
		if uint64(money) > bank.Amount {
			return "你没有那么多钱。"
		}

		bank.Amount -= uint64(money)
		b.setPersonValue("Bank", fromQQ, bank)
		b.setMoney(fromQQ, money)
		return fmt.Sprintf("成功取出%d金镑", money)
	case "查账":
		bank := b.getPersonValue("Bank", fromQQ, &Bank{}).(*Bank)
		bank = calcInterest(bank)
		b.setPersonValue("Bank", fromQQ, bank)
		return fmt.Sprintf("当前存款%d金镑。日利率3％", bank.Amount)
	default:
		return "指令错误" + fmt.Sprintf("%+v", strs)
	}
}
