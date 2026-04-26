import { Box, Flex, Text } from '@radix-ui/themes'
import { FileWarning } from 'lucide-react'

type PdfRendererProps = {
  url?: string
  title: string
}

export function PdfRenderer({ url, title }: PdfRendererProps) {
  if (!url) {
    return (
      <Flex align="center" justify="center" className="h-full w-full bg-gray-50 text-gray-400">
        <Flex direction="column" align="center" gap="2">
          <FileWarning size={48} />
          <Text>No PDF source provided.</Text>
        </Flex>
      </Flex>
    )
  }

  return (
    <Box className="h-full w-full overflow-hidden rounded-lg border border-gray-200">
      <iframe
        src={`${url}#view=FitH&toolbar=1`}
        className="h-full w-full border-0"
        title={title}
      />
    </Box>
  )
}
