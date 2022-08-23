package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//SimpleUsage()
	Usage()
}

// CustomerImapClient 调用NewImapClient
func CustomerImapClient() (*client.Client, error) {
	// 【修改】账号和密码
	return NewImapClient("username", "password")
}

// NewImapClient 创建IMAP客户端
func NewImapClient(username, password string) (*client.Client, error) {
	// 【字符集】  处理us-ascii和utf-8以外的字符集(例如gbk,gb2313等)时,
	//  需要加上这行代码。
	// 【参考】 https://github.com/emersion/go-imap/wiki/Charset-handling
	imap.CharsetReader = charset.Reader

	log.Println("Connecting to server...")

	// 连接邮件服务器
	c, err := client.DialTLS("imap.qq.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// 使用账号密码登录
	if err := c.Login(username, password); err != nil {
		return nil, err
	}

	log.Println("Logged in")

	return c, nil
}

func SimpleUsage() {
	// 连接邮件服务器
	c, err := CustomerImapClient()
	if err != nil {
		log.Fatal(err)
	}
	// Don't forget to logout
	defer c.Logout()

	// 选择收件箱
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// 获取近50封邮件
	from := uint32(50)
	to := mbox.Messages
	if mbox.Messages > 50 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 50
	}
	seqset := new(imap.SeqSet)
	// 设置邮件搜索范围
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		// 抓取邮件消息体传入到messages信道
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	for msg := range messages {
		// 打印邮件标题
		log.Println("* " + msg.Envelope.Subject)
	}

	if err = <-done; err != nil {
		log.Fatal(err)
	}
}

// Usage
// 【处理业务需求】假设需求是找出求以subject开头的标题的最新邮件，并下载附件。
// 【思路】有些邮件包含附件后会变得特别大，如果要遍历的邮件很多，直接遍历处理，每封邮件都获取'RFC822'内容，
// fetch方法执行耗时可能会很长, 因此可以分两次fetch处理，减少处理时长：
// 1)第一次fetch先使用ENVELOP或者RFC822.HEADER获取邮件头信息找到满足业务需求邮件的id
// 2)第二次fetch根据这个邮件id使用'RFC822'获取邮件MIME内容，下载附件
func Usage() {
	// 连接邮件服务器
	c, err := CustomerImapClient()
	if err != nil {
		log.Fatal(err)
	}
	// Don't forget to logout
	defer c.Logout()

	// 选择收件箱
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// 搜索条件实例对象
	criteria := imap.NewSearchCriteria()

	// ALL是默认条件
	// See RFC 3501 section 6.4.4 for a list of searching criteria.
	criteria.WithoutFlags = []string{"ALL"}
	ids, _ := c.Search(criteria)
	var s imap.BodySectionName

	for {
		if len(ids) == 0 {
			break
		}
		id := pop(&ids)

		seqset := new(imap.SeqSet)
		seqset.AddNum(id)
		chanMessage := make(chan *imap.Message, 1)
		go func() {
			// 第一次fetch, 只抓取邮件头，邮件标志，邮件大小等信息，执行速度快
			if err = c.Fetch(seqset,
				[]imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchRFC822Size},
				chanMessage); err != nil {
				// 【实践经验】这里遇到过的err信息是：ENVELOPE doesn't contain 10 fields
				// 原因是对方发送的邮件格式不规范，解析失败
				// 相关的issue: https://github.com/emersion/go-imap/issues/143
				log.Println(seqset, err)
			}
		}()

		message := <-chanMessage
		if message == nil {
			log.Println("Server didn't returned message")
			continue
		}
		fmt.Printf("%v: %v bytes, flags=%v Envelope=%v\n", message.SeqNum, message.Size, message.Flags, message.Envelope)

		if strings.HasPrefix(message.Envelope.Subject, "subject") {
			chanMsg := make(chan *imap.Message, 1)
			go func() {
				// 这里是第二次fetch, 获取邮件MIME内容
				if err = c.Fetch(seqset,
					[]imap.FetchItem{imap.FetchRFC822},
					chanMsg); err != nil {
					log.Println(seqset, err)
				}
			}()

			msg := <-chanMsg
			if msg == nil {
				log.Println("Server didn't returned message")
			}

			section := &s
			r := msg.GetBody(section)
			if r == nil {
				log.Fatal("Server didn't returned message body")
			}

			// Create a new mail reader
			// 创建邮件阅读器
			mr, err := mail.CreateReader(r)
			if err != nil {
				log.Fatal(err)
			}

			// Process each message's part
			// 处理消息体的每个part
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				}

				switch h := p.Header.(type) {
				case *mail.InlineHeader:
					// This is the message's text (can be plain-text or HTML)
					// 获取正文内容, text或者html
					b, _ := ioutil.ReadAll(p.Body)
					log.Println("Got text: ", string(b))
				case *mail.AttachmentHeader:
					// This is an attachment
					// 下载附件
					filename, err := h.Filename()
					if err != nil {
						log.Fatal(err)
					}
					if filename != "" {
						log.Println("Got attachment: ", filename)
						b, _ := ioutil.ReadAll(p.Body)
						file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
						defer file.Close()
						n, err := file.Write(b)
						if err != nil {
							fmt.Println("写入文件异常", err.Error())
						} else {
							fmt.Println("写入Ok：", n)
						}
					}
				}
				fmt.Printf("已找到满足需求的邮件")
				return
			}
		}
	}
}

func pop(list *[]uint32) uint32 {
	length := len(*list)
	lastEle := (*list)[length-1]
	*list = (*list)[:length-1]
	return lastEle
}
