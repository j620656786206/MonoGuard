# 部署指南

本指南涵蓋 MonoGuard 的部署策略、基礎設施需求和營運程序。

## 🏗️ 基礎架構

```
┌─────────────────────────────────────────────────┐
│                   負載平衡器                     │
└─────────────────┬───────────────────────────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        ▼                   ▼
┌──────────────┐    ┌──────────────┐
│     前端     │    │     API      │
│  (Next.js)   │    │    (Go)      │
│   埠號 3000   │    │  埠號 8080   │
└──────┬───────┘    └──────┬───────┘
       │                   │
       └─────────┬─────────┘
                 │
         ┌───────┴───────┐
         │               │
         ▼               ▼
   ┌──────────┐    ┌──────────┐
   │PostgreSQL│    │  Redis   │
   │ 埠號 5432 │    │ 埠號 6379 │
   └──────────┘    └──────────┘
```

## 🚀 部署選項

### 1. Render 部署（推薦）

Render 是現代化的雲端平台，提供簡單的部署流程、自動 HTTPS 和免費方案。

#### 快速部署

1. 將專案推送至 GitHub
2. 在 [Render Dashboard](https://dashboard.render.com) 建立新服務
3. 連接 GitHub 儲存庫
4. 設定環境變數
5. 部署！

#### 優勢
- ✅ **一鍵部署** - GitHub 整合自動部署
- ✅ **免費方案** - 適合小型專案和測試
- ✅ **自動 HTTPS** - 免費 SSL 憑證
- ✅ **PostgreSQL** - 託管資料庫服務
- ✅ **自訂網域** - 支援自訂域名

### 2. Docker Compose（自托管）

適合需要完全控制基礎設施的團隊。

#### 正式環境部署
```bash
# 設定必要的環境變數
export DB_PASSWORD="secure-db-password"
export JWT_SECRET="your-jwt-secret-key"
export API_URL="https://api.monoguard-web.onrender.com"

# 部署
./scripts/prod-deploy.sh
```

#### 手動 Docker Compose
```bash
# 建立正式環境設定檔
cp .env.example .env.production
# 編輯 .env.production 並填入正式環境數值

# 啟動正式環境堆疊
docker-compose -f docker-compose.prod.yml up -d

# 檢查服務健康狀況
docker-compose -f docker-compose.prod.yml ps
docker-compose -f docker-compose.prod.yml logs -f
```

### 3. Kubernetes 部署

建立 Kubernetes 資源清單：

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: monoguard
---
# k8s/postgres.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  namespace: monoguard
spec:
  replicas: 1
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          value: monoguard
        - name: POSTGRES_USER
          value: monoguard
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-data
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-data
        persistentVolumeClaim:
          claimName: postgres-pvc
```

部署至 Kubernetes：
```bash
# 套用資源清單
kubectl apply -f k8s/

# 檢查狀態
kubectl get pods -n monoguard
kubectl logs -f deployment/api -n monoguard
```

### 4. 雲端服務商部署

#### AWS ECS/Fargate
```bash
# 建置並推送映像檔
docker build -t monoguard-api apps/api/
docker tag monoguard-api:latest <account-id>.dkr.ecr.<region>.amazonaws.com/monoguard-api:latest
docker push <account-id>.dkr.ecr.<region>.amazonaws.com/monoguard-api:latest

# 使用 ECS CLI 或 CDK/Terraform 部署
```

#### Google Cloud Run
```bash
# 建置並部署 API
gcloud builds submit --tag gcr.io/$PROJECT_ID/monoguard-api apps/api/
gcloud run deploy monoguard-api --image gcr.io/$PROJECT_ID/monoguard-api --platform managed

# 建置並部署前端
gcloud builds submit --tag gcr.io/$PROJECT_ID/monoguard-frontend apps/frontend/
gcloud run deploy monoguard-frontend --image gcr.io/$PROJECT_ID/monoguard-frontend --platform managed
```

#### Azure Container Instances
```bash
# 建立資源群組
az group create --name MonoGuardRG --location eastus

# 部署容器
az container create --resource-group MonoGuardRG --name monoguard-api --image monoguard-api:latest
az container create --resource-group MonoGuardRG --name monoguard-frontend --image monoguard-frontend:latest
```

## 🔧 環境設定

### 必要環境變數

#### 資料庫設定
```bash
DB_HOST=postgres                    # 資料庫主機
DB_PORT=5432                       # 資料庫埠號
DB_NAME=monoguard                  # 資料庫名稱
DB_USER=monoguard                  # 資料庫使用者
DB_PASSWORD=secure-password        # 資料庫密碼（必填）
DB_SSLMODE=require                 # 正式環境 SSL 模式
```

#### API 設定
```bash
PORT=8080                          # API 伺服器埠號
GIN_MODE=release                   # Gin 框架模式
JWT_SECRET=your-jwt-secret         # JWT 簽章金鑰（必填）
CORS_ORIGINS=https://yourdomain.com # 允許的 CORS 來源
```

#### Redis 設定
```bash
REDIS_HOST=redis                   # Redis 主機
REDIS_PORT=6379                    # Redis 埠號
REDIS_PASSWORD=redis-password      # Redis 密碼
REDIS_DB=0                         # Redis 資料庫編號
```

#### 前端設定
```bash
VITE_API_URL=https://api.yourdomain.com    # API URL
VITE_APP_ENV=production                    # 應用程式環境
```

### 正式環境安全設定

#### SSL/TLS 憑證
使用 Let's Encrypt 自動續期：
```bash
# 使用 Caddy 反向代理
sudo docker run -d \
  --name caddy \
  -p 80:80 -p 443:443 \
  -v caddy_data:/data \
  -v caddy_config:/config \
  -v $PWD/Caddyfile:/etc/caddy/Caddyfile \
  caddy:latest
```

#### Caddyfile 範例：
```
yourdomain.com {
    reverse_proxy frontend:3000
}

api.yourdomain.com {
    reverse_proxy api:8080
}
```

## 📊 監控與日誌

### 健康檢查
所有服務都包含健康檢查端點：

- **API**: `GET /health`
- **前端**: `GET /api/health`
- **資料庫**: 連線檢查
- **Redis**: Ping 檢查

### 日誌設定
```bash
# 設定日誌層級
LOG_LEVEL=info                     # debug, info, warn, error
LOG_FORMAT=json                    # text 或 json

# 日誌聚合（選用）
LOGSTASH_HOST=logstash.yourdomain.com
LOGSTASH_PORT=5044
```

### 監控堆疊（選用）
使用 Prometheus 與 Grafana 部署監控：
```bash
# 加入至 docker-compose.prod.yml
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml
      
  grafana:
    image: grafana/grafana
    ports:
      - "3001:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

## 🔄 資料庫管理

### 資料遷移
```bash
# 執行資料庫遷移
cd apps/api
go run cmd/migrate/main.go up

# 回滾遷移
go run cmd/migrate/main.go down 1
```

### 備份策略
```bash
#!/bin/bash
# 資料庫備份腳本
DATE=$(date +%Y%m%d_%H%M%S)
docker exec postgres pg_dump -U monoguard monoguard > backup_$DATE.sql

# 使用 cron 自動備份
0 2 * * * /path/to/backup-script.sh
```

### 資料庫還原
```bash
# 從備份還原
docker exec -i postgres psql -U monoguard -d monoguard < backup_20240101_020000.sql
```

## 📈 擴展性考量

### 水平擴展
- **前端**: 負載平衡器後方多個 Next.js 實例
- **API**: 負載平衡的多個 Go API 實例
- **資料庫**: 具備讀取副本的 PostgreSQL
- **快取**: 高可用性的 Redis 叢集

### 負載平衡
使用 nginx 或雲端負載平衡器：
```nginx
upstream frontend {
    server frontend-1:3000;
    server frontend-2:3000;
    server frontend-3:3000;
}

upstream api {
    server api-1:8080;
    server api-2:8080;
    server api-3:8080;
}

server {
    listen 80;
    server_name yourdomain.com;
    location / {
        proxy_pass http://frontend;
    }
}

server {
    listen 80;
    server_name api.yourdomain.com;
    location / {
        proxy_pass http://api;
    }
}
```

### 資源需求

#### 最低需求
- **CPU**: 2 核心
- **RAM**: 4GB
- **儲存空間**: 20GB
- **網路**: 100 Mbps

#### 建議正式環境規格
- **CPU**: 4+ 核心
- **RAM**: 8GB+
- **儲存空間**: 100GB+ SSD
- **網路**: 1 Gbps

## 🚨 災難復原

### 備份策略
1. **資料庫備份**: 每日完整備份，每小時增量備份
2. **設定備份**: 版本控制的環境設定檔
3. **磁碟區備份**: Docker 磁碟區與持久性資料
4. **程式碼備份**: Git 儲存庫與發布標籤

### 復原程序
```bash
# 1. 還原資料庫
docker exec -i postgres psql -U monoguard -d monoguard < latest_backup.sql

# 2. 還原設定
cp backup/.env.production .env

# 3. 重啟服務
docker-compose -f docker-compose.prod.yml restart

# 4. 驗證健康狀況
curl -f http://localhost:8080/health
curl -f http://localhost:3000/api/health
```

## 🔍 疑難排解

### 常見問題

#### 服務無法啟動
```bash
# 檢查日誌
docker-compose logs service-name

# 檢查資源使用
docker stats

# 檢查網路連線
docker exec container-name ping other-service
```

#### 資料庫連線問題
```bash
# 測試資料庫連線
docker exec api-container nc -zv postgres 5432

# 檢查資料庫日誌
docker logs postgres-container

# 手動連線資料庫
docker exec -it postgres-container psql -U monoguard -d monoguard
```

#### 效能問題
```bash
# 監控資源使用
docker stats

# 檢查 API 回應時間
curl -w "@curl-format.txt" -o /dev/null -s "http://localhost:8080/health"

# 資料庫效能
docker exec postgres-container pg_stat_activity
```

### 日誌分析
```bash
# 即時檢視所有日誌
docker-compose -f docker-compose.prod.yml logs -f

# 依服務篩選
docker-compose -f docker-compose.prod.yml logs -f api

# 搜尋錯誤
docker-compose -f docker-compose.prod.yml logs | grep ERROR
```

## 📋 部署檢核表

### 部署前
- [ ] 環境變數已設定
- [ ] 金鑰妥善保護
- [ ] SSL 憑證已取得
- [ ] 資料庫已備份
- [ ] 健康檢查已設定
- [ ] 監控系統已建立

### 部署中
- [ ] 本機建置並測試映像檔
- [ ] 部署至測試環境
- [ ] 執行整合測試
- [ ] 部署至正式環境
- [ ] 驗證所有服務健康
- [ ] 測試關鍵使用者流程

### 部署後
- [ ] 監控日誌是否有錯誤
- [ ] 檢查效能指標
- [ ] 驗證備份系統
- [ ] 更新文件
- [ ] 通知團隊部署成功

## 📞 支援與維護

### 定期維護工作
- 每月更新相依套件
- 每季檢視與輪換金鑰
- 監控磁碟使用量並清理日誌
- 每週檢視效能指標
- 每月測試備份/還原程序

### 緊急聯絡人
維護以下聯絡人清單：
- 基礎設施團隊聯絡人
- 資料庫管理員
- 資安團隊
- 雲端服務商支援

### 文件更新
持續更新部署文件，包含：
- 環境變更
- 新設定選項
- 事件處理經驗
- 效能最佳化發現

---

如需其他協助，請參考[開發指南](DEVELOPMENT.md)或在儲存庫中建立議題。