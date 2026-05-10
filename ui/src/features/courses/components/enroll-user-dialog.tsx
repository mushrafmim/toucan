import {useEffect, useState} from 'react'
import {Button, Dialog, Flex, Select, Text,} from '@radix-ui/themes'
import {UserPlus} from 'lucide-react'
import {enrollUser} from '../api/courses-api'
import {fetchUsers} from '@/features/users/api/users-api'
import type {User} from '@/features/users/model/user'

type EnrollUserDialogProps = {
  courseId: string
  onSuccess?: () => void
}

export function EnrollUserDialog({courseId, onSuccess}: EnrollUserDialogProps) {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [users, setUsers] = useState<User[]>([])
  const [selectedUserId, setSelectedUserId] = useState<string>('')
  const [selectedRole, setSelectedRole] = useState<string>('learner')

  useEffect(() => {
    if (open) {
      fetchUsers().then(res => setUsers(res.items))
    }
  }, [open])

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (!selectedUserId) return

    setLoading(true)

    try {
      await enrollUser({
        course_id: courseId,
        user_id: selectedUserId,
        role: selectedRole,
      })
      setOpen(false)
      setSelectedUserId('')
      if (onSuccess) onSuccess()
    } catch (error) {
      console.error('Failed to enroll user:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger>
        <Button variant="solid" color="amber" size="2" className="cursor-pointer">
          <UserPlus size={16}/>
          Enroll User
        </Button>
      </Dialog.Trigger>

      <Dialog.Content maxWidth="450px">
        <Dialog.Title>Enroll User</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Add a student or staff member to this course.
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="4">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Select User
              </Text>
              <Select.Root value={selectedUserId} onValueChange={setSelectedUserId}>
                <Select.Trigger placeholder="Select a user..." style={{width: '100%'}}/>
                <Select.Content>
                  {users.map(user => (
                    <Select.Item key={user.id} value={user.id}>
                      {user.display_name} ({user.email})
                    </Select.Item>
                  ))}
                </Select.Content>
              </Select.Root>
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Course Role
              </Text>
              <Select.Root value={selectedRole} onValueChange={setSelectedRole}>
                <Select.Trigger style={{width: '100%'}}/>
                <Select.Content>
                  <Select.Item value="learner">Learner</Select.Item>
                  <Select.Item value="manager">Manager</Select.Item>
                  <Select.Item value="owner">Owner</Select.Item>
                </Select.Content>
              </Select.Root>
            </label>
          </Flex>

          <Flex gap="3" mt="6" justify="end">
            <Dialog.Close>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </Dialog.Close>
            <Button type="submit" loading={loading} color="amber" disabled={!selectedUserId}>
              Enroll
            </Button>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  )
}
