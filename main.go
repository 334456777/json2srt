package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// --- Struct Definitions ---
// Timestamp represents a timestamp structure containing start and end times
type Timestamp struct {
	From string `json:"from"` // Start time, format "00:00:00,000"
	To   string `json:"to"`   // End time, format "00:00:01,960"
}

// Token represents a single word in transcription
type Token struct {
	Text       string    `json:"text"`       // Token text
	Timestamps Timestamp `json:"timestamps"` // Timestamps for this token
}

// Segment represents a paragraph in transcription
type Segment struct {
	Timestamps Timestamp `json:"timestamps"` // Timestamps for the paragraph
	Text       string    `json:"text"`       // Text content of the paragraph
	Tokens     []Token   `json:"tokens"`     // Token list for this paragraph
}

// WhisperOutput represents the root structure of Whisper transcription output
type WhisperOutput struct {
	Transcription []Segment `json:"transcription"` // List of transcription paragraphs
}

// --- Core Logic ---
// extractTimestamps extracts start and end timestamps from segment's tokens
// Returns startTime, endTime, hasValidTokens
func extractTimestamps(segment Segment) (string, string, bool) {
	var startTime, endTime string
	hasValidTokens := false

	// Find first non-noise token's start time
	for _, token := range segment.Tokens {
		if !strings.HasPrefix(token.Text, "[_") && token.Text != "" {
			// Use token's From as start time
			startTime = token.Timestamps.From
			hasValidTokens = true
			break
		}
	}

	// 查找最后一个非噪音 token 的结束时间
	for j := len(segment.Tokens) - 1; j >= 0; j-- {
		token := segment.Tokens[j]
		if !strings.HasPrefix(token.Text, "[_") && token.Text != "" {
			endTime = token.Timestamps.To
			hasValidTokens = true
			break
		}
	}

	// 如果没有有效 tokens，使用 segment 的时间戳
	if !hasValidTokens || startTime == "" || endTime == "" {
		startTime = segment.Timestamps.From
		endTime = segment.Timestamps.To
	}

	return startTime, endTime, hasValidTokens
}

// processFile 处理单个 JSON 文件，生成 SRT 文件
// ctx 用于支持取消操作
func processFile(ctx context.Context, jsonPath string) error {
	// 检查 context 是否已取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	srtPath := strings.TrimSuffix(jsonPath, filepath.Ext(jsonPath)) + ".srt"
	byteValue, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("读取文件 %s 失败: %v", filepath.Base(jsonPath), err)
	}
	var output WhisperOutput
	if err := json.Unmarshal(byteValue, &output); err != nil {
		return fmt.Errorf("解析 JSON 文件 %s 失败: %v", filepath.Base(jsonPath), err)
	}

	// 验证 JSON 结构完整性
	if output.Transcription == nil {
		return fmt.Errorf("JSON 文件 %s 缺少 transcription 字段", filepath.Base(jsonPath))
	}

	var srtBuilder strings.Builder

	// 顺序处理 segments
	for i, segment := range output.Transcription {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		startTime, endTime, _ := extractTimestamps(segment)
		if startTime == "" || endTime == "" {
			continue // 跳过无效 segment
		}

		text := strings.TrimSpace(segment.Text)
		srtBuilder.WriteString(fmt.Sprintf("%d\n", i+1))
		srtBuilder.WriteString(fmt.Sprintf("%s --> %s\n", startTime, endTime))
		srtBuilder.WriteString(fmt.Sprintf("%s\n\n", text))
	}

	bom := []byte{0xEF, 0xBB, 0xBF}
	content := []byte(srtBuilder.String())
	dataToWrite := append(bom, content...)
	if err := os.WriteFile(srtPath, dataToWrite, 0644); err != nil {
		return fmt.Errorf("写入 SRT 文件 %s 失败: %v", filepath.Base(srtPath), err)
	}
	fmt.Printf("    [成功] 已生成 SRT 文件: %s\n", filepath.Base(srtPath))
	return nil
}

// worker 处理文件的 goroutine
func worker(ctx context.Context, id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case jsonPath, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Printf("Worker %d 开始处理文件: %s\n", id, filepath.Base(jsonPath))
			if err := processFile(ctx, jsonPath); err != nil {
				log.Printf("Worker %d 处理文件 %s 时出错: %v", id, filepath.Base(jsonPath), err)
			}
		}
	}
}

func main() {
	ctx := context.Background() // 创建根 context，可用于取消操作

	// --- MODIFIED: 使用 os.Getwd() 获取当前终端的工作目录 ---
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前工作目录失败: %v", err)
	}

	fmt.Printf("程序正在当前终端目录中运行: %s\n", workDir)

	// MODIFIED: 使用 workDir 变量
	files, err := os.ReadDir(workDir)
	if err != nil {
		log.Fatalf("读取目录失败: %v", err)
	}

	jsonPaths := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			// MODIFIED: 使用 workDir 变量
			jsonPaths = append(jsonPaths, filepath.Join(workDir, file.Name()))
		}
	}

	if len(jsonPaths) == 0 {
		fmt.Println("任务完成：未在当前目录中找到任何 .json 文件。")
		return
	}

	numJobs := len(jsonPaths)
	jobs := make(chan string, numJobs)
	var wg sync.WaitGroup

	numWorkers := runtime.NumCPU()
	fmt.Printf("启动 %d 个 Worker 并行处理 %d 个文件...\n", numWorkers, numJobs)
	fmt.Println(strings.Repeat("=", 40))

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(ctx, w, jobs, &wg)
	}

	for _, path := range jsonPaths {
		jobs <- path
	}
	close(jobs)

	wg.Wait()

	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("任务完成：共处理了 %d 个 JSON 文件。\n", numJobs)
}
