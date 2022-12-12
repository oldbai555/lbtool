package warning

import "testing"

func TestReportToFeishu(t *testing.T) {
	ReportToFeishu("1", "2", "d766aaa2-a20d-4fdb-9d48-ea596f6b1aba", "all")
}
