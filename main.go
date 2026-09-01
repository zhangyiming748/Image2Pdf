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
	clean    bool
	version  = "dev" // 默认版本号,构建时通过 -ldflags 注入
)

var rootCmd = &cobra.Command{
	Use:     "image2pdf",
	Short:   "图片转 PDF 工具",
	Long:    "Image2Pdf 是一个用于将一组图片合并生成单个 PDF 文件的命令行工具",
	Version: version,
}

var singleCmd = &cobra.Command{
	Use:   "single",
	Short: "将单个文件夹内的图片转换为一个 PDF 文件",
	Long:  "给定一个文件夹的绝对路径，将该文件夹下的所有图片文件转换成一个 PDF 文件，并保存到同一文件夹下，且与文件夹同名",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			log.Fatal("请提供文件夹路径参数 -d 或 --dir")
		}
		if err := core.Img2PdfInFolder(dir, compress, clean); err != nil {
			log.Fatal(err)
		}
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
		core.Img2PdfInRoot(dir, compress, clean)
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

var a2jCmd = &cobra.Command{
	Use:   "a2j",
	Short: "将 AVIF 格式图片转换为 JPG 格式",
	Long:  "给定一个根目录路径，将该目录下所有 .avif 文件转换为同名的 .jpg 文件",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			log.Fatal("请提供根目录路径参数 -d 或 --dir")
		}
		decode.Avif2Jpg(dir)
	},
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "清理目录（包括子目录）下的所有图片",
	Long:  "给定一个目录路径，将该目录及所有子目录下的图片逐一深度清理（定向去除扫描竖向黑线并将底灰归一为纯白），清理完成后替换原始图片，文件名、路径和位置保持不变",
	Run: func(cmd *cobra.Command, args []string) {
		if dir == "" {
			log.Fatal("请提供目录路径参数 -d 或 --dir")
		}
		core.CleanImagesInDir(dir)
	},
}

func init() {
	// single 命令的参数
	singleCmd.Flags().StringVarP(&dir, "dir", "d", ".", "包含图片的文件夹绝对路径")
	singleCmd.Flags().BoolVarP(&compress, "compress", "c", true, "是否压缩 PDF 文件大小 (默认 true)")
	singleCmd.Flags().BoolVarP(&clean, "clean", "C", false, "是否启用清理:定向去除扫描竖向黑线并将底灰归一为纯白 (默认 false)")
	singleCmd.MarkFlagRequired("dir")
	/*
		# 默认转换,不做清理
		./Image2Pdf single -d /path/to/images

		# 启用清理:显式写 -C/--clean 或附带 =true
		./Image2Pdf single -d /path/to/images -C
		./Image2Pdf single -d /path/to/images --clean=true

		# 禁用压缩
		./Image2Pdf single -d /path/to/images -c=false
		# 或
		./Image2Pdf single -d /path/to/images --compress=false
	*/

	// multi 命令的参数
	multiCmd.Flags().StringVarP(&dir, "dir", "d", ".", "包含多个子文件夹的根目录绝对路径")
	multiCmd.Flags().BoolVarP(&compress, "compress", "c", true, "是否压缩 PDF 文件大小 (默认 true)")
	multiCmd.Flags().BoolVarP(&clean, "clean", "C", false, "是否启用清理:定向去除扫描竖向黑线并将底灰归一为纯白 (默认 false)")
	/*
		# 默认转换,不做清理
		./Image2Pdf multi -d /path/to/root

		# 启用清理:显式写 -C/--clean 或附带 =true
		./Image2Pdf multi -d /path/to/root -C
		./Image2Pdf multi -d /path/to/root --clean=true

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

	// a2j 命令的参数
	a2jCmd.Flags().StringVarP(&dir, "dir", "d", ".", "包含 AVIF 文件的根目录绝对路径（默认为当前目录）")
	/*
		# 使用默认值（当前目录）
		./Image2Pdf a2j

		# 指定目录
		./Image2Pdf a2j -d /path/to/avif/folder

		# 或使用长格式
		./Image2Pdf a2j --dir /path/to/avif/folder
	*/

	// clean 命令的参数
	cleanCmd.Flags().StringVarP(&dir, "dir", "d", ".", "包含待清理图片的目录绝对路径（默认为当前目录）")
	/*
		# 使用默认值（当前目录）
		./Image2Pdf clean

		# 指定目录,递归清理该目录及所有子目录下的图片
		./Image2Pdf clean -d /path/to/images

		# 或使用长格式
		./Image2Pdf clean --dir /path/to/images
	*/

	// 将子命令添加到根命令
	rootCmd.AddCommand(singleCmd)
	rootCmd.AddCommand(multiCmd)
	rootCmd.AddCommand(decodeCmd)
	rootCmd.AddCommand(a2jCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  "显示当前 Image2Pdf 的版本号",
	Run: func(cmd *cobra.Command, args []string) {
		println("Image2Pdf version", version)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
