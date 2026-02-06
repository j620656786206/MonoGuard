import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { AnalysisResults } from '../components/analysis/AnalysisResults'
import { DependencyGraphViz } from '../components/visualization/DependencyGraph'
import { demoAnalysis, demoCircularDependencies, demoDependencyGraph } from '../lib/demo-data'

export const Route = createFileRoute('/results')({
  component: ResultsPage,
})

function ResultsPage() {
  const navigate = useNavigate()

  return (
    <main className="min-h-screen bg-gray-50 p-8">
      <div className="mx-auto max-w-7xl">
        <AnalysisResults
          analysis={demoAnalysis}
          onNewAnalysis={() => navigate({ to: '/analyze' })}
        />

        <section className="mt-8">
          <h2 className="mb-4 text-xl font-semibold text-gray-900">Dependency Graph</h2>
          <div className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
            <DependencyGraphViz
              data={demoDependencyGraph}
              circularDependencies={demoCircularDependencies}
              height={600}
            />
          </div>
        </section>

        <div className="mt-6">
          <button
            type="button"
            onClick={() => navigate({ to: '/analyze' })}
            className="inline-flex items-center rounded-lg bg-blue-600 px-6 py-3 text-sm font-semibold text-white hover:bg-blue-500"
          >
            New Analysis
          </button>
        </div>
      </div>
    </main>
  )
}
