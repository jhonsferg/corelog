import { App, Card, Typography } from 'antd'
import { useNavigate } from 'react-router-dom'

import { useCreateTicket } from './hooks/useCreateTicket'
import { TicketForm } from './TicketForm'
import styles from './NewTicketPage.module.css'

export function NewTicketPage() {
  const navigate = useNavigate()
  const { message } = App.useApp()
  const { submit, isSubmitting } = useCreateTicket()

  const handleSubmit = async (input: Parameters<typeof submit>[0]) => {
    try {
      const ticket = await submit(input)
      message.success('Ticket created')
      navigate(`/tickets/${ticket.id}`)
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Could not create ticket')
      throw error
    }
  }

  return (
    <Card className={styles.card}>
      <Typography.Title level={3} className={styles.title}>
        New ticket
      </Typography.Title>
      <TicketForm onSubmit={handleSubmit} isSubmitting={isSubmitting} />
    </Card>
  )
}
