package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/monoguard/api/internal/models"
	"github.com/monoguard/api/internal/services"
	"github.com/sirupsen/logrus"
)

// GitHubHandler handles GitHub-related HTTP requests
type GitHubHandler struct {
	integratedAnalysis *services.IntegratedAnalysisService
	uploadService      *services.UploadService
	logger             *logrus.Logger
}

// NewGitHubHandler creates a new GitHub handler
func NewGitHubHandler(
	integratedAnalysis *services.IntegratedAnalysisService,
	uploadService *services.UploadService,
	logger *logrus.Logger,
) *GitHubHandler {
	return &GitHubHandler{
		integratedAnalysis: integratedAnalysis,
		uploadService:      uploadService,
		logger:             logger,
	}
}

// AnalyzeGitHubRequest represents the request body for analyzing GitHub repository
type AnalyzeGitHubRequest struct {
	URL                 string `json:"url" binding:"required"`
	IncludeCircular     bool   `json:"include_circular"`
	IncludeArchitecture bool   `json:"include_architecture"`
	SeverityThreshold   string `json:"severity_threshold"`
	ConfigPath          string `json:"config_path"`
}

// GitHubRepoInfo represents GitHub repository information
type GitHubRepoInfo struct {
	Owner string
	Repo  string
	Ref   string // branch, tag, or commit sha
}

// AnalyzeGitHubRepository handles POST /analysis/github
func (h *GitHubHandler) AnalyzeGitHubRepository(c *gin.Context) {
	var req AnalyzeGitHubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Warning("Invalid request body")
		BadRequest(c, "Invalid request body", err.Error())
		return
	}

	// Set defaults
	if req.SeverityThreshold == "" {
		req.SeverityThreshold = "warning"
	}

	// Parse GitHub URL
	repoInfo, err := h.parseGitHubURL(req.URL)
	if err != nil {
		h.logger.WithError(err).Warning("Invalid GitHub URL")
		BadRequest(c, "Invalid GitHub URL", err.Error())
		return
	}

	h.logger.WithFields(logrus.Fields{
		"owner":                repoInfo.Owner,
		"repo":                 repoInfo.Repo,
		"ref":                  repoInfo.Ref,
		"include_circular":     req.IncludeCircular,
		"include_architecture": req.IncludeArchitecture,
		"severity_threshold":   req.SeverityThreshold,
	}).Info("Starting GitHub repository analysis")

	// Fetch package.json files from GitHub
	processingResult, err := h.fetchPackageJsonFromGitHub(repoInfo)
	if err != nil {
		h.logger.WithError(err).Error("Failed to fetch package.json files from GitHub")
		InternalError(c, "Failed to fetch repository files")
		return
	}

	// Convert to analysis options
	options := services.AnalysisOptions{
		IncludeCircular:     req.IncludeCircular,
		IncludeArchitecture: req.IncludeArchitecture,
		SeverityThreshold:   req.SeverityThreshold,
		ConfigPath:          req.ConfigPath,
	}

	// Start analysis
	analysis, err := h.integratedAnalysis.AnalyzeProcessingResult(
		c.Request.Context(),
		processingResult.ID,
		options,
	)
	if err != nil {
		h.logger.WithError(err).Error("Failed to analyze GitHub repository")
		InternalError(c, "Failed to start analysis")
		return
	}

	// Return 202 Accepted with analysis ID
	c.JSON(202, gin.H{
		"success": true,
		"message": "GitHub repository analysis started successfully",
		"data":    analysis,
	})
}

// parseGitHubURL parses a GitHub URL and extracts repository information
func (h *GitHubHandler) parseGitHubURL(url string) (*GitHubRepoInfo, error) {
	// Regular expression to match GitHub URLs
	// Supports various formats:
	// - https://github.com/owner/repo
	// - https://github.com/owner/repo/tree/branch
	// - https://github.com/owner/repo/blob/branch/path
	re := regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)(?:/(?:tree|blob)/([^/]+))?`)
	matches := re.FindStringSubmatch(url)

	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid GitHub URL format")
	}

	repoInfo := &GitHubRepoInfo{
		Owner: matches[1],
		Repo:  matches[2],
		Ref:   "", // will be determined later by fetching repo info
	}

	// Remove .git suffix if present
	if strings.HasSuffix(repoInfo.Repo, ".git") {
		repoInfo.Repo = strings.TrimSuffix(repoInfo.Repo, ".git")
	}

	// If a specific ref is provided, use it
	if len(matches) > 3 && matches[3] != "" {
		repoInfo.Ref = matches[3]
	}

	return repoInfo, nil
}

// fetchPackageJsonFromGitHub fetches package.json files from GitHub repository
func (h *GitHubHandler) fetchPackageJsonFromGitHub(repoInfo *GitHubRepoInfo) (*models.FileProcessingResult, error) {
	// If ref is not specified, get the default branch
	if repoInfo.Ref == "" {
		defaultBranch, err := h.getDefaultBranch(repoInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to get default branch: %w", err)
		}
		repoInfo.Ref = defaultBranch
	}

	// Get the repository contents to find package.json files
	packageJsonFiles, err := h.findPackageJsonFiles(repoInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to find package.json files: %w", err)
	}

	if len(packageJsonFiles) == 0 {
		return nil, fmt.Errorf("no package.json files found in the repository")
	}

	// Create a file processing result
	processingResult, err := h.uploadService.CreateProcessingResult(packageJsonFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to create processing result: %w", err)
	}

	return processingResult, nil
}

// getDefaultBranch gets the default branch of a GitHub repository
func (h *GitHubHandler) getDefaultBranch(repoInfo *GitHubRepoInfo) (string, error) {
	repoURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", repoInfo.Owner, repoInfo.Repo)

	resp, err := http.Get(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch repository info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var repoResponse GitHubRepoResponse
	if err := json.Unmarshal(body, &repoResponse); err != nil {
		return "", fmt.Errorf("failed to parse repository response: %w", err)
	}

	return repoResponse.DefaultBranch, nil
}

// GitHubTreeResponse represents GitHub Tree API response
type GitHubTreeResponse struct {
	Tree []GitHubTreeItem `json:"tree"`
}

// GitHubTreeItem represents an item in GitHub tree
type GitHubTreeItem struct {
	Path string `json:"path"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
}

