'use client';

import { useEffect, useState } from 'react';
import { apiClient } from '../lib/api/client';
import styles from './page.module.css';

interface HealthStatus {
  status: string;
  service: string;
  version: string;
  environment: string;
  timestamp: string;
  uptime: string;
  checks: {
    database: {
      status: string;
      message: string;
      details?: any;
    };
    redis: {
      status: string;
      message: string;
    };
  };
}

export default function Index() {
  const [healthStatus, setHealthStatus] = useState<HealthStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [projects, setProjects] = useState<any[]>([]);
  const [testResult, setTestResult] = useState<string>('');

  useEffect(() => {
    const testApiConnection = async () => {
      try {
        // Test health endpoint
        const healthResponse = await apiClient.get<HealthStatus>('/health');
        setHealthStatus(healthResponse.data);

        // Test projects endpoint
        const projectsResponse = await apiClient.get('/api/v1/projects');
        setProjects(projectsResponse.data);

        setTestResult('✅ All API endpoints working correctly!');
        setLoading(false);
      } catch (err: any) {
        console.error('API connection failed:', err);
        setError(err.message || 'Failed to connect to API');
        setTestResult('❌ API connection failed');
        setLoading(false);
      }
    };

    testApiConnection();
  }, []);

  const createSampleProject = async () => {
    try {
      setTestResult('Creating sample project...');
      const response = await apiClient.post('/api/v1/projects', {
        name: 'Sample Project',
        description: 'A test project created from the frontend',
        repositoryUrl: 'https://github.com/example/sample-project',
        branch: 'main',
        ownerId: 'user123'
      });
      
      setTestResult('✅ Sample project created successfully!');
      
      // Refresh projects list
      const projectsResponse = await apiClient.get('/api/v1/projects');
      setProjects(projectsResponse.data);
    } catch (err: any) {
      console.error('Failed to create project:', err);
      setTestResult(`❌ Failed to create project: ${err.message}`);
    }
  };

  const testErrorHandling = async () => {
    try {
      setTestResult('Testing error handling...');
      // Try to create a project with invalid data
      await apiClient.post('/api/v1/projects', {
        name: '',
        repositoryUrl: 'invalid-url',
        branch: '',
        ownerId: ''
      });
      setTestResult('❌ Error handling test failed - should have thrown error');
    } catch (err: any) {
      console.log('Expected error caught:', err);
      setTestResult('✅ Error handling working correctly - validation errors caught');
    }
  };

  return (
    <div className={styles.page}>
      <div className="wrapper">
        <div className="container">
          <div id="welcome">
            <h1>
              <span>MonoGuard</span>
              API Integration Test 🧪
            </h1>
          </div>

          <div id="api-test" className="rounded shadow" style={{ margin: '20px 0', padding: '20px', border: '1px solid #ccc' }}>
            <h2>API Connection Test</h2>
            
            {loading && (
              <div style={{ color: '#666' }}>
                Testing API connection...
              </div>
            )}
            
            {error && (
              <div style={{ color: '#ff4444', padding: '10px', backgroundColor: '#ffe6e6', borderRadius: '4px' }}>
                <strong>Error:</strong> {error}
              </div>
            )}
            
            {healthStatus && (
              <div style={{ color: '#28a745', padding: '10px', backgroundColor: '#e6f7e6', borderRadius: '4px' }}>
                <h3>✅ API Connection Successful!</h3>
                <p><strong>Service:</strong> {healthStatus.service}</p>
                <p><strong>Version:</strong> {healthStatus.version}</p>
                <p><strong>Environment:</strong> {healthStatus.environment}</p>
                <p><strong>Status:</strong> {healthStatus.status}</p>
                <p><strong>Uptime:</strong> {healthStatus.uptime}</p>
                
                <h4>Service Health Checks:</h4>
                <ul>
                  <li>
                    <strong>Database:</strong> {healthStatus.checks.database.status} 
                    - {healthStatus.checks.database.message}
                  </li>
                  <li>
                    <strong>Redis:</strong> {healthStatus.checks.redis.status} 
                    - {healthStatus.checks.redis.message}
                  </li>
                </ul>
                
                <h4>Projects API Test:</h4>
                <p>Found {projects.length} projects</p>
                {projects.length === 0 && (
                  <p style={{ fontStyle: 'italic' }}>No projects found (database is empty, which is expected)</p>
                )}
                
                <div style={{ marginTop: '15px' }}>
                  <button 
                    onClick={createSampleProject}
                    style={{
                      padding: '10px 20px',
                      backgroundColor: '#007bff',
                      color: 'white',
                      border: 'none',
                      borderRadius: '4px',
                      cursor: 'pointer',
                      marginRight: '10px'
                    }}
                  >
                    Create Sample Project
                  </button>
                  <button 
                    onClick={testErrorHandling}
                    style={{
                      padding: '10px 20px',
                      backgroundColor: '#dc3545',
                      color: 'white',
                      border: 'none',
                      borderRadius: '4px',
                      cursor: 'pointer'
                    }}
                  >
                    Test Error Handling
                  </button>
                  {testResult && (
                    <p style={{ marginTop: '10px', fontWeight: 'bold' }}>{testResult}</p>
                  )}
                </div>
                
                {projects.length > 0 && (
                  <div style={{ marginTop: '15px' }}>
                    <h5>Projects in Database:</h5>
                    <ul>
                      {projects.map((project: any, index: number) => (
                        <li key={project.id || index}>
                          <strong>{project.name}</strong> - {project.description}
                          <br />
                          <small>Framework: {project.framework}, Language: {project.language}</small>
                        </li>
                      ))}
                    </ul>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};
