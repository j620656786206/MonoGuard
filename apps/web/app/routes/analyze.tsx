import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/analyze')({
  component: AnalyzePage,
})

function AnalyzePage() {
  return (
    <main className="min-h-screen p-8">
      <div className="mx-auto max-w-4xl">
        <h1 className="mb-6 text-3xl font-bold text-gray-900">Analyze Your Workspace</h1>
        <div className="rounded-lg border-2 border-dashed border-gray-300 p-12 text-center">
          <div className="mb-4">
            <svg
              className="mx-auto h-12 w-12 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
          </div>
          <p className="text-lg text-gray-600">Drop your workspace.json here or click to upload</p>
          <p className="mt-2 text-sm text-gray-500">Supports Nx, Lerna, and Turborepo workspaces</p>
          <button className="mt-6 rounded-lg bg-blue-600 px-6 py-3 text-sm font-semibold text-white hover:bg-blue-500">
            Select File
          </button>
        </div>
      </div>
    </main>
  )
}
