import './handlebars-helpers'

document.body.addEventListener('change', event => {
  const target = event.target
  if (target.id !== 'projectSwitcher') return
  if (target.value) location.href = `/board.html?project=${target.value}`
})

document.body.addEventListener('htmx:afterRequest', event => {
  const form = event.detail.elt
  if (!event.detail.successful || form.id !== 'newProjectForm') return
  const project = JSON.parse(event.detail.xhr.responseText)
  location.href = `/board.html?project=${project.slug}`
})
