import { getJSON, postJSON } from '@/shared/api/http'
import type { User, UserListResult, CreateUserRequest } from '../model/user'

export async function fetchUsers(page = 1, pageSize = 10): Promise<UserListResult> {
  return getJSON<UserListResult>(`/api/v1/users?page=${page}&page_size=${pageSize}`)
}

export async function fetchUser(id: string): Promise<User> {
  return getJSON<User>(`/api/v1/users/${id}`)
}

export async function createUser(data: CreateUserRequest): Promise<User> {
  return postJSON<User>('/api/v1/users', data)
}
