import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Docker image uses standalone; Vercel supplies its own output.
  ...(process.env.VERCEL ? {} : { output: "standalone" as const }),
  images: {
    unoptimized: true,
  },
  async rewrites() {
    if (process.env.VERCEL) return [];
    const upstream = (process.env.API_UPSTREAM ?? "http://127.0.0.1:8080").replace(
      /\/$/,
      "",
    );
    return [{ source: "/api/:path*", destination: `${upstream}/api/:path*` }];
  },
};

export default nextConfig;
