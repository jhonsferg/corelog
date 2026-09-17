import type { Team } from '../domain/team'
import type { TeamRepositoryPort } from './ports/team-repository.port'

export class TeamService {
  private readonly repository: TeamRepositoryPort

  constructor(repository: TeamRepositoryPort) {
    this.repository = repository
  }

  list(): Promise<Team[]> {
    return this.repository.list()
  }

  create(name: string): Promise<Team> {
    return this.repository.create(name)
  }

  rename(id: string, name: string): Promise<Team> {
    return this.repository.rename(id, name)
  }

  remove(id: string): Promise<void> {
    return this.repository.remove(id)
  }
}
