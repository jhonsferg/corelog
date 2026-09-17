import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Form, Input, Select } from 'antd'
import { Controller, useForm } from 'react-hook-form'
import { z } from 'zod'

import type { CreateTicketInput } from '@/core/tickets/application/ports/ticket-repository.port'
import styles from './TicketForm.module.css'

const schema = z.object({
  title: z.string().min(3, 'Title is too short'),
  description: z.string().min(3, 'Description is too short'),
  priority: z.enum(['low', 'medium', 'high', 'critical']),
})

type FormValues = z.infer<typeof schema>

interface TicketFormProps {
  onSubmit: (input: CreateTicketInput) => Promise<unknown>
  isSubmitting: boolean
}

export function TicketForm({ onSubmit, isSubmitting }: TicketFormProps) {
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { title: '', description: '', priority: 'medium' },
  })

  const submit = handleSubmit(async (values) => {
    try {
      await onSubmit(values)
      reset()
    } catch {
      return
    }
  })

  return (
    <Form className={styles.form} layout="vertical" onFinish={submit}>
      <Form.Item label="Title" validateStatus={errors.title ? 'error' : ''} help={errors.title?.message}>
        <Controller name="title" control={control} render={({ field }) => <Input {...field} />} />
      </Form.Item>
      <Form.Item
        label="Description"
        validateStatus={errors.description ? 'error' : ''}
        help={errors.description?.message}
      >
        <Controller name="description" control={control} render={({ field }) => <Input.TextArea {...field} rows={4} />} />
      </Form.Item>
      <Form.Item label="Priority" validateStatus={errors.priority ? 'error' : ''} help={errors.priority?.message}>
        <Controller
          name="priority"
          control={control}
          render={({ field }) => (
            <Select
              {...field}
              options={[
                { value: 'low', label: 'Low' },
                { value: 'medium', label: 'Medium' },
                { value: 'high', label: 'High' },
                { value: 'critical', label: 'Critical' },
              ]}
            />
          )}
        />
      </Form.Item>
      <Button type="primary" htmlType="submit" loading={isSubmitting}>
        Create ticket
      </Button>
    </Form>
  )
}
