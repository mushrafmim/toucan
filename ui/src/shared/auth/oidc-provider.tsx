import { useEffect, type ReactNode } from 'react'

import { Button, Card, Flex, Heading, Text } from '@radix-ui/themes'
import { LogIn } from 'lucide-react'
import { AuthProvider, useAuth } from 'react-oidc-context'

import { setAccessTokenSnapshot } from '@/shared/auth/access-token'
import { isOIDCEnabled, oidcConfig } from '@/shared/auth/oidc-config'

export function ToucanAuthProvider({ children }: { children: ReactNode }) {
  if (!isOIDCEnabled || !oidcConfig) {
    return <>{children}</>
  }

  return (
    <AuthProvider {...oidcConfig}>
      <AccessTokenBridge />
      <AuthGate>{children}</AuthGate>
    </AuthProvider>
  )
}

function AccessTokenBridge() {
  const auth = useAuth()

  useEffect(() => {
    setAccessTokenSnapshot(auth.user?.access_token ?? null)
  }, [auth.user?.access_token])

  useEffect(() => {
    return () => setAccessTokenSnapshot(null)
  }, [])

  return null
}

function AuthGate({ children }: { children: ReactNode }) {
  const auth = useAuth()

  if (auth.activeNavigator === 'signinRedirect') {
    return <AuthShell title="Redirecting to sign in..." />
  }

  if (auth.activeNavigator === 'signoutRedirect') {
    return <AuthShell title="Signing out..." />
  }

  if (auth.isLoading) {
    return <AuthShell title="Loading session..." />
  }

  if (auth.error) {
    return (
      <AuthShell
        title="Authentication failed"
        description={auth.error.message}
        actionLabel="Try again"
        onAction={() => void auth.signinRedirect(signinState())}
      />
    )
  }

  if (!auth.isAuthenticated) {
    return (
      <AuthShell
        title="Sign in to Toucan"
        description="Use your organization account to continue."
        actionLabel="Sign in"
        onAction={() => void auth.signinRedirect(signinState())}
      />
    )
  }

  return <>{children}</>
}

function AuthShell({
  title,
  description,
  actionLabel,
  onAction,
}: {
  title: string
  description?: string
  actionLabel?: string
  onAction?: () => void
}) {
  return (
    <Flex align="center" justify="center" className="min-h-dvh px-5 text-[#2f241d]">
      <Card size="4" className="w-full max-w-md bg-[rgba(255,251,245,0.92)] shadow-[0_20px_50px_rgba(106,74,32,0.08)]">
        <Flex direction="column" gap="4">
          <Flex direction="column" gap="2">
            <Heading size="5">{title}</Heading>
            {description ? (
              <Text size="2" color="gray">
                {description}
              </Text>
            ) : null}
          </Flex>
          {actionLabel && onAction ? (
            <Button color="amber" onClick={onAction}>
              <LogIn size={16} />
              {actionLabel}
            </Button>
          ) : null}
        </Flex>
      </Card>
    </Flex>
  )
}

function signinState() {
  return {
    state: {
      returnTo: `${window.location.pathname}${window.location.search}${window.location.hash}`,
    },
  }
}

