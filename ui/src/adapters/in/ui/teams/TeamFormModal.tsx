import { zodResolver } from '@hookform/resolvers/zod'
import { Form, Input, Modal } from 'antd'
import { Controller, useForm } from 'react-hook-form'
import { z } from 'zod'

const schema = z.object({
  name: z.string().min(2, 'Name is too short'),
})

type FormValues = z.infer<typeof schema>

interface TeamFormModalProps {
  open: boolean
  title: string
  initialName: string
  isSubmitting: boolean
  onSubmit: (name: string) => Promise<unknown>
  onCancel: () => void
}

export function TeamFormModal({ open, title, initialName, isSubmitting, onSubmit, onCancel }: TeamFormModalProps) {
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    values: { name: initialName },
  })

  const submit = handleSubmit(async (values) => {
    await onSubmit(values.name)
  })

  return (
    <Modal open={open} title={title} onCancel={onCancel} onOk={submit} confirmLoading={isSubmitting} okText="Save">
      <Form layout="vertical">
        <Form.Item label="Team name" validateStatus={errors.name ? 'error' : ''} help={errors.name?.message}>
          <Controller name="name" control={control} render={({ field }) => <Input {...field} autoFocus />} />
        </Form.Item>
      </Form>
    </Modal>
  )
}
