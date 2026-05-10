import { useEffect, useState } from 'react'
import { Badge, Card, Flex, Grid, Text, Spinner, Callout } from '@radix-ui/themes'
import { Link } from 'react-router-dom'
import { AlertCircle } from 'lucide-react'

import type { Course } from '@/features/courses/model/course'
import { fetchMyCourses } from '@/features/courses/api/courses-api'
import { PageHeader } from '@/components/page-header'
import { useRole } from '@/shared/context/use-role'
import { CreateCourseDialog } from '../components/create-course-dialog'

export function CoursesPage() {
  const { activeRole } = useRole()
  const canCreate = activeRole === 'instructor'

  const [courses, setCourses] = useState<Course[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    let active = true
    setLoading(true)

    fetchMyCourses()
      .then((res) => {
        if (active) {
          setCourses(res)
          setError(null)
        }
      })
      .catch((err) => {
        if (active) setError(err instanceof Error ? err : new Error('Failed to load courses'))
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
          badge="Courses"
          title="Courses"
          description="Select a course to open its detail page, where sections and content are composed from the split domain APIs."
        />
        {canCreate && (
          <Flex pb="2">
            <CreateCourseDialog />
          </Flex>
        )}
      </Flex>

      {error ? (
        <Callout.Root color="red">
          <Callout.Icon>
            <AlertCircle size={16} />
          </Callout.Icon>
          <Callout.Text>
            {error.message}
          </Callout.Text>
        </Callout.Root>
      ) : (
        <Grid columns={{ initial: '1', md: '3' }} gap="4">
          {courses.map((course) => (
            <Card
              key={course.id}
              size="3"
              className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px] transition duration-150 hover:-translate-y-0.5 hover:shadow-[0_24px_60px_rgba(106,74,32,0.14)]"
            >
              <Link to={`/courses/${course.id}`} className="block text-inherit no-underline">
                <Flex direction="column" gap="3">
                  <Flex justify="between" align="center" gap="3">
                    <Text size="2" className="uppercase tracking-[0.14em] text-[#8a6240]">
                      {course.category || 'Course'}
                    </Text>
                    <Badge size="1" color="amber" variant="soft">
                      {course.status}
                    </Badge>
                  </Flex>
                  <Text size="5" weight="bold">
                    {course.title}
                  </Text>
                  <Text size="2" color="gray">
                    {course.summary}
                  </Text>
                  <Text size="2" className="text-[#7a6149]">
                    Level: {course.level}
                  </Text>
                </Flex>
              </Link>
            </Card>
          ))}
        </Grid>
      )}
    </Flex>
  )
}
