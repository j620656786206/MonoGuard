#!/bin/bash

# MonoGuard Railway 資料庫查詢腳本
# 使用方式: railway run ./scripts/query-railway-db.sh

echo "======================================================================"
echo "MonoGuard 使用者數據分析報告"
echo "======================================================================"
echo ""

# 設定 PostgreSQL 連線
export PGPASSWORD="$DB_PASSWORD"
HOST="$DB_HOST"
PORT="${DB_PORT:-5432}"
USER="$DB_USER"
DBNAME="$DB_NAME"

echo "正在連接到 Railway PostgreSQL..."
echo "Host: $HOST"
echo "Database: $DBNAME"
echo ""

# 檢查 psql 是否可用
if ! command -v psql &> /dev/null; then
    echo "❌ 錯誤: psql 未安裝"
    echo "請先安裝 PostgreSQL client"
    exit 1
fi

# 1. 總體統計
echo "======================================================================"
echo "📊 總體統計"
echo "======================================================================"
psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DBNAME" -c "
SELECT
    COUNT(*) as \"總專案數\",
    COUNT(DISTINCT owner_id) as \"獨立使用者\",
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as \"已完成\",
    COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as \"進行中\",
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as \"失敗\",
    ROUND(AVG(health_score), 2) as \"平均健康分數\"
FROM projects;
"

echo ""
echo "======================================================================"
echo "👥 活躍使用者排名 (Top 10)"
echo "======================================================================"
psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DBNAME" -c "
SELECT
    owner_id as \"使用者ID\",
    COUNT(*) as \"專案數量\",
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as \"完成數\",
    ROUND(AVG(health_score), 2) as \"平均分數\",
    MAX(created_at)::date as \"最後活動\"
FROM projects
WHERE owner_id IS NOT NULL AND owner_id != ''
GROUP BY owner_id
ORDER BY COUNT(*) DESC
LIMIT 10;
"

echo ""
echo "======================================================================"
echo "🔍 依賴分析統計"
echo "======================================================================"
psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DBNAME" -c "
SELECT
    COUNT(*) as \"總分析數\",
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as \"已完成\",
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as \"失敗\",
    COUNT(DISTINCT project_id) as \"分析過的專案\"
FROM dependency_analyses;
"

echo ""
echo "======================================================================"
echo "📅 最近 7 天新建專案"
echo "======================================================================"
psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DBNAME" -c "
SELECT
    name as \"專案名稱\",
    status as \"狀態\",
    health_score as \"健康分數\",
    created_at::date as \"建立日期\"
FROM projects
WHERE created_at >= NOW() - INTERVAL '7 days'
ORDER BY created_at DESC
LIMIT 15;
"

echo ""
echo "======================================================================"
echo "🔄 使用者回訪統計"
echo "======================================================================"
psql -h "$HOST" -p "$PORT" -U "$USER" -d "$DBNAME" -c "
SELECT
    CASE
        WHEN project_count = 1 THEN '僅使用 1 次'
        WHEN project_count BETWEEN 2 AND 3 THEN '使用 2-3 次'
        WHEN project_count BETWEEN 4 AND 10 THEN '使用 4-10 次'
        ELSE '使用 10+ 次'
    END as \"使用頻率\",
    COUNT(*) as \"使用者數量\"
FROM (
    SELECT owner_id, COUNT(*) as project_count
    FROM projects
    WHERE owner_id IS NOT NULL AND owner_id != ''
    GROUP BY owner_id
) user_stats
GROUP BY \"使用頻率\"
ORDER BY \"使用者數量\" DESC;
"

echo ""
echo "======================================================================"
echo "✅ 數據提取完成！"
echo "======================================================================"
