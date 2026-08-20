package core

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

const (
	// 黑线的最小竖向连续长度(像素)。
	// 只有竖向延伸达到该长度的深色结构才被认为是"线"并接受中值滤波,
	// 文字笔画的竖向长度远小于该值,因此不会被误伤。
	// 600DPI 下 A4 页高约 7000 像素,扫描折痕/阴影线通常贯穿大半个页面;
	// 若残留较短的黑线可调小该值,若有内容误伤则调大。
	lineMinRun = 300
	// 横向中值滤波窗口宽度,可消除宽度不超过窗口一半的竖线
	medianWindow = 9
)

func checkMagick() {
	// 检查系统中是否存在magick命令
	_, err := exec.LookPath("magick")
	if err != nil {
		switch runtime.GOOS {
		case "darwin":
			log.Println("https://imagemagick.org/script/download.php#macos")
		case "windows":
			log.Println("https://imagemagick.org/script/download.php#windows")
		case "linux":
			log.Println("https://imagemagick.org/script/download.php#linux")
		default:
			log.Println("请检查你的操作系统是否支持ImageMagick")
		}
		log.Fatal("系统中未找到magick命令,请先访问https://imagemagick.org/script/download.php安装ImageMagick")
	}
	log.Println("检测到magick命令,程序正常启动")
}

/*
给定一组 jpg/png 图片路径，生成一个 pdf 文件
magick convert /path/to/image1.jpg /path/to/image2.jpg /path/to/image3.jpg output.pdf
转换前定向消除扫描竖向细线并做白阈值归一,将纸张底灰等近白区域合并为纯白
*/
func Img2Pdf(files []string, dst string, compress bool) error {
	checkMagick()
	if len(files) == 0 {
		log.Fatal("没有提供图片文件!")
	}
	// 逐页去除竖向黑线:遮罩合成(-composite)要求每条命令只处理一张图,
	// 因此先在临时目录生成净化后的页面,再统一合成 pdf
	tmpDir, err := os.MkdirTemp("", "i2p-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	var args []string
	for i, file := range files {
		cleaned := filepath.Join(tmpDir, fmt.Sprintf("page_%05d.png", i))
		if err = removeVerticalLines(file, cleaned); err != nil {
			return err
		}
		args = append(args, cleaned)
	}
	if compress {
		args = append(args, "-quality", "85")
		args = append(args, "-compress", "JPEG")
	} else {
		args = append(args, "-compress", "None")
	}
	args = append(args, "-auto-orient")
	args = append(args, dst)
	cmd := exec.Command("magick", args...)
	log.Printf("执行命令:%v\n", cmd.String())
	b, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	log.Printf("执行结果:%v\n", string(b))
	return nil
}

/*
只消除图中的竖向黑线,不触碰正常文字:
 1. 原图 +clone 一份做横向中值滤波(窗口宽度内的竖线被背景取代),作为"净化版";
 2. 再生成一张遮罩:灰度化取反后二值化,所有深色像素(线+文字)变白;
    用 1xN 核做竖向腐蚀,只有竖向连续长度 ≥lineMinRun 的结构能存活,
    文字笔画远短于该长度全部消失,再竖向膨胀恢复线的完整长度、横向膨胀覆盖线宽;
 3. -composite 按遮罩把净化版合成回原图:遮罩白色处(长线)取净化结果,其余像素原样保留。

等价命令:
magick in.jpg \( +clone -statistic Median 9x1 \) \( +clone -grayscale Rec709Luminance -negate -threshold 25% -morphology Erode 1x300 -morphology Dilate 1x300 -morphology Dilate 9x1 \) -composite out.png
*/
func removeVerticalLines(src, dst string) error {
	verticalKernel := fmt.Sprintf("1x%d", lineMinRun)
	medianKernel := fmt.Sprintf("%dx1", medianWindow)
	args := []string{src,
		// 净化版:横向中值滤波消除竖线
		"(", "+clone", "-statistic", "Median", medianKernel, ")",
		// 遮罩:仅保留超长竖向深色结构
		"(", "+clone",
		"-grayscale", "Rec709Luminance",
		"-negate",
		"-threshold", "25%",
		"-morphology", "Erode", verticalKernel,
		"-morphology", "Dilate", verticalKernel,
		"-morphology", "Dilate", medianKernel,
		")",
		"-composite",
		// 白阈值归一:各通道亮度 ≥85% 的近白底灰合并为纯白,
		// 此时深色文字已被遮罩保护过,不受中值滤波影响
		"-white-threshold", "85%",
		dst,
	}
	cmd := exec.Command("magick", args...)
	log.Printf("执行命令:%v\n", cmd.String())
	b, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("去除竖线失败 %s: %w, %s", src, err, string(b))
	}
	return nil
}

/*
给定一个文件夹的绝对路径
路径下包含的全部图片文件转换成一个pdf文件
并且保存到同一个文件夹下 而且与文件夹同名
*/
func Img2PdfInFolder(srtDir string, compress bool) error {
	imgFiles := finder.FindAllImagesInRoot(srtDir)
	if len(imgFiles) == 0 {
		log.Println("没有找到图片文件!")
		return nil
	}
	log.Printf("找到的图片文件:%v\n", imgFiles)
	baseName := filepath.Base(srtDir)
	//baseName = fmt.Sprintf("%v 共%v页", baseName, len(imgFiles))
	pdfName := strings.Join([]string{baseName, "pdf"}, ".")
	pdfPath := filepath.Join(srtDir, pdfName)
	log.Printf("作为生成pdf的文件名:%v\n", pdfPath)
	return Img2Pdf(imgFiles, pdfPath, compress)
}

/*
给定一个根文件夹的绝对路径
路径下包含多个子文件夹
每个子文件夹下是同一组图片
全部图片文件转换成一个 pdf 文件
并且保存到同一个文件夹下 而且与文件夹同名
*/
func Img2PdfInRoot(root string, compress bool) {
	folders := finder.FindAllFolders(root)
	for _, folder := range folders {
		if err := Img2PdfInFolder(folder, compress); err != nil {
			log.Println(err)
			continue
		}
	}
}
