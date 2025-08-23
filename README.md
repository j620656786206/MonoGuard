# MonoGuard

> 全方位的 Monorepo 架構分析與驗證工具

MonoGuard 是一個強大的平台，專為分析、驗證和維護 monorepo 架構健康狀況而設計。提供相依性分析、架構驗證和即時監控功能，確保您的 monorepo 專案保持可擴展性和可維護性。

## 🏗️ 專案架構

此 monorepo 採用 **整合套件策略**，將所有邏輯保留在各個應用程式內，僅共享型別定義：

```
mono-guard/
├── apps/
│   ├── frontend/        # Next.js 14 網頁應用程式
│   ├── cli/             # Node.js 指令列工具
│   └── api/             # Go API 服務
├── libs/
│   └── shared-types/    # 共用 TypeScript 型別定義
├── tools/               # 開發工具
└── docs/                # 文件
```

## 🚀 快速開始

### 先決條件

- **Node.js** 18+ 與 **pnpm** 8+
- **Go** 1.21+
- **Docker** 與 **Docker Compose**
- **PostgreSQL** 15+ (本機開發不使用 Docker 時需要)
- **Redis** 7+ (本機開發不使用 Docker 時需要)

### 開發環境設定

1. **複製專案並安裝相依套件：**
   ```bash
   git clone <repository-url>
   cd mono-guard
   pnpm install
   ```

2. **啟動開發環境：**
   ```bash
   ./scripts/dev-start.sh
   ```
   此腳本將會：
   - 在 Docker 中啟動 PostgreSQL 和 Redis
   - 從範本建立 `.env` 檔案（請依需求編輯）
   - 建置共用型別定義
   - 顯示後續步驟

3. **啟動應用程式：**
   ```bash
   # 前端 (Next.js)
   pnpm dev:frontend

   # API (Go) - 在另一個終端機
   cd apps/api
   go run cmd/server/main.go

   # CLI (Node.js) - 在另一個終端機
   pnpm dev:cli
   ```

### 正式環境部署

#### Zeabur 部署 (推薦)

```bash
# 準備 Zeabur 部署
./scripts/setup-zeabur.sh
```

