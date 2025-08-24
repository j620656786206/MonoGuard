package services_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/monoguard/api/internal/services"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependencyAnalyzer_AnalyzeMonorepo(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := createTestMonorepo(t)
	defer os.RemoveAll(tempDir)

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress logs during testing

	analyzer := services.NewDependencyAnalyzer(logger)

	results, err := analyzer.AnalyzeMonorepo(context.Background(), tempDir, "test-project-id")

	require.NoError(t, err)
	assert.NotNil(t, results)
	assert.GreaterOrEqual(t, results.Summary.TotalPackages, 1)
	assert.GreaterOrEqual(t, results.Summary.HealthScore, 0.0)
	assert.LessOrEqual(t, results.Summary.HealthScore, 100.0)
}

// createTestMonorepo creates a temporary monorepo structure for testing
func createTestMonorepo(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "monoguard-test-*")
	require.NoError(t, err)

	// Create package.json files
	packageJSONs := map[string]string{
		"package.json": `{
			"name": "root",
			"private": true,
			"workspaces": ["packages/*"],
			"dependencies": {
				"lodash": "^4.17.21"
			},
			"devDependencies": {
				"typescript": "^4.9.5"
			}
		}`,
		"packages/app/package.json": `{
			"name": "@test/app",
			"version": "1.0.0",
			"dependencies": {
				"react": "^18.2.0",
				"lodash": "^4.17.21"
			}
		}`,
		"packages/lib/package.json": `{
			"name": "@test/lib",
			"version": "1.0.0",
			"dependencies": {
				"axios": "^1.4.0",
				"lodash": "^4.17.20"
			}
		}`,
	}

	for path, content := range packageJSONs {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create some TypeScript files to test usage detection
	tsFiles := map[string]string{
		"packages/app/src/index.ts": `
			import React from 'react';
			import _ from 'lodash';
			
			console.log(_.map([1, 2, 3], x => x * 2));
		`,
		"packages/lib/src/index.ts": `
			import axios from 'axios';
			
			export const fetchData = () => axios.get('/api/data');
		`,
	}

	for path, content := range tsFiles {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	return tempDir
}