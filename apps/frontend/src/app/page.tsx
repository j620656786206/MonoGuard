'use client';

import { useState } from 'react';
import { HeroSection } from '../components/landing/HeroSection';
import { FeaturesSection } from '../components/landing/FeaturesSection';
import { SampleResults } from '../components/landing/SampleResults';
import { Footer } from '../components/landing/Footer';
import { EmailSignup } from '../components/landing/EmailSignup';
import { FileUpload } from '../components/common/FileUpload';
import { AnalysisResults } from '../components/analysis/AnalysisResults';
import { FileProcessingResult } from '@monoguard/shared-types';
import { apiClient } from '../lib/api/client';


export default function LandingPage() {
  const [error, setError] = useState<string | null>(null);
  const [uploadResult, setUploadResult] = useState<FileProcessingResult | null>(null);
  const [analysisResult, setAnalysisResult] = useState<any | null>(null);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisProgress, setAnalysisProgress] = useState('');


  const handleUploadComplete = async (result: FileProcessingResult) => {
    setUploadResult(result);
    
    // Trigger analysis after successful upload
    setIsAnalyzing(true);
    setError(null);
    setAnalysisProgress('Starting analysis...');
    
    try {
      const startResponse = await apiClient.post(`/api/v1/analysis/comprehensive/${result.id}`);
      
      const analysisId = startResponse.data?.id;
      if (!analysisId) {
        throw new Error('No analysis ID returned');
      }
      
      // Poll for analysis results
      const pollForResults = async () => {
        const maxAttempts = 30; // 30 attempts = 1 minute max
        let attempts = 0;
        
        while (attempts < maxAttempts) {
          try {
            setAnalysisProgress(`Analyzing... (${attempts + 1}/${maxAttempts})`);
            const resultResponse = await apiClient.get(`/api/v1/analysis/dependencies/${analysisId}`);
            
            if (resultResponse.data && resultResponse.data.status === 'completed') {
              setAnalysisProgress('Analysis completed!');
              
              // Format the data for AnalysisResults component
              
              // Map the API response to the expected format
              const formattedResult = {
                id: resultResponse.data.id,
                status: resultResponse.data.status,
                results: {
                  // Map the actual API data structure correctly
                  healthScore: resultResponse.data.results?.summary?.healthScore || 0,
                  summary: resultResponse.data.results?.summary || {},
                  dependencies: resultResponse.data.results?.dependencies || [],
                  architecture: resultResponse.data.results?.architecture || {},
                  duplicates: resultResponse.data.results?.duplicateDependencies || [],
                  circularDependencies: resultResponse.data.results?.circularDependencies || [],
                  versionConflicts: resultResponse.data.results?.versionConflicts || [],
                  bundleImpact: resultResponse.data.results?.bundleImpact || {},
                  unusedDependencies: resultResponse.data.results?.unusedDependencies || [],
                  // Include all original results data
                  ...resultResponse.data.results
                },
                metadata: resultResponse.data.metadata || {},
                project: resultResponse.data.project || {},
                // Include the original response for debugging
                _originalData: resultResponse
              };
              
              setAnalysisResult(formattedResult);
              return;
            }
            
            // Wait 2 seconds before next poll
            await new Promise(resolve => setTimeout(resolve, 2000));
            attempts++;
            
          } catch (pollError) {
            attempts++;
            await new Promise(resolve => setTimeout(resolve, 2000));
          }
        }
        
        throw new Error('Analysis timed out - please try again');
      };
      
      await pollForResults();
      
    } catch (err: any) {
      console.error('Analysis failed:', err);
      const errorMessage = err.response?.data?.error || err.message || 'Analysis failed';
      setError(`Analysis failed: ${errorMessage}`);
    } finally {
      setIsAnalyzing(false);
    }
  };

  const handleUploadError = (errors: string[]) => {
    setError(errors.join(', '));
  };



  // Show analysis progress while analyzing
  if (isAnalyzing) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {/* Header with back button */}
          <div className="mb-8">
            <button
              onClick={() => {
                setIsAnalyzing(false);
                setUploadResult(null);
                setError(null);
                setAnalysisProgress('');
              }}
              className="inline-flex items-center text-indigo-600 hover:text-indigo-700 font-medium"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Back to Upload
            </button>
          </div>

          {/* Analysis Progress */}
          <div className="bg-white rounded-lg shadow-sm border p-8 text-center">
            <div className="w-16 h-16 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin mx-auto mb-6"></div>
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Analysis in Progress</h2>
            <p className="text-gray-600 text-lg mb-4">{analysisProgress}</p>
            <div className="bg-gray-200 rounded-full h-2 mx-auto max-w-xs">
              <div className="bg-indigo-600 h-2 rounded-full animate-pulse" style={{width: '60%'}}></div>
            </div>
            <p className="text-gray-500 text-sm mt-4">
              We're analyzing your package.json files for dependencies, vulnerabilities, and architecture issues.
            </p>
          </div>
        </div>
      </div>
    );
  }

  // Show analysis results if available
  if (analysisResult) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {/* Header with back button */}
          <div className="mb-8">
            <button
              onClick={() => {
                setAnalysisResult(null);
                setUploadResult(null);
                setError(null);
              }}
              className="inline-flex items-center text-indigo-600 hover:text-indigo-700 font-medium"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Back to Upload
            </button>
          </div>

          {/* Analysis Results */}
          <AnalysisResults 
            analysis={analysisResult}
            onNewAnalysis={() => {
              setAnalysisResult(null);
              setUploadResult(null);
            }}
          />
        </div>
      </div>
    );
  }

  // Main landing page
  return (
    <div className="min-h-screen bg-white">
      {/* Hero Section */}
      <div id="hero-section">
        <HeroSection />
      </div>

      {/* Upload Error Display */}
      {error && (
        <div className="bg-white py-8">
          <div className="max-w-2xl mx-auto px-4 text-center">
            <div className="bg-red-50 border border-red-200 rounded-lg p-6">
              <svg className="w-8 h-8 text-red-600 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h3 className="text-lg font-semibold text-red-900 mb-2">Upload Failed</h3>
              <p className="text-red-700 mb-4">{error}</p>
              <button
                onClick={() => setError(null)}
                className="bg-red-600 hover:bg-red-700 text-white font-medium py-2 px-4 rounded transition-colors"
              >
                Try Again
              </button>
            </div>
          </div>
        </div>
      )}

      {/* File Upload Section */}
      <section id="upload-section" className="py-16 bg-gray-50">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold text-gray-900 mb-4">
              Upload & Analyze Your Files
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              Drag and drop your package.json files or entire monorepo ZIP archives for instant analysis
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-start">
            {/* File Upload Component */}
            <div>
              <FileUpload
                onUploadComplete={handleUploadComplete}
                onUploadError={handleUploadError}
                accept={['.json', '.zip']}
                maxFiles={5}
                className="w-full"
              />
            </div>

            {/* Privacy & Benefits */}
            <div className="space-y-6">
              <div className="bg-white rounded-lg p-6 shadow-sm border">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                  <svg className="w-5 h-5 text-green-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                  Why Choose Local Analysis?
                </h3>
                <ul className="space-y-3 text-gray-700">
                  <li className="flex items-start">
                    <svg className="w-4 h-4 text-green-600 mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span><strong>Privacy Protection</strong> - Files analyzed securely and automatically deleted after processing</span>
                  </li>
                  <li className="flex items-start">
                    <svg className="w-4 h-4 text-green-600 mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span><strong>Instant Analysis</strong> - No waiting, analysis starts immediately after upload</span>
                  </li>
                  <li className="flex items-start">
                    <svg className="w-4 h-4 text-green-600 mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span><strong>No API Limits</strong> - Not affected by GitHub API rate limiting</span>
                  </li>
                  <li className="flex items-start">
                    <svg className="w-4 h-4 text-green-600 mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span><strong>Private Code Friendly</strong> - Analyze internal and private projects safely</span>
                  </li>
                </ul>
              </div>

              {/* Upload Success Message */}
              {uploadResult && (
                <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                  <div className="flex items-center">
                    <svg className="w-5 h-5 text-green-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <h4 className="text-green-800 font-medium">Upload Successful!</h4>
                  </div>
                  <p className="text-green-700 text-sm mt-1">
                    Found {uploadResult.packageJsonFiles.length} package.json {uploadResult.packageJsonFiles.length === 1 ? 'file' : 'files'}.
                    {isAnalyzing ? ' Analyzing...' : ' Analysis complete!'}
                  </p>
                  <button
                    onClick={() => setUploadResult(null)}
                    className="text-green-600 hover:text-green-700 text-sm underline mt-2"
                  >
                    Upload More Files
                  </button>
                </div>
              )}

              {/* Sample Files */}
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h4 className="text-blue-800 font-medium mb-2">Need test files?</h4>
                <p className="text-blue-700 text-sm mb-3">
                  Download our sample files to try MonoGuard's features
                </p>
                <button
                  onClick={() => {
                    // Download sample package.json
                    const samplePackageJson = {
                      "name": "sample-monorepo",
                      "version": "1.0.0",
                      "dependencies": {
                        "react": "^18.2.0",
                        "lodash": "^4.17.21",
                        "axios": "^1.6.0"
                      },
                      "devDependencies": {
                        "@types/react": "^18.2.0",
                        "typescript": "^5.0.0",
                        "lodash": "^4.17.20"
                      }
                    };
                    
                    if (typeof window !== 'undefined') {
                      const blob = new Blob([JSON.stringify(samplePackageJson, null, 2)], { type: 'application/json' });
                      const url = URL.createObjectURL(blob);
                      const a = document.createElement('a');
                      a.href = url;
                      a.download = 'sample-package.json';
                      document.body.appendChild(a);
                      a.click();
                      document.body.removeChild(a);
                      URL.revokeObjectURL(url);
                    }
                  }}
                  className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium py-2 px-4 rounded transition-colors"
                >
                  Download Sample package.json
                </button>
              </div>

            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <FeaturesSection />

      {/* Sample Results Section */}
      <SampleResults />

      {/* Email Signup Section */}
      <section className="py-16 bg-gray-50">
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
          <EmailSignup 
            title="Stay Updated with MonoGuard"
            description="Get the latest features, security insights, and monorepo best practices delivered to your inbox."
            buttonText="Subscribe Now"
            variant="default"
          />
        </div>
      </section>

      {/* Footer */}
      <Footer />
    </div>
  );
}