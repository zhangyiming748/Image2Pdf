# Image2Pdf

<div align="center">

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go)
![License](https://img.shields.io/github/license/zhangyiming748/Image2Pdf?style=flat-square)
![Release](https://img.shields.io/github/v/release/zhangyiming748/Image2Pdf?style=flat-square&label=release)
![Downloads](https://img.shields.io/github/downloads/zhangyiming748/Image2Pdf/total?style=flat-square)

**一个高性能的图片与 PDF 互转命令行工具**

[快速开始](#-快速开始) • [安装指南](#-安装指南) • [使用文档](#-使用文档) • [常见问题](#-常见问题)

</div>

---

## 📖 项目简介

Image2Pdf 是一个基于 Go 语言开发的跨平台命令行工具，专注于图片与 PDF 文件之间的高效转换。它利用 ImageMagick 强大的图像处理能力，提供简洁直观的命令行接口，支持批量处理、智能压缩和自动化工作流。

### ✨ 核心特性

- 🚀 **高性能处理** - 基于 Go 语言编译，执行效率高，资源占用少
- 🔄 **双向转换** - 支持图片转 PDF 和 PDF 转图片两种模式
- 📦 **批量处理** - 支持单文件夹和多文件夹批量转换
- 🗜️ **智能压缩** - 可选的 JPEG 压缩算法，平衡文件大小与质量
- 🌍 **跨平台支持** - 完美支持 Linux、macOS、Windows 三大操作系统
- 🎯 **简单易用** - 清晰的命令行接口，完善的帮助文档
- 🔧 **灵活配置** - 支持压缩开关、自定义输出路径等选项

---

## 📥 快速下载

### 从 GitHub Releases 下载

| 平台 | 架构 | 下载链接 |
|:---:|:---:|:---:|
| Linux | amd64 | [Image2Pdf_linux_amd64](https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_linux_amd64) |
| Linux | arm64 | [Image2Pdf_linux_arm64](https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_linux_arm64) |
| macOS | amd64 | [Image2Pdf_darwin_amd64](https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_darwin_amd64) |
| macOS | arm64 (Apple Silicon) | [Image2Pdf_darwin_arm64](https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_darwin_arm64) |
| Windows | amd64 | [Image2Pdf_windows_amd64.exe](https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_windows_amd64.exe) |
| Windows | arm64 | [Image2Pdf_windows_arm64.exe](https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_windows_arm64.exe) |

### 一键安装脚本

#### Linux / macOS

```bash
# 自动检测系统架构并下载
wget https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/') -O Image2Pdf && chmod +x Image2Pdf

# 移动到系统路径（可选）
sudo mv Image2Pdf /usr/local/bin/
```

#### Windows PowerShell

```powershell
# AMD64 架构
Invoke-WebRequest -Uri "https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_windows_amd64.exe" -OutFile "Image2Pdf.exe"

# ARM64 架构
Invoke-WebRequest -Uri "https://github.com/zhangyiming748/Image2Pdf/releases/latest/download/Image2Pdf_windows_arm64.exe" -OutFile "Image2Pdf.exe"
```

---

## 🛠️ 安装指南

### 前置依赖

Image2Pdf 依赖以下外部工具，请确保已正确安装：

#### 1. ImageMagick（必需）

ImageMagick 是核心的图像处理引擎，用于图片与 PDF 的转换。

**安装方法：**

```bash
# macOS (Homebrew)
brew install imagemagick

# Ubuntu / Debian
sudo apt-get update
sudo apt-get install imagemagick

# CentOS / RHEL
sudo yum install ImageMagick

# Windows
# 访问 https://imagemagick.org/script/download.php#windows 下载安装包
```

**验证安装：**

```bash
magick --version
```

#### 2. Ghostscript（PDF 转图片时必需）

当使用 `decode` 命令将 PDF 转换为图片时，需要安装 Ghostscript。

**安装方法：**

```bash
# macOS (Homebrew)
brew install ghostscript

# Ubuntu / Debian
sudo apt-get install ghostscript

# CentOS / RHEL
sudo yum install ghostscript

# Windows
# 访问 https://www.ghostscript.com/releases/gsdnld.html 下载安装包
```

**验证安装：**

```bash
# macOS / Linux
gs --version

# Windows
gswin64c --version
```

### 从源码构建

如果您希望从源码编译，或进行定制开发：

```bash
# 克隆仓库
git clone https://github.com/zhangyiming748/Image2Pdf.git
cd Image2Pdf

# 编译（默认版本）
go build -o Image2Pdf

# 编译并指定版本号
go build -ldflags "-X main.version=v1.0.0" -o Image2Pdf

# 验证安装
./Image2Pdf version
```

### 查看版本信息

```bash
# 方式 1：使用 version 子命令
./Image2Pdf version

# 方式 2：使用 --version 标志
./Image2Pdf --version

# 方式 3：使用 -v 短标志
./Image2Pdf -v
```

---

## 📚 使用文档

### 命令概览

```bash
Image2Pdf 是一个用于将一组图片合并生成单个 PDF 文件的命令行工具

Usage:
  image2pdf [command]

Available Commands:
  single      将单个文件夹内的图片转换为一个 PDF 文件
  multi       将根目录下的多个子文件夹分别转换为 PDF 文件
  decode      将 PDF 文件解码为图片
  version     显示版本信息
  completion  生成自动补全脚本
  help        获取帮助信息

Flags:
  -h, --help      显示帮助信息
  -v, --version   显示版本信息
```

### 1. single 命令 - 单文件夹图片转 PDF

将指定文件夹中的所有图片文件合并转换为一个 PDF 文件，输出的 PDF 文件名与文件夹名称相同。

#### 命令格式

```bash
image2pdf single [flags]
```

#### 参数说明

| 参数 | 简写 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--dir` | `-d` | string | ✅ | - | 包含图片的文件夹绝对路径 |
| `--compress` | `-c` | bool | ❌ | `true` | 是否启用 PDF 压缩 |

#### 使用示例

```bash
# 基本用法（启用压缩，默认）
./Image2Pdf single -d /path/to/images

# 显式启用压缩
./Image2Pdf single -d /path/to/images -c
./Image2Pdf single -d /path/to/images --compress=true

# 禁用压缩（保持原始质量）
./Image2Pdf single -d /path/to/images -c=false
./Image2Pdf single -d /path/to/images --compress=false

# 使用当前目录
./Image2Pdf single -d .
```

#### 输出说明

- 生成的 PDF 文件保存在源文件夹中
- 文件名格式：`{文件夹名}.pdf`
- 例如：文件夹 `/photos/vacation` → 生成 `/photos/vacation/vacation.pdf`

---

### 2. multi 命令 - 多文件夹批量转换

遍历根目录下的所有子文件夹，将每个子文件夹中的图片分别转换为独立的 PDF 文件。

#### 命令格式

```bash
image2pdf multi [flags]
```

#### 参数说明

| 参数 | 简写 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--dir` | `-d` | string | ✅ | - | 包含多个子文件夹的根目录绝对路径 |
| `--compress` | `-c` | bool | ❌ | `true` | 是否启用 PDF 压缩 |

#### 使用示例

```bash
# 基本用法
./Image2Pdf multi -d /path/to/root

# 启用压缩
./Image2Pdf multi -d /path/to/root -c

# 禁用压缩
./Image2Pdf multi -d /path/to/root -c=false
./Image2Pdf multi -d /path/to/root --compress=false
```

#### 目录结构示例

假设目录结构如下：

```
/books/
├── chapter1/
│   ├── page1.jpg
│   ├── page2.jpg
│   └── page3.jpg
├── chapter2/
│   ├── page1.jpg
│   └── page2.jpg
└── chapter3/
    ├── page1.jpg
    ├── page2.jpg
    └── page3.jpg
```

执行命令：

```bash
./Image2Pdf multi -d /books
```

将生成：

```
/books/
├── chapter1/
│   ├── page1.jpg
│   ├── page2.jpg
│   ├── page3.jpg
│   └── chapter1.pdf    ← 新生成
├── chapter2/
│   ├── page1.jpg
│   ├── page2.jpg
│   └── chapter2.pdf    ← 新生成
└── chapter3/
    ├── page1.jpg
    ├── page2.jpg
    ├── page3.jpg
    └── chapter3.pdf    ← 新生成
```

---

### 3. decode 命令 - PDF 转图片

将指定目录中包含 PDF 的文件夹里的 PDF 文件转换为图片序列。

#### 命令格式

```bash
image2pdf decode [flags]
```

#### 参数说明

| 参数 | 简写 | 类型 | 必填 | 默认值 | 说明 |
|------|------|------|------|--------|------|
| `--dir` | `-d` | string | ❌ | `.` | 包含 PDF 文件的根目录绝对路径 |

#### 使用示例

```bash
# 使用默认值（当前目录）
./Image2Pdf decode

# 指定目录
./Image2Pdf decode -d /path/to/pdf/folder

# 使用长格式
./Image2Pdf decode --dir /path/to/pdf/folder
```

#### 输出说明

- 图片文件保存在 PDF 所在的文件夹中
- 文件名格式：`Scan_0001.jpg`, `Scan_0002.jpg`, `Scan_0003.jpg` ...
- 图片分辨率：300 DPI
- 图片格式：JPEG

---

## ⚙️ 高级配置

### 压缩选项详解

Image2Pdf 提供两种压缩模式：

#### 启用压缩（默认）

```bash
./Image2Pdf single -d /path/to/images -c
```

**特点：**
- 使用 JPEG 压缩算法
- 质量等级：85%
- 显著减小文件体积（通常可减少 50%-80%）
- 适合网络传输和存储节省
- 轻微的质量损失（肉眼难以察觉）

#### 禁用压缩

```bash
./Image2Pdf single -d /path/to/images -c=false
```

**特点：**
- 无压缩处理
- 保持原始图片质量
- 文件体积较大
- 适合印刷出版和高质量归档

### 支持的图片格式

Image2Pdf 支持所有 ImageMagick 兼容的图片格式，包括但不限于：

- **常见格式**：JPG/JPEG, PNG, GIF, BMP, TIFF/TIF
- **专业格式**：WEBP, SVG, PSD, ICO
- **RAW 格式**：CR2, NEF, ARW（需额外配置）

### 性能优化建议

1. **大批量处理**：建议使用 `multi` 命令而非多次调用 `single`
2. **内存优化**：处理超大图片时，确保系统有足够可用内存
3. **并行处理**：当前版本为串行处理，未来版本可能支持并发
4. **磁盘空间**：确保目标文件夹有足够的存储空间

---

## 🔍 常见问题

### Q1: 提示 "未找到 magick 命令"

**原因**：ImageMagick 未安装或未添加到系统 PATH

**解决方案**：

```bash
# 检查是否安装
which magick

# macOS 安装
brew install imagemagick

# Ubuntu 安装
sudo apt-get install imagemagick

# 验证安装
magick --version
```

### Q2: PDF 转图片失败

**原因**：缺少 Ghostscript 依赖

**解决方案**：

```bash
# macOS 安装
brew install ghostscript

# Ubuntu 安装
sudo apt-get install ghostscript

# 验证安装
gs --version
```

### Q3: 生成的 PDF 文件过大

**解决方案**：

1. 确认已启用压缩（默认启用）
   ```bash
   ./Image2Pdf single -d /path -c
   ```

2. 如果仍然过大，考虑在转换前压缩原始图片

3. 使用专业的 PDF 优化工具进行二次压缩

### Q4: 图片顺序不正确

**原因**：文件系统返回的文件列表顺序不确定

**解决方案**：

1. 确保图片文件名按数字顺序命名（如 `001.jpg`, `002.jpg`）
2. 使用字母前缀保证排序正确（如 `page_01.jpg`, `page_02.jpg`）

### Q5: 中文路径或文件名问题

**建议**：

- 尽量使用英文路径和文件名
- 如必须使用中文，确保系统编码设置为 UTF-8
- Windows 用户注意控制台编码设置

### Q6: 如何查看详细的执行日志？

Image2Pdf 默认会输出详细的执行日志，包括：

- 找到的图片文件列表
- 执行的 ImageMagick 命令
- 转换结果和状态

如需更详细的调试信息，可以：

```bash
# 查看完整帮助
./Image2Pdf --help
./Image2Pdf single --help
```

---

## 🏗️ 技术架构

### 核心技术栈

- **编程语言**：Go 1.26+
- **命令行框架**：[Cobra](https://github.com/spf13/cobra) v1.10.2
- **文件查找**：[Finder](https://github.com/zhangyiming748/finder) v0.0.8
- **图像处理**：[ImageMagick](https://imagemagick.org/) 7.x
- **文件类型检测**：[filetype](https://github.com/h2non/filetype) v1.1.3

### 构建与发布

项目使用 [GoReleaser](https://goreleaser.com/) 进行自动化构建和发布：

```yaml
# .goreleaser.yml 关键配置
builds:
  - ldflags:
      - -s -w -X main.version={{.Version}}
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
```

当推送 Git tag 时，GitHub Actions 会自动：

1. 检出代码
2. 编译多平台二进制文件
3. 注入版本号
4. 创建 GitHub Release
5. 上传构建产物

---

## 🤝 贡献指南

欢迎为 Image2Pdf 做出贡献！

### 提交 Issue

- 🐛 **Bug 报告**：详细描述问题现象、复现步骤、环境信息
- 💡 **功能建议**：说明需求场景和期望效果
- 📖 **文档改进**：指出文档错误或提出改进建议

### 提交 PR

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范

- 遵循 [Effective Go](https://go.dev/doc/effective_go) 编码规范
- 添加必要的注释和文档
- 编写单元测试
- 确保 `go fmt` 和 `go vet` 通过

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🙏 致谢

感谢以下开源项目的支持：

- [ImageMagick](https://imagemagick.org/) - 强大的图像处理库
- [Cobra](https://github.com/spf13/cobra) - 优秀的 Go 命令行框架
- [GoReleaser](https://goreleaser.com/) - 便捷的 Go 发布工具

---

## 📬 联系方式

- **作者**：zhangyiming748
- **项目主页**：[https://github.com/zhangyiming748/Image2Pdf](https://github.com/zhangyiming748/Image2Pdf)
- **问题反馈**：[Issues](https://github.com/zhangyiming748/Image2Pdf/issues)

---

<div align="center">

**如果这个项目对您有帮助，请考虑给个 ⭐ Star！**

Made with ❤️ using Go

</div>
