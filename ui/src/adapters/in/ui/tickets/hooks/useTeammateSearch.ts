import { useEffect, useRef, useState } from 'react'

import type { User } from '@/core/auth/domain/user'
import { UserDirectoryService } from '@/core/users/application/user-directory.service'
import { UserDirectoryHttpRepository } from '@/adapters/out/http/user-directory-http-repository'

const directoryService = new UserDirectoryService(new UserDirectoryHttpRepository())

export function useTeammateSearch() {
  const [results, setResults] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const timeoutRef = useRef<ReturnType<typeof setTimeout>>(undefined)

  const search = (query: string) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current)
    }
    if (!query) {
      setResults([])
      return
    }
    setIsLoading(true)
    timeoutRef.current = setTimeout(() => {
      directoryService
        .searchTeammates(query)
        .then(setResults)
        .finally(() => setIsLoading(false))
    }, 300)
  }

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [])

  return { results, isLoading, search }
}
