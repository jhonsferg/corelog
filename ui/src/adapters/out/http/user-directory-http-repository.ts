import type { User } from '@/core/auth/domain/user'
import type { AdminSearchQuery, UserDirectoryPort } from '@/core/users/application/ports/user-directory.port'
import { apiClient } from './api-client'
import { toUser, toUsers, type UserDto } from './user-dto'

export class UserDirectoryHttpRepository implements UserDirectoryPort {
  async searchTeammates(query: string): Promise<User[]> {
    const { data } = await apiClient.get<{ data: UserDto[] }>('/users/search', {
      params: { q: query },
    })
    return toUsers(data.data)
  }

  async search(query: AdminSearchQuery): Promise<User[]> {
    const { data } = await apiClient.get<{ data: UserDto[] }>('/users', {
      params: { q: query.query, team_id: query.teamId },
    })
    return toUsers(data.data)
  }

  async setTeam(userId: string, teamId: string | null): Promise<User> {
    const { data } = await apiClient.patch<{ data: UserDto }>(`/users/${userId}/team`, {
      team_id: teamId,
    })
    return toUser(data.data)
  }
}
