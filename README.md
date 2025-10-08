# json2srt

将 Whisper 格式的 JSON 转换为 SRT 字幕文件。

## 使用方法

1. 将 Whisper 导出的 JSON 文件放在本程序同目录下。
2. 运行程序后自动批量转换为 SRT 文件。

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

## 许可证
MIT
