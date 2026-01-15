import './globals.css';
import { Inter } from 'next/font/google';
import { cn } from '../lib/utils';
import { QueryProvider } from '../lib/providers/QueryProvider';
import { Analytics } from '@vercel/analytics/next';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: {
    default: 'MonoGuard - 30-Second Monorepo Analysis Tool',
    template: '%s | MonoGuard',
  },
  description:
    'Analyze monorepo dependencies, security vulnerabilities, and architecture issues in 30 seconds. Privacy-first with local processing. Perfect for JavaScript/TypeScript development teams.',
  keywords: [
    'monorepo',
    'architecture',
    'dependency-analysis',
    'security-scanning',
    'circular-dependencies',
    'typescript',
    'nextjs',
    'javascript',
    'developer-tools',
    'code-quality',
    'nx-alternative',
    'lerna-alternative',
  ],
  authors: [{ name: 'Alex Yu', url: 'https://mono-guard-frontend.vercel.app' }],
  creator: 'Alex Yu',
  publisher: 'MonoGuard',
  category: 'developer-tools',
  openGraph: {
    type: 'website',
    locale: 'en_US',
    url: 'https://mono-guard-frontend.vercel.app',
    siteName: 'MonoGuard',
    title: 'MonoGuard - 30-Second Monorepo Analysis Tool',
    description:
      'Analyze monorepo dependencies, security vulnerabilities, and architecture issues in 30 seconds. Privacy-first with local processing.',
    images: [],
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  verification: {
    // Add when you have these services
    // google: 'your-google-site-verification',
    // yandex: 'your-yandex-verification',
    // yahoo: 'your-yahoo-site-verification',
  },
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
      <body
        className={cn(
          'bg-background min-h-screen font-sans antialiased',
          inter.className
        )}
      >
        <QueryProvider>
          <div id="root">{children}</div>
        </QueryProvider>
        <Analytics />
      </body>
    </html>
  );
}
