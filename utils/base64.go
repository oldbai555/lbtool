package utils

const (
	HTMLSrcPrefix = "data:image/png;base64," // html直接展示base64的图片
)

const (
	HeaderJPG  = "JPG"
	HeaderPNG  = "PNG"
	HeaderBMP  = "BMP"
	HeaderTIFF = "TIFF"
	HeaderPDF  = "PDF"
	HeaderOFD  = "OFD"
)

var headerMap = map[string]string{
	"/9j": HeaderJPG,
	"iVB": HeaderPNG,
	"Qk0": HeaderBMP,
	"SUk": HeaderTIFF,
	"JVB": HeaderPDF,
	"UEs": HeaderOFD,
}

// Base64FileHeaderMapper 将base64编码后的前三字符传入，返回文件类型
func Base64FileHeaderMapper(fileBase64 string) string {
	if len(fileBase64) > 3 {
		for sub, val := range headerMap {
			if sub == SubStr(fileBase64, 0, 3) {
				return val
			}
		}
	}
	return ""
}
