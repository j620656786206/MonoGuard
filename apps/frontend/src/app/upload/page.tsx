'use client';

import React, { useState, useCallback, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft } from 'lucide-react';
import { PackageJsonFile, ComprehensiveAnalysisResult } from '@monoguard/shared-types';
import { FileProcessingResult } from '@monoguard/shared-types';
import { FileUpload } from '@/components/common/FileUpload';
import { VirtualizedList } from '@/components/ui/VirtualizedList';
import { ProgressSteps, Step } from '@/components/ui/ProgressSteps';
import { AnalysisResults } from '@/components/analysis/AnalysisResults';
import AnalysisService from '@/lib/api/services/analysis';

export default function UploadPage() {
  const router = useRouter();
  const [uploadResult, setUploadResult] = useState<FileProcessingResult | null>(null);
  const [uploadErrors, setUploadErrors] = useState<string[]>([]);
  const [currentStep, setCurrentStep] = useState<'upload' | 'analyze' | 'report'>('upload');
  const [analysisResult, setAnalysisResult] = useState<ComprehensiveAnalysisResult | null>(null);
  const [analysisError, setAnalysisError] = useState<string | null>(null);

  const steps: Step[] = useMemo(() => [
    {
      id: 'upload',
      title: 'File Upload',
      description: 'Upload project files',
      status: uploadResult ? 'completed' : currentStep === 'upload' ? 'current' : 'pending'
    },
    {
      id: 'analyze',
      title: 'Dependency Analysis',
      description: 'Analyze project dependencies',
      status: currentStep === 'analyze' ? 'current' : analysisResult ? 'completed' : 'pending'
    },
    {
      id: 'report',
      title: 'Generate Report',
      description: 'Generate analysis report',
      status: analysisResult ? 'completed' : currentStep === 'report' ? 'current' : 'pending'
    }
  ], [uploadResult, currentStep, analysisResult]);

  const handleUploadComplete = useCallback((result: FileProcessingResult) => {
    setUploadResult(result);
    setUploadErrors([]);
  }, []);

  const handleUploadError = useCallback((errors: string[]) => {
    setUploadErrors(errors);
    setUploadResult(null);
  }, []);

  const handleNewUpload = useCallback(() => {
    setUploadResult(null);
    setUploadErrors([]);
    setAnalysisResult(null);
    setAnalysisError(null);
    setCurrentStep('upload');
  }, []);

  const handleAnalyzeDependencies = useCallback(async () => {
    if (!uploadResult || !uploadResult.uploadId) return;
    
    setCurrentStep('analyze');
    setAnalysisError(null);
    
    try {
      // Start comprehensive analysis
      const uploadId = uploadResult.uploadId.toString();
      const response = await AnalysisService.startComprehensiveAnalysis(uploadId);
      
      // Poll for completion
      if (!response.data.id) {
        throw new Error('Invalid response: missing analysis ID');
      }
      
      const completedAnalysis = await AnalysisService.pollAnalysisCompletion(response.data.id.toString());
      
      setAnalysisResult(completedAnalysis);
      setCurrentStep('report');
    } catch (error) {
      console.error('Analysis failed:', error);
      setAnalysisError(error instanceof Error ? error.message : 'Analysis failed');
      setCurrentStep('upload'); // Go back to upload step on error
    }
  }, [uploadResult]);

  const handleCheckArchitecture = useCallback(async () => {
    // For now, this will do the same as handleAnalyzeDependencies
    // Later we can add separate architecture-only analysis
    await handleAnalyzeDependencies();
  }, [handleAnalyzeDependencies]);

  const handleBackToDashboard = useCallback(() => {
    router.push('/dashboard');
  }, [router]);

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-4 mb-4">
          <button
            onClick={handleBackToDashboard}
            className="flex items-center gap-2 px-3 py-2 text-sm text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            aria-label="Back to Dashboard"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Dashboard
          </button>
        </div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Upload Project Files
        </h1>
        <p className="text-gray-600">
          Upload your project files (.zip archives or package.json files) to analyze dependencies and architecture.
        </p>
      </div>

      {/* Progress Steps */}
      <div className="mb-8">
        <ProgressSteps steps={steps} />
      </div>

      {/* File Upload Section */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 mb-8">
        <FileUpload
          onUploadComplete={handleUploadComplete}
          onUploadError={handleUploadError}
          accept={['.zip', '.json']}
          multiple={true}
          maxFiles={10}
        />
      </div>

      {/* Upload Results */}
      {uploadResult && !analysisResult && (
        <UploadResultsDisplay 
          result={uploadResult} 
          onNewUpload={handleNewUpload}
        />
      )}

      {/* Analysis Results */}
      {analysisResult && (
        <AnalysisResults 
          analysis={analysisResult}
          onNewAnalysis={handleNewUpload}
        />
      )}

      {/* Analysis Error */}
      {analysisError && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 mb-8">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <h3 className="text-sm font-medium text-red-800">Analysis Failed</h3>
              <div className="mt-2 text-sm text-red-700">
                <p>{analysisError}</p>
              </div>
              <div className="mt-4">
                <button
                  onClick={() => setAnalysisError(null)}
                  className="bg-red-100 px-3 py-2 rounded-md text-sm font-medium text-red-800 hover:bg-red-200 transition-colors"
                >
                  Dismiss
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Next Steps */}
      {uploadResult && !analysisResult && (
        <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg p-6 mb-8">
          <h3 className="text-lg font-semibold text-gray-900 mb-3">
            Next Steps
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <button 
              onClick={handleAnalyzeDependencies}
              disabled={!uploadResult || !uploadResult.uploadId || currentStep !== 'upload'}
              className={`text-left p-4 bg-white rounded-lg border transition-all ${
                !uploadResult || !uploadResult.uploadId || currentStep !== 'upload'
                  ? 'border-gray-200 text-gray-400 cursor-not-allowed'
                  : 'border-gray-200 hover:border-blue-300 hover:shadow-sm'
              }`}
            >
              <div className="flex items-center justify-between">
                <div>
                  <div className="font-medium mb-1">
                    Analyze Dependencies
                  </div>
                  <div className="text-sm text-gray-600">
                    Get detailed dependency analysis and health scores
                  </div>
                </div>
                {currentStep === 'analyze' && (
                  <div className="animate-spin h-5 w-5 border-2 border-blue-600 border-t-transparent rounded-full" />
                )}
              </div>
            </button>
            
            <button 
              onClick={handleCheckArchitecture}
              disabled={!uploadResult || !uploadResult.uploadId || currentStep !== 'upload'}
              className={`text-left p-4 bg-white rounded-lg border transition-all ${
                !uploadResult || !uploadResult.uploadId || currentStep !== 'upload'
                  ? 'border-gray-200 text-gray-400 cursor-not-allowed'
                  : 'border-gray-200 hover:border-blue-300 hover:shadow-sm'
              }`}
            >
              <div className="flex items-center justify-between">
                <div>
                  <div className="font-medium mb-1">
                    Check Architecture
                  </div>
                  <div className="text-sm text-gray-600">
                    Validate architecture and detect circular dependencies
                  </div>
                </div>
                {currentStep === 'analyze' && (
                  <div className="animate-spin h-5 w-5 border-2 border-blue-600 border-t-transparent rounded-full" />
                )}
              </div>
            </button>
          </div>
        </div>
      )}

      {/* Instructions */}
      {!uploadResult && (
        <div className="bg-blue-50 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-blue-900 mb-3">
          Upload Instructions
        </h3>
        <ul className="space-y-2 text-blue-800 text-sm">
          <li className="flex items-start">
            <span className="flex-shrink-0 w-2 h-2 bg-blue-400 rounded-full mt-2 mr-3"></span>
            <span>
              <strong>ZIP Files:</strong> Upload compressed project archives. We'll automatically scan for package.json files inside.
            </span>
          </li>
          <li className="flex items-start">
            <span className="flex-shrink-0 w-2 h-2 bg-blue-400 rounded-full mt-2 mr-3"></span>
            <span>
              <strong>Package.json Files:</strong> Upload individual package.json files for quick analysis.
            </span>
          </li>
          <li className="flex items-start">
            <span className="flex-shrink-0 w-2 h-2 bg-blue-400 rounded-full mt-2 mr-3"></span>
            <span>
              <strong>File Size Limit:</strong> Maximum 50MB per file. You can upload up to 10 files at once.
            </span>
          </li>
          <li className="flex items-start">
            <span className="flex-shrink-0 w-2 h-2 bg-blue-400 rounded-full mt-2 mr-3"></span>
            <span>
              <strong>Analysis:</strong> After upload, you can analyze dependencies, detect circular dependencies, and validate architecture.
            </span>
          </li>
        </ul>
        </div>
      )}
    </div>
  );
}

