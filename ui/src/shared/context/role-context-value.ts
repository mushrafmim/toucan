import { createContext } from 'react'

import type { RoleDefinition, UserRole } from '@/shared/auth/roles'

export type RoleContextType = {
  activeRole: UserRole
  setActiveRole: (role: UserRole) => void
  roleDef: RoleDefinition
}

export const RoleContext = createContext<RoleContextType | undefined>(undefined)

