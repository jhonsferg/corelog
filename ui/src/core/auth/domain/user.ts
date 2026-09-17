export type Role = 'admin' | 'agent' | 'requester'

export interface User {
  id: string
  name: string
  email: string
  role: Role
  teamId: string | null
  createdAt: string
}
