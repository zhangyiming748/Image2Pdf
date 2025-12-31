package image2pdf

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zhangyiming748/finder"
)

const (
	DPI float64 = 600
)
const (
	CLOCKWISE        float64 = 90 //顺时针旋转90度
	KEEP             float64 = 0
	COUNTERCLOCKWISE float64 = -90 //逆时针旋转90度
)

func init() {
	// 检查系统中是否存在magick命令
	_, err := exec.LookPath("magick")
	if err != nil {
		log.Fatal("系统中未找到magick命令,请先访问https://imagemagick.org/script/download.php安装ImageMagick")
	}
	log.Println("检测到magick命令,程序正常启动")
}

/*
给定一组 jpg/png 图片路径，生成一个 pdf 文件
magick convert /path/to/image1.jpg /path/to/image2.jpg /path/to/image3.jpg output.pdf
*/
func Img2Pdf(files []string, dst string) {
	if len(files) == 0 {
		log.Fatal("没有提供图片文件!")
	}

	//magick convert /path/to/image1.jpg /path/to/image2.jpg /path/to/image3.jpg output.pdf

	args := []string{"convert"}
	args = append(args, files...)
	args = append(args, dst)
	cmd := exec.Command("magick", args...)
	log.Printf("执行命令:%v\n", cmd.String())
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("执行结果:%v\n", string(b))

}

/*
给定一个文件夹的绝对路径
路径下包含的全部图片文件转换成一个pdf文件
并且保存到同一个文件夹下 而且与文件夹同名
*/
func Img2PdfInFolder(srtDir string) {
	imgFiles := finder.FindAllImagesInRoot(srtDir)
	if len(imgFiles) == 0 {
		log.Println("没有找到图片文件!")
		return
	}
	log.Printf("找到的图片文件:%v\n", imgFiles)
	baseName := filepath.Base(srtDir)
	baseName = fmt.Sprintf("%v 共%v页", baseName, len(imgFiles))
	pdfName := strings.Join([]string{baseName, "pdf"}, ".")
	pdfPath := filepath.Join(srtDir, pdfName)
	log.Printf("作为生成pdf的文件名:%v\n", pdfPath)
	Img2Pdf(imgFiles, pdfPath)
}

/*
给定一个根文件夹的绝对路径
路径下包含多个子文件夹
每个子文件夹下是同一组图片
全部图片文件转换成一个pdf文件
并且保存到同一个文件夹下 而且与文件夹同名
*/
func Img2PdfInRoot(root string) {
	folders := finder.FindAllFolders(root)
	for _, folder := range folders {
		Img2PdfInFolder(folder)
	}
}
