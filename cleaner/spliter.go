package cleaner

import "strings"

type SpliterBy struct {
	By string
}

func (s *SpliterBy) Work(raw []string) []string {
	ret := []string{}
	for _, v := range raw {
		if strings.Contains(v, s.By) {
			ret = append(ret, strings.Split(v, s.By)...)
		}
	}
	return ret
}
