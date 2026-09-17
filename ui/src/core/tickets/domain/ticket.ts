export type TicketStatus = 'open' | 'in_progress' | 'resolved' | 'closed'
export type TicketPriority = 'low' | 'medium' | 'high' | 'critical'

export interface Ticket {
  id: string
  title: string
  description: string
  status: TicketStatus
  priority: TicketPriority
  requesterId: string
  assigneeId: string | null
  createdAt: string
  updatedAt: string
}

const allowedTransitions: Record<TicketStatus, TicketStatus[]> = {
  open: ['in_progress', 'closed'],
  in_progress: ['resolved', 'closed'],
  resolved: ['closed', 'in_progress'],
  closed: [],
}

export function canTransitionTo(current: TicketStatus, next: TicketStatus): boolean {
  return allowedTransitions[current].includes(next)
}
