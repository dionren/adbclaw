import FeatureCard from '../ui/FeatureCard'
import { features } from '../../data/content'

export default function Features() {
  return (
    <section id="features" className="relative py-28">
      {/* Subtle divider */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-48 h-px bg-gradient-to-r from-transparent via-amber-500/20 to-transparent" />

      <div className="mx-auto max-w-6xl px-6">
        <div className="mb-16">
          <span className="inline-block mb-4 text-[11px] font-mono uppercase tracking-[0.2em] text-amber-500/60">
            Capabilities
          </span>
          <h2 className="text-3xl font-display font-bold tracking-tight text-stone-100 sm:text-4xl">
            Everything an AI agent needs
          </h2>
          <p className="mt-4 max-w-xl text-stone-500 leading-relaxed">
            Pure tool layer. No LLM logic. No agent framework. Just reliable, structured commands
            that any AI can call over adb shell.
          </p>
        </div>
        <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature, i) => (
            <FeatureCard key={feature.title} {...feature} index={i} />
          ))}
        </div>
      </div>
    </section>
  )
}
