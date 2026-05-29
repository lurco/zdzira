import './htmx-config'
import './dialog'
import './handlebars-helpers'
import '../styles/main.sass'
import './mode'
import './topbar'
import './board-dnd'
import './board-filter'
import { renderTemplate } from './dialog'
import { PROJECT, refreshBoard } from './project'
import Handlebars from 'handlebars'

// SSE: refresh the board when any agent or other tab mutates the API.
// Debounced to avoid double-render when a local htmx mutation also fires.
let sseRefreshTimer = null
const es = new EventSource('/api/v1/events')
es.onmessage = () => {
  clearTimeout(sseRefreshTimer)
  sseRefreshTimer = setTimeout(refreshBoard, 300)
}

const boardEl = document.getElementById('board')
let currentIssue = null
let currentEpics = []
let currentLanes = []

Handlebars.registerHelper('boardLanes', () => currentLanes)

function boardPath() {
  const epic = new URLSearchParams(location.search).get('epic') || ''
  return epic
    ? `/api/v1/projects/${PROJECT}/board?epic=${encodeURIComponent(epic)}`
    : `/api/v1/projects/${PROJECT}/board`
}

// htmx captures hx-get once at process time, so it can't follow the epic
// filter. Drive the fetch ourselves on every boardUpdated, reading the
// current URL each time.
function loadBoard() {
  if (!boardEl) return
  window.htmx.ajax('GET', boardPath(), { source: boardEl, target: boardEl, swap: 'innerHTML' })
}

if (boardEl) {
  boardEl.setAttribute('hx-ext', 'client-side-templates')
  boardEl.setAttribute('handlebars-template', 'tmpl-board')
  window.htmx.process(boardEl)
  document.body.addEventListener('boardUpdated', loadBoard)
  loadBoard()
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
    panel.innerHTML = renderTemplate('tmpl-issue-edit-form', { ...currentIssue, epics: currentEpics, projectSlug: PROJECT })
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

function updateStatusBar(view) {
  const lanes = view.swimlanes || []
  const total = lanes.reduce((sum, lane) => sum + (lane.issues ? lane.issues.length : 0), 0)
  const totalEl = document.getElementById('totalCount')
  const laneEl = document.getElementById('laneCount')
  if (totalEl) totalEl.textContent = total
  if (laneEl) laneEl.textContent = lanes.length
}

const projectInfoEl = document.getElementById('projectInfo')
if (projectInfoEl) projectInfoEl.textContent = `proj: ${PROJECT}`

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

function wireLaneSelect(issueRef) {
  const select = document.getElementById('issueLaneSelect')
  if (!select) return
  select.addEventListener('change', () => {
    fetch(`/api/v1/projects/${PROJECT}/issues/${issueRef}/move`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ swimlane_id: Number(select.value), position: 0 }),
    }).then(() => refreshBoard())
  })
}

function loadComments(issueRef) {
  const listEl = document.getElementById('commentsList')
  if (!listEl) return

  fetch(`/api/v1/projects/${PROJECT}/issues/${issueRef}/comments`)
    .then(r => r.json())
    .then(comments => {
      listEl.innerHTML = renderTemplate('tmpl-comments', comments)
      listEl.querySelectorAll('[data-comment-id]').forEach(btn => {
        btn.addEventListener('click', () => {
          const id = btn.dataset.commentId
          fetch(`/api/v1/projects/${PROJECT}/issues/${issueRef}/comments/${id}`, { method: 'DELETE' })
            .then(() => loadComments(issueRef))
        })
      })
    })

  const form = document.getElementById('commentForm')
  if (!form || form.dataset.wired) return
  form.dataset.wired = '1'
  form.addEventListener('submit', e => {
    e.preventDefault()
    const textarea = form.querySelector('[name="contents"]')
    const text = textarea?.value.trim()
    if (!text) return
    fetch(`/api/v1/projects/${PROJECT}/issues/${issueRef}/comments`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ contents: text }),
    }).then(() => {
      textarea.value = ''
      loadComments(issueRef)
    })
  })
}

document.body.addEventListener('htmx:afterRequest', event => {
  if (!event.detail.successful) return

  if (event.detail.target?.id === 'issuePanel') {
    try { currentIssue = JSON.parse(event.detail.xhr.responseText) } catch {}
    if (currentIssue?.ref) {
      loadComments(currentIssue.ref)
      wireLaneSelect(currentIssue.ref)
    }
  }

  if (event.detail.target?.id === 'board') {
    try {
      const view = JSON.parse(event.detail.xhr.responseText)
      currentEpics = view.epics || []
      currentLanes = (view.swimlanes || []).map(({ id, name }) => ({ id, name }))
      updateStatusBar(view)
    } catch {}
  }

  const verb = event.detail.requestConfig?.verb
  if (!verb || verb === 'get') return
  refreshBoard()

  const path = event.detail.requestConfig?.path || ''
  if (path.includes('/epics')) window.htmx.trigger(document.body, 'epicsChanged')
})
