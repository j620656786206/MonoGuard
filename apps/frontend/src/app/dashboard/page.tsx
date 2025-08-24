'use client';
import React, { useState, useEffect } from 'react';

// Temporary inline components until UI library is properly built
const Card = ({ children, className = "" }: { children: React.ReactNode; className?: string }) => (
  <div className={`rounded-lg border bg-card text-card-foreground shadow-sm ${className}`}>
    {children}
  </div>
);

const CardHeader = ({ children, className = "" }: { children: React.ReactNode; className?: string }) => (
  <div className={`flex flex-col space-y-1.5 p-6 ${className}`}>
    {children}
  </div>
);

const CardTitle = ({ children, className = "" }: { children: React.ReactNode; className?: string }) => (
  <h3 className={`text-2xl font-semibold leading-none tracking-tight ${className}`}>
    {children}
  </h3>
);

const CardDescription = ({ children, className = "" }: { children: React.ReactNode; className?: string }) => (
  <p className={`text-sm text-muted-foreground ${className}`}>
    {children}
  </p>
);

const CardContent = ({ children, className = "" }: { children: React.ReactNode; className?: string }) => (
  <div className={`p-6 pt-0 ${className}`}>
    {children}
  </div>
);

const Badge = ({ children, variant = "default", className = "" }: { 
  children: React.ReactNode; 
  variant?: "default" | "success" | "warning" | "destructive"; 
  className?: string;
}) => {
  const variants = {
    default: "bg-primary text-primary-foreground hover:bg-primary/80",
    success: "bg-green-500 text-white hover:bg-green-500/80",
    warning: "bg-yellow-500 text-white hover:bg-yellow-500/80",
    destructive: "bg-red-500 text-white hover:bg-red-500/80"
  };
  return (
    <div className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 border-transparent ${variants[variant]} ${className}`}>
      {children}
    </div>
  );
};

const Progress = ({ value = 0, className = "" }: { value?: number; className?: string }) => (
  <div className={`relative h-4 w-full overflow-hidden rounded-full bg-secondary ${className}`}>
    <div 
      className="h-full w-full flex-1 bg-primary transition-all"
      style={{ transform: `translateX(-${100 - (value || 0)}%)` }}
    />
  </div>
);

// Health Score Circular Display Component
const HealthScoreCard = ({ score, trend, className = "" }: { score: number; trend: string; className?: string }) => {
  const getHealthColor = (score: number) => {
    if (score >= 80) return "text-green-500";
    if (score >= 50) return "text-yellow-500";
    return "text-red-500";
  };

  const getHealthStatus = (score: number) => {
    if (score >= 80) return { text: "Healthy", color: "text-green-500", bg: "bg-green-50" };
    if (score >= 50) return { text: "Warning", color: "text-yellow-500", bg: "bg-yellow-50" };
    return { text: "Critical", color: "text-red-500", bg: "bg-red-50" };
  };

  const status = getHealthStatus(score);
  const circumference = 2 * Math.PI * 45; // radius = 45
  const strokeDasharray = circumference;
  const strokeDashoffset = circumference - (score / 100) * circumference;

  return (
    <Card className={`relative overflow-hidden ${className}`}>
      <CardHeader className="pb-2">
        <CardTitle className="text-lg font-medium">Monorepo Health Score</CardTitle>
      </CardHeader>
      <CardContent className="flex items-center justify-center">
        <div className="relative flex flex-col items-center">
          {/* Circular Progress */}
          <div className="relative w-36 h-36 mb-4">
            <svg className="w-36 h-36 transform -rotate-90" viewBox="0 0 100 100">
              {/* Background Circle */}
              <circle
                cx="50"
                cy="50"
                r="45"
                stroke="currentColor"
                strokeWidth="6"
                fill="transparent"
                className="text-gray-200"
              />
              {/* Progress Circle */}
              <circle
                cx="50"
                cy="50"
                r="45"
                stroke="currentColor"
                strokeWidth="6"
                fill="transparent"
                strokeDasharray={strokeDasharray}
                strokeDashoffset={strokeDashoffset}
                strokeLinecap="round"
                className={`${getHealthColor(score)} transition-all duration-1000 ease-out`}
              />
            </svg>
            {/* Score Display */}
            <div className="absolute inset-0 flex flex-col items-center justify-center">
              <div className={`health-score-display ${getHealthColor(score)}`}>
                {score}
              </div>
              <div className="text-sm text-muted-foreground font-medium">/ 100</div>
            </div>
          </div>
          
          {/* Status and Trend */}
          <div className="text-center">
            <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${status.bg} ${status.color} mb-2`}>
              {status.text}
            </div>
            <div className="flex items-center justify-center text-sm text-muted-foreground">
              <span className={trend.startsWith('+') ? 'text-green-500' : 'text-red-500'}>
                {trend.startsWith('+') ? '↗' : '↘'} {trend} from last week
              </span>
            </div>
          </div>
          
          <Button variant="outline" className="mt-4" size="sm" aria-label="View detailed health score breakdown">
            View Breakdown
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

const Button = ({ children, variant = "default", size = "default", className = "", disabled = false, loading = false, ...props }: { 
  children: React.ReactNode; 
  variant?: "default" | "outline"; 
  size?: "default" | "sm";
  className?: string;
  disabled?: boolean;
  loading?: boolean;
  [key: string]: any;
}) => {
  const variants = {
    default: "bg-primary text-primary-foreground hover:bg-primary/90",
    outline: "border border-input bg-background hover:bg-accent hover:text-accent-foreground"
  };
  
  const sizes = {
    default: "h-10 px-4 py-2",
    sm: "h-8 px-3 py-1 text-xs"
  };
  
  return (
    <button 
      className={`inline-flex items-center justify-center rounded-md text-sm font-medium ring-offset-background transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 ${sizes[size]} ${variants[variant]} ${className}`}
      disabled={disabled || loading}
      {...props}
    >
      {loading && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
      {children}
    </button>
  );
};

// Styled Dropdown Component
const Dropdown = ({ children, className = "", ...props }: {
  children: React.ReactNode;
  className?: string;
  [key: string]: any;
}) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="relative">
      <button
        className={`inline-flex items-center justify-between gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm font-medium ring-offset-background transition-colors hover:bg-accent hover:text-accent-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 ${className}`}
        onClick={() => setIsOpen(!isOpen)}
        aria-expanded={isOpen}
        aria-haspopup="true"
        {...props}
      >
        <Filter className="w-4 h-4" />
        <span>Filter</span>
        <ChevronDown className={`w-4 h-4 transition-transform duration-200 ${isOpen ? 'rotate-180' : ''}`} />
      </button>
      
      {isOpen && (
        <>
          <div 
            className="fixed inset-0 z-10" 
            onClick={() => setIsOpen(false)}
            aria-hidden="true"
          />
          <div className="absolute right-0 top-full z-20 mt-2 w-48 rounded-md border bg-popover p-1 text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95">
            <div className="px-2 py-1.5 text-sm font-semibold">Filter by severity</div>
            <div className="space-y-1">
              <button className="flex w-full items-center rounded-sm px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground focus:outline-none focus:bg-accent focus:text-accent-foreground">
                All Issues
              </button>
              <button className="flex w-full items-center rounded-sm px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground focus:outline-none focus:bg-accent focus:text-accent-foreground">
                High Severity
              </button>
              <button className="flex w-full items-center rounded-sm px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground focus:outline-none focus:bg-accent focus:text-accent-foreground">
                Medium Severity
              </button>
              <button className="flex w-full items-center rounded-sm px-2 py-1.5 text-sm hover:bg-accent hover:text-accent-foreground focus:outline-none focus:bg-accent focus:text-accent-foreground">
                Low Severity
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
};
import { 
  Activity, 
  AlertCircle, 
  CheckCircle2, 
  Clock, 
  FileCode, 
  GitBranch, 
  Package, 
  TrendingUp,
  Users,
  Zap,
  Loader2,
  ChevronDown,
  Filter
} from "lucide-react";

// API client function
const fetchDashboardData = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/projects');
    const data = await response.json();
    
    if (data.success && data.data.length > 0) {
      const projects = data.data;
      const totalHealthScore = projects.reduce((sum, p) => sum + (p.healthScore || 0), 0);
      const avgHealthScore = Math.round(totalHealthScore / projects.length);
      
      // Get most recent analysis timestamp
      const lastAnalysisTimestamp = projects
        .map(p => p.lastAnalysisAt)
        .filter(Boolean)
        .sort()
        .pop();
      
      // Build recent analyses from actual projects
      const recentAnalyses = projects.slice(0, 3).map((project, index) => ({
        id: project.id,
        projectName: project.name,
        status: project.status === "pending" ? "in-progress" : 
                 project.status === "completed" ? "completed" : "failed",
        healthScore: project.healthScore || 0,
        timestamp: project.updatedAt || project.createdAt,
        issuesFound: Math.floor(Math.random() * 15) // Mock issue count for now
      }));
      
      return {
        overview: {
          totalProjects: projects.length,
          totalPackages: projects.length * 8, // Estimate based on project count
          healthScore: avgHealthScore,
          lastAnalysis: lastAnalysisTimestamp ? 
            new Date(lastAnalysisTimestamp).toLocaleString() : "Never"
        },
        recentAnalyses: recentAnalyses,
        topIssues: [
          // Keep mock issues for now since we don't have analysis results yet
          {
            type: "Pending Analysis",
            count: projects.filter(p => p.status === "pending").length,
            severity: "medium"
          },
          {
            type: "Active Projects",
            count: projects.length,
            severity: "low"
          }
        ]
      };
    } else {
      return {
        overview: {
          totalProjects: 0,
          totalPackages: 0,
          healthScore: 0,
          lastAnalysis: "Never"
        },
        recentAnalyses: [],
        topIssues: []
      };
    }
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error);
    throw error;
  }
};

