import { createFileRoute } from '@tanstack/react-router'
import { FeaturesSection } from '../components/landing/FeaturesSection'
import { Footer } from '../components/landing/Footer'
import { HeroSection } from '../components/landing/HeroSection'
import { SampleResults } from '../components/landing/SampleResults'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  return (
    <main className="flex min-h-screen flex-col">
      <HeroSection />
      <SampleResults />
      <FeaturesSection />
      <Footer />
    </main>
  )
}
