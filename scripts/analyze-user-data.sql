-- MonoGuard 使用者數據分析查詢
-- 在 Railway PostgreSQL 執行此查詢以了解使用者行為

-- ==========================================
-- 1. 總體使用統計
-- ==========================================
SELECT
    '總體統計' as metric_category,
    COUNT(DISTINCT id) as total_projects,
    COUNT(DISTINCT owner_id) as unique_users,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_projects,
    COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_projects,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_projects,
    ROUND(AVG(health_score), 2) as avg_health_score,
    COUNT(CASE WHEN last_analysis_at IS NOT NULL THEN 1 END) as analyzed_projects
FROM projects;

-- ==========================================
-- 2. 每日活躍趨勢（最近 30 天）
-- ==========================================
SELECT
    DATE(created_at) as date,
    COUNT(*) as projects_created,
    COUNT(DISTINCT owner_id) as unique_users
FROM projects
WHERE created_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- ==========================================
-- 3. 使用者活躍度排名
-- ==========================================
SELECT
    owner_id,
    COUNT(*) as total_projects,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_count,
    ROUND(AVG(health_score), 2) as avg_health_score,
    MAX(created_at) as last_active,
    MIN(created_at) as first_active
FROM projects
WHERE owner_id IS NOT NULL AND owner_id != ''
GROUP BY owner_id
ORDER BY total_projects DESC
LIMIT 20;

-- ==========================================
-- 4. 分析完成率
-- ==========================================
SELECT
    status,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER (), 2) as percentage
FROM projects
GROUP BY status
ORDER BY count DESC;

-- ==========================================
-- 5. 健康分數分布
-- ==========================================
SELECT
    CASE
        WHEN health_score >= 80 THEN '優秀 (80-100)'
        WHEN health_score >= 60 THEN '良好 (60-79)'
        WHEN health_score >= 40 THEN '普通 (40-59)'
        WHEN health_score >= 20 THEN '需改善 (20-39)'
        ELSE '嚴重問題 (0-19)'
    END as health_category,
    COUNT(*) as count,
    ROUND(AVG(health_score), 2) as avg_score
FROM projects
WHERE health_score > 0
GROUP BY health_category
ORDER BY avg_score DESC;

-- ==========================================
-- 6. 依賴分析數據統計
-- ==========================================
SELECT
    '依賴分析統計' as metric,
    COUNT(*) as total_analyses,
    COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
    COUNT(DISTINCT project_id) as analyzed_projects,
    ROUND(AVG(EXTRACT(EPOCH FROM (completed_at - started_at))), 2) as avg_duration_seconds
FROM dependency_analyses
WHERE completed_at IS NOT NULL;

-- ==========================================
-- 7. 最近 7 天的活躍專案
-- ==========================================
SELECT
    p.id,
    p.name,
    p.status,
    p.health_score,
    p.created_at,
    p.last_analysis_at,
    COUNT(da.id) as analysis_count
FROM projects p
LEFT JOIN dependency_analyses da ON p.id = da.project_id
WHERE p.created_at >= NOW() - INTERVAL '7 days'
GROUP BY p.id, p.name, p.status, p.health_score, p.created_at, p.last_analysis_at
ORDER BY p.created_at DESC;

-- ==========================================
-- 8. 發現的問題類型統計（需要查看 JSONB 內容）
-- ==========================================
SELECT
    COUNT(*) as total_completed_analyses,
    COUNT(CASE WHEN results IS NOT NULL THEN 1 END) as analyses_with_results
FROM dependency_analyses
WHERE status = 'completed';

-- ==========================================
-- 9. 回訪率（使用者是否重複使用）
-- ==========================================
SELECT
    CASE
        WHEN project_count = 1 THEN '僅使用 1 次'
        WHEN project_count BETWEEN 2 AND 3 THEN '使用 2-3 次'
        WHEN project_count BETWEEN 4 AND 10 THEN '使用 4-10 次'
        ELSE '使用 10+ 次'
    END as usage_frequency,
    COUNT(*) as user_count
FROM (
    SELECT owner_id, COUNT(*) as project_count
    FROM projects
    WHERE owner_id IS NOT NULL AND owner_id != ''
    GROUP BY owner_id
) user_stats
GROUP BY usage_frequency
ORDER BY user_count DESC;

-- ==========================================
-- 10. 最活躍使用者的專案列表（可用於主動聯繫）
-- ==========================================
SELECT
    p.owner_id,
    p.id as project_id,
    p.name as project_name,
    p.repository_url,
    p.status,
    p.health_score,
    p.created_at,
    COUNT(da.id) as analysis_count
FROM projects p
LEFT JOIN dependency_analyses da ON p.id = da.project_id
WHERE p.owner_id IN (
    SELECT owner_id
    FROM projects
    WHERE owner_id IS NOT NULL AND owner_id != ''
    GROUP BY owner_id
    HAVING COUNT(*) >= 2
)
GROUP BY p.owner_id, p.id, p.name, p.repository_url, p.status, p.health_score, p.created_at
ORDER BY p.owner_id, p.created_at DESC;
