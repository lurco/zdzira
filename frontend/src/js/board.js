import './htmx-config'
import './dialog'
import './handlebars-helpers'
import '../styles/main.sass'
import './mode'
import './topbar'
import './board-dnd'
import { PROJECT, PROJECT_API } from './project'

const boardEl = document.getElementById('board')

if (boardEl) {
  boardEl.setAttribute('hx-get', `${PROJECT_API}/board`)
  boardEl.setAttribute('hx-trigger', 'load, boardUpdated from:body')
  boardEl.setAttribute('hx-ext', 'client-side-templates')
  boardEl.setAttribute('handlebars-template', 'tmpl-board')
  boardEl.setAttribute('hx-target', 'this')
  boardEl.setAttribute('hx-swap', 'innerHTML')
  window.htmx.process(boardEl)
}

document.addEventListener('click', event => {
  const card = event.target.closest('.card[data-card-ref]')
  if (!card) return
  location.href = `/issue.html?project=${PROJECT}&ref=${card.dataset.cardRef}`
})
