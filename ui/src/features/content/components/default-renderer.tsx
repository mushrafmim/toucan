import { Box, Button, Card, Flex, Heading, Text } from '@radix-ui/themes'
import { ExternalLink, FileText } from 'lucide-react'
import type { ContentItem } from '@/features/courses/model/course'

type DefaultRendererProps = {
  item: ContentItem
}

export function DefaultRenderer({ item }: DefaultRendererProps) {
  return (
    <Card size="3" className="h-full w-full bg-white">
      <Flex direction="column" gap="4" className="h-full">
        <Flex align="center" gap="3">
          <Box className="rounded-lg bg-amber-50 p-3 text-amber-600">
            <FileText size={24} />
          </Box>
          <Box>
            <Heading size="4">{item.title}</Heading>
            <Text color="gray" size="2">Type: {item.type}</Text>
          </Box>
        </Flex>

        <Box className="flex-1 overflow-auto rounded-lg bg-gray-50 p-6 border border-dashed border-gray-200">
          {item.body ? (
            <Text size="3" className="whitespace-pre-wrap text-gray-700 leading-relaxed">
              {item.body}
            </Text>
          ) : (
            <Text color="gray" size="2" className="italic">No additional body content available for this item.</Text>
          )}
        </Box>

        {item.source_url && (
          <Flex justify="end">
            <Button variant="soft" color="amber" onClick={() => window.open(item.source_url, '_blank')}>
              <ExternalLink size={16} />
              Open Resource
            </Button>
          </Flex>
        )}
      </Flex>
    </Card>
  )
}
