import type { ContentItem } from '@/features/courses/model/course'
import { VideoRenderer } from './video-renderer'
import { PdfRenderer } from './pdf-renderer'
import { DefaultRenderer } from './default-renderer'

type ContentRendererProps = {
  item: ContentItem
}

export function ContentRenderer({ item }: ContentRendererProps) {
  switch (item.type) {
    case 'video':
      return <VideoRenderer url={item.configs.source_url || ''} title={item.title} />
    case 'pdf':
      return <PdfRenderer url={item.configs.source_url} title={item.title} />
    default:
      return <DefaultRenderer item={item} />
  }
}
