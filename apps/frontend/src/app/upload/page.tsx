'use client';

import React, { useState } from 'react';
import { FileProcessingResult, PackageJsonFile } from '@mono-guard/shared-types';
import { FileUpload } from '@/components/common/FileUpload';

export default function UploadPage() {
  const [uploadResult, setUploadResult] = useState<FileProcessingResult | null>(null);
  const [uploadErrors, setUploadErrors] = useState<string[]>([]);

  const handleUploadComplete = (result: FileProcessingResult) => {
    setUploadResult(result);
    setUploadErrors([]);
  };

  const handleUploadError = (errors: string[]) => {
    setUploadErrors(errors);
    setUploadResult(null);
  };

  const handleNewUpload = () => {
    setUploadResult(null);
    setUploadErrors([]);
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Upload Project Files
        </h1>
        <p className="text-gray-600">
          Upload your project files (.zip archives or package.json files) to analyze dependencies and architecture.
        </p>
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
      {uploadResult && (
        <UploadResultsDisplay 
          result={uploadResult} 
          onNewUpload={handleNewUpload}
        />
      )}

      {/* Instructions */}
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
          <div className="bg-gray-50 rounded-lg p-4">
            <div className="text-2xl font-bold text-gray-900 mb-1">
              {result.files.length}
            </div>
            <div className="text-sm text-gray-600">Files Uploaded</div>
          </div>
          
          <div className="bg-green-50 rounded-lg p-4">
            <div className="text-2xl font-bold text-green-600 mb-1">
              {result.packageJsonFiles.length}
            </div>
            <div className="text-sm text-gray-600">Package.json Found</div>
          </div>
          
          <div className="bg-blue-50 rounded-lg p-4">
            <div className="text-2xl font-bold text-blue-600 mb-1">
              {result.files.reduce((acc, file) => acc + file.fileSize, 0) / (1024 * 1024)}MB
            </div>
            <div className="text-sm text-gray-600">Total Size</div>
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
            {result.packageJsonFiles.map((pkg, index) => (
              <PackageJsonCard key={index} packageJson={pkg} />
            ))}
          </div>
        </div>
      )}

      {/* Next Steps */}
      <div className="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">
          Next Steps
        </h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <button className="text-left p-4 bg-white rounded-lg border border-gray-200 hover:border-blue-300 hover:shadow-sm transition-all">
            <div className="font-medium text-gray-900 mb-1">
              Analyze Dependencies
            </div>
            <div className="text-sm text-gray-600">
              Get detailed dependency analysis and health scores
            </div>
          </button>
          
          <button className="text-left p-4 bg-white rounded-lg border border-gray-200 hover:border-blue-300 hover:shadow-sm transition-all">
            <div className="font-medium text-gray-900 mb-1">
              Check Architecture
            </div>
            <div className="text-sm text-gray-600">
              Validate architecture and detect circular dependencies
            </div>
          </button>
        </div>
      </div>
    </div>
  );
};

// Package.json Card Component
interface PackageJsonCardProps {
  packageJson: PackageJsonFile;
}

const PackageJsonCard: React.FC<PackageJsonCardProps> = ({ packageJson }) => {
  const [isExpanded, setIsExpanded] = useState(false);
  
  const dependencyCount = packageJson.dependencies 
    ? Object.keys(packageJson.dependencies).length 
    : 0;
  
  const devDependencyCount = packageJson.devDependencies 
    ? Object.keys(packageJson.devDependencies).length 
    : 0;

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
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {packageJson.dependencies && Object.keys(packageJson.dependencies).length > 0 && (
              <div>
                <h5 className="font-medium text-gray-900 mb-2">Dependencies</h5>
                <div className="space-y-1 max-h-32 overflow-y-auto">
                  {Object.entries(packageJson.dependencies).map(([name, version]) => (
                    <div key={name} className="flex justify-between text-xs">
                      <span className="text-gray-600">{name}</span>
                      <span className="text-gray-500">{version}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
            
            {packageJson.devDependencies && Object.keys(packageJson.devDependencies).length > 0 && (
              <div>
                <h5 className="font-medium text-gray-900 mb-2">Dev Dependencies</h5>
                <div className="space-y-1 max-h-32 overflow-y-auto">
                  {Object.entries(packageJson.devDependencies).map(([name, version]) => (
                    <div key={name} className="flex justify-between text-xs">
                      <span className="text-gray-600">{name}</span>
                      <span className="text-gray-500">{version}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};