package secret

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
		ret = "恭喜！你的序列晋升了！" + b.allSkillLevelUp(p.QQ)
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

	god := b.getGodFromDb(v.SecretID)
	if god == 0 || god == v.QQ {
		god = 0
		b.setGodToDb(v.SecretID, &god)
	}

	v.SecretID = uint64(value - 1)
	b.setPersonToDb(fromQQ, v)

	b.setMoney(fromQQ, -100)

	if value == 18 {
		b.setSkill(fromQQ, 0, 1)
	}

	return fmt.Sprintf("成功更换到途径：%d", value)
}

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

func (b *Bot) setSkill(fromQQ, skill uint64, level uint64) {
	tree := b.getPersonValue("SkillTree", fromQQ, &SkillTree{}).(*SkillTree)
	find := false
	for i := 0; i < len(tree.Skills); i++ {
		if tree.Skills[i].Name == skillList[skill].Name {
			if level <= tree.Skills[i].MaxLevel {
				tree.Skills[i].Level = level
			} else {
				tree.Skills[i].Level = tree.Skills[i].MaxLevel
			}
			find = true
			break
		}
	}

	if !find {
		tree.Skills = append(tree.Skills, &Skill{ID: skillList[skill].ID, Name: skillList[skill].Name, Level: skillList[skill].Level, MaxLevel: skillList[skill].MaxLevel})
	}
	b.setPersonValue("SkillTree", fromQQ, tree)
}

func (b *Bot) removeSkill(fromQQ, skill uint64) {
	tree := b.getPersonValue("SkillTree", fromQQ, &SkillTree{}).(*SkillTree)
	for i := 0; i < len(tree.Skills); i++ {
		if tree.Skills[i].Name == skillList[skill].Name {
			tree.Skills[i] = tree.Skills[len(tree.Skills)-1] // Copy last element to index i.
			tree.Skills[len(tree.Skills)-1] = nil            // Erase last element (write zero value).
			tree.Skills = tree.Skills[:len(tree.Skills)-1]   // Truncate slice.
			break
		}
	}
	b.setPersonValue("SkillTree", fromQQ, tree)
}

func (b *Bot) clearSkill(fromQQ uint64) {
	b.setPersonValue("SkillTree", fromQQ, &SkillTree{})
}

func (b *Bot) skillLevelUp(fromQQ uint64, skill uint64) {
	tree := b.getPersonValue("SkillTree", fromQQ, &SkillTree{}).(*SkillTree)
	find := false
	for i := 0; i < len(tree.Skills); i++ {
		if tree.Skills[i].ID == skill {
			if tree.Skills[i].Level <= tree.Skills[i].MaxLevel {
				tree.Skills[i].Level++
			} else {
				tree.Skills[i].Level = tree.Skills[i].MaxLevel
			}
			find = true
			break
		}
	}

	if !find {
		tree.Skills = append(tree.Skills, &Skill{ID: skillList[skill].ID, Name: skillList[skill].Name, Level: skillList[skill].Level, MaxLevel: skillList[skill].MaxLevel})
	}
	b.setPersonValue("SkillTree", fromQQ, tree)
}

func (b *Bot) allSkillLevelUp(fromQQ uint64) string {
	tree := b.getPersonValue("SkillTree", fromQQ, &SkillTree{}).(*SkillTree)
	if len(tree.Skills) > 0 {
		for i := 0; i < len(tree.Skills); i++ {
			if tree.Skills[i].Level < tree.Skills[i].MaxLevel {
				tree.Skills[i].Level++
			} else {
				tree.Skills[i].Level = tree.Skills[i].MaxLevel
			}
		}
		b.setPersonValue("SkillTree", fromQQ, tree)
		return "技能等级提示！"
	}

	return ""
}
