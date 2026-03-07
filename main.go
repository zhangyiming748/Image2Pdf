package main

import (
	"log"

	"Image2Pdf/core"
	"github.com/spf13/cobra"
)

func init() {
	log.SetFlags(2 | 16)
}

var (
	dir string
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
		core.Img2PdfInFolder(dir)
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
		core.Img2PdfInRoot(dir)
	},
}

func init() {
	// single 命令的参数
	singleCmd.Flags().StringVarP(&dir, "dir", "d", "", "包含图片的文件夹绝对路径")
	singleCmd.MarkFlagRequired("dir")

	// multi 命令的参数
	multiCmd.Flags().StringVarP(&dir, "dir", "d", "", "包含多个子文件夹的根目录绝对路径")
	multiCmd.MarkFlagRequired("dir")

	// 将子命令添加到根命令
	rootCmd.AddCommand(singleCmd)
	rootCmd.AddCommand(multiCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
