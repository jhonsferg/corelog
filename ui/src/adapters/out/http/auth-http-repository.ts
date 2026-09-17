import type {
  AuthRepositoryPort,
  AuthSession,
  ChangePasswordPayload,
  LoginPayload,
  RegisterPayload,
  UpdateProfilePayload,
} from '@/core/auth/application/ports/auth-repository.port'
import type { User } from '@/core/auth/domain/user'
import { apiClient } from './api-client'
import { toUser, type UserDto } from './user-dto'

interface AuthDto {
  token: string
  user: UserDto
}

export class AuthHttpRepository implements AuthRepositoryPort {
  async register(payload: RegisterPayload): Promise<AuthSession> {
    await apiClient.post<{ data: UserDto }>('/auth/register', payload)
    return this.login({ email: payload.email, password: payload.password })
  }

  async login(payload: LoginPayload): Promise<AuthSession> {
    const { data } = await apiClient.post<{ data: AuthDto }>('/auth/login', payload)
    return { token: data.data.token, user: toUser(data.data.user) }
  }

  async me(token: string): Promise<User> {
    const { data } = await apiClient.get<{ data: UserDto }>('/users/me', {
      headers: { Authorization: `Bearer ${token}` },
    })
    return toUser(data.data)
  }

  async updateProfile(payload: UpdateProfilePayload): Promise<User> {
    const { data } = await apiClient.patch<{ data: UserDto }>('/users/me', payload)
    return toUser(data.data)
  }

  async changePassword(payload: ChangePasswordPayload): Promise<void> {
    await apiClient.post('/users/me/password', {
      current_password: payload.currentPassword,
      new_password: payload.newPassword,
    })
  }
}
