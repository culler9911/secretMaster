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
	b.Update(fromQQ, "空想之喵")
	fmt.Println(b.promotion(fromQQ))
	fmt.Println(b.changeSecretList("更换1", fromQQ))
	b.Update(fromQQ, "空想之喵")
	b.setExp(fromQQ, 101)
	fmt.Println(b.promotion(fromQQ))

}
