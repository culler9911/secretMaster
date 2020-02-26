package secret

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
