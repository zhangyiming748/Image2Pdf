package decode

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zhangyiming748/finder"
	"runtime"
)

func Pdf2Img(root string) {
	path, err := checkGhostscript()
	if err != nil {
		fmt.Printf("未找到 Ghostscript。请先安装 Ghostscript https://www.ghostscript.com/releases/gsdnld.html 并确保它在 PATH 中")
		log.Fatal(err)
	}
	fmt.Println("Ghostscript 路径可用:", path)

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

func checkGhostscript() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// Windows 通常是 gswin64c.exe（推荐控制台版），也可以尝试 gs.exe（部分安装会有）
		if out, err := getGhostscriptVersion("gswin64c.exe"); err == nil {
			return out, nil
		} else {
			fmt.Printf("未找到 Ghostscript。请先安装 Ghostscript 并确保它在 PATH 中")
			return out, err
		}
	case "darwin", "linux":
		// macOS 和 Linux 统一使用 gs
		if out, err := getGhostscriptVersion("gs"); err == nil {
			return out, nil
		} else {
			fmt.Printf("未找到 Ghostscript。请先安装 Ghostscript 并确保它在 PATH 中")
			return out, err
		}
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

// 获取 Ghostscript 版本（推荐额外验证）
func getGhostscriptVersion(cmdName string) (string, error) {
	out, err := exec.Command(cmdName, "--version").Output()
	if err != nil {
		// 部分老版本可能用 -v 或 -h
		out, err = exec.Command(cmdName, "-v").Output()
		if err != nil {
			return "", err
		}
	}
	return string(out), nil
}
