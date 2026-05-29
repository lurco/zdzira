import './htmx-bootstrap'

window.addEventListener('error', event => {
  console.error('window error', event.error || event.message)
})
window.addEventListener('unhandledrejection', event => {
  console.error('unhandled rejection', event.reason)
})

document.body.addEventListener('htmx:responseError', event => {
  const xhr = event.detail.xhr
  console.error('htmx error', xhr.status, xhr.responseText)
})

document.body.addEventListener('htmx:sendError', event => {
  console.error('htmx send error', event.detail)
})

window.addEventListener('popstate', () => {
  const main = document.querySelector('[data-spa-root]')
  if (!main) return
  window.htmx.ajax('GET', location.pathname + location.search, {
    target: main,
    swap: 'innerHTML',
  })
})
