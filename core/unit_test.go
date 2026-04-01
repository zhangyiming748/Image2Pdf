package core

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// go test -v -run TestMagick
func TestMagick(t *testing.T) {
	Img2PdfInFolder("C:\\Users\\zen\\Downloads\\文档", false)
	Img2PdfInRoot("C:\\Users\\zen\\Downloads", false)
}

// go test -v -timeout 100m -run TestPdf2Img
func TestPdf2Img(t *testing.T) {
	Pdf2Img("P:\\Users\\Public\\Documents\\2025档案纸电一致\\基层党委（总支）缴纳党费明细表 打印完")
}
