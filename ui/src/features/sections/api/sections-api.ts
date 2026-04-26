import { postJSON } from '@/shared/api/http'
import type { Section } from '@/features/courses/model/course'

export async function createSection(data: Partial<Section>): Promise<Section> {
  return postJSON<Section>('/api/v1/sections', data)
}
