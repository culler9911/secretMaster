package secret

func (b *Bot) redTheater(fromQQ uint64) string {
	p := b.getPersonFromDb(fromQQ)
	m := b.getMoneyFromDb(fromQQ, p.ChatCount)
	e := b.getExternFromDb(fromQQ)
	if m.Money < 200 {
		return "你站在红剧场售票厅，盯着票价，攥紧了空荡荡的口袋，穷的掩面而去。"
	}

	if e.Magic < 100 {
		return "演出才刚刚开始，你突然一阵头晕目眩，感觉到灵性枯竭预警，立刻退票离开了。"
	}

	m.Money -= 200
	e.Magic -= 100
	p.ChatCount += 120

	b.setMoneyToDb(fromQQ, m)
	b.setExternToDb(fromQQ, e)
	b.setPersonToDb(fromQQ, p)

	return "你看了一场酣畅淋漓的演出，感觉自己对这个世界的了解更加深刻了。经验+120"
}
