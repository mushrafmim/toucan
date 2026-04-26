import { Card, Flex, Grid, Text } from '@radix-ui/themes'

import { PageHeader } from '@/components/page-header'

const sectionCapabilities = [
  'Top-level section queries',
  'Section ordering and editing',
  'Filtering by course_id',
]

export function SectionsPage() {
  return (
    <Flex direction="column" gap="6">
      <PageHeader
        badge="Sections"
        title="Sections have their own route surface now."
        description="This page mirrors the backend boundary you introduced, so the frontend can evolve around sections as a first-class domain."
      />

      <Grid columns={{ initial: '1', md: '3' }} gap="4">
        {sectionCapabilities.map((item) => (
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
