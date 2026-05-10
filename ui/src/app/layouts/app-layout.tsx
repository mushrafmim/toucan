import { useEffect, useRef, useState } from 'react'

import { Box, Flex, Heading, Text } from '@radix-ui/themes'
import {
  BellDot,
  BookOpen,
  House,
  PanelLeftClose,
  PanelLeftOpen,
  Search,
  Users,
  Settings,
  LayoutDashboard,
} from 'lucide-react'
import { NavLink, Outlet, useRevalidator } from 'react-router-dom'
import { RoleSwitcher } from '@/components/role-switcher'
import { AuthControls } from '@/shared/auth/auth-controls'
import type { UserRole } from '@/shared/auth/roles'
import { useRole } from '@/shared/context/use-role'

type NavItem = {
  to: string
  label: string
  icon: typeof House
  end?: boolean
  roles?: UserRole[]
}

const navItems: NavItem[] = [
  { to: '/', label: 'Overview', end: true, icon: House },
  { to: '/courses', label: 'Courses', icon: BookOpen },
  { to: '/curriculum', label: 'Curriculum', icon: LayoutDashboard, roles: ['instructor', 'admin'] },
  { to: '/users', label: 'Users', icon: Users, roles: ['admin'] },
  { to: '/settings', label: 'Settings', icon: Settings, roles: ['admin'] },
]

