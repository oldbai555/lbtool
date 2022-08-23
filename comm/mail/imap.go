package mail

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap-id"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/oldbai555/lb/comm"
	"github.com/oldbai555/lb/log"
	"io"
	"io/ioutil"
	"sort"
	"time"
)

const (
	FlagsAll      = "ALL"      // 获取已读
	FlagsNEW      = "NEW"      // 获取未读
	FlagsSEEN     = "SEEN"     // 邮件已读
	FlagsANSWERED = "ANSWERED" // 邮件已回复
	FlagsFLAGGED  = "ANSWERED" // 邮件标记为紧急或者特别注意
	FlagsDELETED  = "DELETED"  // 邮件为删除状态。
	FlagsDRAFT    = "DRAFT"    // 邮件未写完（标记为草稿状态）。

	DefaultSelectBox        = "INBOX"
	DefaultReadSize  uint32 = 10
)

// NewImapClient 创建IMAP客户端
func NewImapClient(connectType uint32) (*client.Client, error) {
	conn, err := GetConnect(connectType)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	// 【字符集】  处理us-ascii和utf-8以外的字符集(例如gbk,gb2313等)时,
	//  需要加上这行代码。
	// 【参考】 https://github.com/emersion/go-imap/wiki/Charset-handling
	imap.CharsetReader = charset.Reader

	log.Infof("Connecting to server...")

	// 连接邮件服务器
	c, err := client.DialTLS(fmt.Sprintf("%s:%d", conn.ImapHost, conn.ImapPost), nil)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	log.Infof("Connected")
	idClient := id.NewClient(c)
	_, err = idClient.ID(
		id.ID{id.FieldName: "IMAPClient", id.FieldVersion: "1.2.0"}, // 随便定义申明自己身份就行
	)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}
	return c, nil
}

// Login 创建IMAP客户端
func Login(u *user, c *client.Client) error {
	// 【字符集】  处理us-ascii和utf-8以外的字符集(例如gbk,gb2313等)时,
	//  需要加上这行代码。
	// 【参考】 https://github.com/emersion/go-imap/wiki/Charset-handling
	imap.CharsetReader = charset.Reader

	// 使用账号密码登录
	if err := c.Login(u.Username, u.Password); err != nil {
		log.Errorf("err is %v", err)
		return err
	}

	log.Infof("Logged in")

	return nil
}

type ReadMailReq struct {
	// Conn 邮箱连接信息
	C *client.Client `json:"c"`
	// U 收信用户
	U *user `json:"u"`
	// ReadSize 指定查询信件数量,默认十封 -> DefaultReadSize
	ReadSize uint32 `json:"read_size"`
	// SelectBox 指定文件夹名字,默认收件箱 -> DefaultSelectBox
	SelectBox string `json:"select_box"`
	// Flags 构造查询条件
	Flags []string `json:"flags"`
	// FetchItemList 抓取的信件内容，例如 邮件头，邮件标志，邮件大小等信息
	FetchItemList []imap.FetchItem `json:"fetch_item_list"`
}

