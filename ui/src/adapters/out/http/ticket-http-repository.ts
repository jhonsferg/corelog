import type {
  CreateTicketInput,
  ListTicketsFilter,
  TicketRepositoryPort,
} from '@/core/tickets/application/ports/ticket-repository.port'
import type { Ticket, TicketPriority, TicketStatus } from '@/core/tickets/domain/ticket'
import { apiClient } from './api-client'

interface TicketDto {
  id: string
  title: string
  description: string
  status: TicketStatus
  priority: TicketPriority
  requester_id: string
  assignee_id: string | null
  created_at: string
  updated_at: string
}

function toTicket(dto: TicketDto): Ticket {
  return {
    id: dto.id,
    title: dto.title,
    description: dto.description,
    status: dto.status,
    priority: dto.priority,
    requesterId: dto.requester_id,
    assigneeId: dto.assignee_id ?? null,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
  }
}

export class TicketHttpRepository implements TicketRepositoryPort {
  async list(filter: ListTicketsFilter): Promise<Ticket[]> {
    const { data } = await apiClient.get<{ data: TicketDto[] }>('/tickets', {
      params: {
        status: filter.status,
        priority: filter.priority,
        page: filter.page,
        page_size: filter.pageSize,
      },
    })
    return data.data.map(toTicket)
  }

  async getById(id: string): Promise<Ticket> {
    const { data } = await apiClient.get<{ data: TicketDto }>(`/tickets/${id}`)
    return toTicket(data.data)
  }

  async create(input: CreateTicketInput): Promise<Ticket> {
    const { data } = await apiClient.post<{ data: TicketDto }>('/tickets', input)
    return toTicket(data.data)
  }

  async updateStatus(id: string, status: TicketStatus): Promise<Ticket> {
    const { data } = await apiClient.patch<{ data: TicketDto }>(`/tickets/${id}/status`, { status })
    return toTicket(data.data)
  }

  async assign(id: string, assigneeId: string): Promise<Ticket> {
    const { data } = await apiClient.patch<{ data: TicketDto }>(`/tickets/${id}/assign`, {
      assignee_id: assigneeId,
    })
    return toTicket(data.data)
  }
}
