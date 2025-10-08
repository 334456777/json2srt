# json2srt

将 OpenAI Whisper 的 JSON 输出转换为精确到每个字的 SRT 字幕文件。

## 项目目的

OpenAI Whisper 在使用 `--word_timestamps True` 参数时，会生成包含每个单词精确时间戳的 JSON 文件。本工具将这些 JSON 文件转换为 SRT 字幕格式，**保持每个字/词的时间戳精度**，而不是像普通 SRT 那样只有句子级别的时间戳。

这对于需要精确同步的场景特别有用，如：
- 语言学习应用
- 精确的字幕制作
- 音频分析和研究

## 使用方法

1. 使用 Whisper 生成带词级时间戳的 JSON 文件：
   ```bash
   whisper audio.mp3 --word_timestamps True --output_format json
   ```
2. 将生成的 JSON 文件放在本程序同目录下
3. 运行程序后自动批量转换为精确的 SRT 文件

## 特性

- ✨ **精确时间戳**：保持 Whisper 词级别的时间戳精度
- 🚀 **批量处理**：自动处理目录中的所有 JSON 文件
- ⚡ **并发处理**：利用多核 CPU 提高处理速度
- 📝 **BOM 支持**：生成带 BOM 的 SRT 文件，确保中文字符正确显示
- 🔍 **智能过滤**：自动跳过 Whisper 的噪音标记（如 `[_TT_123]`）

## 运行

```bash
go run main.go
```
或编译后运行：
```bash
go build -o json2srt main.go
./json2srt
```

## 依赖
- Go 1.16 及以上

## 输出
- 每个 JSON 文件会生成对应的 SRT 文件，带 BOM。
- SRT 文件中的时间戳精确到每个词的开始和结束时间

## 许可证
MIT
