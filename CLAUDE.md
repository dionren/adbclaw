# adbclaw

Android 设备控制 CLI，供 AI agent 自动化调用。纯工具层，不含 LLM/Agent 逻辑。

## 项目结构

- `src/` — Go 代码根目录（go.mod 在此）
- `src/cmd/` — Cobra CLI 命令
- `src/pkg/` — 核心库（adb / input / observe / output）
- `docs/` — 技术文档（详细方案见 `docs/adbclaw-technical-plan.md`）
- `website/` — 官网

## 构建

```bash
cd src
make build   # 产物 → src/bin/adbclaw
make test
make lint    # go vet
make clean
```

Go 1.24，依赖 cobra v1.10.2。构建产物在 `src/bin/`（已 gitignore）。

## 架构要点

- **Commander 接口** (`pkg/adb/shell.go`) — 所有 pkg 通过接口调用 ADB，测试用 mock
- **JSON Envelope** (`pkg/output/envelope.go`) — 统一 `{ok, command, data, error, duration_ms, timestamp}`
- **UI 树过滤** — 只索引有 text/resource-id/content-desc 或 clickable/scrollable 的节点
- **输入为顶级命令** — `adbclaw tap` 而非 `adbclaw input tap`
- **observe 部分失败容忍** — 截屏和 UI 树并行，互不阻塞

## 命令树

```
adbclaw
├── device list / info
├── observe            # 截屏 + UI树 并行
├── screenshot
├── ui tree / find
├── tap / long-press / swipe / key / type
├── app list / current / launch / stop
├── skill              # 输出 skill.json (go:embed)
└── doctor             # 环境检查
```

## 当前阶段

Phase 1 MVP — 纯 adb shell 命令实现，不含 adbclawd 设备端服务。
