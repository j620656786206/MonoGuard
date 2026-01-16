-- ====================================================================
-- MonoGuard 使用者數據分析查詢
-- 請在 Railway Web Dashboard 的 PostgreSQL Query 頁面執行以下查詢
-- ====================================================================

-- 查詢 1: 總體統計 (最重要！)
-- ====================================================================
SELECT
    COUNT(*) as "總專案數",
    COUNT(DISTINCT owner_id) as "獨立使用者數",
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as "已完成專案",
    COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as "進行中專案",
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as "失敗專案",
    ROUND(AVG(health_score), 2) as "平均健康分數",
    COUNT(CASE WHEN last_analysis_at IS NOT NULL THEN 1 END) as "已分析專案"
FROM projects;

-- 查詢 2: 活躍使用者排名
-- ====================================================================
SELECT
    owner_id as "使用者ID",
    COUNT(*) as "專案數量",
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as "完成數",
    ROUND(AVG(health_score), 2) as "平均健康分數",
    MAX(created_at) as "最後活動時間",
    MIN(created_at) as "首次使用時間"
FROM projects
WHERE owner_id IS NOT NULL AND owner_id != ''
GROUP BY owner_id
ORDER BY COUNT(*) DESC
LIMIT 20;

-- 查詢 3: 最近 7 天新建專案
-- ====================================================================
SELECT
    id as "專案ID",
    name as "專案名稱",
    status as "狀態",
    health_score as "健康分數",
    created_at as "建立時間",
    owner_id as "使用者ID"
FROM projects
WHERE created_at >= NOW() - INTERVAL '7 days'
ORDER BY created_at DESC;

-- 查詢 4: 使用者回訪率分析
-- ====================================================================
SELECT
    CASE
        WHEN project_count = 1 THEN '僅使用 1 次'
        WHEN project_count BETWEEN 2 AND 3 THEN '使用 2-3 次'
        WHEN project_count BETWEEN 4 AND 10 THEN '使用 4-10 次'
        ELSE '使用 10+ 次'
    END as "使用頻率",
    COUNT(*) as "使用者數量",
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as "百分比"
FROM (
    SELECT owner_id, COUNT(*) as project_count
    FROM projects
    WHERE owner_id IS NOT NULL AND owner_id != ''
    GROUP BY owner_id
) user_stats
GROUP BY "使用頻率"
ORDER BY "使用者數量" DESC;

-- 查詢 5: 依賴分析統計
-- ====================================================================
SELECT
    COUNT(*) as "總分析次數",
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as "完成",
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as "失敗",
    COUNT(DISTINCT project_id) as "分析過的專案數"
FROM dependency_analyses;

-- 查詢 6: 專案狀態分布
-- ====================================================================
SELECT
    status as "專案狀態",
    COUNT(*) as "數量",
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as "百分比"
FROM projects
GROUP BY status
ORDER BY COUNT(*) DESC;

-- 查詢 7: 健康分數分布
-- ====================================================================
SELECT
    CASE
        WHEN health_score >= 80 THEN '優秀 (80-100)'
        WHEN health_score >= 60 THEN '良好 (60-79)'
        WHEN health_score >= 40 THEN '普通 (40-59)'
        WHEN health_score >= 20 THEN '需改善 (20-39)'
        WHEN health_score > 0 THEN '嚴重問題 (1-19)'
        ELSE '未評分 (0)'
    END as "健康等級",
    COUNT(*) as "專案數量",
    ROUND(AVG(health_score), 2) as "平均分數"
FROM projects
GROUP BY "健康等級"
ORDER BY "平均分數" DESC;

-- 查詢 8: 每日新增專案趨勢 (最近 30 天)
-- ====================================================================
SELECT
    DATE(created_at) as "日期",
    COUNT(*) as "新增專案數",
    COUNT(DISTINCT owner_id) as "新增使用者數"
FROM projects
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at)
ORDER BY "日期" DESC;
