import { getJSON, postJSON } from '@/shared/api/http'

import type {
  ContentListResult,
  Course,
  CourseDetail,
  CourseListResult,
  CourseSection,
  SectionListResult,
} from '@/features/courses/model/course'

export async function fetchCourses(): Promise<Course[]> {
  const result = await getJSON<CourseListResult>('/api/v1/courses?page_size=100')
  return result.items
}

export async function fetchCourse(courseId: string): Promise<Course> {
  return getJSON<Course>(`/api/v1/courses/${courseId}`)
}

export async function fetchCourseDetail(courseId: string): Promise<CourseDetail> {
  const course = await fetchCourse(courseId)
  const sectionResult = await getJSON<SectionListResult>(
    `/api/v1/sections?course_id=${encodeURIComponent(courseId)}&page_size=100`,
  )

  const sections = await Promise.all(
    sectionResult.items.map(async (section): Promise<CourseSection> => {
      const contentResult = await getJSON<ContentListResult>(
        `/api/v1/content?section_id=${encodeURIComponent(section.id)}&page_size=100`,
      )

      return {
        ...section,
        contentItems: contentResult.items,
      }
    }),
  )

  return {
    course,
    sections,
  }
}

export async function createCourse(data: Partial<Course>): Promise<Course> {
  return postJSON<Course>('/api/v1/courses', data)
}
