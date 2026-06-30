package decode

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/zhangyiming748/finder"
)

func Avif2Jpg(root string) {
	files := finder.FindAllFiles(root)
	for _, file := range files {
		if strings.HasSuffix(file, ".avif") {
			// convert avif to jpg using ImageMagick
			jpgFile := strings.TrimSuffix(file, ".avif") + ".jpg"
			args := []string{file, jpgFile}
			cmd := exec.Command("magick", args...)
			log.Printf("执行命令: %v\n", cmd.String())

			b, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("转换失败 %s: %s\n", file, string(b))
				continue
			}

			log.Printf("转换成功: %s -> %s\n", file, jpgFile)
			os.Remove(file)
		}
	}
}
