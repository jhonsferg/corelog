import { Navigate, Outlet } from 'react-router-dom'

import { useAuthStore } from '@/adapters/in/state/auth.store'

export function ProtectedRoute() {
  const status = useAuthStore((s) => s.status)
  const sessionChecked = useAuthStore((s) => s.sessionChecked)

  if (!sessionChecked) {
    return null
  }

  if (status !== 'authenticated') {
    return <Navigate to="/login" replace />
  }

  return <Outlet />
}
