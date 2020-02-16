package secret

import (
	"fmt"
	"strconv"
	"time"

	"github.com/molin0000/secretMaster/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func getDb() *leveldb.DB {
	if db == nil {
		_db, err := leveldb.OpenFile("secret.db", nil)
		if err != nil {
			fmt.Printf("open db error: %+v", err)
			panic(err)
		}
		db = _db
	}
	return db
}

func (b *Bot) keys(fromQQ uint64) []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10))
}

func (b *Bot) ruleKey(fromQQ uint64) []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10) + "_water")
}

func (b *Bot) moneyKey(fromQQ uint64) []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10) + "_money")
}

func (b *Bot) advKey(fromQQ uint64) []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10) + "_adventure")
}

func (b *Bot) godKey(secretID uint64) []byte {
	return []byte("god_" + strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(secretID), 10))
}

func (b *Bot) rnameKey(fromQQ uint64) []byte {
	return []byte("rname_" + strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10))
}

func (b *Bot) externKey(fromQQ uint64) []byte {
	return []byte("extern_" + strconv.FormatInt(int64(b.Group), 10) + "_" + strconv.FormatInt(int64(fromQQ), 10))
}

func (b *Bot) switchKey() []byte {
	return []byte("switch_" + strconv.FormatInt(int64(b.Group), 10))
}

func (b *Bot) getKeyPrefix() []byte {
	return []byte(strconv.FormatInt(int64(b.Group), 10) + "_")
}

func (b *Bot) getSwitch() bool {
	verify, err := getDb().Get(b.switchKey(), nil)
	if err != nil {
		return true
	}
	var v BotSwitch
	rlp.DecodeBytes(verify, &v)
	return v.Enable
}

func (b *Bot) setSwitch(enable bool) {
	p := &BotSwitch{b.Group, enable}
	buf, _ := rlp.EncodeToBytes(p)
	getDb().Put(b.switchKey(), buf, nil)
}

func (b *Bot) getPersonFromDb(fromQQ uint64) *Person {
	verify, _ := getDb().Get(b.keys(fromQQ), nil)
	var v Person
	rlp.DecodeBytes(verify, &v)
	return &v
}

func (b *Bot) setPersonToDb(fromQQ uint64, p *Person) {
	buf, _ := rlp.EncodeToBytes(p)
	getDb().Put(b.keys(fromQQ), buf, nil)
}

func (b *Bot) getWaterRuleFromDb(fromQQ uint64) *WaterRule {
	ret, err := getDb().Get(b.ruleKey(fromQQ), nil)
	if err != nil {
		return &WaterRule{Group: b.Group, QQ: fromQQ, DayCnt: 1, Days: uint64(time.Now().Unix() / (3600 * 24))}
	}
	var w WaterRule
	rlp.DecodeBytes(ret, &w)
	if w.Days != uint64(time.Now().Unix()/(3600*24)) {
		w.DayCnt = 0
		w.Days = uint64(time.Now().Unix() / (3600 * 24))
	}
	return &w
}

func (b *Bot) setWaterRuleToDb(fromQQ uint64, w *WaterRule) {
	buf, _ := rlp.EncodeToBytes(w)
	getDb().Put(b.ruleKey(fromQQ), buf, nil)
}

func (b *Bot) getMoneyFromDb(fromQQ uint64, chatCnt uint64) *Money {
	ret, err := getDb().Get(b.moneyKey(fromQQ), nil)
	if err != nil {
		return &Money{Group: b.Group, QQ: fromQQ, Money: chatCnt}
	}
	var m Money
	rlp.DecodeBytes(ret, &m)
	return &m
}

func (b *Bot) setMoneyToDb(fromQQ uint64, m *Money) {
	buf, _ := rlp.EncodeToBytes(m)
	getDb().Put(b.moneyKey(fromQQ), buf, nil)
}

func (b *Bot) getAdvFromDb(fromQQ uint64) *Adventure {
	ret, err := getDb().Get(b.advKey(fromQQ), nil)
	if err != nil {
		return &Adventure{Group: b.Group, QQ: fromQQ, Days: uint64(time.Now().Unix() / (3600 * 24)), DayCnt: 0}
	}
	var m Adventure
	rlp.DecodeBytes(ret, &m)
	if m.Days != uint64(time.Now().Unix()/(3600*24)) {
		m.DayCnt = 0
	}
	return &m
}

func (b *Bot) setAdvToDb(fromQQ uint64, m *Adventure) {
	buf, _ := rlp.EncodeToBytes(m)
	getDb().Put(b.advKey(fromQQ), buf, nil)
}

func (b *Bot) getGodFromDb(secretID uint64) uint64 {
	ret, err := getDb().Get(b.godKey(secretID), nil)
	if err != nil {
		return 0
	}
	var god uint64
	rlp.DecodeBytes(ret, &god)
	return god
}

func (b *Bot) setGodToDb(secretID uint64, god *uint64) {
	buf, _ := rlp.EncodeToBytes(god)
	getDb().Put(b.godKey(secretID), buf, nil)
}

func (b *Bot) setRNameToDb(fromQQ uint64, r *RespectName) {
	buf, _ := rlp.EncodeToBytes(r)
	getDb().Put(b.rnameKey(fromQQ), buf, nil)
}

func (b *Bot) getRNameFromDb(fromQQ uint64) string {
	ret, err := getDb().Get(b.rnameKey(fromQQ), nil)
	if err != nil {
		return "无"
	}
	var god RespectName
	rlp.DecodeBytes(ret, &god)
	return "\n" + god.Name
}

func (b *Bot) setExternToDb(fromQQ uint64, r *ExternProperty) {
	buf, _ := rlp.EncodeToBytes(r)
	getDb().Put(b.externKey(fromQQ), buf, nil)
}

func (b *Bot) getExternFromDb(fromQQ uint64) *ExternProperty {
	ret, err := getDb().Get(b.externKey(fromQQ), nil)
	if err != nil {
		return &ExternProperty{Luck: 0, Magic: 200, Days: uint64(time.Now().Unix() / 3600 / 24)}
	}
	var v ExternProperty
	rlp.DecodeBytes(ret, &v)
	if v.Days != uint64(time.Now().Unix()/(3600*24)) {
		v.Magic = 200
		v.Days = uint64(time.Now().Unix() / (3600 * 24))
	}

	p := b.getPersonFromDb(fromQQ)
	if p.SecretID == 17 {
		if v.Luck < 5 {
			v.Luck += 5
		}
	}
	return &v
}
