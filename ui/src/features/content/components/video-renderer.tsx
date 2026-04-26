import { Box } from '@radix-ui/themes'

type VideoRendererProps = {
  url: string
  title: string
}

export function VideoRenderer({ url, title }: VideoRendererProps) {
  // Simple logic to ensure we use embed URL if it's youtube
  const embedUrl = url.replace('watch?v=', 'embed/')

  return (
    <Box className="relative h-full w-full bg-black shadow-2xl">
      <iframe
        className="absolute inset-0 h-full w-full border-0"
        src={embedUrl}
        title={title}
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
        allowFullScreen
      />
    </Box>
  )
}
