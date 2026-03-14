export type UmamiEventData = Record<string, string | number | boolean | null | undefined>

declare global {
  interface Window {
    umami?: {
      track: (eventName: string, eventData?: UmamiEventData) => void
    }
  }
}

export function trackUmamiEvent(eventName: string, eventData: UmamiEventData = {}) {
  if (typeof window === 'undefined') return
  if (!window.umami?.track) return

  try {
    window.umami.track(eventName, eventData)
  } catch {
    // Swallow analytics errors so user actions are never blocked.
  }
}
