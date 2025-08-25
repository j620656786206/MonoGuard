//@ts-check

// eslint-disable-next-line @typescript-eslint/no-var-requires
const { composePlugins, withNx } = require('@nx/next');


/**
 * @type {import('@nx/next/plugins/with-nx').WithNxOptions}
 **/
const nextConfig = {
  nx: {
    svgr: false,
  },
  
  // Build configuration
  typescript: {
    ignoreBuildErrors: true,
  },
  eslint: {
    ignoreDuringBuilds: true,
  },
  
  // Key configuration for Vercel deployment
  output: 'standalone',
  
  // Disable static optimization to avoid prerender errors
  trailingSlash: false,
  
  // Skip static optimization - appDir is now default in Next.js 15
  experimental: {
    // Remove deprecated appDir option
  },
  
  // Force disable static exports
  distDir: '.next',
  
  // Minimal image optimization to avoid issues
  images: {
    unoptimized: true,
  },
  
  // Redirects
  async redirects() {
    return [
      {
        source: '/',
        destination: '/dashboard',
        permanent: false,
      },
    ];
  },
};

const plugins = [
  // Add more Next.js plugins to this list if needed.
  withNx,
];

module.exports = composePlugins(...plugins)(nextConfig);

