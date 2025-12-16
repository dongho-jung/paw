# Workbench

Claude Code 기반의 자동화된 태스크 관리 시스템입니다.

## 개요

- 프로젝트별로 태스크를 관리하고, 파일 변경을 감지하여 Claude Code가 자동으로 처리합니다.
- zellij 세션 기반으로 동작하며, 태스크별로 별도 탭이 생성됩니다.

## 디렉토리 구조

```
workbench/
├── add-project              # 새 프로젝트 추가 스크립트
├── README.md
├── _workbench/              # 전역 설정
│   ├── PROMPT.md            # 전역 시스템 프롬프트
│   ├── start                # start 스크립트 템플릿
│   └── layout.kdl           # zellij 레이아웃
└── projects/
    └── {프로젝트명}/
        ├── start            # 프로젝트 시작 스크립트
        ├── PROMPT.md        # 프로젝트별 프롬프트
        ├── metadata         # 프로젝트 메타데이터
        ├── location/        # 실제 프로젝트 경로 (심볼릭 링크)
        ├── .watcher/        # watcher 상태
        │   ├── watcher.log
        │   └── watcher.pid
        ├── agents/          # 에이전트 작업 공간
        │   └── {태스크명}/
        │       ├── log              # 에이전트 로그
        │       ├── task.json        # 태스크 정보
        │       ├── system-prompt.txt
        │       ├── user-prompt.txt
        │       ├── .tab-created     # 탭 생성 마커
        │       └── worktree/        # git worktree (작업용)
        ├── to-do/           # 대기 중인 태스크
        ├── in-progress/     # 진행 중인 태스크
        ├── in-review/       # 리뷰 대기 중
        ├── done/            # 완료된 태스크
        └── cancelled/       # 취소된 태스크
```

## location 심볼릭 링크

각 프로젝트의 `location/`은 실제 소스 코드 저장소를 가리키는 심볼릭 링크입니다.

- 에이전트는 `location/`을 통해 소스 코드에 접근합니다
- 직접 `location/`에서 작업하지 않고, git worktree를 생성하여 격리된 환경에서 작업합니다
- worktree는 `agents/{태스크명}/worktree/`에 생성됩니다

## 사용법

### 1. 프로젝트 추가

```bash
./add-project
```

- yazi가 열리면 관리할 프로젝트 디렉토리 선택
- 자동으로 `projects/{프로젝트명}/` 생성
- 자동으로 `./start` 실행

### 2. 프로젝트 시작

```bash
cd projects/{프로젝트명}
./start
```

- 기존 세션이 있으면 종료 후 새로 시작
- zellij 세션 시작 + 파일 감시 시작
- zellij 종료 시 fswatch도 자동 정리

### 3. 태스크 생성

```bash
echo "README에 hello 추가해줘" > to-do/add-hello
```

태스크 파일이 `to-do/`에 생성되면:
1. 자동으로 `in-progress/`로 이동
2. 새 zellij 탭 생성: `add-hello/in-progress`
3. 탭에서 수직 분할 (왼쪽: claude, 오른쪽: shell)
4. Claude Code 실행 및 태스크 전달
5. `agents/add-hello/` 디렉토리에 작업 공간 생성

### 4. 태스크 파일 형식

태스크 파일은 구분자로 요청과 결과를 분리합니다:

```
태스크 설명 (사용자 작성)
----------
결과 요약 (에이전트가 작성)
```

- `----------`: 10개의 하이픈으로 구분
- 구분자 위: 사용자의 요청 (에이전트가 수정 금지)
- 구분자 아래: 에이전트가 작업 완료 후 결과 작성

### 5. 태스크 상태 변경

에이전트가 작업 완료 시:
- 성공: `mv in-progress/태스크 done/`
- 리뷰 필요: `mv in-progress/태스크 in-review/`

사용자가 수동으로 상태 변경 가능:
```bash
mv in-progress/fix-bug done/        # 완료로 변경
mv in-progress/fix-bug cancelled/   # 취소
mv in-review/fix-bug to-do/         # 재처리 요청
```

## 설정

### _workbench/PROMPT.md

전역 시스템 프롬프트입니다. 모든 프로젝트의 Claude Code에 적용됩니다.

주요 내용:
- 디렉토리 구조 및 location 심볼릭 링크 설명
- 태스크 파일 형식
- 워크플로우 (to-do → in-progress → done/in-review)
- git worktree 사용법
- 로깅 규칙

### {프로젝트}/PROMPT.md

프로젝트별 프롬프트입니다. 해당 프로젝트에만 적용됩니다.

### _workbench/layout.kdl

zellij 레이아웃 설정:
- 상단: 탭 바
- 하단: 상태 바
- 기본 탭 이름: `_`

## 의존성

```bash
brew install zellij    # 터미널 멀티플렉서
brew install fswatch   # 파일 변경 감시
brew install yazi      # 파일 매니저
```

## 권한

프로젝트 루트 디렉토리는 쓰기 금지됩니다 (`chmod a-w`).
파일은 하위 디렉토리에만 생성할 수 있습니다:
- `to-do/`, `in-progress/`, `in-review/`, `done/`, `cancelled/`, `agents/`

## zellij 단축키

- `Ctrl+O, d`: 세션에서 분리 (detach)
- `Ctrl+O, w`: 탭 전환
- `Ctrl+O, n`: 새 pane
- `Ctrl+O, x`: pane 닫기
