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
	DPI = 200
)



/*
给定一组 jpg/png 图片路径，生成一个 pdf 文件
*/
func Img2Pdf(files []string, dst string) {
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
			log.Printf("打开图片失败 %s: %v", imgPath, err)
			continue
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
			log.Printf("不支持的图片格式 %s: %s", imgPath, ext)
			file.Close()
			continue
		}
		file.Close()

		if err != nil {
			log.Printf("解码失败 %s: %v", imgPath, err)
			continue
		}

		// 计算宽高（按 DPI 换算成 mm）
		bounds := img.Bounds()
		widthPx := float64(bounds.Dx())
		heightPx := float64(bounds.Dy())
		widthMM := widthPx * 25.4 / DPI
		heightMM := heightPx * 25.4 / DPI

		// A4 页面尺寸 210×297 mm
		pageW, pageH := pdf.GetPageSize()

		// 自动旋转：如果图片更宽，就横放（Landscape）
		if widthMM > heightMM && widthMM > pageW {
			pdf.AddPageFormat("L", gofpdf.SizeType{Wd: pageH, Ht: pageW}) // 横版
		} else {
			pdf.AddPage()
		}

		// 居中 + 最大化适应页面（留 10mm 边距）
		pdf.ImageOptions(imgPath, 10, 10, pageW-20, pageH-20, false, gofpdf.ImageOptions{
			ReadDpi: true,
		}, 0, "")
	}

	// 保存 PDF
	err := pdf.OutputFileAndClose(dst)
	if err != nil {
		log.Fatal("保存 PDF 失败:", err)
	}

	fmt.Printf("已生成 %s/%d\n", dst, len(files))
}
