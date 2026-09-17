import type {
  CreateTicketInput,
  ListTicketsFilter,
  TicketRepositoryPort,
} from './ports/ticket-repository.port'
import type { Ticket, TicketStatus } from '../domain/ticket'

export class TicketService {
  private readonly repository: TicketRepositoryPort

  constructor(repository: TicketRepositoryPort) {
    this.repository = repository
  }

  list(filter: ListTicketsFilter = {}): Promise<Ticket[]> {
    return this.repository.list(filter)
  }

  getById(id: string): Promise<Ticket> {
    return this.repository.getById(id)
  }

  create(input: CreateTicketInput): Promise<Ticket> {
    return this.repository.create(input)
  }

  updateStatus(id: string, status: TicketStatus): Promise<Ticket> {
    return this.repository.updateStatus(id, status)
  }

  assign(id: string, assigneeId: string): Promise<Ticket> {
    return this.repository.assign(id, assigneeId)
  }
}
