import type { AuthProviderProps } from 'react-oidc-context'

const authority = import.meta.env.VITE_OIDC_AUTHORITY as string | undefined
const clientId = import.meta.env.VITE_OIDC_CLIENT_ID as string | undefined
const scope = (import.meta.env.VITE_OIDC_SCOPE as string | undefined) ?? 'openid profile email'

export const isOIDCEnabled = Boolean(authority && clientId)

function callbackPath() {
  const configured = import.meta.env.VITE_OIDC_CALLBACK_PATH as string | undefined
  return configured || '/auth/callback'
}

function redirectURI() {
  return `${window.location.origin}${callbackPath()}`
}

function postLogoutRedirectURI() {
  return window.location.origin
}

function callbackReturnTo(user: { state?: unknown } | undefined) {
  if (!user || typeof user.state !== 'object' || user.state === null) {
    return '/'
  }

  const returnTo = (user.state as { returnTo?: unknown }).returnTo
  return typeof returnTo === 'string' && returnTo.startsWith('/') ? returnTo : '/'
}

export const oidcConfig: AuthProviderProps | null = isOIDCEnabled
  ? {
      authority: authority!,
      client_id: clientId!,
      redirect_uri: redirectURI(),
      post_logout_redirect_uri: postLogoutRedirectURI(),
      response_type: 'code',
      scope,
      automaticSilentRenew: true,
      onSigninCallback: (user) => {
        window.history.replaceState({}, document.title, callbackReturnTo(user))
      },
    }
  : null

