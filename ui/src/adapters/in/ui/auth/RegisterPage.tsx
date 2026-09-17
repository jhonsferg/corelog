import { zodResolver } from '@hookform/resolvers/zod'
import { App, Button, Card, Form, Input, Typography } from 'antd'
import { Controller, useForm } from 'react-hook-form'
import { Link, useNavigate } from 'react-router-dom'
import { z } from 'zod'

import { useAuthStore } from '@/adapters/in/state/auth.store'
import styles from './AuthPage.module.css'

const schema = z.object({
  name: z.string().min(2, 'Name is too short'),
  email: z.string().email('Enter a valid email'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
})

type FormValues = z.infer<typeof schema>

export function RegisterPage() {
  const register = useAuthStore((s) => s.register)
  const navigate = useNavigate()
  const { message } = App.useApp()

  const {
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { name: '', email: '', password: '' },
  })

  const onSubmit = handleSubmit(async (values) => {
    try {
      await register(values)
      navigate('/tickets')
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Registration failed')
    }
  })

  return (
    <div className={styles.wrapper}>
      <Card className={styles.card}>
        <Typography.Title level={3} className={styles.title}>
          Create your CoreLog account
        </Typography.Title>
        <Form layout="vertical" onFinish={onSubmit}>
          <Form.Item label="Name" validateStatus={errors.name ? 'error' : ''} help={errors.name?.message}>
            <Controller name="name" control={control} render={({ field }) => <Input {...field} />} />
          </Form.Item>
          <Form.Item label="Email" validateStatus={errors.email ? 'error' : ''} help={errors.email?.message}>
            <Controller name="email" control={control} render={({ field }) => <Input {...field} type="email" />} />
          </Form.Item>
          <Form.Item label="Password" validateStatus={errors.password ? 'error' : ''} help={errors.password?.message}>
            <Controller
              name="password"
              control={control}
              render={({ field }) => <Input.Password {...field} />}
            />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={isSubmitting}>
            Create account
          </Button>
        </Form>
        <Typography.Paragraph className={styles.footer}>
          Already have an account? <Link to="/login">Sign in</Link>
        </Typography.Paragraph>
      </Card>
    </div>
  )
}
