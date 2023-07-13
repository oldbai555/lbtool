package sixhecai

import (
	"fmt"
	"github.com/oldbai555/lbtool/extpkg/pie/pie"
	"github.com/oldbai555/lbtool/log"
	"strconv"
	"strings"
)

func ParseNumber(input string) (string, error) {
	// 点 做区分
	split := strings.Split(input, ".")
	if len(split) != 2 {
		return "", nil
	}
	tm := split[0]
	xz, err := strconv.Atoi(split[1])
	if err != nil {
		log.Errorf("err is %v", err)
		return "", err
	}
	var tmList []uint32
	var tmStrList []string
	for i := 0; i < len(tm); i += 2 {
		if len(tm) < i+2 {
			break
		}
		s := tm[i : i+2]
		atoi, err := strconv.Atoi(s)
		if err != nil {
			log.Errorf("err is %v", err)
			return "", err
		}
		tmList = append(tmList, uint32(atoi))
		tmStrList = append(tmStrList, s)
	}
	for _, tm := range tmList {
		log.Infof("tm is %d", tm)
	}
	return strings.Join(tmStrList, ",") + "各" + fmt.Sprintf("%d", xz), nil
}

var numberList = pie.Strings([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"})

func ParseText(input string) (string, error) {
	lines := strings.Split(input, "\n")
	var result []string
	for _, line := range lines {
		split := strings.Split(line, "各")
		if len(split) != 2 {
			continue
		}

		var xzList []string
		var fisrt bool
		for _, i := range split[1] {
			if numberList.Contains(string(i)) {
				fisrt = true
				xzList = append(xzList, string(i))
			} else {
				if fisrt {
					break
				}
				continue
			}
		}
		xz, err := strconv.Atoi(strings.Join(xzList, ""))
		if err != nil {
			log.Errorf("err is %v", err)
			continue
		}

		var res []string
		for _, s := range split[0] {
			// 过滤生肖
			if SXList.Contains(string(s)) {
				res = append(res, string(s))
			}
		}

		if len(res) == 0 {
			res0 := strings.Split(split[0], ",")
			res1 := strings.Split(split[0], "-")
			res2 := strings.Split(split[0], ".")
			if len(split[0]) == 2 {
				for _, s := range TmStrList {
					if strings.Contains(split[0], s) {
						res = append(res, s)
					}
				}
			} else if len(res) > 1 {
				for _, re := range res0 {
					for _, s := range TmStrList {
						if strings.Contains(re, s) {
							res = append(res, s)
						}
					}
				}
			} else if len(res1) > 1 {
				for _, re := range res1 {
					for _, s := range TmStrList {
						if strings.Contains(re, s) {
							res = append(res, s)
						}
					}
				}
			} else if len(res2) > 1 {
				for _, re := range res2 {
					for _, s := range TmStrList {
						if strings.Contains(re, s) {
							res = append(res, s)
						}
					}
				}
			} else {
				for i := 0; i < len(split[0]); i += 2 {
					if len(split[0]) < i+2 {
						break
					}
					res = append(res, split[0][i:i+2])
				}
			}
		}

		if len(res) > 0 {
			result = append(result, strings.Join(res, ",")+"各"+fmt.Sprintf("%d", xz))
		}
	}
	return strings.Join(result, "\n"), nil
}
