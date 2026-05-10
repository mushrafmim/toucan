import { useEffect, useState } from 'react'
import { Badge, Flex, Table, Text, Spinner, Callout } from '@radix-ui/themes'
import { AlertCircle } from 'lucide-react'
import { PageHeader } from '@/components/page-header'
import type { UserListResult } from '../model/user'
import { CreateUserDialog } from '../components/create-user-dialog'
import { fetchUsers } from '../api/users-api'

export function UsersPage() {
  const [data, setData] = useState<UserListResult | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    let active = true
    setLoading(true)

    fetchUsers()
      .then((res) => {
        if (active) {
          setData(res)
          setError(null)
        }
      })
      .catch((err) => {
        if (active) setError(err instanceof Error ? err : new Error('Failed to load users'))
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [])

  if (loading) {
    return (
      <Flex align="center" justify="center" style={{ height: '60vh' }}>
        <Spinner size="3" />
      </Flex>
    )
  }

  return (
    <Flex direction="column" gap="6">
      <Flex justify="between" align="end" wrap="wrap" gap="4">
        <PageHeader
          badge="Administration"
          title="User Management"
          description="Manage local user profiles and their mapping to external identity providers."
        />
        <Flex pb="2">
          <CreateUserDialog />
        </Flex>
      </Flex>

      {error || !data ? (
        <Callout.Root color="red">
          <Callout.Icon>
            <AlertCircle size={16} />
          </Callout.Icon>
          <Callout.Text>
            {error?.message || 'Failed to load users'}
          </Callout.Text>
        </Callout.Root>
      ) : (
        <Table.Root variant="surface" className="rounded-2xl overflow-hidden bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]">
          <Table.Header>
            <Table.Row>
              <Table.ColumnHeaderCell>Display Name</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>Email</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>External ID</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>Role</Table.ColumnHeaderCell>
              <Table.ColumnHeaderCell>Created</Table.ColumnHeaderCell>
            </Table.Row>
          </Table.Header>

          <Table.Body>
            {data.items.map((user) => (
              <Table.Row key={user.id}>
                <Table.RowHeaderCell>
                  <Text weight="bold">{user.display_name}</Text>
                </Table.RowHeaderCell>
                <Table.Cell>{user.email}</Table.Cell>
                <Table.Cell>
                  <Text size="1" color="gray" className="font-mono">
                    {user.external_subject}
                  </Text>
                </Table.Cell>
                <Table.Cell>
                  <Badge color={getRoleColor(user.roles[0])} variant="soft" className="capitalize">
                    {user.roles[0]}
                  </Badge>
                </Table.Cell>
                <Table.Cell>
                  <Text size="1" color="gray">
                    {new Date(user.created_at).toLocaleDateString()}
                  </Text>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table.Root>
      )}
    </Flex>
  )
}

function getRoleColor(role: string) {
  switch (role) {
    case 'admin':
      return 'red'
    case 'instructor':
      return 'blue'
    default:
      return 'amber'
  }
}
