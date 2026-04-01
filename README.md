# Image2Pdf

一个用于将一组图片合并生成单个 PDF 文件的命令行工具，同时也支持将 PDF 文件解码为图片。

## 依赖

[ImageMagick](https://imagemagick.org/script/download.php)

在运行此工具之前，请先安装 ImageMagick 并确保 `magick` 命令在系统 PATH 中可用。

### 安装 ImageMagick

- **Windows**: 访问 [ImageMagick Windows 下载页面](https://imagemagick.org/script/download.php#windows)
- **macOS**: 访问 [ImageMagick macOS 下载页面](https://imagemagick.org/script/download.php#macos) 或使用 Homebrew: `brew install imagemagick`
- **Linux**: 访问 [ImageMagick Linux 下载页面](https://imagemagick.org/script/download.php#linux) 或使用包管理器: `sudo apt-get install imagemagick`

## 安装

```bash
go build -o Image2Pdf
```

## 使用方法

### 1. single 命令 - 将单个文件夹内的图片转换为一个 PDF 文件

给定一个文件夹的绝对路径，将该文件夹下的所有图片文件转换成一个 PDF 文件，并保存到同一文件夹下，且与文件夹同名。

**参数:**
- `-d, --dir`: 包含图片的文件夹绝对路径（必填）
- `-c, --compress`: 是否压缩 PDF 文件大小，默认为 `true`（可选）

**示例:**

```bash
# 启用压缩（默认）
./Image2Pdf single -d /path/to/images
./Image2Pdf single -d /path/to/images -c

# 禁用压缩
./Image2Pdf single -d /path/to/images -c=false
# 或
./Image2Pdf single -d /path/to/images --compress=false
```

### 2. multi 命令 - 将根目录下的多个子文件夹分别转换为 PDF 文件

给定一个根文件夹的绝对路径，该路径下包含多个子文件夹，每个子文件夹下的图片会分别转换成 PDF 文件。

**参数:**
- `-d, --dir`: 包含多个子文件夹的根目录绝对路径（必填）
- `-c, --compress`: 是否压缩 PDF 文件大小，默认为 `true`（可选）

**示例:**

```bash
# 启用压缩（默认）
./Image2Pdf multi -d /path/to/root
./Image2Pdf multi -d /path/to/root -c

# 禁用压缩
./Image2Pdf multi -d /path/to/root -c=false
# 或
./Image2Pdf multi -d /path/to/root --compress=false
```

### 3. decode 命令 - 将 PDF 文件解码为图片

给定一个根目录路径，将该目录下所有包含 PDF 的文件夹中的 PDF 文件转换为图片，并保存到同一文件夹下。

**参数:**
- `-d, --dir`: 包含 PDF 文件的根目录绝对路径（默认为当前目录）

**示例:**

```bash
# 使用默认值（当前目录）
./Image2Pdf decode

# 指定目录
./Image2Pdf decode -d /path/to/pdf/folder

# 或使用长格式
./Image2Pdf decode --dir /path/to/pdf/folder
```

## 注意事项

1. 图片转 PDF 时，支持常见的图片格式（JPG、PNG 等）
2. PDF 转图片时，输出的图片文件名为 `Scan_0001.jpg`, `Scan_0002.jpg` 等格式
3. 压缩选项可以减小生成的 PDF 文件大小，但可能会略微降低图片质量