// GitHubBlobResponse represents GitHub Blob API response
type GitHubBlobResponse struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// GitHubRepoResponse represents GitHub Repository API response
type GitHubRepoResponse struct {
	DefaultBranch string `json:"default_branch"`
}

// findPackageJsonFiles recursively finds all package.json files in the repository
func (h *GitHubHandler) findPackageJsonFiles(repoInfo *GitHubRepoInfo) ([]models.PackageJsonFile, error) {
	// Get repository tree
	treeURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/%s?recursive=1", 
		repoInfo.Owner, repoInfo.Repo, repoInfo.Ref)

	resp, err := http.Get(treeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository tree: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var treeResponse GitHubTreeResponse
	if err := json.Unmarshal(body, &treeResponse); err != nil {
		return nil, fmt.Errorf("failed to parse tree response: %w", err)
	}

	var packageJsonFiles []models.PackageJsonFile
	
	// Find all package.json files
	for _, item := range treeResponse.Tree {
		if item.Type == "blob" && strings.HasSuffix(item.Path, "package.json") {
			// Fetch the content of this package.json file
			content, err := h.fetchFileContent(repoInfo, item.SHA)
			if err != nil {
				h.logger.WithError(err).WithField("path", item.Path).Warning("Failed to fetch package.json content")
				continue
			}

			// Parse package.json to extract name and version
			var pkgInfo struct {
				Name    string            `json:"name"`
				Version string            `json:"version"`
				Dependencies map[string]string `json:"dependencies"`
				DevDependencies map[string]string `json:"devDependencies"`
			}

			if err := json.Unmarshal([]byte(content), &pkgInfo); err != nil {
				h.logger.WithError(err).WithField("path", item.Path).Warning("Failed to parse package.json")
				continue
			}

			// Convert dependencies to JSON strings
			depsJSON := "{}"
			devDepsJSON := "{}"

			if pkgInfo.Dependencies != nil {
				if depsBytes, err := json.Marshal(pkgInfo.Dependencies); err == nil {
					depsJSON = string(depsBytes)
				}
			}

			if pkgInfo.DevDependencies != nil {
				if devDepsBytes, err := json.Marshal(pkgInfo.DevDependencies); err == nil {
					devDepsJSON = string(devDepsBytes)
				}
			}

			var name *string
			var version *string
			if pkgInfo.Name != "" {
				name = &pkgInfo.Name
			}
			if pkgInfo.Version != "" {
				version = &pkgInfo.Version
			}

			packageJsonFiles = append(packageJsonFiles, models.PackageJsonFile{
				Path:            item.Path,
				Content:         content,
				Name:            name,
				Version:         version,
				Dependencies:    depsJSON,
				DevDependencies: devDepsJSON,
			})
		}
	}

	return packageJsonFiles, nil
}

// fetchFileContent fetches the content of a specific file by its SHA
func (h *GitHubHandler) fetchFileContent(repoInfo *GitHubRepoInfo, sha string) (string, error) {
	blobURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/blobs/%s", 
		repoInfo.Owner, repoInfo.Repo, sha)

	resp, err := http.Get(blobURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch blob: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var blobResponse GitHubBlobResponse
	if err := json.Unmarshal(body, &blobResponse); err != nil {
		return "", fmt.Errorf("failed to parse blob response: %w", err)
	}

	// Decode base64 content
	if blobResponse.Encoding != "base64" {
		return "", fmt.Errorf("unsupported encoding: %s", blobResponse.Encoding)
	}

	content, err := base64.StdEncoding.DecodeString(blobResponse.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return string(content), nil
}