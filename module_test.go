package image2pdf

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

//go test -v -run TestMagick
func TestMagick(t *testing.T) {
	Img2PdfInFolder("C:\\Users\\zen\\Downloads\\文档")
}
