import './htmx-config'
import './dialog'
import './handlebars-helpers'
import '../styles/main.sass'
import './mode'
import './topbar'
import './board-dnd'
import { renderTemplate } from './dialog'
import { PROJECT, refreshBoard } from './project'

const boardEl = document.getElementById('board')

if (boardEl) {
  boardEl.setAttribute('hx-get', `/api/v1/projects/${PROJECT}/board`)
  boardEl.setAttribute('hx-trigger', 'load, boardUpdated from:body')
  boardEl.setAttribute('hx-ext', 'client-side-templates')
  boardEl.setAttribute('handlebars-template', 'tmpl-board')
  boardEl.setAttribute('hx-target', 'this')
  boardEl.setAttribute('hx-swap', 'innerHTML')
  window.htmx.process(boardEl)
}

function closeAllLanePopovers() {
  document.querySelectorAll('[data-lane-popover]').forEach(pop => { pop.hidden = true })
}

document.addEventListener('click', event => {
  const menuBtn = event.target.closest('[data-lane-menu]')
  if (menuBtn) {
    event.stopPropagation()
    const laneId = menuBtn.getAttribute('data-lane-menu')
    const popover = document.querySelector(`[data-lane-popover="${laneId}"]`)
    const willOpen = popover && popover.hidden
    closeAllLanePopovers()
    if (popover && willOpen) popover.hidden = false
    return
  }

  if (event.target.closest('[data-issue-panel-close]')) {
    closeIssuePanel()
    return
  }

  if (event.target.closest('[data-issue-edit]') && currentIssue) {
    const panel = document.getElementById('issuePanel')
    panel.innerHTML = renderTemplate('tmpl-issue-edit-form', { ...currentIssue, projectSlug: PROJECT })
    window.htmx.process(panel)
    return
  }

  if (event.target.closest('[data-issue-edit-cancel]') && currentIssue) {
    const panel = document.getElementById('issuePanel')
    panel.innerHTML = renderTemplate('tmpl-issue-panel', currentIssue)
    window.htmx.process(panel)
    return
  }

  const addBtn = event.target.closest('.add-card-btn')
  if (addBtn) {
    openAddIssueForm(Number(addBtn.dataset.laneId))
    return
  }

  if (event.target.closest('[data-add-issue-cancel]')) {
    document.querySelectorAll('.new-card-form').forEach(form => form.remove())
    return
  }

  if (event.target.closest('[data-lane-popover]')) return
  closeAllLanePopovers()
})

function openAddIssueForm(laneId) {
  const laneBody = document.querySelector(`.lane-body[data-lane-id="${laneId}"]`)
  if (!laneBody) return
  document.querySelectorAll('.new-card-form').forEach(form => form.remove())
  const html = renderTemplate('tmpl-add-issue-form', { laneId, projectSlug: PROJECT })
  const wrapper = document.createElement('div')
  wrapper.innerHTML = html
  const form = wrapper.firstElementChild
  laneBody.appendChild(form)
  window.htmx.process(form)
  form.querySelector('textarea[name="name"]')?.focus()
}

function closeIssuePanel() {
  const panel = document.getElementById('issuePanel')
  if (panel) panel.hidden = true
  const url = new URL(location)
  if (url.searchParams.has('issue')) {
    url.searchParams.delete('issue')
    history.pushState({}, '', url)
  }
}

window.addEventListener('popstate', () => {
  const ref = new URLSearchParams(location.search).get('issue')
  if (ref) openIssuePanel(ref)
  else closeIssuePanel()
})

let currentIssue = null

function openIssuePanel(ref) {
  const panel = document.getElementById('issuePanel')
  if (!panel) return
  window.htmx.ajax('GET', `/api/v1/projects/${PROJECT}/issues/${ref}`, {
    source: panel,
    target: panel,
    swap: 'innerHTML',
  })
}

const initialIssue = new URLSearchParams(location.search).get('issue')
if (initialIssue) openIssuePanel(initialIssue)

document.body.addEventListener('htmx:afterRequest', event => {
  if (!event.detail.successful) return

  if (event.detail.target?.id === 'issuePanel') {
    try { currentIssue = JSON.parse(event.detail.xhr.responseText) } catch {}
  }

  const verb = event.detail.requestConfig?.verb
  if (!verb || verb === 'get') return
  refreshBoard()
})
