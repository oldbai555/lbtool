package warning

import (
	"encoding/json"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/pkg/restysdk"
)

type FeishuBody struct {
	MsgType string        `json:"msg_type"`
	Content FeishuContent `json:"content"`
}

type FeishuContent struct {
	Post PostContent `json:"post"`
}

type PostContent struct {
	ZhCn ZhCnContent `json:"zh_cn"`
}

type ZhCnContent struct {
	Title   string        `json:"title"`
	Content [][]TagDetail `json:"content"`
}

type TagDetail struct {
	Tag      string `json:"tag"`
	UnEscape bool   `json:"un_escape"`
	UserId   string `json:"user_id"`
	Text     string `json:"text"`
}

func ReportToFeishu(title, content, groupId string, atList ...string) {
	url := fmt.Sprintf("https://open.feishu.cn/open-apis/bot/v2/hook/%s", groupId)
	body := &FeishuBody{
		MsgType: "post",
		Content: FeishuContent{
			Post: PostContent{
				ZhCn: ZhCnContent{
					Title:   title,
					Content: [][]TagDetail{},
				},
			},
		},
	}
	var tagDetails []TagDetail
	tagDetails = append(tagDetails, TagDetail{
		Tag:      "text",
		UnEscape: true,
		Text:     content,
	})

	for _, at := range atList {
		tagDetails = append(tagDetails, TagDetail{
			Tag:    "at",
			UserId: at,
		})
	}
	body.Content.Post.ZhCn.Content = append(body.Content.Post.ZhCn.Content, tagDetails)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	headerMap := make(map[string]string)
	headerMap["Content-Type"] = "application/json"
	_, err = restysdk.NewRequest().SetHeaders(headerMap).SetBody(jsonBody).Post(url)
	if err != nil {
		log.Errorf("err: %v", err)
	}
}
