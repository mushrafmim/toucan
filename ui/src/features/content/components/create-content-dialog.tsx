import { useState } from 'react'
import {
  Button,
  Dialog,
  Flex,
  Text,
  TextField,
  TextArea,
  Select,
  Box,
} from '@radix-ui/themes'
import { Plus } from 'lucide-react'
import { createContent } from '../api/content-api'
import { useNavigate } from 'react-router-dom'
import type { ContentType } from '@/features/courses/model/course'

type CreateContentDialogProps = {
  sectionId: string
  onSuccess?: () => void
}

export function CreateContentDialog({ sectionId, onSuccess }: CreateContentDialogProps) {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({
    title: '',
    summary: '',
    type: 'video' as ContentType,
    source_url: '',
    body: '',
  })
  const navigate = useNavigate()

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)

    try {
      const { source_url, body, ...rest } = form
      await createContent({
        ...rest,
        configs: {
          source_url,
          body,
        },
        section_id: sectionId,
        position: 0,
      })
      setOpen(false)
      setForm({
        title: '',
        summary: '',
        type: 'video',
        source_url: '',
        body: '',
      })
      if (onSuccess) {
        onSuccess()
      } else {
        navigate(0)
      }
    } catch (error) {
      console.error('Failed to create content:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger>
        <Button variant="ghost" color="amber" size="1" className="cursor-pointer">
          <Plus size={14} />
          Add Item
        </Button>
      </Dialog.Trigger>

      <Dialog.Content maxWidth="500px">
        <Dialog.Title>Add Content Item</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Add a lesson, video, or document to this section.
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="4">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Item Title
              </Text>
              <TextField.Root
                placeholder="e.g. Introduction to Generics"
                required
                value={form.title}
                onChange={(e) => setForm({ ...form, title: e.target.value })}
              />
            </label>

            <Flex gap="3">
              <Box className="flex-1">
                <Text as="div" size="2" mb="1" weight="bold">
                  Content Type
                </Text>
                <Select.Root
                  value={form.type}
                  onValueChange={(value) => setForm({ ...form, type: value as ContentType })}
                >
                  <Select.Trigger className="w-full" />
                  <Select.Content>
                    <Select.Item value="video">Video</Select.Item>
                    <Select.Item value="pdf">PDF Document</Select.Item>
                    <Select.Item value="rich_text">Rich Text</Select.Item>
                    <Select.Item value="link">External Link</Select.Item>
                    <Select.Item value="embed">Embed Code</Select.Item>
                  </Select.Content>
                </Select.Root>
              </Box>
            </Flex>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Source URL
              </Text>
              <TextField.Root
                placeholder="e.g. https://youtube.com/..."
                value={form.source_url}
                onChange={(e) => setForm({ ...form, source_url: e.target.value })}
              />
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Summary
              </Text>
              <TextArea
                placeholder="Brief description of this item..."
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
              Add Content
            </Button>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  )
}
