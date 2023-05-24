package phone2code

import "testing"

// https://pinfire.feishu.cn/sheets/shtcng6mbYtVDzTTcdjKxw78Ovb
func Test_clearText(t *testing.T) {
	t.Log(GetPhoneCode("+39 339 990 7440"))
}
