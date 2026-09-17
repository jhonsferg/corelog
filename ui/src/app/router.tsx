import { createBrowserRouter, Navigate } from 'react-router-dom'

import { LoginPage } from '@/adapters/in/ui/auth/LoginPage'
import { RegisterPage } from '@/adapters/in/ui/auth/RegisterPage'
import { AdminRoute } from '@/adapters/in/ui/common/AdminRoute'
import { ProtectedRoute } from '@/adapters/in/ui/common/ProtectedRoute'
import { AppLayout } from '@/adapters/in/ui/layout/AppLayout'
import { ProfilePage } from '@/adapters/in/ui/profile/ProfilePage'
import { TeamsPage } from '@/adapters/in/ui/teams/TeamsPage'
import { NewTicketPage } from '@/adapters/in/ui/tickets/NewTicketPage'
import { TicketDetailPage } from '@/adapters/in/ui/tickets/TicketDetailPage'
import { TicketsListPage } from '@/adapters/in/ui/tickets/TicketsListPage'

export const router = createBrowserRouter([
  { path: '/login', element: <LoginPage /> },
  { path: '/register', element: <RegisterPage /> },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <AppLayout />,
        children: [
          { index: true, element: <Navigate to="/tickets" replace /> },
          { path: 'tickets', element: <TicketsListPage /> },
          { path: 'tickets/new', element: <NewTicketPage /> },
          { path: 'tickets/:id', element: <TicketDetailPage /> },
          { path: 'profile', element: <ProfilePage /> },
          {
            element: <AdminRoute />,
            children: [{ path: 'teams', element: <TeamsPage /> }],
          },
        ],
      },
    ],
  },
])
