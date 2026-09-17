import type { User } from '@/core/auth/domain/user'

export interface RegisterPayload {
  name: string
  email: string
  password: string
}

export interface LoginPayload {
  email: string
  password: string
}

export interface AuthSession {
  token: string
  user: User
}

export interface UpdateProfilePayload {
  name: string
  email: string
}

export interface ChangePasswordPayload {
  currentPassword: string
  newPassword: string
}

export interface AuthRepositoryPort {
  register(payload: RegisterPayload): Promise<AuthSession>
  login(payload: LoginPayload): Promise<AuthSession>
  me(token: string): Promise<User>
  updateProfile(payload: UpdateProfilePayload): Promise<User>
  changePassword(payload: ChangePasswordPayload): Promise<void>
}
