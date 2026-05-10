import { postJSON, getJSON } from '@/shared/api/http'
import type { ContentItem } from '@/features/courses/model/course'

export async function createContent(data: Partial<ContentItem>): Promise<ContentItem> {
  return postJSON<ContentItem>('/api/v1/content', data)
}

export async function fetchContentItem(contentId: string): Promise<ContentItem> {
  return getJSON<ContentItem>(`/api/v1/content/${contentId}`)
}
