import { App, Card, Typography } from 'antd'
import { useState } from 'react'

import type { ChangePasswordPayload, UpdateProfilePayload } from '@/core/auth/application/ports/auth-repository.port'
import { useAuthStore } from '@/adapters/in/state/auth.store'
import { extractErrorMessage } from '@/adapters/out/http/api-client'
import { PasswordForm } from './PasswordForm'
import { ProfileForm } from './ProfileForm'
import styles from './ProfilePage.module.css'

export function ProfilePage() {
  const { message } = App.useApp()
  const user = useAuthStore((s) => s.user)
  const updateProfile = useAuthStore((s) => s.updateProfile)
  const changePassword = useAuthStore((s) => s.changePassword)

  const [isSavingProfile, setIsSavingProfile] = useState(false)
  const [isSavingPassword, setIsSavingPassword] = useState(false)

  if (!user) {
    return null
  }

  const handleProfileSubmit = async (payload: UpdateProfilePayload) => {
    setIsSavingProfile(true)
    try {
      await updateProfile(payload)
      message.success('Profile updated')
    } catch (error) {
      message.error(extractErrorMessage(error, 'Could not update profile'))
    } finally {
      setIsSavingProfile(false)
    }
  }

  const handlePasswordSubmit = async (payload: ChangePasswordPayload) => {
    setIsSavingPassword(true)
    try {
      await changePassword(payload)
      message.success('Password changed')
    } catch (error) {
      message.error(extractErrorMessage(error, 'Could not change password'))
      throw error
    } finally {
      setIsSavingPassword(false)
    }
  }

  return (
    <div className={styles.wrapper}>
      <Typography.Title level={3} className={styles.title}>
        My profile
      </Typography.Title>

      <Card title="Profile">
        <ProfileForm
          initialName={user.name}
          initialEmail={user.email}
          isSubmitting={isSavingProfile}
          onSubmit={handleProfileSubmit}
        />
      </Card>

      <Card title="Password">
        <PasswordForm isSubmitting={isSavingPassword} onSubmit={handlePasswordSubmit} />
      </Card>
    </div>
  )
}
