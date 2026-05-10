import { useEffect, useState } from 'react'
import {
  Badge,
  Card,
  Dialog,
  Flex,
  Separator,
  Text,
  Button,
  Box,
  Spinner,
  Callout,
  Tabs,
  Table,
  IconButton, Heading,
} from '@radix-ui/themes'
import { useParams, useNavigate } from 'react-router-dom'
import { ExternalLink, PlayCircle, Plus, AlertCircle, Trash2 } from 'lucide-react'

import type { CourseDetail, ContentItem } from '@/features/courses/model/course'
import { fetchCourseDetail, fetchCourseEnrollments, unenrollUser } from '@/features/courses/api/courses-api'
import { BreadcrumbBar } from '@/components/breadcrumb-bar'
import { ContentIcon } from '@/components/content-icon'
import { PageHeader } from '@/components/page-header'
import { ContentRenderer } from '@/features/content/components/content-renderer'
import { useRole } from '@/shared/context/use-role'
import { CreateSectionDialog } from '@/features/sections/components/create-section-dialog'
import { CreateContentDialog } from '@/features/content/components/create-content-dialog'
import { EnrollUserDialog } from '../components/enroll-user-dialog'

export function CourseDetailPage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { activeRole } = useRole()
  const isInstructor = activeRole === 'instructor' || activeRole === 'admin'

  const [data, setData] = useState<CourseDetail | null>(null)
  const [enrollments, setEnrollments] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  const loadData = async () => {
    if (!id) return
    setLoading(true)
    try {
      const [detail, members] = await Promise.all([
        fetchCourseDetail(id),
        isInstructor ? fetchCourseEnrollments(id) : Promise.resolve([]),
      ])
      setData(detail)
      setEnrollments(members)
      setError(null)
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Failed to load course'))
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadData()
  }, [id, activeRole])

  if (loading) {
    return (
      <Flex align="center" justify="center" style={{ height: '60vh' }}>
        <Spinner size="3" />
      </Flex>
    )
  }

  if (error || !data) {
    return (
      <Box p="4">
        <Callout.Root color="red">
          <Callout.Icon>
            <AlertCircle size={16} />
          </Callout.Icon>
          <Callout.Text>{error?.message || 'Course not found'}</Callout.Text>
        </Callout.Root>
      </Box>
    )
  }

  const { course, sections } = data

  const handleContentClick = (item: ContentItem) => {
    if (item.type === 'video') {
      navigate(`/courses/${course.id}/content/${item.id}`)
    }
  }

  const handleRemoveMember = async (userId: string) => {
    if (!id) return
    if (confirm('Are you sure you want to remove this member?')) {
      try {
        await unenrollUser(id, userId)
        loadData()
      } catch (err) {
        console.error('Failed to remove member:', err)
      }
    }
  }

  return (
    <Flex direction="column" gap="6">
      <BreadcrumbBar
        items={[
          { label: 'Courses', to: '/courses' },
          { label: course.title },
        ]}
        backLabel="Go back to previous page"
      />

      <Flex justify="between" align={{ initial: 'start', md: 'center' }} gap="4" wrap="wrap">
        <PageHeader
          title={course.title}
          description={course.description}
          titleClassName="max-w-none tracking-[-0.04em]"
          descriptionClassName="max-w-4xl"
          afterTitle={
            course.category ? (
              <Badge size="2" color="bronze" variant="soft" className="self-start">
                {course.category}
              </Badge>
            ) : null
          }
        />
        <Flex direction="column" gap="2" align={{ initial: 'start', md: 'end' }}>
          <Badge size="2" color="amber" variant="soft">
            {course.status}
          </Badge>
          <Text size="2" color="gray">
            Level: {course.level}
          </Text>
        </Flex>
      </Flex>

      <Tabs.Root defaultValue="curriculum">
        <Tabs.List size="2">
          <Tabs.Trigger value="curriculum">Curriculum</Tabs.Trigger>
          {isInstructor && <Tabs.Trigger value="members">Members</Tabs.Trigger>}
        </Tabs.List>

        <Box pt="6">
          <Tabs.Content value="curriculum">
            <Flex direction="column" gap="6">
              <Card
                size="3"
                className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
              >
                <Flex direction="column" gap="2">
                  <Text size="2" className="uppercase tracking-[0.14em] text-[#8a6240]">
                    Summary
                  </Text>
                  <Text size="3">{course.summary}</Text>
                </Flex>
              </Card>

              <Flex direction="column" gap="0">
                {isInstructor && <AddSectionWide courseId={course.id} position={1} />}

                {sections.map((section) => (
                  <Box key={section.id}>
                    <Card
                      size="3"
                      className="relative z-10 bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
                    >
                      <Flex direction="column" gap="4">
                        <Flex direction="column">
                          <Flex justify="between" align="start" gap="4" wrap="wrap">
                            <Text size="5" weight="bold">
                              {section.title}
                            </Text>
                            <Flex align="center" gap="3">
                              {isInstructor && <CreateContentDialog sectionId={section.id} />}
                              <Badge size="2" color="bronze" variant="soft">
                                {section.contentItems.length} items
                              </Badge>
                            </Flex>
                          </Flex>
                          <Text size="2" color="gray">
                            {section.summary}
                          </Text>
                        </Flex>

                        <Separator size="4" />

                        <Flex direction="column" gap="3">
                          {section.contentItems.map((item) => (
                            <Dialog.Root key={item.id}>
                              {item.type === 'pdf' ? (
                                <Dialog.Trigger>
                                  <Card
                                    size="2"
                                    className="cursor-pointer bg-[rgba(255,252,247,0.92)] transition hover:bg-[rgba(255,248,238,1)] hover:shadow-md"
                                  >
                                    <ContentItemSummary item={item} />
                                  </Card>
                                </Dialog.Trigger>
                              ) : (
                                <Card
                                  size="2"
                                  className="cursor-pointer bg-[rgba(255,252,247,0.92)] transition hover:bg-[rgba(255,248,238,1)] hover:shadow-md"
                                  onClick={() => handleContentClick(item)}
                                >
                                  <ContentItemSummary item={item} />
                                </Card>
                              )}

                              <Dialog.Content maxWidth="90vw" style={{ height: '85vh' }}>
                                <Dialog.Title>{item.title}</Dialog.Title>
                                <Dialog.Description size="2" mb="4">
                                  {item.summary}
                                </Dialog.Description>

                                <Box className="h-[calc(100%-120px)] w-full">
                                  <ContentRenderer item={item} />
                                </Box>

                                <Flex gap="3" mt="4" justify="end">
                                  <Dialog.Close>
                                    <Button variant="soft" color="gray">
                                      Close
                                    </Button>
                                  </Dialog.Close>
                                  {item.source_url && (
                                    <Button onClick={() => window.open(item.source_url, '_blank')}>
                                      <ExternalLink size={14} />
                                      Open in New Tab
                                    </Button>
                                  )}
                                </Flex>
                              </Dialog.Content>
                            </Dialog.Root>
                          ))}
                        </Flex>
                      </Flex>
                    </Card>
                    {isInstructor && (
                      <Box my="-2">
                        <AddSectionWide courseId={course.id} position={section.position + 1} />
                      </Box>
                    )}
                  </Box>
                ))}
              </Flex>
            </Flex>
          </Tabs.Content>

          {isInstructor && (
            <Tabs.Content value="members">
              <Flex direction="column" gap="4">
                <Flex justify="between" align="center">
                  <Heading size="4">Course Members</Heading>
                  <EnrollUserDialog courseId={course.id} onSuccess={loadData} />
                </Flex>

                <Card size="2">
                  <Table.Root variant="surface">
                    <Table.Header>
                      <Table.Row>
                        <Table.ColumnHeaderCell>User ID</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Role</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell>Enrolled</Table.ColumnHeaderCell>
                        <Table.ColumnHeaderCell></Table.ColumnHeaderCell>
                      </Table.Row>
                    </Table.Header>
                    <Table.Body>
                      {enrollments.map((member) => (
                        <Table.Row key={member.user_id}>
                          <Table.Cell>{member.user_id}</Table.Cell>
                          <Table.Cell>
                            <Badge variant="soft" color={member.role === 'owner' ? 'red' : member.role === 'manager' ? 'blue' : 'amber'}>
                              {member.role}
                            </Badge>
                          </Table.Cell>
                          <Table.Cell>
                            <Text size="1" color="gray">
                              {new Date(member.created_at).toLocaleDateString()}
                            </Text>
                          </Table.Cell>
                          <Table.Cell>
                            {member.role !== 'owner' && (
                              <IconButton
                                size="1"
                                variant="ghost"
                                color="red"
                                onClick={() => handleRemoveMember(member.user_id)}
                              >
                                <Trash2 size={14} />
                              </IconButton>
                            )}
                          </Table.Cell>
                        </Table.Row>
                      ))}
                    </Table.Body>
                  </Table.Root>
                </Card>
              </Flex>
            </Tabs.Content>
          )}
        </Box>
      </Tabs.Root>
    </Flex>
  )
}

