# json2srt

Convert OpenAI Whisper JSON output with word-level timestamps to precise SRT subtitle files.

## Purpose

OpenAI Whisper generates JSON files with word-level timestamps when using `--word_timestamps True`. This tool converts those JSON files to SRT subtitle format, **preserving per-character/word timing accuracy** rather than sentence-level timestamps.

This is particularly useful for scenarios requiring precise synchronization:
- Language learning applications
- Precise subtitle production
- Audio analysis and research

## Usage

1. Generate JSON file with word-level timestamps using Whisper:
   ```bash
   whisper audio.mp3 --word_timestamps True --output_format json
   ```
2. Place generated JSON file in same directory as program
3. Run program to auto-convert to precise SRT files

## Features

- ✨ **Precise timestamps**: Maintains Whisper's word-level timing precision
- 🚀 **Batch processing**: Auto-processes all JSON files in directory
- ⚡ **Concurrent processing**: Uses multi-core CPU for faster processing
- 📝 **BOM support**: Generates SRT files with BOM for correct Chinese character display
- 🔍 **Smart filtering**: Auto-skips Whisper noise markers (like `[_TT_123]`)

## Run

```bash
go run main.go
```
Or compile and run:
```bash
go build -o json2srt main.go
./json2srt
```

## Dependencies

- Go 1.16+

## Output

- Each JSON file generates corresponding SRT file with BOM
- SRT file timestamps are precise to each word's start and end time

## License

MIT
