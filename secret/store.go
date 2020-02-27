package secret

import "time"

func (b *Bot) redTheater(fromQQ uint64) string {
	if b.getMoney(fromQQ) < 200 {
		return "你站在红剧场售票厅，盯着票价，攥紧了空荡荡的口袋，穷的掩面而去。"
	}

	if b.getMagic(fromQQ) < 100 {
		return "演出才刚刚开始，你突然一阵头晕目眩，感觉到灵性枯竭预警，立刻退票离开了。"
	}

	b.setMoney(fromQQ, -200)
	b.setMagic(fromQQ, -99)
	b.setExp(fromQQ, 120)
	return "你看了一场酣畅淋漓的演出，感觉自己对这个世界的了解更加深刻了。经验+120"
}

func (b *Bot) buyMagicPotion(fromQQ uint64) string {
	potion := b.getPersonValue("Potion", fromQQ, &Potion{0, uint64(time.Now().Unix() / 3600 / 24)})

	if potion.(*Potion).Days != uint64(time.Now().Unix()/3600/24) {
		potion.(*Potion).DayCnt = 0
		potion.(*Potion).Days = uint64(time.Now().Unix() / 3600 / 24)
	}

	if b.getMoney(fromQQ) < 70 {
		return "你盯着商品报价，攥紧了空荡荡的钱包，穷的掩面而去。"
	}

	if potion.(*Potion).DayCnt >= 3 && potion.(*Potion).Days == uint64(time.Now().Unix()/3600/24) {
		return "你今天喝了太多，实在是喝不下了。"
	}

	potion.(*Potion).DayCnt++

	b.setMagic(fromQQ, 30)
	b.setMoney(fromQQ, -70)
	b.setPersonValue("Potion", fromQQ, potion)

	return "喝下这瓶药剂后，你满嘴苦涩，灵性似乎恢复了一点。"
}
