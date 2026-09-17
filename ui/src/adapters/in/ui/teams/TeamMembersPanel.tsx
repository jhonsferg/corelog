import { UserDeleteOutlined } from '@ant-design/icons'
import { App, Button, Card, List, Select, Typography } from 'antd'
import { useEffect, useRef, useState } from 'react'

import type { User } from '@/core/auth/domain/user'
import type { Team } from '@/core/teams/domain/team'
import { UserDirectoryService } from '@/core/users/application/user-directory.service'
import { useTeamsStore } from '@/adapters/in/state/teams.store'
import { UserDirectoryHttpRepository } from '@/adapters/out/http/user-directory-http-repository'
import { extractErrorMessage } from '@/adapters/out/http/api-client'
import styles from './TeamMembersPanel.module.css'

const directoryService = new UserDirectoryService(new UserDirectoryHttpRepository())

interface TeamMembersPanelProps {
  team: Team
}

export function TeamMembersPanel({ team }: TeamMembersPanelProps) {
  const { message } = App.useApp()
  const members = useTeamsStore((s) => s.members)
  const membersStatus = useTeamsStore((s) => s.membersStatus)
  const fetchMembers = useTeamsStore((s) => s.fetchMembers)
  const addMember = useTeamsStore((s) => s.addMember)
  const removeMember = useTeamsStore((s) => s.removeMember)

  const [searchResults, setSearchResults] = useState<User[]>([])
  const [isSearching, setIsSearching] = useState(false)
  const timeoutRef = useRef<ReturnType<typeof setTimeout>>(undefined)

  useEffect(() => {
    fetchMembers(team.id)
  }, [team.id, fetchMembers])

  const handleSearch = (query: string) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    if (!query) {
      setSearchResults([])
      return
    }
    setIsSearching(true)
    timeoutRef.current = setTimeout(() => {
      directoryService
        .search({ query })
        .then(setSearchResults)
        .finally(() => setIsSearching(false))
    }, 300)
  }

  const handleAdd = async (userId: string) => {
    try {
      await addMember(team.id, userId)
      message.success('Member added')
      setSearchResults([])
    } catch (error) {
      message.error(extractErrorMessage(error, 'Could not add member'))
    }
  }

  const handleRemove = async (userId: string) => {
    try {
      await removeMember(team.id, userId)
      message.success('Member removed')
    } catch (error) {
      message.error(extractErrorMessage(error, 'Could not remove member'))
    }
  }

  return (
    <Card title={`Members of ${team.name}`}>
      <Select<string>
        showSearch
        filterOption={false}
        placeholder="Search a person by name or email to add"
        className={styles.search}
        loading={isSearching}
        onSearch={handleSearch}
        onSelect={handleAdd}
        options={searchResults.map((u) => ({ value: u.id, label: `${u.name} (${u.email})` }))}
      />
      <List
        loading={membersStatus === 'loading'}
        dataSource={members}
        locale={{ emptyText: 'No members yet' }}
        renderItem={(member) => (
          <List.Item
            actions={[
              <Button key="remove" icon={<UserDeleteOutlined />} onClick={() => handleRemove(member.id)}>
                Remove
              </Button>,
            ]}
          >
            <Typography.Text>{member.name}</Typography.Text>
            <Typography.Text type="secondary">&nbsp;{member.email}</Typography.Text>
          </List.Item>
        )}
      />
    </Card>
  )
}
