import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useCallback, useEffect, useRef, useState } from 'react'
import { AnalysisResults } from '../components/analysis/AnalysisResults'
import { DependencyGraphViz } from '../components/visualization/DependencyGraph'
import { demoAnalysis, demoCircularDependencies, demoDependencyGraph } from '../lib/demo-data'

export const Route = createFileRoute('/analyze')({
  component: AnalyzePage,
})

type Phase = 'upload' | 'analyzing' | 'results'

const ANALYSIS_STEPS = [
  'Parsing workspace configuration...',
  'Building dependency graph...',
  'Detecting circular dependencies...',
  'Validating architecture layers...',
  'Computing health score...',
] as const

const STEP_DURATION_MS = 500
const TOTAL_DURATION_MS = STEP_DURATION_MS * ANALYSIS_STEPS.length

function AnalyzePage() {
  const navigate = useNavigate()
  const [phase, setPhase] = useState<Phase>('upload')
  const [progress, setProgress] = useState(0)
  const [currentStepIndex, setCurrentStepIndex] = useState(0)
  const animationRef = useRef<number | null>(null)
  const startTimeRef = useRef<number>(0)

  const startAnalysis = useCallback(() => {
    setPhase('analyzing')
    setProgress(0)
    setCurrentStepIndex(0)
    startTimeRef.current = Date.now()

    const tick = () => {
      const elapsed = Date.now() - startTimeRef.current
      const pct = Math.min((elapsed / TOTAL_DURATION_MS) * 100, 100)
      const stepIdx = Math.min(Math.floor(elapsed / STEP_DURATION_MS), ANALYSIS_STEPS.length - 1)

      setProgress(pct)
      setCurrentStepIndex(stepIdx)

      if (elapsed < TOTAL_DURATION_MS) {
        animationRef.current = requestAnimationFrame(tick)
      } else {
        setProgress(100)
        setCurrentStepIndex(ANALYSIS_STEPS.length - 1)
        setTimeout(() => setPhase('results'), 300)
      }
    }

    animationRef.current = requestAnimationFrame(tick)
  }, [])

  useEffect(() => {
    return () => {
      if (animationRef.current !== null) {
        cancelAnimationFrame(animationRef.current)
      }
    }
  }, [])

  if (phase === 'results') {
    return (
      <main className="min-h-screen bg-gray-50 p-8">
        <div className="mx-auto max-w-7xl">
          <AnalysisResults
            analysis={demoAnalysis}
            onNewAnalysis={() => {
              setPhase('upload')
              setProgress(0)
              setCurrentStepIndex(0)
            }}
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
        </div>
      </main>
    )
  }

  if (phase === 'analyzing') {
    return (
      <main className="min-h-screen p-8">
        <div className="mx-auto max-w-2xl pt-24">
          <div className="text-center">
            <div className="mb-6 inline-flex h-16 w-16 items-center justify-center rounded-full bg-blue-100">
              <svg
                className="h-8 w-8 animate-spin text-blue-600"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <title>Loading</title>
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                />
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                />
              </svg>
            </div>
            <h1 className="mb-2 text-2xl font-bold text-gray-900">Analyzing Workspace</h1>
            <p className="mb-8 text-gray-500">{ANALYSIS_STEPS[currentStepIndex]}</p>
          </div>

          <div className="mb-4 h-3 overflow-hidden rounded-full bg-gray-200">
            <div
              className="h-full rounded-full bg-blue-600 transition-all duration-150 ease-linear"
              style={{ width: `${progress}%` }}
            />
          </div>

          <div className="flex justify-between text-sm text-gray-500">
            <span>
              Step {currentStepIndex + 1} of {ANALYSIS_STEPS.length}
            </span>
            <span>{Math.round(progress)}%</span>
          </div>

          <div className="mt-8 space-y-2">
            {ANALYSIS_STEPS.map((step, i) => (
              <div key={step} className="flex items-center gap-3 text-sm">
                {i < currentStepIndex ? (
                  <svg
                    className="h-5 w-5 text-green-500"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <title>Completed</title>
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M5 13l4 4L19 7"
                    />
                  </svg>
                ) : i === currentStepIndex ? (
                  <div className="h-5 w-5 animate-pulse rounded-full bg-blue-500" />
                ) : (
                  <div className="h-5 w-5 rounded-full bg-gray-200" />
                )}
                <span className={i <= currentStepIndex ? 'text-gray-900' : 'text-gray-400'}>
                  {step}
                </span>
              </div>
            ))}
          </div>
        </div>
      </main>
    )
  }

  // Phase: upload
  return (
    <main className="min-h-screen p-8">
      <div className="mx-auto max-w-4xl">
        <h1 className="mb-2 text-3xl font-bold text-gray-900">Analyze Your Workspace</h1>
        <p className="mb-8 text-gray-600">
          Upload a workspace configuration or try the demo analysis with sample data.
        </p>

        <div className="rounded-lg border-2 border-dashed border-gray-300 bg-white p-12 text-center transition-colors hover:border-blue-400 hover:bg-blue-50/50">
          <div className="mb-4">
            <svg
              className="mx-auto h-12 w-12 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <title>Upload</title>
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

          <div className="mt-6 flex flex-col items-center gap-3">
            <button
              type="button"
              onClick={startAnalysis}
              className="rounded-lg bg-blue-600 px-8 py-3 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
            >
              Start Demo Analysis
            </button>
            <span className="text-xs text-gray-400">
              Uses sample data &mdash; no file upload required
            </span>
          </div>
        </div>

        <div className="mt-6 text-center">
          <button
            type="button"
            onClick={() => navigate({ to: '/results' })}
            className="text-sm text-blue-600 hover:text-blue-500 hover:underline"
          >
            Skip to results &rarr;
          </button>
        </div>
      </div>
    </main>
  )
}
