import { PlusOutlined } from '@ant-design/icons'
import { Button, Select, Typography } from 'antd'
import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'

import type { Ticket, TicketPriority, TicketStatus } from '@/core/tickets/domain/ticket'
import { useTickets } from './hooks/useTickets'
import { TicketList } from './TicketList'
import styles from './TicketsListPage.module.css'

const statusOptions: { value: TicketStatus; label: string }[] = [
  { value: 'open', label: 'Open' },
  { value: 'in_progress', label: 'In progress' },
  { value: 'resolved', label: 'Resolved' },
  { value: 'closed', label: 'Closed' },
]

const priorityOptions: { value: TicketPriority; label: string }[] = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
  { value: 'critical', label: 'Critical' },
]

export function TicketsListPage() {
  const navigate = useNavigate()
  const [statusFilter, setStatusFilter] = useState<TicketStatus | undefined>()
  const [priorityFilter, setPriorityFilter] = useState<TicketPriority | undefined>()

  const { tickets, isLoading } = useTickets({ status: statusFilter, priority: priorityFilter })

  const handleSelect = (ticket: Ticket) => navigate(`/tickets/${ticket.id}`)

  return (
    <div>
      <div className={styles.header}>
        <Typography.Title level={3} className={styles.title}>
          Tickets
        </Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/tickets/new')}>
          New ticket
        </Button>
      </div>

      <div className={styles.filters}>
        <Select
          allowClear
          placeholder="Filter by status"
          className={styles.filterSelect}
          options={statusOptions}
          value={statusFilter}
          onChange={setStatusFilter}
        />
        <Select
          allowClear
          placeholder="Filter by priority"
          className={styles.filterSelect}
          options={priorityOptions}
          value={priorityFilter}
          onChange={setPriorityFilter}
        />
      </div>

      <TicketList tickets={tickets} isLoading={isLoading} onSelect={handleSelect} />

      <Typography.Paragraph type="secondary" className={styles.hint}>
        <Link to="/tickets/new">Create the first ticket</Link> if the list above is empty.
      </Typography.Paragraph>
    </div>
  )
}
