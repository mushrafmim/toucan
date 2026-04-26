import { Text } from '@radix-ui/themes'
import { ArrowLeft } from 'lucide-react'
import { Link, useNavigate } from 'react-router-dom'

type BreadcrumbItem = {
  label: string
  to?: string
}

type BreadcrumbBarProps = {
  items: BreadcrumbItem[]
  backLabel?: string
}

export function BreadcrumbBar({
  items,
  backLabel = 'Go back',
}: BreadcrumbBarProps) {
  const navigate = useNavigate()

  return (
    <div className="flex flex-wrap items-center gap-2">
      <button
        type="button"
        onClick={() => navigate(-1)}
        aria-label={backLabel}
        title={backLabel}
        className="inline-flex h-8 w-8 items-center justify-center rounded-full border-0 bg-transparent p-0 text-[#8a6240] transition hover:bg-[rgba(210,174,126,0.16)] hover:text-[#3d2b1b] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[rgba(121,75,25,0.5)]"
      >
        <ArrowLeft size={16} />
      </button>

      {items.map((item, index) => {
        const isLast = index === items.length-1

        return (
          <div key={`${item.label}-${index}`} className="flex items-center gap-2">
            <Text size="2" className={isLast ? 'uppercase tracking-[0.14em] text-[#5e4731]' : 'uppercase tracking-[0.14em] text-[#8a6240]'}>
              {item.to && !isLast ? (
                <Link to={item.to} className="text-inherit no-underline hover:text-[#3d2b1b]">
                  {item.label}
                </Link>
              ) : (
                item.label
              )}
            </Text>
            {!isLast ? (
              <Text size="2" color="gray">
                /
              </Text>
            ) : null}
          </div>
        )
      })}
    </div>
  )
}
