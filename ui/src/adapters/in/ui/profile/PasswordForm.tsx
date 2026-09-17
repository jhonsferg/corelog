import { zodResolver } from '@hookform/resolvers/zod'
import { Button, Form, Input } from 'antd'
import { Controller, useForm } from 'react-hook-form'
import { z } from 'zod'

import type { ChangePasswordPayload } from '@/core/auth/application/ports/auth-repository.port'
import styles from './ProfileForm.module.css'

const schema = z.object({
  currentPassword: z.string().min(1, 'Current password is required'),
  newPassword: z.string().min(8, 'New password must be at least 8 characters'),
})

type FormValues = z.infer<typeof schema>

interface PasswordFormProps {
  isSubmitting: boolean
  onSubmit: (payload: ChangePasswordPayload) => Promise<unknown>
}

export function PasswordForm({ isSubmitting, onSubmit }: PasswordFormProps) {
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { currentPassword: '', newPassword: '' },
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
      <Form.Item
        label="Current password"
        validateStatus={errors.currentPassword ? 'error' : ''}
        help={errors.currentPassword?.message}
      >
        <Controller name="currentPassword" control={control} render={({ field }) => <Input.Password {...field} />} />
      </Form.Item>
      <Form.Item
        label="New password"
        validateStatus={errors.newPassword ? 'error' : ''}
        help={errors.newPassword?.message}
      >
        <Controller name="newPassword" control={control} render={({ field }) => <Input.Password {...field} />} />
      </Form.Item>
      <Button type="primary" htmlType="submit" loading={isSubmitting}>
        Change password
      </Button>
    </Form>
  )
}
