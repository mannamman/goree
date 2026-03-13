package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FileEntry struct {
	Name     string
	IsDir    bool
	Depth    int
	IsRoot   bool
	IsLast   bool
	Children []FileEntry
}

const (
	Reset = "\033[0m"
	Blue  = "\033[1;36m"
)

const indentWidth int = 3
const verticalBranchChar string = "├"
const verticalNormalChar string = "│"
const verticalEndChar string = "└"
const horizontalLine string = "──"

func colorizeBlue(text string) string {
	return fmt.Sprintf("%s%s%s", Blue, text, Reset)
}

func formatTreeLine(entry *FileEntry) string {
	lineStr := ""
	for i := 0; i < entry.Depth; i++ {
		if i == entry.Depth-1 {
			// 현재 깊이일 경우
			if entry.IsLast {
				// 마지막일 경우 세로 마지막 문자사용
				lineStr = lineStr + verticalEndChar
			} else {
				// 마지막이 아닌 경우 분기처리 문자사용
				lineStr = lineStr + verticalBranchChar
			}
		} else {
			// 현재 깊이보다 상위 깊이면 세로선만 추가
			lineStr = lineStr + verticalNormalChar + strings.Repeat(" ", indentWidth)
		}
	}

	if entry.IsDir {
		// 폴더의 경우 파란색 세팅
		lineStr = lineStr + horizontalLine + " " + colorizeBlue(entry.Name)
	} else {
		lineStr = lineStr + horizontalLine + " " + entry.Name
	}

	return lineStr
}

func printTree(entry *FileEntry) {
	if entry.IsRoot {
		// 입력한 최상위 경로의 경우 색만 표기
		fmt.Printf("%s\n", colorizeBlue(entry.Name))
	} else {
		// 입력한 최상위 경로가 아니라면 깊이에 맞게 출력
		fmt.Println(formatTreeLine(entry))
	}

	for _, child := range entry.Children {
		// 자식 노드 순회하면서 출력
		printTree(&child)
	}
}

func buildTree(dirPath string, parent *FileEntry) error {
	// 전달받은 경로의 파일 혹은 폴더 확인
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for idx, entry := range entries {
		// 현재 경로의 파일 혹은 폴더를 열기 위해 경로 세팅
		entryPath := filepath.Join(dirPath, entry.Name())
		child := FileEntry{
			Name:     entry.Name(),
			Depth:    parent.Depth + 1,
			Children: make([]FileEntry, 0),
			IsDir:    false,
			IsLast:   idx == len(entries)-1,
		}

		if entry.IsDir() {
			// 폴더일 경우 플래그 세팅 후 재귀 호출
			child.IsDir = true
			if err := buildTree(entryPath, &child); err != nil {
				return err
			}
		}
		// 파일 혹은 폴더를 부모 노드의 자식 배열에 추가
		parent.Children = append(parent.Children, child)
	}

	return nil
}

func main() {
	rootPath := "."
	if len(os.Args) == 2 {
		rootPath = os.Args[1]
	}

	if _, err := os.Stat(rootPath); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	root := FileEntry{
		Name:     rootPath,
		IsDir:    true,
		Children: make([]FileEntry, 0),
		IsRoot:   true,
		Depth:    0,
	}

	if err := buildTree(rootPath, &root); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		return
	}

	printTree(&root)
}
