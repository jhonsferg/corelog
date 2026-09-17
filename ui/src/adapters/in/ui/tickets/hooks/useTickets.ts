import { useEffect } from 'react'

import type { TicketPriority, TicketStatus } from '@/core/tickets/domain/ticket'
import { useTicketsStore } from '@/adapters/in/state/tickets.store'

interface UseTicketsFilter {
  status?: TicketStatus
  priority?: TicketPriority
}

export function useTickets(filter: UseTicketsFilter) {
  const { status: statusFilter, priority: priorityFilter } = filter
  const tickets = useTicketsStore((s) => s.tickets)
  const status = useTicketsStore((s) => s.status)
  const error = useTicketsStore((s) => s.error)
  const fetchTickets = useTicketsStore((s) => s.fetchTickets)

  useEffect(() => {
    fetchTickets({ status: statusFilter, priority: priorityFilter })
  }, [fetchTickets, statusFilter, priorityFilter])

  return { tickets, isLoading: status === 'loading', error }
}
