package core

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zhangyiming748/finder"
)

func Pdf2Img(root string) {
	folders := finder.FindAllFolders(root)
	for _, folder := range folders {
		if hasPdfInFolder(folder) {
			Pdf2ImgInFolder(folder)
		}
	}
}

func hasPdfInFolder(root string) bool {
	files := finder.FindAllFiles(root)
	for _, file := range files {
		if strings.HasSuffix(file, ".pdf") {
			return true
		}
	}
	return false

}
func Pdf2ImgInFolder(pdf string) {
	files := finder.FindAllFiles(pdf)
	var pdfFile string
	for _, file := range files {
		if strings.HasSuffix(file, ".pdf") {
			pdfFile = file
			break
		}
	}

	if pdfFile == "" {
		log.Println("未找到 PDF 文件!")
		return
	}

	log.Printf("开始转换 PDF 文件:%v\n", pdfFile)

	// 使用 ImageMagick 将 PDF 转换为图片
	// magick -density 600 input.pdf output.png
	outputPattern := filepath.Join(filepath.Dir(pdfFile), "Scan_%d.jpg")
	args := []string{}
	args = append(args, "-density", "300")
	args = append(args, pdfFile)
	args = append(args, outputPattern)
	cmd := exec.Command("magick", args...)
	log.Printf("执行命令:%v\n", cmd.String())

	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("转换失败:", string(b))
		log.Fatal(err)
	}

	log.Printf("转换完成，结果:%v\n", string(b))
}
