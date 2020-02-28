package secret

import (
	"fmt"
	"testing"
)

func TestSkill(t *testing.T) {
	fromQQ := uint64(111)
	b := NewSecretBot(1234, 4567, "aa", false, &debugInteract{})
	b.clearSkill(fromQQ)

	fmt.Println(b.getSkill(fromQQ))

	b.setSkill(fromQQ, 0, 1)
	b.setSkill(fromQQ, 1, 1)
	b.setSkill(fromQQ, 2, 1)

	fmt.Println(b.getSkill(fromQQ))

	b.skillLevelUp(fromQQ, 0)
	b.skillLevelUp(fromQQ, 1)
	b.skillLevelUp(fromQQ, 2)
	b.skillLevelUp(fromQQ, 0)
	b.skillLevelUp(fromQQ, 0)
	b.skillLevelUp(fromQQ, 0)
	b.skillLevelUp(fromQQ, 0)
	b.skillLevelUp(fromQQ, 0)

	fmt.Println(b.getSkill(fromQQ))
	fmt.Println(b.allSkillLevelUp(fromQQ))
	fmt.Println(b.getSkill(fromQQ))

	b.clearSkill(fromQQ)
	b.skillLevelUp(fromQQ, 1)

	fmt.Println(b.getSkill(fromQQ))
}

func TestPromotion(t *testing.T) {
	fromQQ := uint64(111)
	b := NewSecretBot(1234, 4567, "aa", false, &debugInteract{})
	fmt.Println(b.deletePerson(fromQQ))
	b.Update(fromQQ, "ThinkCat")
	fmt.Println(b.promotion(fromQQ))
	fmt.Println(b.changeSecretList("更换1", fromQQ))
	b.Update(fromQQ, "ThinkCat")
	b.setExp(fromQQ, 101)
	fmt.Println(b.promotion(fromQQ))
	b.setMoney(fromQQ, 200)
	b.buyMagicItem(fromQQ)
	fmt.Println(b.promotion(fromQQ))
	b.setExp(fromQQ, 101)
	b.setMoney(fromQQ, 200)
	b.buyMagicItem(fromQQ)
	fmt.Println(b.getMoney(fromQQ), b.getExp(fromQQ))
	fmt.Println(b.promotion(fromQQ))

	for i := 0; i < 1000; i++ {
		b.setMoney(fromQQ, 200)
		b.setExp(fromQQ, 101)

		fmt.Println(b.getMoney(fromQQ), b.getExp(fromQQ))
		fmt.Println(b.promotion(fromQQ))
	}
}
