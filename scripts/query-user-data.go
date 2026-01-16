package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type ProjectStats struct {
	TotalProjects      int     `json:"total_projects"`
	UniqueUsers        int     `json:"unique_users"`
	CompletedProjects  int     `json:"completed_projects"`
	InProgressProjects int     `json:"in_progress_projects"`
	FailedProjects     int     `json:"failed_projects"`
	AvgHealthScore     float64 `json:"avg_health_score"`
	AnalyzedProjects   int     `json:"analyzed_projects"`
}

type UserActivity struct {
	OwnerID        string  `json:"owner_id"`
	TotalProjects  int     `json:"total_projects"`
	CompletedCount int     `json:"completed_count"`
	AvgHealthScore float64 `json:"avg_health_score"`
	LastActive     string  `json:"last_active"`
	FirstActive    string  `json:"first_active"`
}

type AnalysisStats struct {
	TotalAnalyses     int     `json:"total_analyses"`
	Completed         int     `json:"completed"`
	Failed            int     `json:"failed"`
	AnalyzedProjects  int     `json:"analyzed_projects"`
	AvgDurationSec    float64 `json:"avg_duration_seconds"`
}

func main() {
	// 從環境變數獲取 DATABASE_URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// 連接資料庫
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 測試連線
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✅ Successfully connected to Railway PostgreSQL")
	fmt.Println("\n" + "="*60)

	// 1. 總體統計
	fmt.Println("\n📊 總體統計")
	fmt.Println("-" * 60)
	stats := getProjectStats(db)
	printJSON(stats)

	// 2. 使用者活躍度
	fmt.Println("\n👥 活躍使用者排名 (Top 10)")
	fmt.Println("-" * 60)
	users := getUserActivity(db)
	for i, user := range users {
		fmt.Printf("%d. Owner: %s | Projects: %d | Completed: %d | Avg Health: %.1f | Last Active: %s\n",
			i+1, user.OwnerID, user.TotalProjects, user.CompletedCount, user.AvgHealthScore, user.LastActive)
	}

	// 3. 分析統計
	fmt.Println("\n🔍 依賴分析統計")
	fmt.Println("-" * 60)
	analysisStats := getAnalysisStats(db)
	printJSON(analysisStats)

	// 4. 最近 7 天活動
	fmt.Println("\n📅 最近 7 天新建專案")
	fmt.Println("-" * 60)
	recentProjects := getRecentProjects(db)
	for _, p := range recentProjects {
		fmt.Printf("- %s | Status: %s | Health: %d | Created: %s\n",
			p["name"], p["status"], p["health_score"], p["created_at"])
	}

	// 5. 回訪率
	fmt.Println("\n🔄 使用者回訪統計")
	fmt.Println("-" * 60)
	retention := getRetentionStats(db)
	for _, r := range retention {
		fmt.Printf("%s: %d users\n", r["frequency"], r["count"])
	}

	fmt.Println("\n" + "="*60)
	fmt.Println("✅ 數據提取完成！")
}

func getProjectStats(db *sql.DB) ProjectStats {
	var stats ProjectStats
	query := `
		SELECT
			COUNT(*) as total_projects,
			COUNT(DISTINCT owner_id) as unique_users,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_projects,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress_projects,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_projects,
			COALESCE(ROUND(AVG(health_score), 2), 0) as avg_health_score,
			COUNT(CASE WHEN last_analysis_at IS NOT NULL THEN 1 END) as analyzed_projects
		FROM projects
	`
	err := db.QueryRow(query).Scan(
		&stats.TotalProjects,
		&stats.UniqueUsers,
		&stats.CompletedProjects,
		&stats.InProgressProjects,
		&stats.FailedProjects,
		&stats.AvgHealthScore,
		&stats.AnalyzedProjects,
	)
	if err != nil {
		log.Printf("Error querying project stats: %v", err)
	}
	return stats
}

func getUserActivity(db *sql.DB) []UserActivity {
	query := `
		SELECT
			owner_id,
			COUNT(*) as total_projects,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_count,
			COALESCE(ROUND(AVG(health_score), 2), 0) as avg_health_score,
			MAX(created_at)::text as last_active,
			MIN(created_at)::text as first_active
		FROM projects
		WHERE owner_id IS NOT NULL AND owner_id != ''
		GROUP BY owner_id
		ORDER BY total_projects DESC
		LIMIT 10
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying user activity: %v", err)
		return nil
	}
	defer rows.Close()

	var users []UserActivity
	for rows.Next() {
		var user UserActivity
		err := rows.Scan(&user.OwnerID, &user.TotalProjects, &user.CompletedCount,
			&user.AvgHealthScore, &user.LastActive, &user.FirstActive)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}
		users = append(users, user)
	}
	return users
}

func getAnalysisStats(db *sql.DB) AnalysisStats {
	var stats AnalysisStats
	query := `
		SELECT
			COUNT(*) as total_analyses,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
			COUNT(DISTINCT project_id) as analyzed_projects,
			COALESCE(ROUND(AVG(EXTRACT(EPOCH FROM (completed_at - started_at))), 2), 0) as avg_duration_seconds
		FROM dependency_analyses
		WHERE completed_at IS NOT NULL
	`
	err := db.QueryRow(query).Scan(
		&stats.TotalAnalyses,
		&stats.Completed,
		&stats.Failed,
		&stats.AnalyzedProjects,
		&stats.AvgDurationSec,
	)
	if err != nil {
		log.Printf("Error querying analysis stats: %v", err)
	}
	return stats
}

func getRecentProjects(db *sql.DB) []map[string]interface{} {
	query := `
		SELECT
			name,
			status,
			health_score,
			created_at::text
		FROM projects
		WHERE created_at >= NOW() - INTERVAL '7 days'
		ORDER BY created_at DESC
		LIMIT 10
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying recent projects: %v", err)
		return nil
	}
	defer rows.Close()

	var projects []map[string]interface{}
	for rows.Next() {
		var name, status, createdAt string
		var healthScore int
		err := rows.Scan(&name, &status, &healthScore, &createdAt)
		if err != nil {
			log.Printf("Error scanning project row: %v", err)
			continue
		}
		projects = append(projects, map[string]interface{}{
			"name":         name,
			"status":       status,
			"health_score": healthScore,
			"created_at":   createdAt,
		})
	}
	return projects
}

func getRetentionStats(db *sql.DB) []map[string]interface{} {
	query := `
		SELECT
			CASE
				WHEN project_count = 1 THEN '僅使用 1 次'
				WHEN project_count BETWEEN 2 AND 3 THEN '使用 2-3 次'
				WHEN project_count BETWEEN 4 AND 10 THEN '使用 4-10 次'
				ELSE '使用 10+ 次'
			END as frequency,
			COUNT(*) as count
		FROM (
			SELECT owner_id, COUNT(*) as project_count
			FROM projects
			WHERE owner_id IS NOT NULL AND owner_id != ''
			GROUP BY owner_id
		) user_stats
		GROUP BY frequency
		ORDER BY count DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error querying retention stats: %v", err)
		return nil
	}
	defer rows.Close()

	var stats []map[string]interface{}
	for rows.Next() {
		var frequency string
		var count int
		err := rows.Scan(&frequency, &count)
		if err != nil {
			log.Printf("Error scanning retention row: %v", err)
			continue
		}
		stats = append(stats, map[string]interface{}{
			"frequency": frequency,
			"count":     count,
		})
	}
	return stats
}

func printJSON(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
