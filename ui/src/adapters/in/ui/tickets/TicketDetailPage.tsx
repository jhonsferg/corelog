import { App, Button, Card, Descriptions, Select, Skeleton, Space, Typography } from 'antd'
import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'

import { canTransitionTo, type TicketStatus } from '@/core/tickets/domain/ticket'
import { useTicketsStore } from '@/adapters/in/state/tickets.store'
import { useTeammateSearch } from './hooks/useTeammateSearch'
import { TicketBadge } from './TicketBadge'
import styles from './TicketDetailPage.module.css'

const allStatuses: TicketStatus[] = ['open', 'in_progress', 'resolved', 'closed']

export function TicketDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { message } = App.useApp()

  const selected = useTicketsStore((s) => s.selected)
  const status = useTicketsStore((s) => s.status)
  const fetchTicket = useTicketsStore((s) => s.fetchTicket)
  const updateStatus = useTicketsStore((s) => s.updateStatus)
  const assign = useTicketsStore((s) => s.assign)

  const [assigneeId, setAssigneeId] = useState('')
  const teammateSearch = useTeammateSearch()

  useEffect(() => {
    if (id) fetchTicket(id)
  }, [id, fetchTicket])

  if (status === 'loading' || !selected) {
    return <Skeleton active />
  }

  const nextStatuses = allStatuses.filter((s) => canTransitionTo(selected.status, s))

  const handleStatusChange = async (next: TicketStatus) => {
    try {
      await updateStatus(selected.id, next)
      message.success('Status updated')
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Could not update status')
    }
  }

  const handleAssign = async () => {
    if (!assigneeId) return
    try {
      await assign(selected.id, assigneeId)
      message.success('Ticket assigned')
      setAssigneeId('')
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Could not assign ticket')
    }
  }

  return (
    <div className={styles.wrapper}>
      <Button className={styles.backButton} onClick={() => navigate('/tickets')}>
        Back to tickets
      </Button>

      <Card>
        <Typography.Title level={3} className={styles.title}>
          {selected.title}
        </Typography.Title>
        <Typography.Paragraph>{selected.description}</Typography.Paragraph>

        <Descriptions column={1} bordered size="small">
          <Descriptions.Item label="Status">
            <TicketBadge kind="status" value={selected.status} />
          </Descriptions.Item>
          <Descriptions.Item label="Priority">
            <TicketBadge kind="priority" value={selected.priority} />
          </Descriptions.Item>
          <Descriptions.Item label="Requester">{selected.requesterId}</Descriptions.Item>
          <Descriptions.Item label="Assignee">{selected.assigneeId ?? 'Unassigned'}</Descriptions.Item>
          <Descriptions.Item label="Created">{new Date(selected.createdAt).toLocaleString()}</Descriptions.Item>
          <Descriptions.Item label="Updated">{new Date(selected.updatedAt).toLocaleString()}</Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="Change status">
        <Space>
          <Select
            className={styles.statusSelect}
            placeholder="Select next status"
            disabled={nextStatuses.length === 0}
            options={nextStatuses.map((s) => ({ value: s, label: s.replace('_', ' ') }))}
            onChange={handleStatusChange}
          />
          {nextStatuses.length === 0 && <Typography.Text type="secondary">No transitions available</Typography.Text>}
        </Space>
      </Card>

      <Card title="Assign">
        <Space>
          <Select
            showSearch
            filterOption={false}
            placeholder="Search a teammate by name or email"
            className={styles.assigneeInput}
            loading={teammateSearch.isLoading}
            notFoundContent={teammateSearch.isLoading ? 'Searching...' : 'No teammates found'}
            onSearch={teammateSearch.search}
            onChange={(value) => setAssigneeId(value)}
            value={assigneeId || undefined}
            options={teammateSearch.results.map((u) => ({ value: u.id, label: `${u.name} (${u.email})` }))}
          />
          <Button type="primary" onClick={handleAssign} disabled={!assigneeId}>
            Assign
          </Button>
        </Space>
      </Card>
    </div>
  )
}
