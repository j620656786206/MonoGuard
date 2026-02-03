<div align="center">

[English](README.md) | **[繁體中文](README.zh-TW.md)**

# MonoGuard

**幾秒內分析你的 Monorepo 健康狀況**

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0-blue.svg)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.21-00ADD8.svg)](https://go.dev/)
[![pnpm](https://img.shields.io/badge/pnpm-9.0-F69220.svg)](https://pnpm.io/)

[Demo](https://monoguard-web.onrender.com) · [文件](docs/) · [回報問題](https://github.com/user/monoguard/issues)

</div>

---

## MonoGuard 是什麼？

MonoGuard 是一個開源工具，幫助你了解並改善 monorepo 的相依性健康狀況。上傳你的 workspace 設定檔，即可獲得即時分析：

- **循環相依性偵測** — 找出造成建置問題和緊密耦合的相依性循環
- **相依性圖表視覺化** — 互動式 D3.js 圖表顯示套件關係
- **健康分數** — 單一指標（0-100）總結你的 monorepo 整體健康狀況
- **架構驗證** — 驗證分層相依性是否符合你的規則
- **Bundle 影響分析** — 找出讓 bundle 膨脹的重複相依性

<div align="center">

![MonoGuard Demo](docs/assets/demo.gif)

*互動式相依性圖表，支援循環相依性高亮顯示*

</div>

## 快速開始

### 試用 Demo

造訪 [monoguard-web.onrender.com](https://monoguard-web.onrender.com) 並點擊 **「開始 Demo 分析」** — 無需註冊。

### 本機執行

```bash
# 複製專案
git clone https://github.com/user/monoguard.git
cd monoguard

# 安裝相依套件
pnpm install

# 啟動網頁應用程式
pnpm dev:web
```

開啟 [http://localhost:5173](http://localhost:5173) 開始分析。

## 功能

### 循環相依性偵測

MonoGuard 能識別 monorepo 中的直接與間接循環相依性：

- **直接循環**: `A → B → A`
- **間接循環**: `A → B → C → A`

每個循環都包含嚴重程度評級、影響評估和修復建議。

### 互動式相依性圖表

視覺化你的整個套件相依性結構：

- 力導向圖表佈局
- 縮放、平移和小地圖導航
- 大型圖表的節點展開/收合
- 循環相依性路徑高亮
- 匯出為 SVG/PNG

### 健康分數

獲得單一數字（0-100）代表你的 monorepo 健康狀況，細分為：

- 相依性健康（重複、衝突）
- 架構合規性
- 可維護性指標
- 安全性考量

### 報告匯出

以多種格式匯出分析結果：

- **HTML** — 可分享的獨立報告
- **JSON** — 機器可讀，適合 CI 整合
- **Markdown** — 適合 PR 描述

## 技術棧

| 層級 | 技術 |
|------|------|
| 前端 | React 19, TanStack Router, Tailwind CSS |
| 視覺化 | D3.js（混合 SVG/Canvas 渲染）|
| 後端 | Go 1.21, Gin framework |
| 型別 | 共用 TypeScript 型別（`@monoguard/types`）|
| 建置 | pnpm workspaces, Vite |

## 專案結構

```
monoguard/
├── apps/
│   ├── web/          # React 網頁應用程式
│   └── api/          # Go API 伺服器
├── packages/
│   └── types/        # 共用 TypeScript 型別
└── docs/             # 文件
```

## 開發指南

詳細的開發環境設定、部署說明和環境變數設定，請參考：

- [開發指南](docs/DEVELOPMENT.md) — 完整開發環境設定
- [部署指南](docs/DEPLOYMENT.md) — Render 和 Docker 部署說明

## 路線圖

- [x] 循環相依性偵測
- [x] D3.js 相依性圖表視覺化
- [x] 健康分數計算
- [x] 報告匯出（HTML/JSON/Markdown）
- [ ] WebAssembly 分析器（純前端分析）
- [ ] GitHub 整合（從 repo URL 分析）
- [ ] CI/CD 整合（GitHub Actions 等）
- [ ] VS Code 擴充套件

完整路線圖請見 [ROADMAP.zh-TW.md](ROADMAP.zh-TW.md)。

## 貢獻

歡迎貢獻！請參考 [CONTRIBUTING.zh-TW.md](CONTRIBUTING.zh-TW.md) 了解貢獻指南。

```bash
# 開發環境設定
pnpm install
pnpm dev:web

# 執行測試
pnpm test

# 型別檢查
pnpm type-check
```

## 授權條款

MIT License — 詳見 [LICENSE](LICENSE)。

---

<div align="center">

**[立即試用 MonoGuard →](https://monoguard-web.onrender.com)**

由 MonoGuard 團隊用心打造

</div>
