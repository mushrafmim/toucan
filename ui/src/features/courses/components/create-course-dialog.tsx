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
import { createCourse } from '../api/courses-api'
import { useNavigate } from 'react-router-dom'
import type { CourseLevel } from '../model/course'

type FormState = {
  title: string
  category: string
  level: CourseLevel
  summary: string
  description: string
}

const initialFormState: FormState = {
  title: '',
  category: '',
  level: 'beginner',
  summary: '',
  description: '',
}

export function CreateCourseDialog() {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState<FormState>(initialFormState)
  const navigate = useNavigate()

  const updateField = (field: keyof FormState, value: string) => {
    setForm((prev) => ({ ...prev, [field]: value }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setLoading(true)

    try {
      const course = await createCourse(form)
      setOpen(false)
      setForm(initialFormState)
      navigate(`/courses/${course.id}`)
    } catch (error) {
      console.error('Failed to create course:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger>
        <Button variant="solid" color="amber" size="3" className="cursor-pointer">
          <Plus size={18} />
          Create New Course
        </Button>
      </Dialog.Trigger>

      <Dialog.Content maxWidth="500px">
        <Dialog.Title>Create New Course</Dialog.Title>
        <Dialog.Description size="2" mb="4">
          Fill in the details below to start a new learning journey.
        </Dialog.Description>

        <form onSubmit={handleSubmit}>
          <Flex direction="column" gap="4">
            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Course Title
              </Text>
              <TextField.Root
                placeholder="e.g. Advanced Go Patterns"
                required
                value={form.title}
                onChange={(e) => updateField('title', e.target.value)}
              />
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Category
              </Text>
              <TextField.Root
                placeholder="e.g. Engineering"
                required
                value={form.category}
                onChange={(e) => updateField('category', e.target.value)}
              />
            </label>

            <Flex gap="3" width="100%">
              <Box className="flex-1">
                <Text as="div" size="2" mb="1" weight="bold">
                  Difficulty Level
                </Text>
                <Select.Root
                  value={form.level}
                  onValueChange={(value) => updateField('level', value)}
                >
                  <Select.Trigger className="w-full" />
                  <Select.Content>
                    <Select.Item value="beginner">Beginner</Select.Item>
                    <Select.Item value="intermediate">Intermediate</Select.Item>
                    <Select.Item value="advanced">Advanced</Select.Item>
                  </Select.Content>
                </Select.Root>
              </Box>
            </Flex>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Short Summary
              </Text>
              <TextArea
                placeholder="Briefly describe what this course is about..."
                required
                value={form.summary}
                onChange={(e) => updateField('summary', e.target.value)}
              />
            </label>

            <label>
              <Text as="div" size="2" mb="1" weight="bold">
                Full Description
              </Text>
              <TextArea
                placeholder="Provide a detailed overview of the curriculum..."
                style={{ height: 100 }}
                value={form.description}
                onChange={(e) => updateField('description', e.target.value)}
              />
            </label>
          </Flex>

          <Flex gap="3" mt="6" justify="end">
            <Dialog.Close>
              <Button
                variant="soft"
                color="gray"
                onClick={() => setForm(initialFormState)}
              >
                Cancel
              </Button>
            </Dialog.Close>
            <Button type="submit" loading={loading} color="amber">
              Create Course
            </Button>
          </Flex>
        </form>
      </Dialog.Content>
    </Dialog.Root>
  )
}
