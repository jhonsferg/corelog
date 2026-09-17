import type { User } from '@/core/auth/domain/user'
import type { AdminSearchQuery, UserDirectoryPort } from './ports/user-directory.port'

export class UserDirectoryService {
  private readonly repository: UserDirectoryPort

  constructor(repository: UserDirectoryPort) {
    this.repository = repository
  }

  searchTeammates(query: string): Promise<User[]> {
    return this.repository.searchTeammates(query)
  }

  search(query: AdminSearchQuery): Promise<User[]> {
    return this.repository.search(query)
  }

  setTeam(userId: string, teamId: string | null): Promise<User> {
    return this.repository.setTeam(userId, teamId)
  }
}