// Upload Results Display Component
interface UploadResultsDisplayProps {
  result: FileProcessingResult;
  onNewUpload: () => void;
}

const UploadResultsDisplay: React.FC<UploadResultsDisplayProps> = ({ 
  result, 
  onNewUpload 
}) => {
  return (
    <div className="space-y-6 mb-8">
      {/* Summary Card */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-semibold text-gray-900">
            Upload Summary
          </h2>
          <button
            onClick={onNewUpload}
            className="text-sm bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
          >
            Upload New Files
          </button>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-gradient-to-br from-gray-50 to-gray-100 rounded-lg p-4 border border-gray-200 hover:shadow-sm transition-shadow">
            <div className="flex items-center justify-between">
              <div>
                <div className="text-3xl font-bold text-gray-900 mb-1">
                  {result.files.length}
                </div>
                <div className="text-sm text-gray-600 font-medium">Files Uploaded</div>
              </div>
              <div className="p-3 bg-gray-200 rounded-full">
                <svg className="w-6 h-6 text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
            </div>
          </div>
          
          <div className="bg-gradient-to-br from-green-50 to-green-100 rounded-lg p-4 border border-green-200 hover:shadow-sm transition-shadow">
            <div className="flex items-center justify-between">
              <div>
                <div className="text-3xl font-bold text-green-700 mb-1">
                  {result.packageJsonFiles.length}
                </div>
                <div className="text-sm text-green-600 font-medium">Package.json Found</div>
              </div>
              <div className="p-3 bg-green-200 rounded-full">
                <svg className="w-6 h-6 text-green-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                </svg>
              </div>
            </div>
          </div>
          
          <div className="bg-gradient-to-br from-blue-50 to-blue-100 rounded-lg p-4 border border-blue-200 hover:shadow-sm transition-shadow">
            <div className="flex items-center justify-between">
              <div>
                <div className="text-3xl font-bold text-blue-700 mb-1">
                  {(result.files.reduce((acc: number, file: any) => acc + file.fileSize, 0) / (1024 * 1024)).toFixed(2)}MB
                </div>
                <div className="text-sm text-blue-600 font-medium">Total Size</div>
              </div>
              <div className="p-3 bg-blue-200 rounded-full">
                <svg className="w-6 h-6 text-blue-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 3.79 4 8.5 4s8.5-1.79 8.5-4V7M4 7c0 2.21 3.79 4 8.5 4s8.5-1.79 8.5-4M4 7c0-2.21 3.79-4 8.5-4s8.5 1.79 8.5 4" />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Package.json Files */}
      {result.packageJsonFiles.length > 0 && (
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            Discovered Package.json Files
          </h3>
          
          <div className="space-y-4">
            {result.packageJsonFiles.map((pkg: any, index: number) => (
              <PackageJsonCard key={index} packageJson={pkg} />
            ))}
          </div>
        </div>
      )}

    </div>
  );
};

