# 貢獻指南

感謝您有興趣為 MonoGuard 做出貢獻！本文件提供貢獻者所需的指南和資訊。

## 行為準則

請在所有互動中保持尊重和建設性。我們正在一起建構一個專案。

## 如何貢獻

### 回報錯誤

1. 搜尋[現有的 Issues](https://github.com/j620656786206/MonoGuard/issues) 以避免重複
2. 如果沒有找到，[開啟新的 Issue](https://github.com/j620656786206/MonoGuard/issues/new)，包含：
   - 清晰、描述性的標題
   - 重現步驟
   - 預期行為與實際行為
   - 環境資訊（作業系統、Node 版本、瀏覽器）

### 建議功能

1. 查看 [Roadmap](README.md#roadmap) 了解已規劃的功能
2. 開啟帶有 `enhancement` 標籤的 Issue
3. 描述使用情境和建議的解決方案

### Pull Request

1. Fork 此儲存庫
2. 建立功能分支：`git checkout -b feat/your-feature`
3. 進行變更
4. 根據需要撰寫/更新測試
5. 確保所有檢查通過：`pnpm lint && pnpm type-check && pnpm test`
6. 使用 Conventional Commits 格式提交（見下方）
7. Push 並開啟 PR

## 開發環境設定

### 先決條件

- Node.js 20+
- pnpm 9+
- Go 1.21+（用於 API 開發）

### 開始使用

```bash
# Clone 您的 fork
git clone https://github.com/YOUR_USERNAME/monoguard.git
cd monoguard

# 安裝相依套件
pnpm install

# 啟動開發伺服器
pnpm dev:web
```

### 專案結構

```
monoguard/
├── apps/
│   ├── web/          # React 網頁應用程式 (Vite + TanStack Router)
│   └── api/          # Go API 伺服器 (Gin)
├── packages/
│   └── types/        # 共用 TypeScript 型別
├── docs/             # 文件
└── scripts/          # 建置和部署腳本
```

### 主要指令

```bash
# 開發
pnpm dev:web          # 在 localhost:3000 啟動網頁應用程式

# 品質檢查
pnpm lint             # 執行 Biome 檢查
pnpm type-check       # TypeScript 型別檢查
pnpm test             # 執行測試

# 建置
pnpm build            # 建置所有套件
```

## 程式碼規範

### TypeScript

- 所有新程式碼使用 TypeScript
- 優先使用 `type` 而非 `interface` 定義物件型別
- 從 `@monoguard/types` 套件匯出型別

### React

- 使用函數式元件和 hooks
- 遵循 `apps/web/app/components/` 中現有的元件模式
- 使用 TanStack Router 進行路由

### 樣式

- 使用 Tailwind CSS 工具類別
- 遵循現有的設計模式
- 維持響應式設計

### 測試

- 為新功能撰寫測試
- 維持 >80% 的程式碼覆蓋率
- 使用 `src/__tests__/factories/` 中現有的測試工廠

## 提交訊息規範

我們使用 [Conventional Commits](https://www.conventionalcommits.org/)：

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### 類型

- `feat`: 新功能
- `fix`: 錯誤修正
- `docs`: 僅文件變更
- `style`: 格式調整，無程式碼變更
- `refactor`: 重構，既不修正錯誤也不新增功能
- `test`: 新增或更新測試
- `chore`: 維護任務

### 範例

```
feat(web): add circular dependency highlighting
fix(api): resolve memory leak in graph traversal
docs: update README with new features
test(web): add integration tests for analysis flow
```

## Pull Request 流程

1. 如有需要，更新文件
2. 為新功能新增測試
3. 確保 CI 通過
4. 請求維護者審查
5. 回應審查意見
6. 審核通過後 squash and merge

## 有問題嗎？

- 開啟 [Discussion](https://github.com/j620656786206/MonoGuard/discussions) 詢問一般問題
- 如果遇到阻礙，在 Issue 中標記維護者

---

感謝您的貢獻！
