import { useCallback, useMemo, useState, type ReactNode } from 'react'

import {
  ROLES,
  getActiveRoleSnapshot,
  persistActiveRole,
  type UserRole,
} from '@/shared/auth/roles'
import { RoleContext } from '@/shared/context/role-context-value'

export function RoleProvider({ children }: { children: ReactNode }) {
  const [activeRoleState, setActiveRoleState] = useState<UserRole>(() => getActiveRoleSnapshot())

  const setActiveRole = useCallback((role: UserRole) => {
    persistActiveRole(role)
    setActiveRoleState(role)
  }, [])

  const value = useMemo(
    () => ({
      activeRole: activeRoleState,
      setActiveRole,
      roleDef: ROLES[activeRoleState],
    }),
    [activeRoleState, setActiveRole],
  )

  return (
    <RoleContext.Provider value={value}>{children}</RoleContext.Provider>
  )
}
