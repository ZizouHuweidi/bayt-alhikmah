import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/')({ component: App })

function App() {
  const hello = 'hello, world!'

  return <div className="space-x-10">{hello}</div>
}
