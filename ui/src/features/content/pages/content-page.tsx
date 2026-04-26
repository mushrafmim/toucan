import { Card, Flex, Grid, Text } from '@radix-ui/themes'

import { PageHeader } from '@/components/page-header'

const contentCapabilities = [
  'Typed content item management',
  'Filtering by section_id',
  'Video, PDF, and rich text authoring flows',
]

export function ContentPage() {
  return (
    <Flex direction="column" gap="6">
      <PageHeader
        badge="Content"
        title="Content items are routed as a separate top-level domain."
        description="This route is where content-specific forms, previews, and future type-aware editors should live."
      />

      <Grid columns={{ initial: '1', md: '3' }} gap="4">
        {contentCapabilities.map((item) => (
          <Card
            key={item}
            size="3"
            className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
          >
            <Text size="3">{item}</Text>
          </Card>
        ))}
      </Grid>
    </Flex>
  )
}
