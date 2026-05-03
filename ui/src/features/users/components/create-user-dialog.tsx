import { useState } from 'react'
import { Button, Dialog, Flex, Select, Text, TextField } from '@radix-ui/themes'
import { UserPlus } from 'lucide-react'
import { useRevalidator } from 'react-router-dom'
import { createUser } from '../api/users-api'
import type { CreateUserRequest, UserRole } from '../model/user'

export function CreateUserDialog() {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const revalidator = useRevalidator()

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setLoading(true)
    setError(null)

    const formData = new FormData(event.currentTarget)
    const data: CreateUserRequest = {
      external_subject: formData.get('external_subject') as string,
      email: formData.get('email') as string,
      displayName: formData.get('displayName') as string,
      roles: [formData.get('role') as UserRole], // For now, still picking one from Select, but sending as array
    }

    try {
      await createUser(data)
      setOpen(false)
      revalidator.revalidate()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create user')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger>
        <Button size="3" className="rounded-2xl cursor-pointer">
          <UserPlus size={18} />
          Add User
        </Button>
      </Dialog.Trigger>

      <Dialog.Content maxWidth="450px" className="rounded-[24px]">
        <Dialog.Title>Add New User</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Manually provision a local user mapping for an external identity.
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="3">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                External Subject (IDP ID)
              </Text>
              <TextField.Root
                name="external_subject"
                placeholder="e.g. asgardeo user id"
                required
              />
            </label>
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Email
              </Text>
              <TextField.Root
                name="email"
                type="email"
                placeholder="user@example.com"
                required
              />
            </label>
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Display Name
              </Text>
              <TextField.Root
                name="displayName"
                placeholder="John Doe"
                required
              />
            </label>
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Role
              </Text>
              <Select.Root name="role" defaultValue="learner">
                <Select.Trigger className="w-full" />
                <Select.Content>
                  <Select.Item value="admin">Admin</Select.Item>
                  <Select.Item value="instructor">Instructor</Select.Item>
                  <Select.Item value="learner">Learner</Select.Item>
                </Select.Content>
              </Select.Root>
            </label>

            {error && (
              <Text color="red" size="2">
                {error}
              </Text>
            )}

            <Flex gap="3" mt="4" justify="end">
              <Dialog.Close>
                <Button variant="soft" color="gray">
                  Cancel
                </Button>
              </Dialog.Close>
              <Button type="submit" loading={loading}>
                Create User
              </Button>
            </Flex>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  )
}
