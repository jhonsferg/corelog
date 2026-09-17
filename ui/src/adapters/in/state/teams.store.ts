import { create } from 'zustand'

import type { User } from '@/core/auth/domain/user'
import { TeamService } from '@/core/teams/application/team.service'
import type { Team } from '@/core/teams/domain/team'
import { UserDirectoryService } from '@/core/users/application/user-directory.service'
import { TeamHttpRepository } from '@/adapters/out/http/team-http-repository'
import { UserDirectoryHttpRepository } from '@/adapters/out/http/user-directory-http-repository'
import { extractErrorMessage } from '@/adapters/out/http/api-client'

const teamService = new TeamService(new TeamHttpRepository())
const directoryService = new UserDirectoryService(new UserDirectoryHttpRepository())

type Status = 'idle' | 'loading' | 'error'

interface TeamsState {
  teams: Team[]
  status: Status
  error: string | null
  members: User[]
  membersStatus: Status
  fetchTeams: () => Promise<void>
  createTeam: (name: string) => Promise<Team>
  renameTeam: (id: string, name: string) => Promise<void>
  removeTeam: (id: string) => Promise<void>
  fetchMembers: (teamId: string) => Promise<void>
  addMember: (teamId: string, userId: string) => Promise<void>
  removeMember: (teamId: string, userId: string) => Promise<void>
}

export const useTeamsStore = create<TeamsState>((set, get) => ({
  teams: [],
  status: 'idle',
  error: null,
  members: [],
  membersStatus: 'idle',

  fetchTeams: async () => {
    set({ status: 'loading', error: null })
    try {
      const teams = await teamService.list()
      set({ teams, status: 'idle' })
    } catch (error) {
      set({ status: 'error', error: extractErrorMessage(error, 'Could not load teams') })
    }
  },

  createTeam: async (name) => {
    const team = await teamService.create(name)
    set({ teams: [...get().teams, team].sort((a, b) => a.name.localeCompare(b.name)) })
    return team
  },

  renameTeam: async (id, name) => {
    const updated = await teamService.rename(id, name)
    set({ teams: get().teams.map((t) => (t.id === id ? updated : t)) })
  },

  removeTeam: async (id) => {
    await teamService.remove(id)
    set({ teams: get().teams.filter((t) => t.id !== id) })
  },

  fetchMembers: async (teamId) => {
    set({ membersStatus: 'loading' })
    try {
      const members = await directoryService.search({ query: '', teamId })
      set({ members, membersStatus: 'idle' })
    } catch (error) {
      set({ membersStatus: 'error', error: extractErrorMessage(error, 'Could not load team members') })
    }
  },

  addMember: async (teamId, userId) => {
    await directoryService.setTeam(userId, teamId)
    await get().fetchMembers(teamId)
  },

  removeMember: async (teamId, userId) => {
    await directoryService.setTeam(userId, null)
    await get().fetchMembers(teamId)
  },
}))
