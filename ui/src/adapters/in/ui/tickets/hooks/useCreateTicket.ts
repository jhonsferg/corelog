import { useState } from 'react'

import type { CreateTicketInput } from '@/core/tickets/application/ports/ticket-repository.port'
import type { Ticket } from '@/core/tickets/domain/ticket'
import { extractErrorMessage } from '@/adapters/out/http/api-client'
import { useTicketsStore } from '@/adapters/in/state/tickets.store'

export function useCreateTicket() {
  const createTicket = useTicketsStore((s) => s.createTicket)
  const [isSubmitting, setSubmitting] = useState(false)

  const submit = async (input: CreateTicketInput): Promise<Ticket> => {
    setSubmitting(true)
    try {
      return await createTicket(input)
    } catch (error) {
      throw new Error(extractErrorMessage(error, 'Could not create ticket'), { cause: error })
    } finally {
      setSubmitting(false)
    }
  }

  return { submit, isSubmitting }
}
