import { onUnmounted, watch, type Ref } from 'vue'

let activeLocks = 0
let savedScrollY = 0
let savedBodyStyle = ''
let savedHtmlStyle = ''

const lockBody = () => {
  if (!import.meta.client || activeLocks++ > 0) return

  const body = document.body
  const html = document.documentElement
  savedScrollY = window.scrollY
  savedBodyStyle = body.getAttribute('style') || ''
  savedHtmlStyle = html.getAttribute('style') || ''

  const scrollbarWidth = window.innerWidth - html.clientWidth
  const currentPadding = Number.parseFloat(window.getComputedStyle(body).paddingRight) || 0

  html.style.overflow = 'hidden'
  body.style.position = 'fixed'
  body.style.inset = `-${savedScrollY}px 0 0`
  body.style.width = '100%'
  body.style.overflow = 'hidden'
  if (scrollbarWidth > 0) body.style.paddingRight = `${currentPadding + scrollbarWidth}px`
}

const unlockBody = () => {
  if (!import.meta.client || activeLocks === 0 || --activeLocks > 0) return

  if (savedBodyStyle) document.body.setAttribute('style', savedBodyStyle)
  else document.body.removeAttribute('style')
  if (savedHtmlStyle) document.documentElement.setAttribute('style', savedHtmlStyle)
  else document.documentElement.removeAttribute('style')
  window.scrollTo(0, savedScrollY)
}

/** Locks page scrolling while a modal is open and safely supports nested modals. */
export const useBodyScrollLock = (isOpen: Ref<boolean>) => {
  let ownsLock = false

  const stop = watch(isOpen, (open) => {
    if (open && !ownsLock) {
      lockBody()
      ownsLock = true
    } else if (!open && ownsLock) {
      unlockBody()
      ownsLock = false
    }
  }, { immediate: true })

  onUnmounted(() => {
    stop()
    if (ownsLock) unlockBody()
  })
}
