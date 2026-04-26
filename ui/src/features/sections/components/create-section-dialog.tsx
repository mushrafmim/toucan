import { useState } from 'react'
import {
  Button,
  Dialog,
  Flex,
  Text,
  TextField,
  TextArea,
} from '@radix-ui/themes'
import { Plus } from 'lucide-react'
import { createSection } from '../api/sections-api'
import { useNavigate } from 'react-router-dom'

type CreateSectionDialogProps = {
  courseId: string
  onSuccess?: () => void
}

export function CreateSectionDialog({ courseId, onSuccess }: CreateSectionDialogProps) {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({
    title: '',
    summary: '',
  })
  const navigate = useNavigate()

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)

    try {
      await createSection({
        ...form,
        course_id: courseId,
        position: 0, // Simplified for now
      })
      setOpen(false)
      setForm({ title: '', summary: '' })
      if (onSuccess) {
        onSuccess()
      } else {
        // Refresh the page to show the new section
        navigate(0)
      }
    } catch (error) {
      console.error('Failed to create section:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger>
        <Button variant="outline" color="amber" size="2" className="cursor-pointer">
          <Plus size={16} />
          Add Section
        </Button>
      </Dialog.Trigger>

      <Dialog.Content maxWidth="450px">
        <Dialog.Title>Add New Section</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Group your lessons into modules or topics.
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="4">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Section Title
              </Text>
              <TextField.Root
                placeholder="e.g. Getting Started"
                required
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
              />
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Summary
              </Text>
              <TextArea
                placeholder="What will students learn in this section?"
                value={form.summary}
                onChange={(e) => setForm({ ...form, summary: e.target.value })}
              />
            </label>
          </Flex>

          <Flex gap="3" mt="6" justify="end">
            <Dialog.Close>
              <Button variant="soft" color="gray">
                Cancel
              </Button>
            </Dialog.Close>
            <Button type="submit" loading={loading} color="amber">
              Add Section
            </Button>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  )
}