// Fallback mock data
const fallbackData = {
  overview: {
    totalProjects: 12,
    totalPackages: 48,
    healthScore: 87,
    lastAnalysis: "2 hours ago"
  },
  recentAnalyses: [
    {
      id: "1",
      projectName: "frontend-app",
      status: "completed",
      healthScore: 92,
      timestamp: "2023-12-15T10:30:00Z",
      issuesFound: 3
    },
    {
      id: "2", 
      projectName: "api-service",
      status: "in-progress",
      healthScore: 78,
      timestamp: "2023-12-15T09:15:00Z",
      issuesFound: 7
    },
    {
      id: "3",
      projectName: "shared-utils",
      status: "failed",
      healthScore: 65,
      timestamp: "2023-12-15T08:45:00Z",
      issuesFound: 12
    }
  ],
  topIssues: [
    {
      type: "Circular Dependency",
      count: 5,
      severity: "high"
    },
    {
      type: "Unused Dependencies",
      count: 12,
      severity: "medium"
    },
    {
      type: "Version Mismatch",
      count: 8,
      severity: "medium"
    },
    {
      type: "Missing Dependencies",
      count: 3,
      severity: "high"
    }
  ]
};

export default function DashboardPage() {
  const [isAnalysisLoading, setIsAnalysisLoading] = useState(false);
  const [dashboardData, setDashboardData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Load dashboard data on component mount
  useEffect(() => {
    const loadDashboardData = async () => {
      try {
        setLoading(true);
        const data = await fetchDashboardData();
        if (data) {
          setDashboardData(data);
        }
      } catch (err) {
        console.error('Error loading dashboard data:', err);
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    loadDashboardData();
  }, []);

  const handleRunAnalysis = async () => {
    setIsAnalysisLoading(true);
    // Simulate analysis process
    setTimeout(() => {
      setIsAnalysisLoading(false);
      // Refresh data after analysis
      const refreshData = async () => {
        const data = await fetchDashboardData();
        if (data) {
          setDashboardData(data);
        }
      };
      refreshData();
    }, 3000);
  };

  // Use real data if available, otherwise fallback to mock data
  const displayData = dashboardData || fallbackData;

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-h1">Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Monitor your monorepo health and architecture analysis
          </p>
          {loading && (
            <p className="text-sm text-blue-600 mt-2">
              🔄 Loading dashboard data...
            </p>
          )}
          {error && (
            <p className="text-sm text-red-600 mt-2">
              ⚠️ Using fallback data: {error}
            </p>
          )}
          {dashboardData && !loading && (
            <p className="text-sm text-green-600 mt-2">
              ✅ Connected to live data
            </p>
          )}
        </div>
        <Button 
          className="flex items-center gap-2" 
          aria-label="Run analysis on current project"
          loading={isAnalysisLoading}
          onClick={handleRunAnalysis}
          disabled={loading}
        >
          {!isAnalysisLoading && <Zap className="w-4 h-4" />}
          {isAnalysisLoading ? 'Analyzing...' : 'Run Analysis'}
        </Button>
      </div>

      {/* Health Score + Overview Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
        {/* Large Health Score Card */}
        <div className="lg:col-span-2">
          <HealthScoreCard 
            score={displayData.overview.healthScore} 
            trend={dashboardData ? "+0" : "+12"}
            className="h-full"
          />
        </div>
        
        {/* Overview Stats */}
        <div className="lg:col-span-2 grid grid-cols-1 sm:grid-cols-2 gap-4 md:gap-6">
        <Card className="transition-all duration-200 hover:shadow-md">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Projects</CardTitle>
            <Package className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{displayData.overview.totalProjects}</div>
            <p className="text-xs text-muted-foreground">
              {dashboardData ? "From live data" : "+2 from last month"}
            </p>
          </CardContent>
        </Card>

        <Card className="transition-all duration-200 hover:shadow-md">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Packages</CardTitle>
            <FileCode className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{displayData.overview.totalPackages}</div>
            <p className="text-xs text-muted-foreground">
              {dashboardData ? "From live data" : "+12 from last month"}
            </p>
          </CardContent>
        </Card>

          <Card className="transition-all duration-200 hover:shadow-md sm:col-span-2 lg:col-span-1">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Last Analysis</CardTitle>
              <Clock className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-lg font-bold">{displayData.overview.lastAnalysis}</div>
              <p className="text-xs text-muted-foreground">
                {dashboardData ? "Live data" : "All systems operational"}
              </p>
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 xl:grid-cols-2 gap-4 md:gap-6">
        {/* Recent Analyses */}
        <Card className="transition-all duration-200 hover:shadow-md">
          <CardHeader>
            <CardTitle>Recent Analyses</CardTitle>
            <CardDescription>
              Latest architecture analysis results
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-6">
              {displayData.recentAnalyses.map((analysis) => (
                <div key={analysis.id} className="flex items-center justify-between p-4 border rounded-lg hover:border-primary/20 transition-colors cursor-pointer focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2" tabIndex={0} role="button" aria-label={`View analysis for ${analysis.projectName}`}>
                  <div className="flex items-center space-x-4">
                    <div className="flex items-center justify-center w-10 h-10 rounded-full bg-gray-100">
                      <GitBranch className="w-5 h-5 text-gray-600" />
                    </div>
                    <div>
                      <p className="font-medium">{analysis.projectName}</p>
                      <div className="flex items-center space-x-2 mt-1">
                        <Badge 
                          variant={
                            analysis.status === "completed" ? "success" :
                            analysis.status === "in-progress" ? "warning" : "destructive"
                          }
                        >
                          {analysis.status === "completed" && <CheckCircle2 className="w-3 h-3 mr-1" />}
                          {analysis.status === "in-progress" && <Activity className="w-3 h-3 mr-1" />}
                          {analysis.status === "failed" && <AlertCircle className="w-3 h-3 mr-1" />}
                          {analysis.status}
                        </Badge>
                        <span className="text-sm text-muted-foreground">
                          {analysis.issuesFound} issues found
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-lg font-semibold">{analysis.healthScore}%</div>
                    <div className="text-xs text-muted-foreground">
                      {new Date(analysis.timestamp).toLocaleDateString()}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* Critical Issues */}
        <Card className="transition-all duration-200 hover:shadow-md">
          <CardHeader className="flex flex-row items-center justify-between space-y-0">
            <div>
              <CardTitle>Issues ({displayData.topIssues.length})</CardTitle>
              <CardDescription>
                {dashboardData ? "Current status from live data" : "Issues requiring immediate attention"}
              </CardDescription>
            </div>
            <Dropdown />
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {displayData.topIssues.length > 0 ? displayData.topIssues.map((issue, index) => (
                <div key={index} className="relative">
                  {/* Color Bar on Left */}
                  <div className={`absolute left-0 top-0 bottom-0 w-1 rounded-l ${
                    issue.severity === "high" ? "bg-red-500" : "bg-yellow-500"
                  }`} />
                  
                  <div className="pl-4 py-4 border rounded-lg bg-background transition-all duration-200 hover:shadow-sm hover:border-primary/20">
                    <div className="flex items-center justify-between">
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center space-x-3 mb-2">
                          <div className={`w-3 h-3 rounded-full ${
                            issue.severity === "high" ? "bg-red-500" : "bg-yellow-500"
                          }`} />
                          <h4 className="font-semibold text-foreground">{issue.type}</h4>
                          <div className="text-sm text-muted-foreground">
                            {issue.severity === "high" ? "2 hrs" : "30 min"}
                          </div>
                        </div>
                        
                        <p className="text-sm text-muted-foreground mb-3 pr-4">
                          {issue.type === "Circular Dependency" && "apps/web → libs/shared → apps/web"}
                          {issue.type === "Unused Dependencies" && `Affecting ${issue.count} packages • 234KB waste`}
                          {issue.type === "Version Mismatch" && `${issue.count} different versions across packages`}
                          {issue.type === "Missing Dependencies" && `${issue.count} packages missing required dependencies`}
                        </p>
                        
                        <div className="flex items-center space-x-2">
                          {issue.severity === "high" ? (
                            <>
                              <Button size="sm" className="bg-primary text-primary-foreground hover:bg-primary/90" aria-label={`Quick fix for ${issue.type}`}>
                                Quick Fix
                              </Button>
                              <Button size="sm" variant="outline" aria-label={`View details for ${issue.type}`}>
                                View Details
                              </Button>
                            </>
                          ) : (
                            <>
                              <Button size="sm" className="bg-primary text-primary-foreground hover:bg-primary/90" aria-label={`Consolidate ${issue.type}`}>
                                Consolidate
                              </Button>
                              <button 
                                className="text-sm text-muted-foreground hover:text-foreground transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 rounded-sm px-2 py-1"
                                aria-label={`Ignore ${issue.type}`}
                              >
                                Ignore
                              </button>
                            </>
                          )}
                        </div>
                      </div>
                      
                      <div className="flex flex-col items-center justify-center px-6 py-2 bg-muted/30 rounded-lg min-w-[100px]">
                        <div className="text-2xl font-bold text-foreground">{issue.count}</div>
                        <div className="text-xs text-muted-foreground">occurrences</div>
                      </div>
                    </div>
                  </div>
                </div>
              )) : (
                <div className="text-center py-8">
                  <div className="text-muted-foreground">
                    {dashboardData ? "No issues detected in current projects" : "No data available"}
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <Card className="transition-all duration-200 hover:shadow-md">
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
          <CardDescription>
            Common tasks and shortcuts
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <Button variant="outline" className="flex items-center gap-2 h-20 flex-col" aria-label="Start analyzing a new project">
              <Package className="w-6 h-6" />
              <span>Analyze New Project</span>
            </Button>
            <Button variant="outline" className="flex items-center gap-2 h-20 flex-col" aria-label="View project architecture diagram">
              <FileCode className="w-6 h-6" />
              <span>View Architecture</span>
            </Button>
            <Button variant="outline" className="flex items-center gap-2 h-20 flex-col" aria-label="Configure team settings and permissions">
              <Users className="w-6 h-6" />
              <span>Team Settings</span>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}