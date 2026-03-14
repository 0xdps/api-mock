'use client'

import { useEffect } from 'react'
import { usePathname, useSearchParams } from 'next/navigation'
import Script from 'next/script'
import { trackUmamiEvent } from '@/lib/analytics'

const UMAMI_SCRIPT_SRC = 'https://manage.anately.sh/script.js'
const UMAMI_WEBSITE_ID = '6535c752-d8b4-4758-8710-98a1a5f6d751'

function getElementLabel(element: HTMLElement): string {
  const explicitLabel = element.getAttribute('data-umami-label') || element.getAttribute('aria-label')
  if (explicitLabel) return explicitLabel

  const text = element.textContent?.trim().replace(/\s+/g, ' ')
  if (text) return text.slice(0, 80)

  return 'unknown'
}

export function UmamiAnalytics() {
  const pathname = usePathname()
  const searchParams = useSearchParams()

  useEffect(() => {
    const query = searchParams.toString()
    const fullPath = query ? `${pathname}?${query}` : pathname

    trackUmamiEvent('pv', {
      path: fullPath,
      route: pathname,
    })
  }, [pathname, searchParams])

  useEffect(() => {
    const onClick = (event: MouseEvent) => {
      const target = event.target as HTMLElement | null
      if (!target) return

      const clickable = target.closest<HTMLElement>('[data-umami-event], a, button')
      if (!clickable) return

      const explicitEvent = clickable.getAttribute('data-umami-event')
      const eventName = explicitEvent || 'cta'

      const href = clickable instanceof HTMLAnchorElement ? clickable.href : undefined

      trackUmamiEvent(eventName, {
        label: getElementLabel(clickable),
        element: clickable.tagName.toLowerCase(),
        href,
        path: window.location.pathname,
      })
    }

    const onSubmit = (event: SubmitEvent) => {
      const form = event.target as HTMLFormElement | null
      if (!form) return

      trackUmamiEvent('fs', {
        id: form.id || undefined,
        action: form.action || undefined,
        path: window.location.pathname,
      })
    }

    document.addEventListener('click', onClick)
    document.addEventListener('submit', onSubmit)

    return () => {
      document.removeEventListener('click', onClick)
      document.removeEventListener('submit', onSubmit)
    }
  }, [])

  return (
    <Script
      src={UMAMI_SCRIPT_SRC}
      data-website-id={UMAMI_WEBSITE_ID}
      strategy="afterInteractive"
    />
  )
}
