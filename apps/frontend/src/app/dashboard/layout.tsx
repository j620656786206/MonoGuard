'use client';
import { ReactNode, useState } from 'react';
import { 
  BarChart3,
  GitBranch, 
  LayoutDashboard,
  Package2,
  Building2,
  FileText,
  TrendingUp,
  ChevronRight,
  Menu,
  X
} from 'lucide-react';

interface DashboardLayoutProps {
  children: ReactNode;
}

// Navigation items based on design spec
const navigationItems = [
  { 
    name: 'Dashboard', 
    href: '/dashboard', 
    icon: LayoutDashboard, 
    current: true 
  },
  { 
    name: 'Dependencies', 
    href: '/dashboard/dependencies', 
    icon: GitBranch, 
    current: false 
  },
  { 
    name: 'Architecture', 
    href: '/dashboard/architecture', 
    icon: Building2, 
    current: false 
  },
  { 
    name: 'Reports', 
    href: '/dashboard/reports', 
    icon: FileText, 
    current: false 
  },
];

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = () => setIsSidebarOpen(!isSidebarOpen);
  const closeSidebar = () => setIsSidebarOpen(false);

  return (
    <div className="flex min-h-screen bg-muted/10">
      {/* Header Bar */}
      <div className="fixed top-0 left-0 right-0 z-50 bg-background border-b border-border h-16">
        <div className="flex items-center justify-between px-6 h-full">
          {/* Mobile Menu Button + Logo and Project Selector */}
          <div className="flex items-center space-x-4">
            {/* Mobile Menu Button */}
            <button
              onClick={toggleSidebar}
              className="lg:hidden p-2 rounded-md hover:bg-muted/80 transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
              aria-label="Toggle navigation menu"
            >
              {isSidebarOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
            
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
                <span className="text-primary-foreground font-bold text-sm">MG</span>
              </div>
              <h1 className="text-h3 font-bold hidden sm:block">MonoGuard</h1>
            </div>
            
            {/* Project Selector - Hidden on small mobile */}
            <div className="hidden sm:flex items-center space-x-2 px-3 py-1 bg-muted rounded-md cursor-pointer hover:bg-muted/80 transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2" tabIndex={0}>
              <Package2 className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm font-medium">mono-guard</span>
              <ChevronRight className="w-4 h-4 text-muted-foreground" />
            </div>
          </div>
          
          {/* Global Health Indicator */}
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2 px-3 py-1 bg-green-50 text-green-700 rounded-full">
              <div className="w-2 h-2 bg-green-500 rounded-full"></div>
              <span className="text-sm font-medium hidden sm:inline">Health: </span>
              <span className="text-sm font-medium">87</span>
              <TrendingUp className="w-4 h-4 hidden sm:block" />
            </div>
          </div>
        </div>
      </div>
      
      {/* Mobile Overlay */}
      {isSidebarOpen && (
        <div 
          className="fixed inset-0 z-30 bg-black/20 backdrop-blur-sm lg:hidden" 
          onClick={closeSidebar}
          aria-hidden="true"
        />
      )}

      {/* Sidebar */}
      <aside className={`
        fixed top-16 left-0 z-40 w-64 h-[calc(100vh-4rem)] bg-background border-r border-border
        transform transition-transform duration-200 ease-in-out
        lg:translate-x-0 ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}
      `}>
        <nav className="p-4 space-y-2">
          {navigationItems.map((item) => {
            const Icon = item.icon;
            return (
              <a
                key={item.name}
                href={item.href}
                onClick={closeSidebar}
                className={`
                  flex items-center px-3 py-2.5 text-sm font-medium rounded-lg transition-all duration-200
                  focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2
                  ${item.current 
                    ? 'bg-primary text-primary-foreground shadow-sm' 
                    : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
                  }
                `}
              >
                <Icon className="w-5 h-5 mr-3" />
                {item.name}
              </a>
            );
          })}
        </nav>
        
        {/* Navigation Sections */}
        <div className="px-4 pt-6">
          <div className="text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-3">
            Quick Actions
          </div>
          <div className="space-y-1">
            <button 
              className="w-full flex items-center px-3 py-2 text-sm text-muted-foreground hover:text-foreground hover:bg-muted/50 rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
              onClick={closeSidebar}
            >
              <BarChart3 className="w-4 h-4 mr-3" />
              Run Analysis
            </button>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="lg:ml-64 pt-16 flex-1 min-h-screen">
        <div className="p-4 sm:p-6 max-w-7xl mx-auto">
          {children}
        </div>
      </main>
    </div>
  );
}