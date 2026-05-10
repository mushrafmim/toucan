export type CourseStatus = 'draft' | 'published' | 'archived'
export type CourseLevel = 'beginner' | 'intermediate' | 'advanced'
export type ContentType =
  | 'video'
  | 'pdf'
  | 'rich_text'
  | 'file'
  | 'link'
  | 'embed'

export type Course = {
  id: string
  title: string
  slug: string
  summary: string
  description: string
  category: string
  level: CourseLevel
  tags: string[]
  status: CourseStatus
  created_at: string
  updated_at: string
  published_at?: string
}

export type CourseListResult = {
  items: Course[]
  page: number
  page_size: number
  total: number
}

export type Section = {
  id: string
  course_id: string
  title: string
  summary: string
  position: number
  created_at: string
  updated_at: string
}

export type SectionListResult = {
  items: Section[]
  page: number
  page_size: number
  total: number
}

export type ContentItem = {
  id: string
  section_id: string
  title: string
  summary: string
  type: ContentType
  position: number
  configs: Record<string, any>
  created_at: string
  updated_at: string
}

export type ContentListResult = {
  items: ContentItem[]
  page: number
  page_size: number
  total: number
}

export type CourseSection = Section & {
  contentItems: ContentItem[]
}

export type CourseDetail = {
  course: Course
  sections: CourseSection[]
  enrollment?: {
    course_id: string
    user_id: string
    role: 'owner' | 'manager' | 'learner'
  } | null
}
