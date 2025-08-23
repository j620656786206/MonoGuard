# 開發指南

本指南涵蓋 MonoGuard 的開發工作流程、架構決策和最佳實務。

## 🏗️ 架構概覽

MonoGuard 採用 **整合套件策略**，旨在最大化可維護性和部署彈性：

### 核心原則
1. **應用程式包含所有邏輯** - 業務邏輯保留在應用程式界限內
2. **最小化共用程式庫** - 僅共享必要的型別定義
3. **獨立部署** - 每個應用程式可獨立部署和擴展
4. **明確服務邊界** - 服務間具備完整定義的介面

### 目錄結構
```
mono-guard/
├── apps/
│   ├── frontend/           # Next.js 14 網頁應用程式
│   │   ├── src/
│   │   │   ├── app/        # Next.js App Router 頁面
│   │   │   ├── components/ # React 元件（按領域組織）
│   │   │   ├── hooks/      # 自訂 React hooks
│   │   │   ├── lib/        # 工具函數與設定
│   │   │   ├── store/      # Zustand 狀態管理
│   │   │   └── types/      # 應用程式專用型別定義
│   │   ├── tests/          # 測試檔案
│   │   └── public/         # 靜態資源
│   │
│   ├── cli/                # Node.js CLI 工具
│   │   ├── src/
│   │   │   ├── commands/   # CLI 指令實作
│   │   │   ├── lib/        # 核心 CLI 邏輯
│   │   │   └── utils/      # CLI 工具
│   │   └── tests/          # CLI 測試
│   │
│   └── api/                # Go API 服務
│       ├── cmd/server/     # 應用程式進入點
│       ├── internal/       # 私有應用程式程式碼
│       │   ├── handlers/   # HTTP 處理器
│       │   ├── services/   # 業務邏輯
│       │   ├── models/     # 資料模型
│       │   └── config/     # 設定
│       └── pkg/            # 公用套件
│
├── libs/
│   └── shared-types/       # TypeScript 型別定義
│       └── src/
│           ├── api.ts      # API 合約型別
│           ├── domain.ts   # 領域模型型別
│           ├── auth.ts     # 認證型別
│           └── common.ts   # 通用工具型別
│
├── tools/                  # 開發工具與腳本
├── scripts/                # 建置與部署腳本
└── docs/                   # 文件
```

## 🛠️ 開發環境設定

### 先決條件
- Node.js 18+ 與 pnpm 8+
- Go 1.21+
- Docker Desktop
- Git

### 初始設定

1. **複製專案並安裝：**
   ```bash
   git clone <repo-url>
   cd mono-guard
   pnpm install
   ```

2. **環境設定：**
   ```bash
   cp .env.example .env
   # 編輯 .env 檔案設定本機環境
   
   # 如果計劃部署至 Zeabur，可參考 Zeabur 格式設定
   # 詳見 .env.example 中的 Zeabur 設定區塊
   ```

3. **啟動開發基礎設施：**
   ```bash
   ./scripts/dev-start.sh
   ```

### 開發工作流程

#### 前端開發
```bash
# 啟動前端開發伺服器
pnpm dev:frontend

# 執行測試
pnpm nx test frontend
pnpm nx test frontend --watch

# 端對端測試
pnpm nx e2e frontend-e2e

# 型別檢查
pnpm nx type-check frontend

# 程式碼檢查
pnpm nx lint frontend --fix
```

#### API 開發
```bash
# 啟動 API 伺服器
cd apps/api
go run cmd/server/main.go

# 執行測試
go test ./...

# 建置
go build -o bin/server cmd/server/main.go

# 格式化程式碼
go fmt ./...
```

#### CLI 開發
```bash
# 以監看模式啟動 CLI
pnpm dev:cli

# 本機測試 CLI
pnpm nx build cli
node dist/apps/cli/main.js --help

# 執行測試
pnpm nx test cli
```

## 📝 程式碼標準

### TypeScript/JavaScript
- 所有新程式碼使用 TypeScript
- 遵循 ESLint 設定
- 偏好函式型元件和 hooks
- 使用適當的型別（避免使用 `any`）
- 為公用 API 撰寫 JSDoc 註解

### Go
- 遵循 Go 最佳實務和慣用法
- 使用 `gofmt` 保持一致的格式
- 撰寫完整的測試
- 使用有意義的套件和函數名稱
- 明確處理錯誤

### 一般準則
- 撰寫自我記錄的程式碼
- 為新功能添加測試
- 需要時更新文件
- 使用慣例式提交訊息
- 保持函數簡短且專注

## 🧪 測試策略

### 前端測試
```bash
# 使用 Jest + React Testing Library 進行單元測試
pnpm nx test frontend

# 使用 Playwright 進行端對端測試
pnpm nx e2e frontend-e2e --ui

# 元件測試
# 測試檔案與元件放在一起或在 __tests__ 資料夾中
```

### API 測試
```bash
# 單元測試
go test ./internal/...

# 整合測試
go test ./tests/integration/...

# 帶覆蓋率的測試
go test -cover ./...
```

### CLI 測試
```bash
# 單元測試
pnpm nx test cli

# 整合測試
pnpm nx test cli --testPathPattern=integration
```

## 🔧 建置與部署

