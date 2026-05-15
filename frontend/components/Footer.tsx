import Link from 'next/link'

export function Footer() {
  return (
    <footer className="border-t border-slate-800/50 bg-gradient-to-b from-slate-900/50 to-slate-950/80 backdrop-blur-md mt-20">
      <div className="container mx-auto px-6 py-8">
        <div className="grid grid-cols-[50%_25%_25%] gap-12 text-sm">
          {/* Column 1: Copyright & Info */}
          <div className="space-y-2 text-slate-400 flex flex-col justify-center">
            <p className="text-slate-300 leading-relaxed">
              © 2025 Mockly. Free, schema-driven mock API service for developers.
            </p>
            <p className="text-sm">
              Powered by{' '}
              <a 
                href="https://dps.codes" 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-400 hover:text-blue-300 transition-colors font-medium underline decoration-blue-400/30 hover:decoration-blue-300"
              >
                0xdps
              </a>
            </p>
          </div>
          
          {/* Column 2: Resources */}
          <div>
            <h3 className="text-white font-semibold mb-3 text-sm uppercase tracking-wider">Resources</h3>
            <ul className="space-y-2">
              <li>
                <Link 
                  href="/resources" 
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">API Resources</span>
                </Link>
              </li>
              <li>
                <Link 
                  href="/playground" 
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">Playground</span>
                </Link>
              </li>
              <li>
                <a 
                  href="https://github.com/0xdps/api-mockly" 
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">GitHub</span>
                </a>
              </li>
            </ul>
          </div>
          
          {/* Column 3: Projects */}
          <div>
            <h3 className="text-white font-semibold mb-3 text-sm uppercase tracking-wider">Projects</h3>
            <ul className="space-y-2">
              <li>
                <a 
                  href="https://www.mockly.codes/" 
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">Mockly</span>
                </a>
              </li>
              <li>
                <a 
                  href="https://pinboard-gpt.dps.codes/" 
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">Pinboard GPT</span>
                </a>
              </li>
              <li>
                <a 
                  href="https://devutil.dps.codes/" 
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">DevUtil</span>
                </a>
              </li>
              <li>
                <a 
                  href="https://www.pingpong.codes/" 
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">PingPong</span>
                </a>
              </li>
              <li>
                <a 
                  href="https://fake-stack.readthedocs.io/" 
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-white transition-colors duration-200 inline-flex items-center group"
                >
                  <span className="group-hover:translate-x-0.5 transition-transform">Fake Stack</span>
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  )
}
