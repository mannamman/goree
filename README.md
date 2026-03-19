# goree

go언어 공부 겸 tree를 모사하기 위해 만들었습니다.
 - 디렉토리 구조를 트리 형태로 출력하는 CLI 도구입니다.

## 사용법

```bash
# 현재 디렉토리
goree

# 특정 디렉토리
goree -path /path/to/dir

# 숨김 파일/폴더 포함 조회
goree -all

# 특정 폴더 제외
goree -ignore node_modules
goree -ignore node_modules,.git

# 최대 탐색 깊이 세팅
goree -depth 2

# 조합 예시
goree -path /path/to/dir -all -ignore node_modules,.git -depth 3
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

## TODO
### 더 생각나면 추가 예정

| 완료 여부 | 작업 내용 | 생성일자 | 완료일자 | 비고 |
|:---------:|-----------|----------|----------|------|
| [x] | 숨김 파일/폴더 조회 선택 추가 (기본 숨김 파일/폴더 미조회) | 2026-03-13 | 2026-03-15 | |
| [x] | 특정 폴더 탐색 거부 | 2026-03-13 | 2026-03-16 | |
| [x] | 최대 탐색 깊이 설정 | 2026-03-13 | 2026-03-16 | |
| [-] | ~~특정 내용의 파일 찾기~~ | 2026-03-16 | 2026-03-19 | tree명령어와 어울리지 않는 것 같음 |
| [x] | 파일 타입 색 지정하기 | 2026-03-16 | 2026-03-18 | |
| [ ] | 디렉토리만 조회 플래그 | 2026-03-19 | |  |

---

> 이 README는 [Claude](https://claude.ai)의 도움으로 작성되었습니다.