### 本機建置
```bash
# 建置所有應用程式
pnpm build

# 建置特定應用程式
pnpm nx build frontend
pnpm nx build cli
cd apps/api && go build -o bin/server cmd/server/main.go
```

### Docker 開發
```bash
# 以開發模式啟動所有服務
docker-compose up -d

# 查看日誌
docker-compose logs -f frontend
docker-compose logs -f api

# 重建服務
docker-compose build frontend
```

### 正式環境部署

#### 推薦：Zeabur 部署
```bash
# 準備 Zeabur 部署
./scripts/setup-zeabur.sh

# 推送至 GitHub 並連接到 Zeabur
# 詳細步驟請參考 docs/ZEABUR_DEPLOYMENT.md
```

#### 替代：Docker 部署
```bash
# 建置正式環境映像檔
docker-compose -f docker-compose.prod.yml build

# 部署至正式環境
./scripts/prod-deploy.sh
```

## 🎯 狀態管理

### 前端狀態管理
- **Zustand** 用於全域狀態（使用者、UI 偏好設定）
- **React Query** 用於伺服器狀態（API 資料、快取）
- **React hooks** 用於本機元件狀態
- **Context** 節制使用，主要用於主題

### 狀態組織
```typescript
// store/auth.ts
interface AuthStore {
  user: User | null;
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => void;
}

// store/ui.ts
interface UIStore {
  theme: 'light' | 'dark';
  sidebarOpen: boolean;
  toggleTheme: () => void;
}
```

## 🛡️ 安全性考量

### 認證與授權
- 使用 JWT 權杖進行 API 認證
- 使用 NextAuth.js 進行前端認證
- 角色型存取控制（RBAC）
- 安全的權杖儲存

### API 安全性
- 使用 Zod 綱要進行輸入驗證
- SQL 注入防護
- 速率限制
- CORS 設定
- 安全標頭

### 環境安全性
- 環境變數驗證
- 金鑰管理
- 程式碼中不寫死金鑰
- 安全的 Docker 映像檔實務

## 🚀 效能最佳化

### 前端效能
- Next.js 自動最佳化
- 按路由進行程式碼分割
- 使用 Next.js Image 進行圖像最佳化
- 使用 webpack-bundle-analyzer 進行套件分析
- 對昂貴的元件使用 React.memo

### API 效能
- 資料庫連線池
- 對頻繁查詢使用 Redis 快取
- 使用高效的 Go routines 進行並行操作
- 查詢最佳化
- 回應壓縮

### CLI 效能
- 高效的檔案解析
- 適當時使用平行處理
- 長時間操作的進度指示器
- 記憶體效率的資料結構

## 🔍 除錯

### 前端除錯
- React Developer Tools
- Next.js 內建除錯
- 瀏覽器開發工具
- VS Code 除錯設定

### API 除錯
- Go 除錯器（delve）
- 結構化日誌記錄
- 健康檢查端點
- 指標收集

### CLI 除錯
- Node.js 除錯
- 詳細日誌標誌
- 進度指示器
- 錯誤處理與回報

## 📊 監控與可觀測性

### 健康檢查
- API：`/health` 端點
- 前端：自訂健康檢查頁面
- 資料庫：連線監控
- Redis：Ping 檢查

### 日誌記錄
- 結構化 JSON 日誌
- 日誌層級（debug、info、warn、error）
- 請求/回應日誌
- 錯誤追蹤

### 指標（未來）
- 應用程式效能指標
- 業務指標（分析次數等）
- 基礎設施指標
- 使用者行為分析

## 🤝 貢獻準則

### 工作流程
1. 從 `main` 建立功能分支
2. 實作功能並添加測試
3. 確保所有 CI 檢查通過
4. 提交 pull request
5. 程式碼審查並合併

### 提交訊息
使用慣例式提交：
```
feat: 新增依賴分析儀表板
fix: 修正循環依賴偵測
docs: 更新 API 文件
test: 新增 CLI 整合測試
refactor: 簡化認證流程
```

### Pull Request 流程
1. 填寫 PR 範本
2. 為新功能加入測試
3. 需要時更新文件
4. 確保 CI 通過
5. 向維護者請求審查

### 程式碼審查檢核表
- [ ] 程式碼遵循專案標準
- [ ] 包含測試且通過
- [ ] 文件已更新
- [ ] 無安全性漏洞
- [ ] 已考慮效能問題
- [ ] 錯誤處理適當

## 🗂️ 專案管理

### 議題範本
- Bug 回報
- 功能請求
- 文件改進
- 效能問題

### 標籤
- `bug` - 功能無法正常運作
- `enhancement` - 新功能或請求
- `documentation` - 文件改進
- `good first issue` - 適合新手
- `help wanted` - 需要額外關注

### 里程碑
按版本發布組織，具備明確的功能集和時程估計。

## 🔗 實用連結

- [Next.js 文件](https://nextjs.org/docs)
- [Go 文件](https://golang.org/doc/)
- [Nx 文件](https://nx.dev)
- [Docker 文件](https://docs.docker.com)
- [PostgreSQL 文件](https://www.postgresql.org/docs/)
- [Redis 文件](https://redis.io/documentation)

---

需要協助嗎？查看 [README](../README.md) 或建立議題！