import {
  Code,
  File,
  FileText,
  Link as LinkIcon,
  Type,
  Video,
  type LucideProps,
} from 'lucide-react'

import type { ContentType } from '@/features/courses/model/course'

type ContentIconProps = LucideProps & {
  type: ContentType
}

export function ContentIcon({ type, ...props }: ContentIconProps) {
  switch (type) {
    case 'video':
      return <Video {...props} />
    case 'pdf':
      return <FileText {...props} />
    case 'rich_text':
      return <Type {...props} />
    case 'file':
      return <File {...props} />
    case 'link':
      return <LinkIcon {...props} />
    case 'embed':
      return <Code {...props} />
    default:
      return <File {...props} />
  }
}
