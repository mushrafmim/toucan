import { Card, Flex, Grid, Inset, Separator, Text } from '@radix-ui/themes'

import { PageHeader } from '@/components/page-header'

const dashboardHighlights = [
  {
    title: 'Courses',
    value: '12',
    detail: 'Core learning objects managed from a single catalog.',
  },
  {
    title: 'Sections',
    value: '37',
    detail: 'Top-level structure ready for authoring workflows.',
  },
  {
    title: 'Content Items',
    value: '128',
    detail: 'Videos, PDFs, links, and extensible content formats.',
  },
]

const dashboardNextSteps = [
  'Connect the dashboard to the Go APIs under `/api/v1`.',
  'Introduce feature routes for authoring and learner views.',
  'Add data loading patterns and API error boundaries.',
]

export function DashboardPage() {
  return (
    <Flex direction="column" gap="6">
      <Grid columns={{ initial: '1', md: '2' }} gap="6" align="center">
        <Card
          size="4"
          className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
        >
          <PageHeader
            badge="Admin Workspace"
            title="A deliberate starting point for the learning platform UI."
            description="The app now runs on a proper route tree with a shared shell, so feature pages can be added without reworking the entrypoint again."
          />
        </Card>

        <Card
          size="3"
          className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
        >
          <Inset clip="padding-box" side="top" pb="current">
            <div className="relative min-h-48 overflow-hidden rounded-2xl bg-[linear-gradient(135deg,rgba(77,47,15,0.82),rgba(189,123,37,0.76)),linear-gradient(180deg,rgba(255,255,255,0.06),rgba(255,255,255,0))] before:absolute before:-top-32 before:right-[-6rem] before:h-72 before:w-72 before:rounded-full before:bg-[rgba(255,244,227,0.26)] before:content-[''] after:absolute after:-bottom-12 after:left-[-2rem] after:h-40 after:w-40 after:rounded-full after:bg-[rgba(255,244,227,0.26)] after:content-['']" />
          </Inset>
          <Flex direction="column" gap="3">
            <Text size="2" weight="bold" className="uppercase tracking-[0.18em]">
              Project Setup
            </Text>
            <Text size="3">
              Radix Themes and React Router are both active. The app is ready
              for feature modules and API-backed navigation.
            </Text>
            <Separator size="4" />
            <Flex direction="column" gap="2">
              {dashboardNextSteps.map((item) => (
                <Text key={item} size="2">
                  <span className="mr-1 text-[#9a641d]">•</span>
                  {item}
                </Text>
              ))}
            </Flex>
          </Flex>
        </Card>
      </Grid>

      <Grid columns={{ initial: '1', sm: '3' }} gap="4">
        {dashboardHighlights.map((highlight) => (
          <Card
            key={highlight.title}
            size="3"
            className="bg-[rgba(255,251,245,0.82)] shadow-[0_20px_50px_rgba(106,74,32,0.08)] backdrop-blur-[14px]"
          >
            <Flex direction="column" gap="3">
              <Text size="2" className="uppercase tracking-[0.14em] text-[#8a6240]">
                {highlight.title}
              </Text>
              <Text size="8" weight="bold">
                {highlight.value}
              </Text>
              <Text size="2" color="gray">
                {highlight.detail}
              </Text>
            </Flex>
          </Card>
        ))}
      </Grid>
    </Flex>
  )
}
