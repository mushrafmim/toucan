import type { ReactNode } from 'react'

import { Badge, Flex, Heading, Text } from '@radix-ui/themes'

type PageHeaderProps = {
  badge?: string
  title: string
  description: string
  titleClassName?: string
  descriptionClassName?: string
  afterTitle?: ReactNode
}

export function PageHeader({
  badge,
  title,
  description,
  titleClassName,
  descriptionClassName,
  afterTitle,
}: PageHeaderProps) {
  return (
    <Flex direction="column" gap="3">
      <Heading
        size="8"
        className={
          titleClassName ?? 'max-w-[14ch] tracking-[-0.04em] max-md:max-w-none'
        }
      >
        {title}
      </Heading>
      {badge ? (
        <Badge size="2" color="bronze" variant="soft">
          {badge}
        </Badge>
      ) : null}
      {afterTitle}
      <Text size="4" className={descriptionClassName ?? 'max-w-3xl'}>
        {description}
      </Text>
    </Flex>
  )
}
