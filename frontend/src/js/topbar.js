import './handlebars-helpers'

document.body.addEventListener('change', event => {
  const target = event.target
  if (target.id !== 'projectSwitcher') return
  if (target.value) location.href = `/board.html?project=${target.value}`
})
