package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FileEntry struct {
	Name     string
	IsDir    bool
	Depth    uint
	IsRoot   bool
	IsLast   bool
	Children []FileEntry
}

type Config struct {
	RootPath           string
	IsSearchHiddenPath bool
	IgnoreDirArray     []string
	MaxDepth           uint
}

const (
	Reset              string = "\033[0m"
	BoldCyan           string = "\033[1;36m"
	indentWidth        int    = 3
	verticalBranchChar string = "├"
	verticalNormalChar string = "│"
	verticalEndChar    string = "└"
	horizontalLine     string = "──"
)

func isHiddenPath(entryName string) bool {
	return strings.HasPrefix(entryName, ".")
}

func colorize(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

func formatTreeLine(entry *FileEntry) string {
	var strBuilder strings.Builder
	for i := uint(0); i < entry.Depth; i++ {
		if i == entry.Depth-1 {
			// 현재 깊이일 경우
			if entry.IsLast {
				// 마지막일 경우 세로 마지막 문자사용
				strBuilder.WriteString(verticalEndChar)
			} else {
				// 마지막이 아닌 경우 분기처리 문자사용
				strBuilder.WriteString(verticalBranchChar)
			}
		} else {
			// 현재 깊이보다 상위 깊이면 세로선만 추가
			strBuilder.WriteString(verticalNormalChar)
			strBuilder.WriteString(strings.Repeat(" ", indentWidth))
		}
	}

	if entry.IsDir {
		// 폴더의 경우 파란색 세팅
		strBuilder.WriteString(horizontalLine)
		strBuilder.WriteString(" ")
		strBuilder.WriteString(colorize(entry.Name, BoldCyan))
	} else {
		strBuilder.WriteString(horizontalLine)
		strBuilder.WriteString(" ")
		strBuilder.WriteString(entry.Name)
	}

	return strBuilder.String()
}

func printTree(entry *FileEntry) {
	if entry.IsRoot {
		// 입력한 최상위 경로의 경우 색만 표기
		fmt.Println(colorize(entry.Name, BoldCyan))
	} else {
		// 입력한 최상위 경로가 아니라면 깊이에 맞게 출력
		fmt.Println(formatTreeLine(entry))
	}

	for _, child := range entry.Children {
		// 자식 노드 순회하면서 출력
		printTree(&child)
	}
}

func buildTree(dirPath string, parent *FileEntry, cfg *Config) error {
	// 전달받은 경로의 파일 혹은 폴더 확인
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// 현재 경로의 파일 혹은 폴더를 열기 위해 경로 세팅
		entryPath := filepath.Join(dirPath, entry.Name())
		// 숨김 경로일 경우 건너뜀
		if !cfg.IsSearchHiddenPath && isHiddenPath(entry.Name()) {
			continue
		}
		// 무시할 폴더일 경우 건너뜀
		if entry.IsDir() && len(cfg.IgnoreDirArray) > 0 {
			if slices.Contains(cfg.IgnoreDirArray, entry.Name()) {
				continue
			}
		}
		// maxDepth값이 설정되었고, 부모의 깊이가 같거나 크면 탐색 종료
		if cfg.MaxDepth > 0 && parent.Depth >= cfg.MaxDepth {
			continue
		}
		child := FileEntry{
			Name:     entry.Name(),
			Depth:    parent.Depth + 1,
			Children: make([]FileEntry, 0),
			IsDir:    false,
			IsLast:   false,
		}

		if entry.IsDir() {
			// 폴더일 경우 플래그 세팅 후 재귀 호출
			child.IsDir = true
			if err := buildTree(entryPath, &child, cfg); err != nil {
				return err
			}
		}
		// 파일 혹은 폴더를 부모 노드의 자식 배열에 추가
		parent.Children = append(parent.Children, child)
	}

	// 마지막에 추가된 자식에 직접 세팅
	if len(parent.Children) > 0 {
		parent.Children[len(parent.Children)-1].IsLast = true
	}

	return nil
}

func initFlag() *Config {
	rootPathRef := flag.String("path", ".", "root directory to search")
	// -all 만 사용시 bool의 경우 true로 세팅됨
	isSearchHiddenPathRef := flag.Bool("all", false, "include hidden files (default: false)")
	ignoreDirStringRef := flag.String("ignore", "", "comma-separated list of directories to ignore")
	maxDepthRef := flag.Uint("depth", 0, "max search depth (0: unlimited)")

	flag.Parse()

	// 설정값 세팅
	var config Config
	config.RootPath = *rootPathRef
	config.IsSearchHiddenPath = *isSearchHiddenPathRef
	ignoreDirString := *ignoreDirStringRef
	if ignoreDirString != "" {
		config.IgnoreDirArray = strings.Split(ignoreDirString, ",")
	}
	config.MaxDepth = *maxDepthRef

	return &config
}

func main() {

	// 플래그 초기화 및 설정값 로드
	cfg := initFlag()

	if _, err := os.Stat(cfg.RootPath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	root := FileEntry{
		Name:     cfg.RootPath,
		IsDir:    true,
		Children: make([]FileEntry, 0),
		IsRoot:   true,
		Depth:    0,
	}

	if err := buildTree(cfg.RootPath, &root, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}

	printTree(&root)
}
