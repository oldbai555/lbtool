package orm

import "strings"

type tagRules struct {
	ruleMap  map[string]string
	ruleList [][]string
}

// parseModelTag
func parseModelTag(dbDef string) *tagRules {
	d := &tagRules{
		ruleMap: map[string]string{},
	}
	items := strings.Split(dbDef, ";")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		ab := strings.SplitN(item, ":", 2)
		if len(ab) == 0 {
			continue
		}
		k := strings.TrimSpace(ab[0])
		var v string
		if len(ab) > 1 {
			v = strings.TrimSpace(ab[1])
		}
		d.ruleMap[k] = v
		d.ruleList = append(d.ruleList, []string{k, v})
	}
	return d
}
