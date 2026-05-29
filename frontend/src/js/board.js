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
let currentIssue = null
let currentEpics = []

function syncBoardURL() {
  const epic = new URLSearchParams(location.search).get('epic') || ''
  const path = epic
    ? `/api/v1/projects/${PROJECT}/board?epic=${encodeURIComponent(epic)}`
    : `/api/v1/projects/${PROJECT}/board`
  boardEl.setAttribute('hx-get', path)
}

if (boardEl) {
  syncBoardURL()
  boardEl.setAttribute('hx-trigger', 'boardUpdated from:body')
  boardEl.setAttribute('hx-ext', 'client-side-templates')
  boardEl.setAttribute('handlebars-template', 'tmpl-board')
  boardEl.setAttribute('hx-target', 'this')
  boardEl.setAttribute('hx-swap', 'innerHTML')
  window.htmx.process(boardEl)
  refreshBoard()
}

const epicFilterEl = document.getElementById('epicFilter')
if (epicFilterEl) {
  epicFilterEl.setAttribute('hx-get', `/api/v1/projects/${PROJECT}/epics`)
  epicFilterEl.setAttribute('hx-trigger', 'epicsChanged from:body')
  window.htmx.process(epicFilterEl)
  window.htmx.trigger(epicFilterEl, 'epicsChanged')
  epicFilterEl.addEventListener('change', () => {
    const url = new URL(location)
    if (epicFilterEl.value) url.searchParams.set('epic', epicFilterEl.value)
    else url.searchParams.delete('epic')
    history.replaceState({}, '', url)
    syncBoardURL()
    refreshBoard()
  })
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

  if (event.target.closest('[data-issue-deleted]')) {
    closeIssuePanel()
    currentIssue = null
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
    const lane = addBtn.closest('.lane')
    openAddIssueDialog(Number(addBtn.dataset.laneId), lane?.dataset.laneName || '')
    return
  }

  if (event.target.closest('#addIssueBtn')) {
    const firstLane = document.querySelector('.lane[data-lane-id]')
    if (firstLane) {
      openAddIssueDialog(Number(firstLane.dataset.laneId), firstLane.dataset.laneName)
    }
    return
  }

  if (event.target.closest('[data-open-epics-manager]')) {
    window.openDialog('tmpl-epics-manager', { epics: currentEpics, projectSlug: PROJECT })
    return
  }

  const epicDetailBtn = event.target.closest('[data-open-epic-detail]')
  if (epicDetailBtn) {
    openEpicDetail(epicDetailBtn.getAttribute('data-epic-ref'))
    return
  }

  const epicEditBtn = event.target.closest('[data-open-epic-edit]')
  if (epicEditBtn) {
    window.openDialog('tmpl-epic-edit', {
      ref: epicEditBtn.getAttribute('data-epic-ref'),
      name: epicEditBtn.getAttribute('data-epic-name'),
      description: epicEditBtn.getAttribute('data-epic-description') || '',
      projectSlug: PROJECT,
    })
    return
  }

  if (event.target.closest('[data-lane-popover]')) return
  closeAllLanePopovers()
})

document.body.addEventListener('epicsChanged', () => {
  fetch(`/api/v1/projects/${PROJECT}/epics`).then(r => r.json()).then(epics => {
    currentEpics = epics
    const body = document.getElementById('sharedModalBody')
    if (body && body.querySelector('.epics-manager')) {
      body.innerHTML = renderTemplate('tmpl-epics-manager', { epics, projectSlug: PROJECT })
      window.htmx.process(body)
    }
  })
})

function openEpicDetail(ref) {
  const base = `/api/v1/projects/${PROJECT}`
  fetch(`${base}/epics/${ref}`)
    .then(r => r.json())
    .then(epic =>
      fetch(`${base}/issues?epic_id=${epic.id}`)
        .then(r => r.json())
        .then(issues => window.openDialog('tmpl-epic-detail', { ...epic, issues, projectSlug: PROJECT })),
    )
}

function openAddIssueDialog(laneId, laneName) {
  window.openDialog('tmpl-add-issue-form', {
    laneId,
    laneName,
    projectSlug: PROJECT,
    epics: currentEpics,
  })
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

function openIssuePanel(ref) {
  const panel = document.getElementById('issuePanel')
  if (!panel) return
  window.htmx.ajax('GET', `/api/v1/projects/${PROJECT}/issues/${ref}`, {
    source: panel,
    target: panel,
    swap: 'innerHTML',
  })
}

window.addEventListener('popstate', () => {
  const ref = new URLSearchParams(location.search).get('issue')
  if (ref) openIssuePanel(ref)
  else closeIssuePanel()
})

const initialIssue = new URLSearchParams(location.search).get('issue')
if (initialIssue) openIssuePanel(initialIssue)

document.body.addEventListener('htmx:afterRequest', event => {
  if (!event.detail.successful) return

  if (event.detail.target?.id === 'issuePanel') {
    try { currentIssue = JSON.parse(event.detail.xhr.responseText) } catch {}
  }

  if (event.detail.target?.id === 'board') {
    try {
      const view = JSON.parse(event.detail.xhr.responseText)
      currentEpics = view.epics || []
    } catch {}
  }

  const verb = event.detail.requestConfig?.verb
  if (!verb || verb === 'get') return
  refreshBoard()

  const path = event.detail.requestConfig?.path || ''
  if (path.includes('/epics')) window.htmx.trigger(document.body, 'epicsChanged')
})
