/**
 * @Author: zjj
 * @Date: 2025/2/18
 * @Desc: 解析导出的html书签
**/

package parsehtmlbookbookmark

import (
	"bufio"
	"fmt"
	"github.com/oldbai555/lbtool/log"
	"github.com/oldbai555/lbtool/utils"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func Parse(filePath string) {
	// 输入书签文件路径
	outputFile := fmt.Sprintf("%d.txt", utils.TimeNow()) // 输出文件路径

	// 打开书签文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorf("无法打开书签文件: %v\n", err)
		return
	}
	defer file.Close()

	// 解析 HTML
	doc, err := html.Parse(file)
	if err != nil {
		log.Errorf("解析 HTML 失败: %v\n", err)
		return
	}

	// 打开输出文件
	output, err := os.Create(outputFile)
	if err != nil {
		log.Errorf("无法创建输出文件: %v\n", err)
		return
	}
	defer output.Close()

	writer := bufio.NewWriter(output)
	defer writer.Flush()

	// 递归遍历 HTML 节点
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// 提取 URL 和标题
			var url, title string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url = attr.Val
				}
			}
			if n.FirstChild != nil {
				title = strings.TrimSpace(n.FirstChild.Data)
			}

			// 写入文件
			if url != "" && title != "" {
				line := fmt.Sprintf("[%s]======>%s\n", title, url)
				_, err := writer.WriteString(line)
				if err != nil {
					log.Errorf("写入文件失败: %v\n", err)
					return
				}
			}
		}

		// 递归遍历子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	// 开始遍历
	traverse(doc)

	log.Infof("解析完成，结果已保存到 %s\n", outputFile)
}
