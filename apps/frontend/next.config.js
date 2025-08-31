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
  
  trailingSlash: false,
  
  // Build configuration - fix issues instead of ignoring
  typescript: {
    ignoreBuildErrors: true, // Keep for now, will fix gradually
  },
  eslint: {
    ignoreDuringBuilds: true, // Keep for now, will fix gradually
  },
  
  // Re-enable image optimization
  images: {
    unoptimized: false,
  },
  
  // Clean experimental config - let Next.js handle defaults
  experimental: {
    // Keep empty or add valid Next.js 15 experimental features only
  },
  
  // Redirects - MVP stage: only allow landing page access
  async redirects() {
    return [
      {
        source: '/dashboard/:path*',
        destination: '/',
        permanent: false,
      },
      {
        source: '/upload/:path*',
        destination: '/',
        permanent: false,
      },
      {
        source: '/projects/:path*',
        destination: '/',
        permanent: false,
      },
      {
        source: '/analysis/:path*',
        destination: '/',
        permanent: false,
      },
      {
        source: '/settings/:path*',
        destination: '/',
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

