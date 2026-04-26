import { BookOpen, GraduationCap, ShieldCheck } from 'lucide-react'

export type UserRole = 'admin' | 'teacher' | 'student'

export const TOUCAN_ROLE_HEADER = 'X-Toucan-Role'
export const DEFAULT_ACTIVE_ROLE: UserRole = 'student'

export type RoleDefinition = {
  id: UserRole
  label: string
  description: string
  icon: typeof ShieldCheck
}

export const ROLES: Record<UserRole, RoleDefinition> = {
  admin: {
    id: 'admin',
    label: 'Administrator',
    description: 'Manage platform, tenants, and system-wide settings.',
    icon: ShieldCheck,
  },
  teacher: {
    id: 'teacher',
    label: 'Instructor',
    description: 'Design curriculum, manage content, and track student progress.',
    icon: GraduationCap,
  },
  student: {
    id: 'student',
    label: 'Student',
    description: 'Access learning materials, complete lessons, and view progress.',
    icon: BookOpen,
  },
}

const ACTIVE_ROLE_STORAGE_KEY = 'toucan.activeRole'

export function isUserRole(value: string): value is UserRole {
  return value === 'admin' || value === 'teacher' || value === 'student'
}

export function normalizeRole(value: string | null): UserRole | null {
  if (value === 'learner') {
    return 'student'
  }
  if (value && isUserRole(value)) {
    return value
  }
  return null
}

export function getActiveRoleSnapshot(): UserRole {
  if (typeof window === 'undefined') {
    return DEFAULT_ACTIVE_ROLE
  }
  return normalizeRole(window.localStorage.getItem(ACTIVE_ROLE_STORAGE_KEY)) ?? DEFAULT_ACTIVE_ROLE
}

export function persistActiveRole(role: UserRole) {
  if (typeof window === 'undefined') {
    return
  }
  window.localStorage.setItem(ACTIVE_ROLE_STORAGE_KEY, role)
}

export function rolesFromClaims(claims: Record<string, unknown> | undefined): UserRole[] {
  if (!claims) {
    return [DEFAULT_ACTIVE_ROLE]
  }

  const rawRoles = [
    ...claimValues(claims.roles),
    ...claimValues(claims.role),
    ...claimValues(claims.groups),
    ...claimValues(claims['http://wso2.org/claims/role']),
  ]

  const roles = rawRoles
    .map((role) => normalizeRole(role.trim().toLowerCase()))
    .filter((role): role is UserRole => Boolean(role))

  return uniqueRoles(roles.length > 0 ? roles : [DEFAULT_ACTIVE_ROLE])
}

export function allRoles(): UserRole[] {
  return Object.keys(ROLES) as UserRole[]
}

function claimValues(value: unknown): string[] {
  if (typeof value === 'string') {
    return value.split(/[,\s]+/).filter(Boolean)
  }

  if (Array.isArray(value)) {
    return value.flatMap((item) => (typeof item === 'string' ? claimValues(item) : []))
  }

  return []
}

function uniqueRoles(roles: UserRole[]) {
  return Array.from(new Set(roles))
}
