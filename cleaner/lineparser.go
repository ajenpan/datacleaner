package cleaner

import (
	"regexp"
	"sort"
	"strings"
)

type LineParser interface {
	Parse(line string) (map[string]interface{}, bool)
	// String() string
}

var AllParser = map[string]LineParser{
	"SpliterBy----": &SpliterBy{
		By:   "----",
		Slot: []string{"account", "passwd"},
	},
	"SpliterByTab": &SpliterBy{
		By:   "\t",
		Slot: []string{"account", "passwd"},
	},
	"SpliterByBlank": &SpliterBy{
		By:   " ",
		Slot: []string{"account", "passwd"},
	},
	"SpliterBy|": &SpliterBy{
		By:   "|",
		Slot: []string{"account", "passwd"},
	},
}

var NamePasswd = regexp.MustCompile(`^([a-zA-Z0-9_-]{4,})----(.+)$`)

type ParserScore struct {
	Parser string `json:"parser"`
	Score  int    `json:"score"`
}

func GetParserScore(Lines []string) []*ParserScore {
	var res []*ParserScore
	for k, v := range AllParser {
		score := 0
		for _, line := range Lines {
			if _, done := v.Parse(line); done {
				score++
			}
		}
		res = append(res, &ParserScore{k, score})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Score > res[j].Score
	})
	return res
}

func BestParser(lines []string) *ParserScore {
	res := GetParserScore(lines)
	if len(res) == 0 {
		return nil
	}
	return res[0]
}

type SpliterBy struct {
	By   string
	Slot []string

	// Name string
}

// func (s *SpliterBy) String() string {
// 	return s.Name
// }

func (s *SpliterBy) Parse(line string) (map[string]interface{}, bool) {
	if len(line) < len(s.By)*(len(s.Slot)-1) {
		return nil, false
	}
	temp := strings.Split(line, s.By)

	res := make([]string, 0, len(s.Slot))

	for _, v := range temp {
		vv := strings.TrimSpace(v)
		if len(vv) != 0 {
			res = append(res, vv)
		}
	}

	if len(res) != len(s.Slot) {
		return nil, false
	}

	resMap := make(map[string]interface{})
	for i, v := range s.Slot {
		resMap[v] = res[i]
	}
	return resMap, true
}