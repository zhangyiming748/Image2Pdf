package decode

import (
	"testing"
)

// go test -v -timeout 100m -run TestPdf2Img
func TestPdf2Img(t *testing.T) {
	Pdf2Img("P:\\Users\\Public\\Documents\\2025档案纸电一致\\基层党委（总支）缴纳党费明细表 打印完 明细表有问题\\8\\test")
}


func TestAvif2Jpg(t *testing.T) {
	Avif2Jpg("C:\\Users\\zhang\\Desktop\\2022会议纪要及合同已复核\\12-35巨龙pdf缺第一页和第九页 - 副本")
}
