import { Table } from 'antd'
import type { ColumnsType } from 'antd/es/table'

import type { Ticket, TicketPriority, TicketStatus } from '@/core/tickets/domain/ticket'
import { TicketBadge } from './TicketBadge'
import styles from './TicketList.module.css'

interface TicketListProps {
  tickets: Ticket[]
  isLoading: boolean
  onSelect: (ticket: Ticket) => void
}

export function TicketList({ tickets, isLoading, onSelect }: TicketListProps) {
  const columns: ColumnsType<Ticket> = [
    { title: 'Title', dataIndex: 'title', key: 'title' },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (value: TicketStatus) => <TicketBadge kind="status" value={value} />,
    },
    {
      title: 'Priority',
      dataIndex: 'priority',
      key: 'priority',
      render: (value: TicketPriority) => <TicketBadge kind="priority" value={value} />,
    },
    {
      title: 'Created',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (value: string) => new Date(value).toLocaleString(),
    },
  ]

  return (
    <Table<Ticket>
      rowKey="id"
      loading={isLoading}
      dataSource={tickets}
      columns={columns}
      onRow={(record) => ({ onClick: () => onSelect(record) })}
      rowClassName={() => styles.row}
      locale={{ emptyText: <div className={styles.empty}>No tickets yet</div> }}
    />
  )
}
