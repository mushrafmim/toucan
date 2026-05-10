import { useEffect, useState } from 'react'
import { Badge, Box, Button, Flex, Heading, Text, Spinner, Callout } from '@radix-ui/themes'
import { ArrowLeft, Maximize2, AlertCircle } from 'lucide-react'
import { useNavigate, useParams } from 'react-router-dom'
import { VideoRenderer } from '../components/video-renderer'
import { fetchContentItem } from '../api/content-api'
import type { ContentItem } from '@/features/courses/model/course'

export function FocusedVideoPage() {
  const navigate = useNavigate()
  const { courseId, contentId } = useParams()
  
  const [content, setContent] = useState<ContentItem | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    if (!contentId) return

    let active = true
    setLoading(true)

    fetchContentItem(contentId)
      .then((item) => {
        if (active) {
          setContent(item)
          setError(null)
        }
      })
      .catch((err) => {
        if (active) setError(err instanceof Error ? err : new Error('Failed to load video'))
      })
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [contentId])

  if (loading) {
    return (
      <Box className="fixed inset-0 z-50 bg-[#0f0a08] text-white">
        <Flex align="center" justify="center" className="h-full">
          <Spinner size="3" />
        </Flex>
      </Box>
    )
  }

  if (error || !content) {
    return (
      <Box className="fixed inset-0 z-50 bg-[#0f0a08] text-white p-6">
        <Callout.Root color="red">
          <Callout.Icon>
            <AlertCircle size={16} />
          </Callout.Icon>
          <Callout.Text>
            {error?.message || 'Video not found'}
          </Callout.Text>
        </Callout.Root>
        <Button mt="4" onClick={() => navigate(`/courses/${courseId}`)}>
          Go Back
        </Button>
      </Box>
    )
  }

  return (
    <Box className="fixed inset-0 z-50 bg-[#0f0a08] text-white">
      <Flex direction="column" className="h-full">
        {/* Top Bar */}
        <Flex justify="between" align="center" className="bg-[#1a1412] px-6 py-4 shadow-xl">
          <Flex align="center" gap="4">
            <Button
              variant="ghost"
              color="gray"
              onClick={() => navigate(`/courses/${courseId}`)}
              className="cursor-pointer transition-transform hover:-translate-x-1"
            >
              <ArrowLeft size={20} />
            </Button>
            <Box>
              <Badge color="amber" variant="soft" className="mb-1">VIDEO LESSON</Badge>
              <Heading size="4" className="tracking-tight">{content.title}</Heading>
            </Box>
          </Flex>
          <Button variant="soft" color="gray">
            <Maximize2 size={16} />
            Theater Mode
          </Button>
        </Flex>

        {/* Video Player Section */}
        <Box className="relative flex-1 bg-black">
          <VideoRenderer url={content.source_url || ''} title={content.title} />
        </Box>

        {/* Bottom Content / Controls */}
        <Box className="bg-[#1a1412] p-6">
          <Flex justify="between" align="start">
            <Box className="max-w-2xl">
              <Text color="gray" size="2" className="mb-2 block uppercase tracking-widest">Description</Text>
              <Text size="3" className="text-gray-300">
                {content.description || 'This lesson covers the core foundations. Take your time to understand the concepts before moving to the next section.'}
              </Text>
            </Box>
            <Flex gap="3">
              <Button size="3" variant="outline" color="gray">Previous</Button>
              <Button size="3" color="amber">Mark as Complete</Button>
            </Flex>
          </Flex>
        </Box>
      </Flex>
    </Box>
  )
}
