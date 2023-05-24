package proto_emicklei

import (
	"os"
)

func MustCloseFile(f *os.File) {
	Must(f.Close())
}

func MustWriteToFile(filename, content string) {
	f, err := os.Create(filename)
	Must(err)
	defer MustCloseFile(f)
	_, err = f.WriteString(content)
	Must(err)
}

func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
