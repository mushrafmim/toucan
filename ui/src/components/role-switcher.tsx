import { useEffect } from 'react'

import { Button, DropdownMenu, Flex, Text, Box } from '@radix-ui/themes'
import { ChevronDown } from 'lucide-react'
import { useAuth } from 'react-oidc-context'

import { allRoles, ROLES, rolesFromClaims, type UserRole } from '@/shared/auth/roles'
import { isOIDCEnabled } from '@/shared/auth/oidc-config'
import { useRole } from '@/shared/context/use-role'

export function RoleSwitcher() {
  if (!isOIDCEnabled) {
    return <RoleSwitcherContent availableRoles={allRoles()} />
  }

  return <OIDCRoleSwitcher />
}

function OIDCRoleSwitcher() {
  const auth = useAuth()
  return <RoleSwitcherContent availableRoles={rolesFromClaims(auth.user?.profile)} />
}

function RoleSwitcherContent({ availableRoles }: { availableRoles: UserRole[] }) {
  const { activeRole, setActiveRole, roleDef } = useRole()
  const activeRoleDef = availableRoles.includes(activeRole) ? roleDef : ROLES[availableRoles[0]]
  const CurrentIcon = activeRoleDef.icon

  useEffect(() => {
    if (!availableRoles.includes(activeRole)) {
      setActiveRole(availableRoles[0])
    }
  }, [activeRole, availableRoles, setActiveRole])

  if (availableRoles.length === 1) {
    return (
      <div className="inline-flex min-h-10 w-[160px] items-center gap-2 rounded-full border border-[rgba(134,99,57,0.1)] bg-[rgba(255,248,238,0.88)] px-3 py-2 text-[#6e5842]">
        <CurrentIcon size={16} className="shrink-0" />
        <span className="flex-1 truncate text-left text-sm">{activeRoleDef.label}</span>
      </div>
    )
  }

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        <Button variant="soft" color="amber" size="2" className="w-[160px] cursor-pointer justify-start gap-2">
          <CurrentIcon size={16} className="shrink-0" />
          <span className="flex-1 text-left truncate">{activeRoleDef.label}</span>
          <ChevronDown size={14} className="shrink-0 opacity-60" />
        </Button>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content variant="soft" size="2" align="end" className="min-w-[280px]">
        {availableRoles.map((roleID) => {
          const role = ROLES[roleID]
          const Icon = role.icon
          const isActive = activeRole === role.id

          return (
            <DropdownMenu.Item
              key={role.id}
              onClick={() => setActiveRole(role.id)}
              className={`p-3 transition-colors ${isActive ? 'bg-amber-50' : ''}`}
            >
              <Flex gap="3" align="start">
                <Box className={`mt-0.5 rounded-md p-1.5 ${isActive ? 'bg-amber-600 text-white' : 'bg-gray-100 text-gray-500'}`}>
                  <Icon size={18} />
                </Box>
                <Flex direction="column" gap="0">
                  <Text weight="bold" size="2" color={isActive ? 'amber' : undefined}>
                    {role.label}
                  </Text>
                  <Text size="1" color="gray" className="leading-tight">
                    {role.description}
                  </Text>
                </Flex>
              </Flex>
            </DropdownMenu.Item>
          )
        })}
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  )
}
