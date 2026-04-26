import {Badge, Card, Dialog, Flex, Separator, Text, Button, Box} from '@radix-ui/themes'
import { useLoaderData, useNavigate } from 'react-router-dom'
import { ExternalLink, PlayCircle } from 'lucide-react'

import type { CourseDetail, ContentItem } from '@/features/courses/model/course'
import { BreadcrumbBar } from '@/components/breadcrumb-bar'
import { ContentIcon } from '@/components/content-icon'
import { PageHeader } from '@/components/page-header'
import { ContentRenderer } from '@/features/content/components/content-renderer'
import { useRole } from '@/shared/context/use-role'
import { CreateSectionDialog } from '@/features/sections/components/create-section-dialog'
import { CreateContentDialog } from '@/features/content/components/create-content-dialog'

export function CourseDetailPage() {
  const { course, sections } = useLoaderData() as CourseDetail
  const navigate = useNavigate()
  const { activeRole } = useRole()
  const isTeacher = activeRole === 'teacher'

  const handleContentClick = (item: ContentItem) => {
    if (item.type === 'video') {
      navigate(`/courses/${course.id}/content/${item.id}`)
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

      <Card
        size="3"
        className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
      >
        <Flex justify="between" align="center" gap="4" wrap="wrap">
          <Flex direction="column" gap="2">
            <Text size="2" className="uppercase tracking-[0.14em] text-[#8a6240]">
              Summary
            </Text>
            <Text size="3">{course.summary}</Text>
          </Flex>
          {isTeacher && <CreateSectionDialog courseId={course.id} />}
        </Flex>
      </Card>

      <Flex direction="column" gap="4">
        {sections.map((section) => (
          <Card
            key={section.id}
            size="3"
            className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
          >
            <Flex direction="column" gap="4">
              <Flex direction="column">
                <Flex justify="between" align="start" gap="4" wrap="wrap">
                  <Text size="5" weight="bold">
                    {section.title}
                  </Text>
                  <Flex align="center" gap="3">
                    {isTeacher && <CreateContentDialog sectionId={section.id} />}
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
                          <Button variant="soft" color="gray">Close</Button>
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
        ))}
      </Flex>
    </Flex>
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
