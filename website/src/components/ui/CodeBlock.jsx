export default function CodeBlock({ commands, title }) {
  return (
    <div className="group rounded-xl border border-stone-800/60 bg-surface-900/40 overflow-hidden transition-all duration-300 hover:border-amber-500/15">
      {title && (
        <div className="flex items-center justify-between px-5 py-3 border-b border-stone-800/40">
          <span className="text-xs text-stone-400 font-mono tracking-wide">{title}</span>
          <div className="flex gap-1.5">
            <span className="w-2 h-2 rounded-full bg-stone-800" />
            <span className="w-2 h-2 rounded-full bg-stone-800" />
            <span className="w-2 h-2 rounded-full bg-stone-800" />
          </div>
        </div>
      )}
      <div className="p-5 overflow-x-auto scanline">
        <pre className="text-[13px] font-mono leading-[1.9]">
          {commands.map((line, i) => (
            <div key={i} className="flex gap-2 group/line">
              <span className="text-amber-500/50 select-none shrink-0">$</span>
              <span className="text-stone-300 group-hover/line:text-stone-100 transition-colors">{line.cmd}</span>
              {line.comment && (
                <span className="text-stone-700 shrink-0 ml-1"># {line.comment}</span>
              )}
            </div>
          ))}
        </pre>
      </div>
    </div>
  )
}