// ReadMail 读取指定前 x 封邮件，默认十封
func ReadMail(req *ReadMailReq, f func(msgList []*ReceiveMailMessage) error) error {
	if req.U == nil {
		log.Errorf("ReadMail request U not found")
		return fmt.Errorf("req.U is nil")
	}

	err := Login(req.U, req.C)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}

	// Don't forget to logout
	defer func() {
		logoutErr := req.C.Logout()
		if logoutErr != nil {
			log.Errorf("err is %v", logoutErr)
		}
	}()

	// 选择收件箱
	if req.SelectBox == "" {
		req.SelectBox = DefaultSelectBox
	}
	_, err = req.C.Select(req.SelectBox, false)
	if err != nil {
		log.Errorf("err is %v", err)
		return err
	}

	// 搜索条件实例对象
	criteria := imap.NewSearchCriteria()
	// See RFC 3501 section 6.4.4 for a list of searching criteria.
	if len(req.Flags) == 0 {
		// 默认拿所有
		req.Flags = append(req.Flags, FlagsAll)
	}
	criteria.WithFlags = req.Flags
	seqNums, _ := req.C.Search(criteria)
	var msgList MessageList
	if len(seqNums) == 0 {
		log.Warnf("not seqNums: %v", seqNums)
		err = f(msgList)
		if err != nil {
			log.Errorf("err is %v", err)
			return err
		}
		return nil
	}

	// 获取邮件的数量以及设置邮件头的标志
	if req.ReadSize == 0 {
		req.ReadSize = DefaultReadSize
	}
	if req.ReadSize < uint32(len(seqNums)) {
		start := uint32(len(seqNums)) - req.ReadSize
		seqNums = seqNums[start:]
	}
	seqset := new(imap.SeqSet)
	for _, num := range seqNums {
		seqset.AddNum(num)
	}
	if len(req.FetchItemList) == 0 {
		// 默认只抓取邮件头，邮件标志，邮件大小，正文附件等信息
		req.FetchItemList = append(req.FetchItemList, imap.FetchEnvelope, imap.FetchFlags, imap.FetchRFC822Size, imap.FetchRFC822)
	}

	// 开始去拿邮件列表
	chanMessage := make(chan *imap.Message, DefaultReadSize)
	go fetchMail(req.C, seqset, chanMessage, req.FetchItemList)

	var getAdressList = func(imapAdrList []*imap.Address) (list []*Address) {
		for _, address := range imapAdrList {
			list = append(list, &Address{
				PersonalName: address.PersonalName,
				AtDomainList: address.AtDomainList,
				MailboxName:  address.MailboxName,
				HostName:     address.HostName,
			})
		}
		return
	}

	for msg := range chanMessage {
		if msg.Envelope == nil {
			log.Warnf("val.Envelope is nil ,%v", msg)
			continue
		}
		envelope := msg.Envelope
		receMsg := &ReceiveMailMessage{
			SeqNum:    msg.SeqNum,
			Size:      msg.Size,
			MessageId: envelope.MessageId,
			Envelope: &Envelope{
				Date:        uint32(envelope.Date.Unix()),
				Subject:     envelope.Subject,
				InReplyTo:   envelope.InReplyTo,
				MessageId:   envelope.MessageId,
				FromList:    getAdressList(envelope.From),
				SenderList:  getAdressList(envelope.Sender),
				ReplyToList: getAdressList(envelope.ReplyTo),
				ToList:      getAdressList(envelope.To),
				CcList:      getAdressList(envelope.Cc),
				BccList:     getAdressList(envelope.Bcc),
			},
		}

		// 解析文本内容和附件
		section := &imap.BodySectionName{}
		r := msg.GetBody(section)
		if r == nil {
			log.Warnf("Server didn't returned message body")
			continue
		}

		// Create a new mail reader
		// 创建邮件阅读器
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		// Process each message's part
		// 处理消息体的每个part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Errorf("err is %v", err)
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				// 获取正文内容, text或者html
				b, _ := ioutil.ReadAll(p.Body)
				// log.Debugf("Got text: ", string(b))
				receMsg.MsgBodyList = append(receMsg.MsgBodyList, b)
			case *mail.AttachmentHeader:
				// This is an attachment
				// 下载附件
				filename, err := h.Filename()
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				if filename != "" {
					log.Debugf("Got attachment: %s", filename)
					b, _ := ioutil.ReadAll(p.Body)
					receMsg.AttachList = append(receMsg.AttachList, &Attach{
						Buf:      b,
						FileName: filename,
					})
				}
			}
		}
		log.Debugf("find end")
		err = setSeenFlag(req.C, msg.SeqNum)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}
		log.Debugf("find end , msg-id is %s , seqNum is %d", msg.Envelope.MessageId, msg.SeqNum)
		msgList = append(msgList, receMsg)
	}

	sort.Sort(msgList)

	err = f(msgList)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

//fetchMail
func fetchMail(c *client.Client, seqset *imap.SeqSet, chanMsg chan *imap.Message, fetchItemList []imap.FetchItem) {
	// 第一次fetch, 只抓取邮件头，邮件标志，邮件大小等信息，执行速度快
	if err := c.Fetch(seqset,
		fetchItemList,
		chanMsg); err != nil {
		// 【实践经验】这里遇到过的err信息是：ENVELOPE doesn't contain 10 fields
		// 原因是对方发送的邮件格式不规范，解析失败
		// 相关的issue: https://github.com/emersion/go-imap/issues/143
		log.Errorf("err is %v ，seqset is %v", err, seqset)
	}
}

// GetAllMailBoxesList 获取所有邮件文件夹
func GetAllMailBoxesList(c *client.Client) ([]*imap.MailboxInfo, error) {
	var list []*imap.MailboxInfo
	mailboxs := make(chan *imap.MailboxInfo, 10)
	err := c.List("", "*", mailboxs)
	if err != nil {
		log.Errorf("err is %v", err)
		return nil, err
	}

	for m := range mailboxs {
		log.Debugf("m.Name is %s", m.Name)
		list = append(list, m)
	}
	return list, nil
}

// setSeenFlag 标记邮件已读
func setSeenFlag(c *client.Client, seqNum uint32) error {
	if c == nil {
		log.Errorf("client is nil")
		return fmt.Errorf("[setSeenFlag] para c is nil")
	}
	if seqNum < 0 {
		log.Errorf("seqNum is %zero", seqNum)
		return fmt.Errorf("[setSeenFlag] para msgSeq is invalid")
	}

	seenSet, err := imap.ParseSeqSet(fmt.Sprintf("%d", seqNum))
	if err != nil {
		log.Errorf("err is %v", err)
		return fmt.Errorf("err is %v", err)
	}
	done := make(chan error, 1)

	go func() {
		done <- c.Store(seenSet, imap.AddFlags, []interface{}{imap.SeenFlag}, nil)
	}()

	if err != nil {
		log.Errorf("err is %v", err)
		return fmt.Errorf("err is %v", err)
	}
	// there can be bug set time out 2s
	select {
	case newErr := <-done:
		if newErr != nil {
			log.Errorf("err is %v", newErr)
			return fmt.Errorf("err is %v", newErr)
		}

	case <-time.After(2 * time.Second):
		log.Errorf("set flag time out(2s)")
		return fmt.Errorf("set flag time out(2s)")
	}
	return nil
}

// MessageList 消息对象列表
type MessageList []*ReceiveMailMessage

// Len 重写 Len() 方法
func (a MessageList) Len() int {
	return len(a)
}

// Swap 重写 Swap() 方法
func (a MessageList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Less 重写 Less() 方法
func (a MessageList) Less(i, j int) bool {
	return comm.ReflectCompareFieldDesc(a[j], a[i], "SeqNum")
}
