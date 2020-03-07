package mission

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/molin0000/secretMaster/secret"
)

func TestJsonRlp(t *testing.T) {
	a := -100
	b := uint64(a)

	fmt.Printf("%d %d\n\n", b, int64(b))

	var ms Mission
	str := `{
		"title": "简单怪兽",
		"author": "空想之喵",
		"date": "2020-03-06T14:48:29.080Z",
		"events": [
		  {
			"key": 0,
			"desc": "你收到一封信",
			"options": [
			  {
				"key": 0,
				"desc": "拆开",
				"results": [
				  {
					"key": 0,
					"type": "跳转",
					"data": 1
				  }
				]
			  },
			  {
				"key": 1,
				"desc": "退回",
				"results": [
				  {
					"key": 0,
					"type": "结束",
					"data": "无事发生"
				  }
				]
			  }
			]
		  },
		  {
			"key": 1,
			"desc": "信中是一副地图，上面标注了一个位置",
			"options": [
			  {
				"key": 0,
				"desc": "前往",
				"results": [
				  {
					"key": 0,
					"type": "跳转",
					"data": 2
				  }
				]
			  },
			  {
				"key": 1,
				"desc": "不去",
				"results": [
				  {
					"key": 0,
					"type": "结束",
					"data": "无事发生"
				  }
				]
			  }
			]
		  },
		  {
			"key": 2,
			"desc": "你在地图地点反复搜索，找到一个山洞。",
			"options": [
			  {
				"key": 0,
				"desc": "进入",
				"results": [
				  {
					"key": 0,
					"type": ".ra",
					"data": "80"
				  },
				  {
					"key": 1,
					"type": "ra成功",
					"data": "3"
				  },
				  {
					"key": 2,
					"type": "ra失败",
					"data": "4"
				  }
				]
			  },
			  {
				"key": 1,
				"desc": "返回",
				"results": [
				  {
					"key": 0,
					"type": "结束",
					"data": "一无所获"
				  }
				]
			  }
			]
		  },
		  {
			"key": 3,
			"desc": "你遇到了怪兽。",
			"options": [
			  {
				"key": 0,
				"desc": "战斗",
				"results": [
				  {
					"key": 0,
					"type": "描述",
					"data": "你打败了怪兽"
				  },
				  {
					"key": 1,
					"type": "经验",
					"data": "50"
				  },
				  {
					"key": 2,
					"type": "金镑",
					"data": "-50"
				  },
				  {
					"key": 3,
					"type": "结束",
					"data": "你满载而归"
				  }
				]
			  },
			  {
				"key": 1,
				"desc": "逃离",
				"results": [
				  {
					"key": 0,
					"type": "结束",
					"data": "你艰难的逃走了。"
				  }
				]
			  }
			]
		  },
		  {
			"key": 4,
			"desc": "你一无所获。"
		  }
		]
	  }`
	json.Unmarshal([]byte(str), &ms)
	fmt.Printf("%+v\n\n\n", ms)

	secret.SetGlobalValue("Test", &ms)
	new := secret.GetGlobalValue("Test", &Mission{}).(*Mission)
	fmt.Printf("%+v", new)
}

func TestRandomMission(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 100; i++ {
		ms := NewRandomMission("/Users/molin/Downloads/mission")
		fmt.Printf("%+v\n\n", ms)
	}
}

func TestRunMission(t *testing.T) {
	for m := 0; m < 10; m++ {
		ms := NewRandomMission("/Users/molin/Downloads/mission")
		fmt.Println(ms.ShowEvent(0))
		for i := 0; i < 100; i++ {
			msg, ret := ms.SelectOption(int(ms.Event), rand.Intn(len(ms.Ms.Events[ms.Event].Options)+1))
			fmt.Println(msg)
			if ret {
				break
			}
		}
		fmt.Println(ms.Finish())
	}
}
