import { Avatar, Button, DropdownMenu, Flex, Text } from '@radix-ui/themes'
import { LogIn, LogOut, Settings, UserCircle } from 'lucide-react'
import { useAuth } from 'react-oidc-context'

import { isOIDCEnabled } from '@/shared/auth/oidc-config'

export function AuthControls() {
  if (!isOIDCEnabled) {
    return null
  }

  return <OIDCAuthControls />
}

function OIDCAuthControls() {
  const auth = useAuth()
  const name = auth.user?.profile.name
  const email = auth.user?.profile.email
  const subject = auth.user?.profile.sub
  const label = email ?? name ?? subject ?? 'User'
  const fallback = avatarFallback(name ?? email ?? subject)

  if (!auth.isAuthenticated) {
    return (
      <Button color="amber" variant="soft" onClick={() => void auth.signinRedirect()}>
        <LogIn size={16} />
        Sign in
      </Button>
    )
  }

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        <button
          type="button"
          className="inline-flex h-10 w-10 items-center justify-center rounded-full border border-[rgba(134,99,57,0.1)] bg-[rgba(255,248,238,0.88)] text-[#6e5842] transition duration-150 hover:bg-[rgba(230,196,150,0.3)] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[rgba(121,75,25,0.5)]"
          aria-label="User settings"
        >
          <Avatar fallback={fallback} size="2" radius="full" color="amber" />
        </button>
      </DropdownMenu.Trigger>
      <DropdownMenu.Content variant="soft" size="2" align="end" className="min-w-[260px]">
        <Flex direction="column" gap="1" className="px-3 py-2">
          <Flex align="center" gap="2">
            <UserCircle size={16} className="text-[#8a6240]" />
            <Text size="2" weight="bold" className="max-w-[22ch] truncate">
              {name ?? label}
            </Text>
          </Flex>
          {email ? (
            <Text size="1" color="gray" className="max-w-[26ch] truncate">
              {email}
            </Text>
          ) : null}
        </Flex>
        <DropdownMenu.Separator />
        <DropdownMenu.Item disabled>
          <Settings size={16} />
          Account settings
        </DropdownMenu.Item>
        <DropdownMenu.Item color="red" onClick={() => void auth.signoutRedirect()}>
          <LogOut size={16} />
          Sign out
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  )
}

function avatarFallback(value: string | undefined) {
  if (!value) {
    return 'U'
  }

  const parts = value
    .replace(/@.*/, '')
    .split(/[\s._-]+/)
    .filter(Boolean)

  const first = parts[0]?.[0] ?? value[0]
  const second = parts.length > 1 ? parts[1]?.[0] : undefined
  return `${first ?? 'U'}${second ?? ''}`.toUpperCase()
}