// Package.json Card Component
interface PackageJsonCardProps {
  packageJson: PackageJsonFile;
}

const PackageJsonCard: React.FC<PackageJsonCardProps> = ({ packageJson }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [activeTab, setActiveTab] = useState<'all' | 'dependencies' | 'devDependencies'>('all');
  
  const dependencyCount = packageJson.dependencies 
    ? Object.keys(packageJson.dependencies).length 
    : 0;
  
  const devDependencyCount = packageJson.devDependencies 
    ? Object.keys(packageJson.devDependencies).length 
    : 0;

  // Filter dependencies based on search term and active tab
  const filterDependencies = (deps: Record<string, string>, searchTerm: string) => {
    if (!deps) return [];
    
    return Object.entries(deps).filter(([name, version]) =>
      name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      version.toLowerCase().includes(searchTerm.toLowerCase())
    );
  };

  const filteredDependencies = filterDependencies(packageJson.dependencies || {}, searchTerm);
  const filteredDevDependencies = filterDependencies(packageJson.devDependencies || {}, searchTerm);

  const allDependencies = [
    ...Object.entries(packageJson.dependencies || {}).map(([name, version]) => ({ name, version, type: 'dependency' as const })),
    ...Object.entries(packageJson.devDependencies || {}).map(([name, version]) => ({ name, version, type: 'devDependency' as const }))
  ].filter(dep =>
    dep.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    dep.version.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getDisplayDependencies = () => {
    switch (activeTab) {
      case 'dependencies':
        return filteredDependencies.map(([name, version]) => ({ name, version, type: 'dependency' as const }));
      case 'devDependencies':
        return filteredDevDependencies.map(([name, version]) => ({ name, version, type: 'devDependency' as const }));
      default:
        return allDependencies;
    }
  };

  const displayDependencies = getDisplayDependencies();

  return (
    <div className="border border-gray-200 rounded-lg p-4">
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-3 mb-2">
            <h4 className="font-medium text-gray-900">
              {packageJson.name || 'Unnamed Package'}
            </h4>
            {packageJson.version && (
              <span className="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded">
                v{packageJson.version}
              </span>
            )}
          </div>
          
          <div className="text-sm text-gray-600 mb-2">
            Path: {packageJson.path}
          </div>
          
          <div className="flex gap-4 text-sm text-gray-500">
            {dependencyCount > 0 && (
              <span>{dependencyCount} dependencies</span>
            )}
            {devDependencyCount > 0 && (
              <span>{devDependencyCount} dev dependencies</span>
            )}
          </div>
        </div>
        
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="text-sm text-blue-600 hover:text-blue-700"
        >
          {isExpanded ? 'Show Less' : 'Show More'}
        </button>
      </div>
      
      {isExpanded && (
        <div className="mt-4 pt-4 border-t border-gray-200">
          {/* Search and Filter Controls */}
          <div className="mb-4 space-y-3">
            {/* Search Bar */}
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input
                type="text"
                placeholder="Search dependencies..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
              {searchTerm && (
                <button
                  onClick={() => setSearchTerm('')}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  <svg className="h-4 w-4 text-gray-400 hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              )}
            </div>

            {/* Tab Filter */}
            <div className="flex space-x-1 bg-gray-100 rounded-lg p-1">
              <button
                onClick={() => setActiveTab('all')}
                className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
                  activeTab === 'all'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                All ({allDependencies.length})
              </button>
              <button
                onClick={() => setActiveTab('dependencies')}
                className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
                  activeTab === 'dependencies'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                Dependencies ({dependencyCount})
              </button>
              <button
                onClick={() => setActiveTab('devDependencies')}
                className={`px-3 py-1 text-sm font-medium rounded-md transition-colors ${
                  activeTab === 'devDependencies'
                    ? 'bg-white text-blue-600 shadow-sm'
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                Dev Dependencies ({devDependencyCount})
              </button>
            </div>
          </div>

          {/* Dependencies List */}
          <div className="space-y-2">
            {displayDependencies.length > 0 ? (
              <>
                <div className="text-sm text-gray-500 mb-2">
                  Showing {displayDependencies.length} dependencies
                  {searchTerm && ` matching "${searchTerm}"`}
                </div>
                {displayDependencies.length > 50 ? (
                  <VirtualizedList
                    items={displayDependencies}
                    itemHeight={32}
                    containerHeight={256}
                    className="border border-gray-200 rounded-lg"
                    renderItem={({ name, version, type }) => (
                      <div className="flex items-center justify-between py-1 px-2 hover:bg-gray-50 h-8">
                        <div className="flex items-center space-x-2">
                          <span className={`inline-block w-2 h-2 rounded-full ${
                            type === 'dependency' ? 'bg-green-500' : 'bg-blue-500'
                          }`}></span>
                          <span className="text-gray-700 text-sm">{name}</span>
                        </div>
                        <span className="text-gray-500 text-xs font-mono">{version}</span>
                      </div>
                    )}
                  />
                ) : (
                  <div className="max-h-64 overflow-y-auto space-y-1">
                    {displayDependencies.map(({ name, version, type }) => (
                      <div key={`${type}-${name}`} className="flex items-center justify-between py-1 px-2 rounded hover:bg-gray-50">
                        <div className="flex items-center space-x-2">
                          <span className={`inline-block w-2 h-2 rounded-full ${
                            type === 'dependency' ? 'bg-green-500' : 'bg-blue-500'
                          }`}></span>
                          <span className="text-gray-700 text-sm">{name}</span>
                        </div>
                        <span className="text-gray-500 text-xs font-mono">{version}</span>
                      </div>
                    ))}
                  </div>
                )}
              </>
            ) : (
              <div className="text-center py-8 text-gray-500">
                <svg className="mx-auto h-8 w-8 text-gray-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                <p className="text-sm">No dependencies found matching "{searchTerm}"</p>
                <button
                  onClick={() => setSearchTerm('')}
                  className="text-blue-600 hover:text-blue-700 text-sm underline mt-1"
                >
                  Clear search
                </button>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};