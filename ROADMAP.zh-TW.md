# MonoGuard 發展路線圖

[English](ROADMAP.md) | [繁體中文](ROADMAP.zh-TW.md)

本文件概述 MonoGuard 的計劃功能和改進項目。

## 圖例

- ✅ 已完成
- 🚧 進行中
- 📋 已規劃
- 💡 評估中

---

## 第一階段：核心分析（已完成）

✅ **相依性圖解析**
- 解析 monorepo 工作區的 package.json 檔案
- 支援 npm、yarn、pnpm 和 Nx 工作區
- 建立記憶體內相依性圖

✅ **循環相依性偵測**
- 偵測直接循環（A → B → A）
- 偵測間接循環（A → B → C → A）
- 嚴重性分類（critical、warning、info）
- 影響評估和修復建議

✅ **健康分數計算**
- 整體健康分數（0-100）
- 按類別分解（相依性、架構、可維護性）
- 趨勢追蹤

✅ **D3.js 視覺化**
- 力導向圖佈局
- 循環相依性高亮顯示
- 縮放、平移和小地圖導航
- 節點展開/收合（適用於大型圖）
- 混合 SVG/Canvas 渲染以提升效能

✅ **報告匯出**
- HTML 獨立報告
- JSON 用於 CI 整合
- Markdown 用於 PR 描述

---

## 第二階段：增強分析（進行中）

🚧 **WebAssembly 分析器**
- 使用 Go 編譯為 WASM 的客戶端分析
- 基本分析不需要伺服器
- 隱私優先：檔案永不離開瀏覽器

📋 **架構驗證**
- 定義層級規則（domain、application、infrastructure）
- 偵測層級違規
- 自訂規則設定

📋 **Bundle 影響分析**
- 識別重複相依性
- 計算浪費的 bundle 大小
- 建議整合策略

📋 **版本衝突偵測**
- 找出衝突的相依性版本
- 衝突風險評估
- 解決方案建議

---

## 第三階段：整合（已規劃）

📋 **GitHub 整合**
- 直接從 GitHub URL 分析儲存庫
- PR 留言顯示分析結果
- CI/CD 狀態檢查

📋 **CI/CD 整合**
- GitHub Actions 工作流程
- GitLab CI 範本
- 可設定的門檻和品質閘門

📋 **CLI 工具**
- 從命令列執行本機分析
- JSON 輸出供腳本使用
- 開發用監看模式

---

## 第四階段：進階功能（未來）

💡 **VS Code 擴充功能**
- 即時循環相依性警告
- 內嵌視覺化
- 快速修復

💡 **歷史追蹤**
- 追蹤健康分數的時間變化
- 退化警報
- 趨勢報告

💡 **團隊協作**
- 共享工作區
- 留言和註解
- 問題指派

💡 **自訂規則引擎**
- 定義自訂驗證規則
- 規則市集
- 匯入/匯出設定

---

## 貢獻

我們歡迎貢獻！如果您想參與任何這些功能的開發，請：

1. 檢查是否已有相關 issue
2. 開啟新 issue 討論您的方法
3. 提交參考該 issue 的 PR

請參閱 [CONTRIBUTING.zh-TW.md](CONTRIBUTING.zh-TW.md) 了解指南。

---

## 意見回饋

有新功能的想法嗎？請開啟 [GitHub Discussion](https://github.com/user/monoguard/discussions) 或 [Issue](https://github.com/user/monoguard/issues)。
