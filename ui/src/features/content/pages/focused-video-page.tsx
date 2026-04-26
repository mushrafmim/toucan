import { Badge, Box, Button, Flex, Heading, Text } from '@radix-ui/themes'
import { ArrowLeft, Maximize2 } from 'lucide-react'
import { useNavigate, useParams } from 'react-router-dom'
import { VideoRenderer } from '../components/video-renderer'

// Note: In a real app, we'd fetch the content detail here via a loader.
// For this implementation, we'll assume the URL provides context or we'd add the loader later.

export function FocusedVideoPage() {
  const navigate = useNavigate()
  const { courseId } = useParams()

  // Placeholder for content fetching logic
  const videoTitle = "Introduction to the Course"
  const videoUrl = "https://www.youtube.com/embed/dQw4w9WgXcQ" // Example embed URL

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
              <Heading size="4" className="tracking-tight">{videoTitle}</Heading>
            </Box>
          </Flex>
          <Button variant="soft" color="gray">
            <Maximize2 size={16} />
            Theater Mode
          </Button>
        </Flex>

        {/* Video Player Section */}
        <Box className="relative flex-1 bg-black">
          <VideoRenderer url={videoUrl} title={videoTitle} />
        </Box>

        {/* Bottom Content / Controls */}
        <Box className="bg-[#1a1412] p-6">
          <Flex justify="between" align="start">
            <Box className="max-w-2xl">
              <Text color="gray" size="2" className="mb-2 block uppercase tracking-widest">Description</Text>
              <Text size="3" className="text-gray-300">
                This lesson covers the core foundations. Take your time to understand the concepts
                before moving to the next section.
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
