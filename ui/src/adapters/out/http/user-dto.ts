import type { Role, User } from '@/core/auth/domain/user'

export interface UserDto {
  id: string
  name: string
  email: string
  role: Role
  team_id?: string
  created_at: string
}

export function toUser(dto: UserDto): User {
  return {
    id: dto.id,
    name: dto.name,
    email: dto.email,
    role: dto.role,
    teamId: dto.team_id ?? null,
    createdAt: dto.created_at,
  }
}

export function toUsers(dtos: UserDto[]): User[] {
  return dtos.map(toUser)
}
