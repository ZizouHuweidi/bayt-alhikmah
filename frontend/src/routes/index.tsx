import { Link } from '@tanstack/react-router'
import {
  Book,
  BookOpen,
  Brain,
  Library,
  Quote,
  Scroll,
  Search,
  Sparkles,
  Star,
  Users,
} from 'lucide-react'
import { Button } from '@/components/ui/button'

export const Route = createFileRoute('/')({
  component: LandingPage,
})

function LandingPage() {
  const features = [
    {
      icon: <Library className="w-8 h-8 text-emerald-600" />,
      title: 'Organize Your Knowledge',
      description:
        'Track books, papers, podcasts, videos, and articles in one unified library. Build your personal collection with rich metadata.',
    },
    {
      icon: <BookOpen className="w-8 h-8 text-emerald-600" />,
      title: 'Smart Annotations',
      description:
        'Create structured notes and annotations tied to specific pages, timestamps, or sections. Link your thoughts to source materials.',
    },
    {
      icon: <Search className="w-8 h-8 text-emerald-600" />,
      title: 'Discovery Engine',
      description:
        'Get AI-enhanced recommendations based on your reading history, interests, and notes. Discover hidden gems aligned with your curiosity.',
    },
    {
      icon: <Users className="w-8 h-8 text-emerald-600" />,
      title: 'Build Your Identity',
      description:
        'Create a public or private knowledge profile. Share reviews, ratings, and curated lists with the community.',
    },
    {
      icon: <Brain className="w-8 h-8 text-emerald-600" />,
      title: 'Connect Ideas',
      description:
        'Use tags, topics, and taxonomies to create meaningful connections between different sources and your personal insights.',
    },
    {
      icon: <Star className="w-8 h-8 text-emerald-600" />,
      title: 'Track & Review',
      description:
        'Rate sources based on usefulness, clarity, and depth. Build a timeline of your learning journey and growth.',
    },
  ]

  const supportedTypes = [
    { icon: <Book className="w-6 h-6" />, label: 'Books' },
    { icon: <Scroll className="w-6 h-6" />, label: 'Academic Papers' },
    { icon: <Quote className="w-6 h-6" />, label: 'Articles' },
    { icon: <Sparkles className="w-6 h-6" />, label: 'Podcasts & Videos' },
  ]

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-50 to-slate-100">
      {/* Hero Section */}
      <section className="relative px-6 py-24 overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-br from-emerald-50 via-teal-50 to-cyan-50 opacity-70" />
        <div className="absolute top-0 right-0 w-96 h-96 bg-emerald-200 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-pulse" />
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-teal-200 rounded-full mix-blend-multiply filter blur-3xl opacity-30 animate-pulse" />

        <div className="relative max-w-6xl mx-auto text-center">
          <div className="inline-flex items-center gap-2 px-4 py-2 mb-8 text-sm font-medium text-emerald-800 bg-emerald-100 rounded-full">
            <Sparkles className="w-4 h-4" />
            <span>Inspired by the House of Wisdom</span>
          </div>

          <h1 className="mb-6 text-5xl font-bold tracking-tight text-slate-900 md:text-7xl">
            <span className="block">Bayt al</span>
            <span className="block text-transparent bg-gradient-to-r from-emerald-600 to-teal-600 bg-clip-text">
              Hikmah
            </span>
          </h1>

          <p className="max-w-2xl mx-auto mb-4 text-xl text-slate-600">
            A modern, AI-powered platform for organizing, engaging with, and
            tracking knowledge sources across all media.
          </p>

          <p className="max-w-xl mx-auto mb-10 text-slate-500">
            Inspired by the historical center of learning from the Islamic Golden
            Age, reimagined for the 21st century.
          </p>

          <div className="flex flex-col items-center justify-center gap-4 sm:flex-row">
            <Link to="/registration">
              <Button
                size="lg"
                className="px-8 py-6 text-lg font-semibold text-white transition-all shadow-lg bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-700 hover:to-teal-700 hover:shadow-xl"
              >
                Start Your Journey
              </Button>
            </Link>
            <Link to="/login">
              <Button
                variant="outline"
                size="lg"
                className="px-8 py-6 text-lg font-semibold border-2 border-slate-300 hover:border-emerald-500 hover:text-emerald-700"
              >
                Sign In
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Supported Types */}
      <section className="px-6 py-12 border-y border-slate-200 bg-white/50">
        <div className="max-w-6xl mx-auto">
          <p className="mb-8 text-sm font-semibold tracking-wider text-center text-slate-500 uppercase">
            Track Everything You Learn From
          </p>
          <div className="flex flex-wrap items-center justify-center gap-8">
            {supportedTypes.map((type) => (
              <div
                key={type.label}
                className="flex items-center gap-3 px-6 py-3 bg-white rounded-full shadow-sm"
              >
                <span className="text-emerald-600">{type.icon}</span>
                <span className="font-medium text-slate-700">{type.label}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Grid */}
      <section className="px-6 py-24">
        <div className="max-w-6xl mx-auto">
          <div className="mb-16 text-center">
            <h2 className="mb-4 text-3xl font-bold text-slate-900">
              Your Knowledge, Organized
            </h2>
            <p className="max-w-2xl mx-auto text-lg text-slate-600">
              Build a personal library that grows with you. Track what you read,
              watch, and listen to — and connect it all together.
            </p>
          </div>

          <div className="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
            {features.map((feature) => (
              <div
                key={feature.title}
                className="p-8 transition-all duration-300 bg-white border border-slate-200 rounded-2xl hover:shadow-lg hover:border-emerald-200 group"
              >
                <div className="p-3 mb-6 transition-colors rounded-xl bg-slate-50 w-fit group-hover:bg-emerald-50">
                  {feature.icon}
                </div>
                <h3 className="mb-3 text-xl font-semibold text-slate-900">
                  {feature.title}
                </h3>
                <p className="leading-relaxed text-slate-600">
                  {feature.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="px-6 py-24">
        <div className="max-w-4xl mx-auto text-center">
          <div className="p-12 bg-gradient-to-br from-emerald-600 to-teal-700 rounded-3xl">
            <h2 className="mb-4 text-3xl font-bold text-white">
              Ready to Build Your Library?
            </h2>
            <p className="mb-8 text-lg text-emerald-100">
              Join Bayt al Hikmah and start tracking your knowledge journey today.
              Free to get started.
            </p>
            <Link to="/registration">
              <Button
                size="lg"
                variant="secondary"
                className="px-8 py-6 text-lg font-semibold text-emerald-700 bg-white hover:bg-emerald-50"
              >
                Create Free Account
              </Button>
            </Link>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="px-6 py-12 border-t border-slate-200">
        <div className="max-w-6xl mx-auto">
          <div className="flex flex-col items-center justify-between gap-4 md:flex-row">
            <div className="flex items-center gap-3">
              <div className="p-2 rounded-lg bg-emerald-600">
                <Library className="w-6 h-6 text-white" />
              </div>
              <span className="text-xl font-bold text-slate-900">Bayt al Hikmah</span>
            </div>
            <p className="text-sm text-slate-500">
              House of Wisdom — Reimagined for the Modern Age
            </p>
          </div>
        </div>
      </footer>
    </div>
  )
}

// Need to import createFileRoute for the Route export
import { createFileRoute } from '@tanstack/react-router'
