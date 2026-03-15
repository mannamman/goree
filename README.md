# goree

go언어 공부 겸 tree를 모사하기 위해 만들었습니다.
 - 디렉토리 구조를 트리 형태로 출력하는 CLI 도구입니다.

## 사용법

```bash
# 현재 디렉토리
goree

# 특정 디렉토리
goree /path/to/dir
```

## 빌드

```bash
# 현재 플랫폼
go build -o goree

# macOS - Apple Silicon (M1/M2/M3)
GOOS=darwin GOARCH=arm64 go build -o goree-darwin-arm64

# macOS - Intel
GOOS=darwin GOARCH=amd64 go build -o goree-darwin-amd64

# Linux - AMD/Intel
GOOS=linux GOARCH=amd64 go build -o goree-linux-amd64

# Linux - ARM
GOOS=linux GOARCH=arm64 go build -o goree-linux-arm64

# Windows - AMD/Intel
GOOS=windows GOARCH=amd64 go build -o goree-amd64.exe

# Windows - ARM
GOOS=windows GOARCH=arm64 go build -o goree-arm64.exe
```

---

> 이 README는 [Claude](https://claude.ai)의 도움으로 작성되었습니다.

## TODO
### 더 생각나면 추가 예정

- [V] 숨김 파일/폴더 조회 선택 추가 (기본 숨김 파일/폴더 미조회)
- [ ] 특정 폴더 탐색 거부
- [ ] 최대 탐색 깊이 설정