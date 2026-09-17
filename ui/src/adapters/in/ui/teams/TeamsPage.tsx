import { PlusOutlined } from '@ant-design/icons'
import { App, Button, Empty, Flex, Typography } from 'antd'
import { useEffect, useState } from 'react'

import type { Team } from '@/core/teams/domain/team'
import { useTeamsStore } from '@/adapters/in/state/teams.store'
import { extractErrorMessage } from '@/adapters/out/http/api-client'
import { TeamFormModal } from './TeamFormModal'
import { TeamList } from './TeamList'
import { TeamMembersPanel } from './TeamMembersPanel'
import styles from './TeamsPage.module.css'

type ModalState = { kind: 'create' } | { kind: 'rename'; team: Team } | { kind: 'closed' }

export function TeamsPage() {
  const { message } = App.useApp()
  const teams = useTeamsStore((s) => s.teams)
  const status = useTeamsStore((s) => s.status)
  const fetchTeams = useTeamsStore((s) => s.fetchTeams)
  const createTeam = useTeamsStore((s) => s.createTeam)
  const renameTeam = useTeamsStore((s) => s.renameTeam)
  const removeTeam = useTeamsStore((s) => s.removeTeam)

  const [selectedTeam, setSelectedTeam] = useState<Team | null>(null)
  const [modal, setModal] = useState<ModalState>({ kind: 'closed' })
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    fetchTeams()
  }, [fetchTeams])

  const handleSubmitModal = async (name: string) => {
    setIsSubmitting(true)
    try {
      if (modal.kind === 'create') {
        await createTeam(name)
        message.success('Team created')
      } else if (modal.kind === 'rename') {
        await renameTeam(modal.team.id, name)
        message.success('Team renamed')
      }
      setModal({ kind: 'closed' })
    } catch (error) {
      message.error(extractErrorMessage(error, 'Could not save team'))
      throw error
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleDelete = async (team: Team) => {
    try {
      await removeTeam(team.id)
      message.success('Team deleted')
      if (selectedTeam?.id === team.id) {
        setSelectedTeam(null)
      }
    } catch (error) {
      message.error(extractErrorMessage(error, 'Could not delete team'))
    }
  }

  return (
    <div>
      <Flex justify="space-between" align="center" className={styles.header}>
        <Typography.Title level={3} className={styles.title}>
          Teams
        </Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setModal({ kind: 'create' })}>
          New team
        </Button>
      </Flex>

      <TeamList
        teams={teams}
        isLoading={status === 'loading'}
        selectedTeamId={selectedTeam?.id ?? null}
        onSelect={setSelectedTeam}
        onRename={(team) => setModal({ kind: 'rename', team })}
        onDelete={handleDelete}
      />

      <div className={styles.membersSection}>
        {selectedTeam ? (
          <TeamMembersPanel key={selectedTeam.id} team={selectedTeam} />
        ) : (
          <Empty description="Select a team to manage its members" />
        )}
      </div>

      <TeamFormModal
        open={modal.kind !== 'closed'}
        title={modal.kind === 'rename' ? 'Rename team' : 'New team'}
        initialName={modal.kind === 'rename' ? modal.team.name : ''}
        isSubmitting={isSubmitting}
        onSubmit={handleSubmitModal}
        onCancel={() => setModal({ kind: 'closed' })}
      />
    </div>
  )
}
