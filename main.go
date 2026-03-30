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
	Name      string
	Depth     uint
	IsLast    bool
	Children  []FileEntry
	EntryType FileType
	LinkDst   string
}

type FileType int

type Config struct {
	RootPath           string
	IsSearchHiddenPath bool
	IgnoreDirArray     []string
	MaxDepth           uint
	OnlyDir            bool
}

type Result struct {
	dirCnt  uint
	fileCnt uint
}

// color
const (
	Reset    string = "\033[0m"
	BoldCyan string = "\033[1;36m"
	Red      string = "\033[31m"
	Magenta  string = "\033[35m"
	Yellow   string = "\033[33m"
	Green    string = "\033[32m"
)

// fileType
const (
	FileTypeNone FileType = iota
	FileTypeDir
	FileTypeExec
	FileTypePipe
	FileTypeSymlink
	FileTypeSocket
)

// line
const (
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
	if color == "" {
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

func printTree(entry *FileEntry, prefix string) {
	nextPrefix := prefix

	if entry.Depth == 0 {
		// 입력한 최상위 경로의 경우 색만 표기
		fmt.Println(colorize(entry.Name, BoldCyan))
	} else {
		var strBuilder strings.Builder

		// 현재 깊이일 경우
		if entry.IsLast {
			// 마지막일 경우 세로 마지막 문자사용
			strBuilder.WriteString(verticalEndChar)
		} else {
			// 마지막이 아닌 경우 분기처리 문자사용
			strBuilder.WriteString(verticalBranchChar)
		}

		var color string
		switch entry.EntryType {
		case FileTypeDir:
			color = BoldCyan
		case FileTypeExec:
			color = Red
		case FileTypePipe:
			color = Yellow
		case FileTypeSymlink:
			color = Magenta
		case FileTypeSocket:
			color = Green
		}

		strBuilder.WriteString(horizontalLine)
		strBuilder.WriteString(" ")
		strBuilder.WriteString(colorize(entry.Name, color))

		if entry.EntryType == FileTypeSymlink {
			strBuilder.WriteString(fmt.Sprintf(" -> %s", entry.LinkDst))
		}

		if !entry.IsLast {
			// 현재노드가 마지막이 아니라면 세로선 추가해서 자식에게 전달
			nextPrefix += verticalNormalChar + strings.Repeat(" ", indentWidth)
		} else {
			// 현재노드가 마지막이라면 빈 공간 자식에게 전달
			nextPrefix += strings.Repeat(" ", indentWidth+1)
		}

		fmt.Println(prefix + strBuilder.String())
	}

	for _, child := range entry.Children {
		// 자식 노드 순회하면서 출력
		printTree(&child, nextPrefix)
	}
}

func buildTree(dirPath string, parent *FileEntry, cfg *Config, result *Result) error {
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
		// onlyDir
		if cfg.OnlyDir && !entry.IsDir() {
			continue
		}
		child := FileEntry{
			Name:      entry.Name(),
			Depth:     parent.Depth + 1,
			Children:  make([]FileEntry, 0),
			EntryType: FileTypeNone,
		}

		info, err := os.Lstat(entryPath)
		if err != nil {
			return err
		}

		mode := info.Mode()
		if mode&os.ModeSymlink != 0 {
			// 심볼릭 링크
			child.EntryType = FileTypeSymlink
			dst, err := os.Readlink(entryPath)
			if err != nil {
				return err
			}
			child.LinkDst = dst
		} else if mode&os.ModeDir != 0 {
			// 폴더
			child.EntryType = FileTypeDir
		} else if mode.Perm()&0111 != 0 {
			// 실행파일
			child.EntryType = FileTypeExec
		} else if mode&os.ModeNamedPipe != 0 {
			// 파이프
			child.EntryType = FileTypePipe
		} else if mode&os.ModeSocket != 0 {
			// 소켓
			child.EntryType = FileTypeSocket
		}

		if child.EntryType == FileTypeDir {
			// 폴더일 경우 플래그 세팅 후 재귀 호출
			if err := buildTree(entryPath, &child, cfg, result); err != nil {
				return err
			}
			result.dirCnt += 1
		} else {
			result.fileCnt += 1
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
	onlyDirRef := flag.Bool("dir", false, "show only directories (default: false)")

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
	config.OnlyDir = *onlyDirRef
	return &config
}

func main() {

	// 플래그 초기화 및 설정값 로드
	cfg := initFlag()
	result := &Result{
		// tree의 경우 탐색시작 디렉토리도 포함하므로 1으로 시작
		dirCnt:  1,
		fileCnt: 0,
	}

	rootInfo, err := os.Stat(cfg.RootPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if !rootInfo.IsDir() {
		fmt.Fprintln(os.Stderr, fmt.Errorf("input is not directory"))
		return
	}

	root := FileEntry{
		Name:      cfg.RootPath,
		EntryType: FileTypeDir,
		Children:  make([]FileEntry, 0),
		Depth:     0,
	}

	if err := buildTree(cfg.RootPath, &root, cfg, result); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}

	printTree(&root, "")

	if cfg.OnlyDir {
		fmt.Printf("\n%d directories\n", result.dirCnt)
	} else {
		fmt.Printf("\n%d directories, %d files\n", result.dirCnt, result.fileCnt)
	}
}
