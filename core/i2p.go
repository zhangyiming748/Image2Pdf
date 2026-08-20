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
clean 为 true 时转换前定向消除扫描竖向细线并做白阈值归一,将纸张底灰等近白区域合并为纯白;
为 false 时保持默认转换参数,不做任何清理
*/
func Img2Pdf(files []string, dst string, compress bool, clean bool) error {
	checkMagick()
	if len(files) == 0 {
		log.Fatal("没有提供图片文件!")
	}
	var args []string
	if clean {
		// 逐页去除竖向黑线:遮罩合成(-composite)要求每条命令只处理一张图,
		// 因此先在临时目录生成净化后的页面,再统一合成 pdf
		tmpDir, err := os.MkdirTemp("", "i2p-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)
		for i, file := range files {
			cleaned := filepath.Join(tmpDir, fmt.Sprintf("page_%05d.png", i))
			if err = removeVerticalLines(file, cleaned); err != nil {
				return err
			}
			args = append(args, cleaned)
		}
	} else {
		args = append(args, files...)
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
magick in.jpg \( +clone -statistic Median 9x1 \) \( +clone -grayscale Rec709Luminance -negate -threshold 25% -morphology Erode Rectangle:1x300 -morphology Dilate Rectangle:1x300 -morphology Dilate Rectangle:9x1 \) -composite out.png

注意:-morphology 的核字符串必须带形状前缀(如 Rectangle:1x300),裸写 1x300 会报
"unable to parse kernel string";-statistic Median 走的是几何参数解析,裸写 9x1 即可
*/
func removeVerticalLines(src, dst string) error {
	// 形态学核需带 Rectangle: 前缀;中值滤波的窗口参数则保持裸几何格式
	verticalKernel := fmt.Sprintf("Rectangle:1x%d", lineMinRun)
	morphDilateKernel := fmt.Sprintf("Rectangle:%dx1", medianWindow)
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
		"-morphology", "Dilate", morphDilateKernel,
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
给定一个目录的绝对路径,将该目录下(包括所有子目录)的全部图片逐一深度清理:
定向去除扫描竖向黑线并将底灰归一为纯白。
每张图先在原目录生成同后缀的临时文件,清理成功后删除原始图片,
再把临时文件重命名为原始文件名,路径和位置保持不变;单张失败时跳过继续处理其余图片
*/
func CleanImagesInDir(root string) {
	checkMagick()
	files := finder.FindAllImages(root)
	if len(files) == 0 {
		log.Println("没有找到图片文件!")
		return
	}
	log.Printf("找到的图片文件:%v\n", files)
	for _, file := range files {
		if err := cleanImageInPlace(file); err != nil {
			log.Printf("清理失败 %s: %v\n", file, err)
			continue
		}
		log.Printf("清理完成: %s\n", file)
	}
}

/*
原地清理单张图片:在与原图同目录下创建同后缀临时文件承接清理结果,
成功后删除原图并将临时文件重命名为原文件名,保留原文件权限;
任一步骤失败时清理临时文件并返回错误,原图不受影响
*/
func cleanImageInPlace(src string) error {
	ext := filepath.Ext(src)
	tmp, err := os.CreateTemp(filepath.Dir(src), ".i2p-clean-*"+ext)
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	tmp.Close()
	if err = removeVerticalLines(src, tmpName); err != nil {
		os.Remove(tmpName)
		return err
	}
	info, err := os.Stat(src)
	if err != nil {
		os.Remove(tmpName)
		return err
	}
	if err = os.Chmod(tmpName, info.Mode().Perm()); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err = os.Remove(src); err != nil {
		os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, src)
}

/*
给定一个文件夹的绝对路径
路径下包含的全部图片文件转换成一个pdf文件
并且保存到同一个文件夹下 而且与文件夹同名
*/
func Img2PdfInFolder(srtDir string, compress bool, clean bool) error {
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
	return Img2Pdf(imgFiles, pdfPath, compress, clean)
}

/*
给定一个根文件夹的绝对路径
路径下包含多个子文件夹
每个子文件夹下是同一组图片
全部图片文件转换成一个 pdf 文件
并且保存到同一个文件夹下 而且与文件夹同名
*/
func Img2PdfInRoot(root string, compress bool, clean bool) {
	folders := finder.FindAllFolders(root)
	for _, folder := range folders {
		if err := Img2PdfInFolder(folder, compress, clean); err != nil {
			log.Println(err)
			continue
		}
	}
}
