let accessTokenSnapshot: string | null = null

export function setAccessTokenSnapshot(token: string | null) {
  accessTokenSnapshot = token
}

export function getAccessTokenSnapshot() {
  return accessTokenSnapshot
}

