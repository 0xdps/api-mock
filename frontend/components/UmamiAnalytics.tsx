'use client'

import { useEffect } from 'react'
import { usePathname } from 'next/navigation'
import Script from 'next/script'
import { trackUmamiEvent } from '@/lib/analytics'

const UMAMI_SCRIPT_SRC = 'https://manage.anately.sh/script.js'
const UMAMI_WEBSITE_ID = 'cf979932-e5c7-4c0b-af26-66828359c028'

function getElementLabel(element: HTMLElement): string {
  const explicitLabel = element.getAttribute('data-umami-label') || element.getAttribute('aria-label')
  if (explicitLabel) return explicitLabel

  const text = element.textContent?.trim().replace(/\s+/g, ' ')
  if (text) return text.slice(0, 80)

  return 'unknown'
}

export function UmamiAnalytics() {
  const pathname = usePathname()

  useEffect(() => {
    const query = typeof window !== 'undefined' ? window.location.search : ''
    const fullPath = query ? `${pathname}${query}` : pathname

    trackUmamiEvent('pv', {
      path: fullPath,
      route: pathname,
    })
  }, [pathname])

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
