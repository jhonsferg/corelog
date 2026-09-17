import { useEffect } from 'react'
import { RouterProvider } from 'react-router-dom'

import { useAuthStore } from '@/adapters/in/state/auth.store'
import { AppProviders } from './providers/AppProviders'
import { router } from './router'

export function App() {
  const restoreSession = useAuthStore((s) => s.restoreSession)

  useEffect(() => {
    restoreSession()
  }, [restoreSession])

  return (
    <AppProviders>
      <RouterProvider router={router} />
    </AppProviders>
  )
}
