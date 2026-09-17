import type { Ticket, TicketPriority, TicketStatus } from '@/core/tickets/domain/ticket'

export interface CreateTicketInput {
  title: string
  description: string
  priority: TicketPriority
}

export interface ListTicketsFilter {
  status?: TicketStatus
  priority?: TicketPriority
  page?: number
  pageSize?: number
}

export interface TicketRepositoryPort {
  list(filter: ListTicketsFilter): Promise<Ticket[]>
  getById(id: string): Promise<Ticket>
  create(input: CreateTicketInput): Promise<Ticket>
  updateStatus(id: string, status: TicketStatus): Promise<Ticket>
  assign(id: string, assigneeId: string): Promise<Ticket>
}