export function AppLayout() {
  const [collapsed, setCollapsed] = useState(false)
  const { activeRole } = useRole()
  const revalidator = useRevalidator()
  const previousRole = useRef(activeRole)

  useEffect(() => {
    if (previousRole.current === activeRole) {
      return
    }
    previousRole.current = activeRole
    revalidator.revalidate()
  }, [activeRole, revalidator])

  const visibleNavItems = navItems.filter(
    (item) => !item.roles || item.roles.includes(activeRole)
  )

  return (
    <Box className="h-dvh overflow-hidden text-[#2f241d]">
      <header className="sticky top-0 z-10 border-b border-[rgba(134,99,57,0.14)] bg-[linear-gradient(180deg,rgba(250,243,232,0.98),rgba(244,234,218,0.94))] backdrop-blur-xl">
        <Flex
          justify="between"
          align="center"
          gap="4"
          wrap="wrap"
          className="min-h-[4.75rem] w-full box-border px-6 py-4"
        >
          <Flex align="center" gap="3">
            <Box className="h-11 w-11 rounded-2xl bg-[linear-gradient(135deg,rgba(113,67,22,0.92),rgba(244,167,67,0.92))] shadow-[0_18px_40px_rgba(120,82,31,0.18),inset_0_1px_0_rgba(255,255,255,0.35)]" />
            <Box>
              <Text size="1" weight="bold" className="uppercase tracking-[0.18em]">
                Toucan LMS
              </Text>
              <Heading size="4">Learning Workspace</Heading>
            </Box>
          </Flex>

          <Flex align="center" gap="3" wrap="wrap">
            <RoleSwitcher />
            <div className="inline-flex min-h-10 items-center justify-center gap-2 rounded-full border border-[rgba(134,99,57,0.1)] bg-[rgba(255,248,238,0.88)] px-3 py-2 text-[#6e5842] md:px-3.5">
              <Search size={16} />
            </div>
            <div className="inline-flex min-h-10 items-center justify-center gap-2 rounded-full border border-[rgba(134,99,57,0.1)] bg-[rgba(255,248,238,0.88)] px-3 py-2 text-[#6e5842] md:px-3.5">
              <BellDot size={16} />
            </div>
            <AuthControls />
          </Flex>
        </Flex>
      </header>

      <div
        className={`grid h-[calc(100dvh-4.75rem)] min-h-0 grid-cols-1 ${
          collapsed ? 'md:grid-cols-[5.75rem_minmax(0,1fr)]' : 'md:grid-cols-[18rem_minmax(0,1fr)]'
        }`}
      >
        <aside className="hidden h-full min-h-0 overflow-hidden border-r border-r-[rgba(134,99,57,0.18)] bg-[linear-gradient(180deg,rgba(255,247,234,0.98),rgba(240,228,209,0.96))] px-4 py-4 shadow-[inset_-1px_0_0_rgba(255,255,255,0.45)] backdrop-blur-xl md:block">
          <Flex direction="column" className="h-auto min-h-0 gap-4 md:h-full">
            <Flex direction="column" gap="4">
              <nav aria-label="Primary" className="flex flex-col gap-2">
                {visibleNavItems.map((item) => (
                  <NavLink
                    key={item.label}
                    to={item.to}
                    end={item.end}
                    className={({ isActive }) =>
                      [
                        'inline-flex min-h-12 items-center gap-3 rounded-2xl px-3.5 py-3 no-underline transition duration-150 md:flex-none',
                        collapsed ? 'justify-center px-3' : 'justify-start',
                        isActive
                          ? 'bg-[linear-gradient(135deg,rgba(121,75,25,0.92),rgba(211,148,55,0.92))] text-[#fffaf2] shadow-[0_12px_28px_rgba(113,67,22,0.2)]'
                          : 'text-[#6e5842] hover:bg-[rgba(210,174,126,0.18)] hover:text-[#3d2b1b] md:hover:translate-x-0.5',
                      ].join(' ')
                    }
                    aria-label={item.label}
                    title={collapsed ? item.label : undefined}
                  >
                    <item.icon size={18} className="shrink-0" />
                    {!collapsed ? <span className="whitespace-nowrap">{item.label}</span> : null}
                  </NavLink>
                ))}
              </nav>
            </Flex>

            <Flex direction="column" gap="3" className="mt-4 md:mt-auto">
              <button
                type="button"
                className={[
                  'inline-flex min-h-10 items-center gap-2 rounded-2xl border-0 bg-[rgba(255,251,245,0.82)] text-[#6e5842] shadow-[0_10px_30px_rgba(106,74,32,0.08)] transition duration-150 hover:bg-[rgba(230,196,150,0.3)] hover:text-[#3d2b1b] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[rgba(121,75,25,0.5)]',
                  collapsed ? 'justify-center px-3 py-3' : 'justify-start px-3.5 py-3',
                ].join(' ')}
                onClick={() => setCollapsed((value) => !value)}
                aria-label={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
              >
                {collapsed ? <PanelLeftOpen size={18} /> : <PanelLeftClose size={18} />}
                {!collapsed ? <span className="whitespace-nowrap">Collapse</span> : null}
              </button>
            </Flex>
          </Flex>
        </aside>

        <Box className="min-h-0 min-w-0 overflow-auto">
          <div className="mx-auto w-full max-w-7xl px-5 py-5 pb-24 md:pb-5">
            <Outlet />
          </div>
        </Box>
      </div>

      <nav
        aria-label="Mobile"
        className="fixed inset-x-0 bottom-0 z-20 border-t border-[rgba(134,99,57,0.14)] bg-[linear-gradient(180deg,rgba(250,243,232,0.98),rgba(244,234,218,0.96))] px-3 py-2 backdrop-blur-xl md:hidden"
      >
        <div className="mx-auto flex max-w-md justify-around gap-2">
          {visibleNavItems.map((item) => (
            <NavLink
              key={item.label}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                [
                  'inline-flex min-h-14 flex-col items-center justify-center gap-1.5 rounded-2xl px-3 py-2 no-underline transition duration-150 flex-1',
                  isActive
                    ? 'bg-[linear-gradient(135deg,rgba(121,75,25,0.92),rgba(211,148,55,0.92))] text-[#fffaf2] shadow-[0_12px_28px_rgba(113,67,22,0.2)]'
                    : 'text-[#6e5842] hover:bg-[rgba(210,174,126,0.18)] hover:text-[#3d2b1b]',
                ].join(' ')
              }
            >
              <item.icon size={18} className="shrink-0" />
              <span className="text-[10px] font-medium truncate">{item.label}</span>
            </NavLink>
          ))}
        </div>
      </nav>
    </Box>
  )
}
