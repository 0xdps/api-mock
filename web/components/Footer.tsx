import Link from 'next/link'

export function Footer() {
  return (
    <footer className="border-t border-slate-700 bg-slate-900/50 backdrop-blur mt-20">
      <div className="container mx-auto px-4 py-8">
        <div className="grid md:grid-cols-3 gap-8">
          <div>
            <h3 className="text-white font-semibold mb-4">apimock.codes</h3>
            <p className="text-slate-400 text-sm">
              Free, schema-driven mock API service for developers.
            </p>
          </div>
          
          <div>
            <h3 className="text-white font-semibold mb-4">Resources</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <Link href="/docs" className="text-slate-400 hover:text-white transition">
                  Documentation
                </Link>
              </li>
              <li>
                <Link href="/playground" className="text-slate-400 hover:text-white transition">
                  API Playground
                </Link>
              </li>
              <li>
                <a 
                  href="https://github.com/0xdps/fake-stack" 
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-slate-400 hover:text-white transition"
                >
                  GitHub Repository
                </a>
              </li>
            </ul>
          </div>
          
          <div>
            <h3 className="text-white font-semibold mb-4">Related Projects</h3>
            <ul className="space-y-2 text-sm">
              <li>
                <a 
                  href="https://fakestack.readthedocs.io" 
                  className="text-slate-400 hover:text-white transition"
                >
                  FakeStack (Python/Go/Node)
                </a>
              </li>
            </ul>
          </div>
        </div>
        
        <div className="mt-8 pt-8 border-t border-slate-700 text-center text-slate-400 text-sm">
          <p>© 2025 apimock.codes. Open source and free to use.</p>
        </div>
      </div>
    </footer>
  )
}
