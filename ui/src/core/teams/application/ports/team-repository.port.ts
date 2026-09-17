import type { Team } from '@/core/teams/domain/team'

export interface TeamRepositoryPort {
  list(): Promise<Team[]>
  create(name: string): Promise<Team>
  rename(id: string, name: string): Promise<Team>
  remove(id: string): Promise<void>
}
