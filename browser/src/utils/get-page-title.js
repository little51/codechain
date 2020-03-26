import defaultSettings from '@/settings'

const title = defaultSettings.title || 'Code Chain'

export default function getPageTitle(pageTitle) {
  if (pageTitle) {
    return `${pageTitle} - ${title}`
  }
  return `${title}`
}
