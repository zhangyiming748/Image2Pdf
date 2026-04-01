package main

import (
	"log"

	"Image2Pdf/core"
	"Image2Pdf/decode"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(2 | 16)
}

var (
	dir      string
	compress bool
)

var rootCmd = &cobra.Command{
	Use:   "image2pdf",
	Short: "图片转 PDF 工具",
	Long:  "Image2Pdf 是一个用于将一组图片合并生成单个 PDF 文件的命令行工具",
}

var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "将单个文件夹内的图片转换为一个 PDF 文件",
	Long:  "给定一个文件夹的绝对路径，将该文件夹下的所有图片文件转换成一个 PDF 文件，并保存到同一文件夹下，且与文件夹同名",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			log.Fatal("请提供文件夹路径参数 -d 或 --dir")
		}
		core.Img2PdfInFolder(dir, compress)
	},
}

var multiCmd = &cobra.Command{
	Use:   "multi",
	Short: "将根目录下的多个子文件夹分别转换为 PDF 文件",
	Long:  "给定一个根文件夹的绝对路径，该路径下包含多个子文件夹，每个子文件夹下的图片会分别转换成 PDF 文件",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			log.Fatal("请提供根目录路径参数 -d 或 --dir")
		}
		core.Img2PdfInRoot(dir, compress)
	},
}

var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "将 PDF 文件解码为图片",
	Long:  "给定一个根目录路径，将该目录下所有包含 PDF 的文件夹中的 PDF 文件转换为图片，并保存到同一文件夹下",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			log.Fatal("请提供根目录路径参数 -d 或 --dir")
		}
		decode.Pdf2Img(dir)
	},
}

func init() {
	// single 命令的参数
	singleCmd.Flags().StringVarP(&dir, "dir", "d", ".", "包含图片的文件夹绝对路径")
	singleCmd.Flags().BoolVarP(&compress, "compress", "c", true, "是否压缩 PDF 文件大小 (默认 true)")
	singleCmd.MarkFlagRequired("dir")
	/*
	# 启用压缩（默认）
	./Image2Pdf single -d /path/to/images
	./Image2Pdf single -d /path/to/images -c

	# 禁用压缩
	./Image2Pdf single -d /path/to/images -c=false
	# 或
	./Image2Pdf single -d /path/to/images --compress=false
	*/

	// multi 命令的参数
	multiCmd.Flags().StringVarP(&dir, "dir", "d", ".", "包含多个子文件夹的根目录绝对路径")
	multiCmd.Flags().BoolVarP(&compress, "compress", "c", true, "是否压缩 PDF 文件大小 (默认 true)")
	multiCmd.MarkFlagRequired("dir")
	/*
	# 启用压缩（默认）
	./Image2Pdf multi -d /path/to/root
	./Image2Pdf multi -d /path/to/root -c

	# 禁用压缩
	./Image2Pdf multi -d /path/to/root -c=false
	# 或
	./Image2Pdf multi -d /path/to/root --compress=false
	*/

	// decode 命令的参数
	decodeCmd.Flags().StringVarP(&dir, "dir", "d", ".", "包含 PDF 文件的根目录绝对路径（默认为当前目录）")
	/*
	# 使用默认值（当前目录）
	./Image2Pdf decode

	# 指定目录
	./Image2Pdf decode -d /path/to/pdf/folder

	# 或使用长格式
	./Image2Pdf decode --dir /path/to/pdf/folder
	*/

	// 将子命令添加到根命令
	rootCmd.AddCommand(singleCmd)
	rootCmd.AddCommand(multiCmd)
	rootCmd.AddCommand(decodeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
