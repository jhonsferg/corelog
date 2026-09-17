import type { TicketPriority, TicketStatus } from '@/core/tickets/domain/ticket'
import styles from './TicketBadge.module.css'

type Tone = 'neutral' | 'blue' | 'gold' | 'green' | 'orange' | 'red'

const statusMap: Record<TicketStatus, { label: string; tone: Tone }> = {
  open: { label: 'Open', tone: 'blue' },
  in_progress: { label: 'In progress', tone: 'gold' },
  resolved: { label: 'Resolved', tone: 'green' },
  closed: { label: 'Closed', tone: 'neutral' },
}

const priorityMap: Record<TicketPriority, { label: string; tone: Tone }> = {
  low: { label: 'Low', tone: 'neutral' },
  medium: { label: 'Medium', tone: 'blue' },
  high: { label: 'High', tone: 'orange' },
  critical: { label: 'Critical', tone: 'red' },
}

type TicketBadgeProps = { kind: 'status'; value: TicketStatus } | { kind: 'priority'; value: TicketPriority }

export function TicketBadge(props: TicketBadgeProps) {
  const { label, tone } = props.kind === 'status' ? statusMap[props.value] : priorityMap[props.value]

  return <span className={`${styles.badge} ${styles[tone]}`}>{label}</span>
}
