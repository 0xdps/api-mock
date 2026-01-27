'use client'

import { DocsSidebar } from './DocsSidebar'
import { usePathname } from 'next/navigation'

interface DocsSidebarWrapperProps {
  groups: Record<string, string[]>
}

export function DocsSidebarWrapper({ groups }: DocsSidebarWrapperProps) {
  const pathname = usePathname()
  
  // Extract resource from pathname: /resources/users -> users
  const currentResource = pathname?.startsWith('/resources/') 
    ? pathname.split('/').pop() || ''
    : ''
  
  return <DocsSidebar groups={groups} selectedResource={currentResource} />
}