然後將專案推送到 GitHub 並連接到 [Zeabur](https://zeabur.com)。詳細步驟請參見 [docs/ZEABUR_DEPLOYMENT.md](docs/ZEABUR_DEPLOYMENT.md)。

#### Docker 部署

```bash
# 設定必要的環境變數
export DB_PASSWORD="your-secure-password"
export JWT_SECRET="your-jwt-secret"
export NEXTAUTH_SECRET="your-nextauth-secret"

# 部署至正式環境
./scripts/prod-deploy.sh
```

## 📱 應用程式

### 🌐 前端 (`apps/frontend`)
- **技術：** Next.js 14 with App Router
- **樣式：** Tailwind CSS + Shadcn/UI
- **狀態管理：** Zustand + React Query
- **測試：** Jest + Playwright
- **埠號：** 3000

**主要功能：**
- 專案儀表板與健康指標
- 互動式相依性分析
- 架構驗證報告
- 即時監控

### ⚡ 指令列工具 (`apps/cli`)
- **技術：** Node.js with TypeScript
- **框架：** Commander.js
- **埠號：** N/A (指令列工具)

**可用指令：**
```bash
monoguard init        # 初始化專案設定
monoguard analyze     # 執行相依性分析
monoguard validate    # 驗證架構規則
```

### 🚀 API (`apps/api`)
- **技術：** Go with Gin framework
- **資料庫：** PostgreSQL
- **快取：** Redis
- **埠號：** 8080

**主要端點：**
- `GET /health` - 服務健康檢查
- `GET /api/v1/projects` - 專案管理
- `GET /api/v1/analysis` - 分析結果
- `GET /api/v1/dependencies` - 相依性資料

## 🛠️ 開發指令

### 工作區指令
```bash
# 開發
pnpm dev                    # 啟動所有開發伺服器
pnpm dev:frontend          # 僅啟動前端
pnpm dev:cli               # 以監看模式啟動 CLI

# 建置
pnpm build                 # 建置所有應用程式
pnpm build:frontend        # 僅建置前端
pnpm build:cli            # 僅建置 CLI

# 測試
pnpm test                  # 執行所有測試
pnpm test:watch           # 以監看模式執行測試
pnpm test:coverage        # 執行測試並產生覆蓋率報告
pnpm test:e2e             # 執行端對端測試

# 程式碼檢查與格式化
pnpm lint                  # 檢查所有程式碼
pnpm lint:fix             # 修正程式碼檢查問題
pnpm type-check           # TypeScript 型別檢查

# 相依性管理
pnpm clean                # 清理所有建置產物
pnpm graph                # 顯示相依性圖表
```

### Docker 指令
```bash
# 開發環境
docker-compose up -d              # 啟動基礎設施
docker-compose down               # 停止基礎設施
docker-compose logs -f            # 查看日誌

# 正式環境
docker-compose -f docker-compose.prod.yml up -d    # 啟動正式環境
docker-compose -f docker-compose.prod.yml down     # 停止正式環境
```

## 📊 可用服務

### 開發環境
| 服務 | URL | 說明 |
|---------|-----|-------------|
| 前端 | http://localhost:3000 | Next.js 網頁應用程式 |
| API | http://localhost:8080 | Go API 伺服器 |
| 資料庫 | localhost:5432 | PostgreSQL 資料庫 |
| Redis | localhost:6379 | Redis 快取 |
| Adminer | http://localhost:8081 | 資料庫管理工具 |

### 正式環境
所有服務都在可設定的埠號上運行，並具備適當的健康檢查和監控。

## 🔧 設定

### 環境變數
複製 `.env.example` 為 `.env` 並進行設定：

```bash
# 資料庫
DATABASE_URL=postgresql://monoguard:password@localhost:5432/monoguard
REDIS_URL=redis://localhost:6379

# API
API_URL=http://localhost:8080
JWT_SECRET=your-secret-key

# 前端
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXTAUTH_SECRET=your-nextauth-secret

# 選用：GitHub OAuth
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
```

## 🏛️ 架構原則

### 整合套件策略
- **應用程式包含所有邏輯** - 無共用業務邏輯程式庫
- **最小化共用程式碼** - 僅在 `libs/shared-types` 中共用型別定義
- **獨立部署** - 每個應用程式可獨立部署
- **明確界限** - 服務間具備完整定義的介面

### 技術選擇
- **前端：** Next.js 用於現代化 React 開發，支援 SSR/SSG
- **後端：** Go 提供高效能與併發處理
- **CLI：** Node.js 具備生態系統相容性
- **資料庫：** PostgreSQL 提供可靠性與進階功能
- **快取：** Redis 用於會話與分析快取
- **Monorepo：** Nx 用於建置最佳化與開發者體驗

## 🧪 測試策略

### 前端測試
- **單元測試：** Jest + React Testing Library
- **端對端測試：** Playwright
- **覆蓋率目標：** 80%+

### API 測試
- **單元測試：** Go 內建測試
- **整合測試：** 使用測試資料庫
- **負載測試：** 計劃於 v1.0 實作

### CLI 測試
- **單元測試：** Jest
- **整合測試：** 針對實際專案
- **指令測試：** 所有 CLI 指令

## 📈 效能考量

- **前端：** 程式碼分割、圖像最佳化、快取
- **API：** 連線池、查詢最佳化、快取
- **CLI：** 高效檔案處理、平行分析
- **基礎設施：** Docker 健康檢查、資源限制

## 🤝 貢獻

1. Fork 此儲存庫
2. 建立功能分支
3. 進行變更
4. 為新功能添加測試
5. 確保所有測試通過
6. 提交 pull request

### 開發指引
- 遵循 TypeScript 最佳實務
- 撰寫完整的測試
- 使用慣例式提交訊息
- 適時更新文件
- 確保程式碼通過檢查與型別檢查

## 📝 授權條款

MIT License - 詳見 [LICENSE](LICENSE) 檔案。

## 🆘 支援

- **文件：** 查看 `docs/` 目錄
- **問題回報：** 使用 GitHub Issues 回報 bug
- **討論：** 使用 GitHub Discussions 進行提問

---

由 MonoGuard 團隊用 ❤️ 打造