import './htmx-config'
import './dialog'
import '../styles/main.sass'
import './mode'
import './topbar'

document.body.addEventListener('htmx:afterRequest', event => {
  const form = event.detail.elt
  if (!event.detail.successful) return
  if (form.id !== 'newProjectForm') return

  const project = JSON.parse(event.detail.xhr.responseText)
  location.href = `/board.html?project=${project.slug}`
})
