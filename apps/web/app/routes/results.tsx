import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/results')({
  component: ResultsPage,
});

function ResultsPage() {
  return (
    <main className="min-h-screen p-8">
      <div className="mx-auto max-w-6xl">
        <h1 className="mb-6 text-3xl font-bold text-gray-900">
          Analysis Results
        </h1>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="mb-2 text-lg font-semibold text-gray-900">
              Health Score
            </h2>
            <p className="text-4xl font-bold text-green-600">--</p>
            <p className="mt-2 text-sm text-gray-500">
              No analysis data available
            </p>
          </div>
          <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="mb-2 text-lg font-semibold text-gray-900">
              Circular Dependencies
            </h2>
            <p className="text-4xl font-bold text-gray-400">--</p>
            <p className="mt-2 text-sm text-gray-500">
              Run analysis to detect cycles
            </p>
          </div>
          <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
            <h2 className="mb-2 text-lg font-semibold text-gray-900">
              Total Packages
            </h2>
            <p className="text-4xl font-bold text-gray-400">--</p>
            <p className="mt-2 text-sm text-gray-500">
              Upload workspace to count
            </p>
          </div>
        </div>
        <div className="mt-8 rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
          <h2 className="mb-4 text-lg font-semibold text-gray-900">
            Dependency Graph
          </h2>
          <div className="flex h-64 items-center justify-center rounded-lg bg-gray-50">
            <p className="text-gray-500">
              Run analysis to visualize dependencies
            </p>
          </div>
        </div>
        <div className="mt-6">
          <a
            href="/analyze"
            className="inline-flex items-center rounded-lg bg-blue-600 px-6 py-3 text-sm font-semibold text-white hover:bg-blue-500"
          >
            Start New Analysis
          </a>
        </div>
      </div>
    </main>
  );
}
