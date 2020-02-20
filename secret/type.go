package secret

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
	Rank  []uint64
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

type Adventure struct {
	Group  uint64
	QQ     uint64
	DayCnt uint64
	Days   uint64
}

type RespectName struct {
	Group uint64
	QQ    uint64
	Name  string
}

type ExternProperty struct {
	Luck  uint64
	Magic uint64
	Days  uint64
}

type AdventureEvent struct {
	// 类型：1）-钱＆经验5%；2）-经验10%；3）-钱10%；4）+钱＆经验10%；5）+钱15%；6）+经验15%；7）无事发生35%
	Type int
	// 概率：百分比1~100
	Probability int
	// 文案
	Messages []string
}

type BotSwitch struct {
	Group  uint64
	Enable bool
}

type MoneyBind struct {
	IniPath    string
	IniSection string
	IniKey     string
	HasUpdate  bool
}
