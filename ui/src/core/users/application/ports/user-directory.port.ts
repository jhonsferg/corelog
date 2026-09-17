import type { User } from '@/core/auth/domain/user'

export interface AdminSearchQuery {
  query: string
  teamId?: string
}

export interface UserDirectoryPort {
  searchTeammates(query: string): Promise<User[]>
  search(query: AdminSearchQuery): Promise<User[]>
  setTeam(userId: string, teamId: string | null): Promise<User>
}
