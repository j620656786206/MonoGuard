import './globals.css';
import { Inter } from 'next/font/google';
import { cn } from '../lib/utils';
import { QueryProvider } from '../lib/providers/QueryProvider';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'MonoGuard - Comprehensive Monorepo Architecture Analysis',
  description: 'Comprehensive monorepo architecture analysis and validation tool for modern development teams',
  keywords: ['monorepo', 'architecture', 'dependency-analysis', 'typescript', 'nextjs'],
  authors: [{ name: 'MonoGuard Team' }],
  creator: 'MonoGuard Team',
  publisher: 'MonoGuard Team',
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={cn(
        'min-h-screen bg-background font-sans antialiased',
        inter.className
      )}>
        <QueryProvider>
          <div id="root">
            {children}
          </div>
        </QueryProvider>
      </body>
    </html>
  );
}
