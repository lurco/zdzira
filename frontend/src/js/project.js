export const PROJECT = new URLSearchParams(location.search).get('project') || 'main'

export const PROJECT_API = `/api/v1/projects/${PROJECT}`

export function refreshBoard() {
  window.htmx.trigger(document.body, 'boardUpdated')
}
