import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Form, Input } from 'antd'
import { Controller, useForm } from 'react-hook-form'
import { z } from 'zod'

import type { UpdateProfilePayload } from '@/core/auth/application/ports/auth-repository.port'
import styles from './ProfileForm.module.css'

const schema = z.object({
  name: z.string().min(2, 'Name is too short'),
  email: z.string().email('Enter a valid email'),
})

type FormValues = z.infer<typeof schema>

interface ProfileFormProps {
  initialName: string
  initialEmail: string
  isSubmitting: boolean
  onSubmit: (payload: UpdateProfilePayload) => Promise<unknown>
}

export function ProfileForm({ initialName, initialEmail, isSubmitting, onSubmit }: ProfileFormProps) {
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    values: { name: initialName, email: initialEmail },
  })

  const submit = handleSubmit(async (values) => {
    await onSubmit(values)
  })

  return (
    <Form className={styles.form} layout="vertical" onFinish={submit}>
      <Form.Item label="Name" validateStatus={errors.name ? 'error' : ''} help={errors.name?.message}>
        <Controller name="name" control={control} render={({ field }) => <Input {...field} />} />
      </Form.Item>
      <Form.Item label="Email" validateStatus={errors.email ? 'error' : ''} help={errors.email?.message}>
        <Controller name="email" control={control} render={({ field }) => <Input {...field} type="email" />} />
      </Form.Item>
      <Button type="primary" htmlType="submit" loading={isSubmitting}>
        Save profile
      </Button>
    </Form>
  )
}
