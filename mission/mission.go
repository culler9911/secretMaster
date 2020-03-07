package mission

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
)

type Result struct {
	Key  json.Number `json:"key"`
	Type string      `json:"type"`
	Data json.Number `json:"data"`
}

type Option struct {
	Key     json.Number `json:"key"`
	Desc    string      `json:"desc"`
	Results []Result    `json:"results"`
}

type Event struct {
	Key     json.Number `json:"key"`
	Desc    string      `json:"desc"`
	Options []Option    `json:"options"`
}

type Mission struct {
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Date   string  `json:"date"`
	Events []Event `json:"events"`
}

type MissionGame struct {
	Ms    *Mission
	State uint64 //0: not start, 1: running, 2: finish
	Event uint64
	Money uint64
	Exp   uint64
	Magic uint64
}

var missionList = []*Mission{}

func loadMissions(jsonPath string) {
	files, _ := ioutil.ReadDir(jsonPath)
	for _, f := range files {
		fmt.Println(f.Name())
		str := readJSONFile(path.Join(jsonPath, f.Name()))
		var ms Mission
		err := json.Unmarshal([]byte(str), &ms)
		if err != nil {
			fmt.Println(err)
			continue
		}
		missionList = append(missionList, &ms)
	}
}

func readJSONFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return ""
	}
	return string(fd)
}

// NewRandomMission 随机获取一个新的任务副本
func NewRandomMission(jsonPath string) *MissionGame {
	if len(missionList) == 0 {
		loadMissions(jsonPath)
	}

	idx := rand.Intn(len(missionList))
	mg := &MissionGame{Ms: missionList[idx]}
	return mg
}

func Run(mg *MissionGame, selection uint64) string {
	if mg.State == 0 {
		return mg.Ms.Events[0].Desc
	}
	if selection > len(mg.Ms.Events[mg.Event].Options) {
		return "选项错误"
	}

}

func GetInfo(mg *MissionGame) string {
	return fmt.Sprintf("副本：%s, 作者：%s, 时间：%s", mg.Ms.Title, mg.Ms.Author, mg.Ms.Date)
}
