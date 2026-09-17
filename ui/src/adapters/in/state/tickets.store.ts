import { create } from 'zustand'

import { TicketService } from '@/core/tickets/application/ticket.service'
import type {
  CreateTicketInput,
  ListTicketsFilter,
} from '@/core/tickets/application/ports/ticket-repository.port'
import type { Ticket, TicketStatus } from '@/core/tickets/domain/ticket'
import { TicketHttpRepository } from '@/adapters/out/http/ticket-http-repository'
import { extractErrorMessage } from '@/adapters/out/http/api-client'

const ticketService = new TicketService(new TicketHttpRepository())

type Status = 'idle' | 'loading' | 'error'

interface TicketsState {
  tickets: Ticket[]
  selected: Ticket | null
  status: Status
  error: string | null
  fetchTickets: (filter?: ListTicketsFilter) => Promise<void>
  fetchTicket: (id: string) => Promise<void>
  createTicket: (input: CreateTicketInput) => Promise<Ticket>
  updateStatus: (id: string, status: TicketStatus) => Promise<void>
  assign: (id: string, assigneeId: string) => Promise<void>
}

export const useTicketsStore = create<TicketsState>((set, get) => ({
  tickets: [],
  selected: null,
  status: 'idle',
  error: null,

  fetchTickets: async (filter = {}) => {
    set({ status: 'loading', error: null })
    try {
      const tickets = await ticketService.list(filter)
      set({ tickets, status: 'idle' })
    } catch (error) {
      set({ status: 'error', error: extractErrorMessage(error, 'Could not load tickets') })
    }
  },

  fetchTicket: async (id) => {
    set({ status: 'loading', error: null })
    try {
      const selected = await ticketService.getById(id)
      set({ selected, status: 'idle' })
    } catch (error) {
      set({ status: 'error', error: extractErrorMessage(error, 'Could not load ticket') })
    }
  },

  createTicket: async (input) => {
    const ticket = await ticketService.create(input)
    set({ tickets: [ticket, ...get().tickets] })
    return ticket
  },

  updateStatus: async (id, status) => {
    const updated = await ticketService.updateStatus(id, status)
    applyUpdate(set, get, updated)
  },

  assign: async (id, assigneeId) => {
    const updated = await ticketService.assign(id, assigneeId)
    applyUpdate(set, get, updated)
  },
}))

function applyUpdate(
  set: (partial: Partial<TicketsState>) => void,
  get: () => TicketsState,
  updated: Ticket,
) {
  set({
    tickets: get().tickets.map((t) => (t.id === updated.id ? updated : t)),
    selected: get().selected?.id === updated.id ? updated : get().selected,
  })
}
