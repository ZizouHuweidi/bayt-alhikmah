import { Button } from "~/components/ui/button";
import { cn } from "~/lib/utils";

export default function Landing() {
  return (
    <div className="flex flex-col min-h-screen">
      {/* Navigation */}
      <header className="container mx-auto py-6 px-4 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <img 
            src="/logo-light.svg" 
            alt="Bayt Alhikmah Logo" 
            className="h-8 w-auto hidden dark:block" 
          />
          <img 
            src="/logo-dark.svg" 
            alt="Bayt Alhikmah Logo" 
            className="h-8 w-auto block dark:hidden" 
          />
        </div>
        <div className="flex items-center gap-4">
          <Button variant="ghost" className="text-sm font-medium">
            Login
          </Button>
          <Button className="text-sm font-medium">
            Get Started
          </Button>
        </div>
      </header>

      {/* Hero Section */}
      <section className="container mx-auto px-4 py-20 flex flex-col items-center text-center">
        <h1 className="text-4xl md:text-6xl font-bold tracking-tight mb-6">
          Your Knowledge, <span className="text-primary dark:text-primary">Organized</span>
        </h1>
        <p className="text-xl text-muted-foreground max-w-[600px] mb-10">
          Bayt Alhikmah helps you capture, organize, and share your knowledge in one beautiful, accessible place.
        </p>
        <div className="flex flex-col sm:flex-row gap-4 w-full justify-center">
          <Button size="lg" className="font-medium">
            Start for Free
          </Button>
          <Button size="lg" variant="outline" className="font-medium">
            Take a Tour
          </Button>
        </div>
        <div className="mt-16 relative w-full max-w-5xl">
          <div className="rounded-xl overflow-hidden border shadow-lg">
            <div className="aspect-video bg-muted/30 w-full"></div>
          </div>
          <div className="absolute -bottom-6 -right-6 -z-10 h-full w-full rounded-xl bg-muted/50"></div>
        </div>
      </section>

      {/* Features Section */}
      <section className="bg-muted/30 py-24">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl font-bold mb-4">Powerful Knowledge Management</h2>
            <p className="text-xl text-muted-foreground max-w-[600px] mx-auto">
              Everything you need to organize and access your information efficiently.
            </p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {/* Feature 1 */}
            <div className="bg-background rounded-xl p-6 shadow-sm border">
              <div className="size-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-primary">
                  <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
                  <polyline points="14 2 14 8 20 8" />
                  <path d="M12 18v-6" />
                  <path d="M8 18v-1" />
                  <path d="M16 18v-3" />
                </svg>
              </div>
              <h3 className="text-xl font-bold mb-2">Smart Organization</h3>
              <p className="text-muted-foreground">
                Organize your knowledge with powerful categorization, tagging, and search capabilities.
              </p>
            </div>

            {/* Feature 2 */}
            <div className="bg-background rounded-xl p-6 shadow-sm border">
              <div className="size-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-primary">
                  <circle cx="12" cy="12" r="10" />
                  <line x1="2" x2="22" y1="12" y2="12" />
                  <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold mb-2">Access Anywhere</h3>
              <p className="text-muted-foreground">
                Access your knowledge base from any device, anytime, with seamless synchronization.
              </p>
            </div>

            {/* Feature 3 */}
            <div className="bg-background rounded-xl p-6 shadow-sm border">
              <div className="size-12 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-primary">
                  <path d="M17 6.1H3" />
                  <path d="M21 12.1H3" />
                  <path d="M15.1 18H3" />
                </svg>
              </div>
              <h3 className="text-xl font-bold mb-2">Powerful Search</h3>
              <p className="text-muted-foreground">
                Find exactly what you need with our advanced search and filtering capabilities.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works Section */}
      <section className="container mx-auto px-4 py-24">
        <div className="text-center mb-16">
          <h2 className="text-3xl font-bold mb-4">How It Works</h2>
          <p className="text-xl text-muted-foreground max-w-[600px] mx-auto">
            Get started with Bayt Alhikmah in just a few simple steps.
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-4xl mx-auto">
          {/* Step 1 */}
          <div className="flex flex-col items-center text-center">
            <div className="size-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 relative">
              <span className="text-xl font-bold text-primary">1</span>
              <div className="absolute h-px w-full bg-primary/20 right-0 top-1/2 -z-10 md:w-full md:left-1/2"></div>
            </div>
            <h3 className="text-xl font-bold mb-2">Create an Account</h3>
            <p className="text-muted-foreground">
              Sign up for free and set up your personal knowledge base.
            </p>
          </div>

          {/* Step 2 */}
          <div className="flex flex-col items-center text-center">
            <div className="size-12 rounded-full bg-primary/10 flex items-center justify-center mb-4 relative">
              <span className="text-xl font-bold text-primary">2</span>
              <div className="absolute h-px w-full bg-primary/20 right-0 top-1/2 -z-10 md:w-full md:left-1/2"></div>
            </div>
            <h3 className="text-xl font-bold mb-2">Add Your Knowledge</h3>
            <p className="text-muted-foreground">
              Import or create content and organize it your way.
            </p>
          </div>

          {/* Step 3 */}
          <div className="flex flex-col items-center text-center">
            <div className="size-12 rounded-full bg-primary/10 flex items-center justify-center mb-4">
              <span className="text-xl font-bold text-primary">3</span>
            </div>
            <h3 className="text-xl font-bold mb-2">Access Anywhere</h3>
            <p className="text-muted-foreground">
              Use your knowledge base on any device, anytime.
            </p>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="bg-primary/5 py-24">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold mb-4">Ready to Get Started?</h2>
          <p className="text-xl text-muted-foreground max-w-[600px] mx-auto mb-8">
            Join thousands of users who are already organizing their knowledge with Bayt Alhikmah.
          </p>
          <Button size="lg" className="font-medium">
            Start for Free
          </Button>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-muted/30 py-12">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <div className="mb-6 md:mb-0">
              <img 
                src="/logo-light.svg" 
                alt="Bayt Alhikmah Logo" 
                className="h-6 w-auto hidden dark:block" 
              />
              <img 
                src="/logo-dark.svg" 
                alt="Bayt Alhikmah Logo" 
                className="h-6 w-auto block dark:hidden" 
              />
              <p className="text-sm text-muted-foreground mt-2">
                © {new Date().getFullYear()} Bayt Alhikmah. All rights reserved.
              </p>
            </div>
            <div className="flex gap-6">
              <a href="#" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Privacy
              </a>
              <a href="#" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Terms
              </a>
              <a href="#" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                Contact
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
