import './htmx-bootstrap'

document.body.addEventListener('htmx:responseError', event => {
  const xhr = event.detail.xhr
  console.error('htmx error', xhr.status, xhr.responseText)
})

window.addEventListener('popstate', () => {
  const main = document.querySelector('[data-spa-root]')
  if (!main) return
  window.htmx.ajax('GET', location.pathname + location.search, {
    target: main,
    swap: 'innerHTML',
  })
})
