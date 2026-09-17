import type {
  AuthRepositoryPort,
  LoginPayload,
  RegisterPayload,
  AuthSession,
  UpdateProfilePayload,
  ChangePasswordPayload,
} from './ports/auth-repository.port'
import type { User } from '../domain/user'

export class AuthService {
  private readonly repository: AuthRepositoryPort

  constructor(repository: AuthRepositoryPort) {
    this.repository = repository
  }

  register(payload: RegisterPayload): Promise<AuthSession> {
    return this.repository.register(payload)
  }

  login(payload: LoginPayload): Promise<AuthSession> {
    return this.repository.login(payload)
  }

  currentUser(token: string): Promise<User> {
    return this.repository.me(token)
  }

  updateProfile(payload: UpdateProfilePayload): Promise<User> {
    return this.repository.updateProfile(payload)
  }

  changePassword(payload: ChangePasswordPayload): Promise<void> {
    return this.repository.changePassword(payload)
  }
}
