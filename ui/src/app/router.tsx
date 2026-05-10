import { createBrowserRouter } from 'react-router-dom'

import { AppLayout } from '@/app/layouts/app-layout'
import { ContentPage } from '@/features/content/pages/content-page'
import { FocusedVideoPage } from '@/features/content/pages/focused-video-page'
import { CourseDetailPage } from '@/features/courses/pages/course-detail-page'
import { CoursesPage } from '@/features/courses/pages/courses-page'
import { DashboardPage } from '@/features/dashboard/pages/dashboard-page'
import { SectionsPage } from '@/features/sections/pages/sections-page'
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
      },
      {
        path: 'courses/:id',
        element: <CourseDetailPage />,
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
      },
      {
        path: 'courses/:courseId/content/:contentId',
        element: <FocusedVideoPage />,
      },
    ],
  },
])
