package image2pdf

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/phpdave11/gofpdf"
)

const (
	DPI float64 = 600
)
const (
	CLOCKWISE        float64 = 90  //顺时针旋转90度
	KEEP             float64 = 0
	COUNTERCLOCKWISE float64 = -90 //逆时针旋转90度
)

/*
给定一组 jpg/png 图片路径，生成一个 pdf 文件
*/
func Img2Pdf(files []string, dst string, rotate float64) error {
	if len(files) == 0 {
		log.Fatal("没有提供图片文件!")
	}

	// 创建 PDF
	pdf := gofpdf.New("P", "mm", "A4", "") // 竖版 A4
	pdf.SetCompression(true)

	for i, imgPath := range files {
		fmt.Printf("[%d/%d] 正在添加: %s\n", i+1, len(files), filepath.Base(imgPath))
		// 打开图片，获取真实尺寸
		file, err := os.Open(imgPath)
		if err != nil {
			return fmt.Errorf("打开图片失败 %s: %v", imgPath, err)
		}

		// 根据扩展名选择解码器
		var img image.Image
		ext := filepath.Ext(imgPath)
		switch ext {
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(file)
		case ".png":
			img, err = png.Decode(file)
		default:
			file.Close()
			return fmt.Errorf("不支持的图片格式 %s: %s", imgPath, ext)
		}
		file.Close()

		if err != nil {
			return fmt.Errorf("解码失败 %s: %v", imgPath, err)
		}

		// 计算宽高（按 DPI 换算成 mm）
		bounds := img.Bounds()
		widthPx := float64(bounds.Dx())
		heightPx := float64(bounds.Dy())
		widthMM := widthPx * 25.4 / DPI
		heightMM := heightPx * 25.4 / DPI

		// A4 页面尺寸 210×297 mm
		pageW, pageH := pdf.GetPageSize()

		// 总是使用纵向页面
		pdf.AddPage()

		// 显示图片，保持纵横比，不拉伸变形，无边距
		availableW := pageW
		availableH := pageH

		// 计算缩放比例以保持纵横比
		scaleW := availableW / widthMM
		scaleH := availableH / heightMM
		scale := scaleW
		if scaleH < scaleW {
			scale = scaleH
		}

		// 计算实际绘制尺寸
		drawW := widthMM * scale
		drawH := heightMM * scale

		// 计算居中位置
		x := (pageW - drawW) / 2
		y := (pageH - drawH) / 2

		// 应用旋转变换
		if rotate != KEEP {
			// 开始变换
			pdf.TransformBegin()
			// 将旋转中心设为图片中心点
			centerX := x + drawW/2
			centerY := y + drawH/2
			// 旋转指定角度
			pdf.TransformRotate(rotate, centerX, centerY)
			// 绘制图片
			pdf.ImageOptions(imgPath, x, y, drawW, drawH, false, gofpdf.ImageOptions{
				ReadDpi: true,
			}, 0, "")
			// 结束变换
			pdf.TransformEnd()
		} else {
			// 直接绘制图片
			pdf.ImageOptions(imgPath, x, y, drawW, drawH, false, gofpdf.ImageOptions{
				ReadDpi: true,
			}, 0, "")
		}
	}
	// 保存 PDF
	err := pdf.OutputFileAndClose(dst)
	if err != nil {
		return fmt.Errorf("保存 PDF 失败:%v", err)
	}

	fmt.Printf("已生成 %s\t共%d页\n", dst, len(files))
	return nil
}