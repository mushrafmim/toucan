import { createBrowserRouter } from 'react-router-dom'

import { AppLayout } from '@/app/layouts/app-layout'
import { ContentPage } from '@/features/content/pages/content-page'
import { FocusedVideoPage } from '@/features/content/pages/focused-video-page'
import { fetchCourseDetail, fetchCourses } from '@/features/courses/api/courses-api'
import { CourseDetailPage } from '@/features/courses/pages/course-detail-page'
import { CoursesPage } from '@/features/courses/pages/courses-page'
import { DashboardPage } from '@/features/dashboard/pages/dashboard-page'
import { SectionsPage } from '@/features/sections/pages/sections-page'
import { fetchUsers } from '@/features/users/api/users-api'
import { UsersPage } from '@/features/users/pages/users-page'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      {
        index: true,
        element: <DashboardPage />,
      },
      {
        path: 'courses',
        element: <CoursesPage />,
        loader: () => fetchCourses(),
      },
      {
        path: 'courses/:courseId',
        element: <CourseDetailPage />,
        loader: async ({ params }) => {
          const courseId = params.courseId
          if (!courseId) {
            throw new Response('Course not found', { status: 404 })
          }
          return fetchCourseDetail(courseId)
        },
      },
      {
        path: 'sections',
        element: <SectionsPage />,
      },
      {
        path: 'content',
        element: <ContentPage />,
      },
      {
        path: 'users',
        element: <UsersPage />,
        loader: () => fetchUsers(),
      },
      {
        path: 'courses/:courseId/content/:contentId',
        element: <FocusedVideoPage />,
      },
    ],
  },
])
