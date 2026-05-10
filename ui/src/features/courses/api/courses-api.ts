import { getJSON, postJSON, deleteItem } from '@/shared/api/http'

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

export async function fetchMyCourses(): Promise<Course[]> {
  const result = await getJSON<CourseListResult>('/api/v1/me/courses?page_size=100')
  return result.items
}

export async function fetchCourse(courseId: string): Promise<Course> {
  return getJSON<Course>(`/api/v1/courses/${courseId}`)
}

export async function fetchMyEnrollment(courseId: string): Promise<any> {
  try {
    return await getJSON<any>(`/api/v1/courses/${courseId}/member/me`)
  } catch (error) {
    return null
  }
}

export async function fetchCourseDetail(courseId: string): Promise<CourseDetail> {
  const course = await fetchCourse(courseId)
  const enrollment = await fetchMyEnrollment(courseId)
  
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
    enrollment,
  }
}

export async function fetchCourseEnrollments(courseId: string): Promise<any[]> {
  return getJSON<any[]>(`/api/v1/courses/${courseId}/enrollments`)
}

export async function enrollUser(data: { course_id: string; user_id: string; role: string }): Promise<any> {
  return postJSON<any>('/api/v1/enrollments', data)
}

export async function unenrollUser(courseId: string, userId: string): Promise<void> {
  await deleteItem(`/api/v1/enrollments?course_id=${encodeURIComponent(courseId)}&user_id=${encodeURIComponent(userId)}`)
}

export async function createCourse(data: Partial<Course>): Promise<Course> {
  return postJSON<Course>('/api/v1/courses', data)
}
