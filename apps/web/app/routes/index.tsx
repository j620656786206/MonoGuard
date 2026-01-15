import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/')({
  component: HomePage,
});

function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="max-w-4xl text-center">
        <h1 className="mb-4 text-4xl font-bold tracking-tight text-gray-900 sm:text-6xl">
          MonoGuard
        </h1>
        <p className="mb-8 text-lg text-gray-600">
          Analyze your monorepo dependencies and detect circular dependencies
          with powerful visualization tools.
        </p>
        <div className="flex justify-center gap-4">
          <a
            href="/analyze"
            className="rounded-lg bg-blue-600 px-6 py-3 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
          >
            Start Analysis
          </a>
          <a
            href="/results"
            className="rounded-lg border border-gray-300 px-6 py-3 text-sm font-semibold text-gray-700 shadow-sm hover:bg-gray-50"
          >
            View Results
          </a>
        </div>
      </div>
    </main>
  );
}
