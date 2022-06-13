package filter

import (
	"regexp"
	"sort"
)

var AllParser = map[string]Filter{
	"SpliterBy----": &SplitBy{
		By:   "----",
		Slot: []string{"account", "passwd"},
	},
	"SpliterByTab": &SplitBy{
		By:   "\t",
		Slot: []string{"account", "passwd"},
	},
	"SpliterByBlank": &SplitBy{
		By:   " ",
		Slot: []string{"account", "passwd"},
	},
	"SpliterBy|": &SplitBy{
		By:   "|",
		Slot: []string{"account", "passwd"},
		// Field: "_raw",
	},
}

var NamePasswd = regexp.MustCompile(`^([a-zA-Z0-9_-]{4,})----(.+)$`)

type ParserScore struct {
	Parser string `json:"parser"`
	Score  int    `json:"score"`
}

func GetParserScore(lines []string) []*ParserScore {
	var res []*ParserScore
	for k, v := range AllParser {
		score := 0
		for _, line := range lines {
			e := NewElement(line)
			if _, ok := v.Do(e); ok {
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
