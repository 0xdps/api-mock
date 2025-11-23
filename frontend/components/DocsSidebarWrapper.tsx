'use client'

import { DocsSidebar } from './DocsSidebar'
import { usePathname } from 'next/navigation'

interface DocsSidebarWrapperProps {
  groups: Record<string, string[]>
}

export function DocsSidebarWrapper({ groups }: DocsSidebarWrapperProps) {
  const pathname = usePathname()
  
  // Extract resource from pathname: /docs/users -> users
  const currentResource = pathname?.startsWith('/docs/') 
    ? pathname.split('/').pop() || ''
    : ''
  
  return <DocsSidebar groups={groups} selectedResource={currentResource} />
}