function AddSectionWide({ courseId, position }: { courseId: string; position: number }) {
  return (
    <CreateSectionDialog
      courseId={courseId}
      position={position}
      trigger={
        <button className="group relative flex w-full cursor-pointer items-center justify-center py-2 outline-none border-none bg-transparent">
          <div className="absolute inset-x-0 h-px bg-transparent transition-colors group-hover:bg-amber-300" />
          <div className="z-10 flex scale-90 items-center gap-2 rounded-full border border-amber-200 bg-white px-3 py-1 text-[10px] font-bold uppercase tracking-wider text-amber-700 opacity-0 shadow-sm transition-all group-hover:scale-100 group-hover:opacity-100">
            <Plus size={12} />
            Add Section
          </div>
        </button>
      }
    />
  )
}

function ContentItemSummary({ item }: { item: ContentItem }) {
  return (
    <Flex direction="column" gap="2">
      <Flex justify="between" align="center" gap="3">
        <Flex align="center" gap="2">
          <ContentIcon type={item.type} size={18} className="text-[#8a6240]" />
          <Text size="3" weight="bold">
            {item.title}
          </Text>
        </Flex>
        <Flex align="center" gap="2">
          {item.type === 'video' && <PlayCircle size={16} className="text-amber-600" />}
          <Badge size="1" color="gray" variant="soft">
            {item.type}
          </Badge>
        </Flex>
      </Flex>
      <Text size="2" color="gray">
        {item.summary}
      </Text>
    </Flex>
  )
}
