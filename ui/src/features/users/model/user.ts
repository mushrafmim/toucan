export type UserRole = 'admin' | 'instructor' | 'learner'

export type User = {
  id: string
  external_subject: string
  email: string
  display_name: string
  roles: UserRole[]
  created_at: string
  updated_at: string
}

export type UserListResult = {
  items: User[]
  page: number
  page_size: number
  total: number
}

export type CreateUserRequest = {
  external_subject: string
  email: string
  displayName: string
  roles: UserRole[]
}
