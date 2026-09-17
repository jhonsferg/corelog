import type { TeamRepositoryPort } from '@/core/teams/application/ports/team-repository.port'
import type { Team } from '@/core/teams/domain/team'
import { apiClient } from './api-client'

interface TeamDto {
  id: string
  name: string
  created_at: string
  updated_at: string
}

function toTeam(dto: TeamDto): Team {
  return {
    id: dto.id,
    name: dto.name,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
  }
}

export class TeamHttpRepository implements TeamRepositoryPort {
  async list(): Promise<Team[]> {
    const { data } = await apiClient.get<{ data: TeamDto[] }>('/teams')
    return data.data.map(toTeam)
  }

  async create(name: string): Promise<Team> {
    const { data } = await apiClient.post<{ data: TeamDto }>('/teams', { name })
    return toTeam(data.data)
  }

  async rename(id: string, name: string): Promise<Team> {
    const { data } = await apiClient.patch<{ data: TeamDto }>(`/teams/${id}`, { name })
    return toTeam(data.data)
  }

  async remove(id: string): Promise<void> {
    await apiClient.delete(`/teams/${id}`)
  }
}
