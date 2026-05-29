import '../styles/main.sass'
import './mode'
import './topbar'

// ── config ────────────────────────────────────────────────────────────────────
const params  = new URLSearchParams(location.search)
const PROJECT = params.get('project') || 'main'
const API     = `/api/v1/projects/${PROJECT}`

const PRIORITY_CLASS = { IMMEDIATE: 'p0', HIGH: 'p1', LOW: 'p3' }
const LANE_COLORS = ['#2A6FDB','#F5D547','#7A4FD6','#E63946','#2A9D5A','#6B655A','#0A0908']

// ── state ─────────────────────────────────────────────────────────────────────
let swimlanes = []   // [{ id, name, color, position }]
let issues    = []   // [{ ref, swimlane_id, name, type, priority }]
let epics     = []   // [{ id, ref, name }]

// ── DOM refs ──────────────────────────────────────────────────────────────────
const boardEl      = document.getElementById('board')
const totalCountEl = document.getElementById('totalCount')
const laneCountEl  = document.getElementById('laneCount')

// ── filters ───────────────────────────────────────────────────────────────────
const filters = { q: '', priority: new Set(), type: new Set() }

function matchesFilter(issue) {
  if (filters.q) {
    const hay = `${issue.ref} ${issue.name}`.toLowerCase()
    if (!hay.includes(filters.q.toLowerCase())) return false
  }
  if (filters.priority.size && !filters.priority.has(issue.priority)) return false
  if (filters.type.size && !filters.type.has(issue.type)) return false
  return true
}

// ── API helpers ───────────────────────────────────────────────────────────────
async function api(path, opts = {}) {
  const res = await fetch(`/api/v1${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...opts,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  })
  if (!res.ok) throw new Error(`${opts.method || 'GET'} ${path} → ${res.status}`)
  if (res.status === 204) return null
  return res.json()
}

async function loadBoard() {
  const [sw, iss, ep] = await Promise.all([
    api(`/projects/${PROJECT}/swimlanes`),
    api(`/projects/${PROJECT}/issues`),
    api(`/projects/${PROJECT}/epics`),
  ])
  swimlanes = sw ?? []
  issues    = iss ?? []
  epics     = ep ?? []
}

// ── render ────────────────────────────────────────────────────────────────────
function render() {
  boardEl.innerHTML = ''

  swimlanes.forEach(lane => {
    const laneIssues   = issues.filter(i => i.swimlane_id === lane.id)
    const visibleIssues = laneIssues.filter(matchesFilter)

    const el = document.createElement('div')
    el.className = 'lane'
    el.dataset.laneId = lane.id

    // header
    const hdr = document.createElement('div')
    hdr.className = 'lane-header'
    hdr.draggable = true
    hdr.dataset.dragKind = 'lane'
    hdr.dataset.laneId = lane.id
    hdr.innerHTML = `
      <span class="lane-swatch" style="background:${escHtml(lane.color || '#6B655A')}"></span>
      <input class="lane-title" value="${escHtml(lane.name)}" maxlength="40" />
      <span class="lane-count">${visibleIssues.length}</span>
      <button class="lane-menu" aria-label="Lane menu">⋯</button>
    `
    el.appendChild(hdr)

    // body
    const body = document.createElement('div')
    body.className = 'lane-body'
    body.dataset.laneId = lane.id
    if (visibleIssues.length === 0) {
      const empty = document.createElement('div')
      empty.className = 'lane-empty'
      empty.textContent = laneIssues.length === 0 ? 'No issues' : 'No matches'
      body.appendChild(empty)
    } else {
      visibleIssues.forEach(issue => body.appendChild(renderCard(issue, lane)))
    }
    el.appendChild(body)

    // footer
    const footer = document.createElement('div')
    footer.className = 'lane-footer'
    const addBtn = document.createElement('button')
    addBtn.className = 'add-card-btn'
    addBtn.textContent = '+ Add issue'
    addBtn.addEventListener('click', () => beginAddCard(lane.id, body, footer))
    footer.appendChild(addBtn)
    el.appendChild(footer)

    boardEl.appendChild(el)
  })

  // add-lane tile
  const addLaneBtn = document.createElement('button')
  addLaneBtn.className = 'add-lane'
  addLaneBtn.innerHTML = '<span class="plus">+</span> Add lane'
  addLaneBtn.addEventListener('click', addLane)
  boardEl.appendChild(addLaneBtn)

  totalCountEl.textContent = issues.length
  laneCountEl.textContent  = swimlanes.length

  wireHandlers()
}

function renderCard(issue, lane) {
  const el = document.createElement('div')
  el.className = 'card'
  el.draggable = true
  el.dataset.dragKind = 'card'
  el.dataset.cardRef  = issue.ref
  el.dataset.laneId   = lane.id

  const strip = document.createElement('div')
  strip.className = 'card-strip'
  strip.style.background = lane.color || '#2A6FDB'
  el.appendChild(strip)

  const top = document.createElement('div')
  top.className = 'card-top'
  top.innerHTML = `
    <span class="card-id mono">${escHtml(issue.ref)}</span>
    <span class="badge ${PRIORITY_CLASS[issue.priority] || 'p3'}">${escHtml(issue.priority)}</span>
  `
  el.appendChild(top)

  const title = document.createElement('div')
  title.className = 'card-title'
  title.textContent = issue.name
  el.appendChild(title)

  const meta = document.createElement('div')
  meta.className = 'card-meta'
  const typeTag = document.createElement('span')
  typeTag.className = `tag tag-${(issue.type || '').toLowerCase()}`
  typeTag.textContent = issue.type || ''
  meta.appendChild(typeTag)
  el.appendChild(meta)

  el.addEventListener('click', e => {
    if (e.shiftKey) deleteIssue(issue.ref, el)
    else location.href = `/issue.html?project=${PROJECT}&ref=${issue.ref}`
  })

  return el
}

function escHtml(s) {
  return String(s).replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))
}

// ── inline add-card form ──────────────────────────────────────────────────────
function beginAddCard(laneId, bodyEl, footerEl) {
  document.querySelectorAll('.new-card-form').forEach(n => n.remove())
  document.querySelectorAll('.add-card-btn').forEach(b => b.style.display = '')

  const form = document.createElement('div')
  form.className = 'new-card-form'
  const epicOptions = epics.length
    ? epics.map(e => `<option value="${escHtml(e.ref)}">${escHtml(e.ref)} ${escHtml(e.name)}</option>`).join('')
    : ''

  form.innerHTML = `
    <textarea class="new-card-form__textarea" placeholder="What needs doing?"></textarea>
    <div class="new-card-form__row">
      <select class="mini-select" data-k="type">
        <option value="TASK" selected>Task</option>
        <option value="BUG">Bug</option>
        <option value="STORY">Story</option>
      </select>
      <select class="mini-select" data-k="priority">
        <option value="HIGH" selected>High</option>
        <option value="IMMEDIATE">Immediate</option>
        <option value="LOW">Low</option>
      </select>
      ${epicOptions ? `<select class="mini-select" data-k="epic"><option value="">No epic</option>${epicOptions}</select>` : ''}
      <span style="flex:1"></span>
      <button class="btn primary" data-act="save" style="padding:6px 10px;font-size:11px">Add</button>
      <button class="btn ghost"   data-act="cancel" style="padding:6px 10px;font-size:11px">Esc</button>
    </div>
  `

  const addBtn = footerEl.querySelector('.add-card-btn')
  if (addBtn) addBtn.style.display = 'none'
  bodyEl.appendChild(form)
  const ta = form.querySelector('.new-card-form__textarea')
  ta.focus()

  async function commit() {
    const name = ta.value.trim()
    if (!name) { cancel(); return }
    const type     = form.querySelector('[data-k="type"]').value
    const priority = form.querySelector('[data-k="priority"]').value
    const epicRef  = form.querySelector('[data-k="epic"]')?.value || ''
    try {
      const body = { name, type, priority, swimlane_id: laneId }
      if (epicRef) body.epic_ref = epicRef
      const issue = await api(`/projects/${PROJECT}/issues`, {
        method: 'POST',
        body,
      })
      issues.push(issue)
      render()
    } catch (e) {
      alert('Failed to create issue: ' + e.message)
      cancel()
    }
  }

  function cancel() {
    form.remove()
    if (addBtn) addBtn.style.display = ''
  }

  form.querySelector('[data-act="save"]').addEventListener('click', commit)
  form.querySelector('[data-act="cancel"]').addEventListener('click', cancel)
  ta.addEventListener('keydown', e => {
    if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) { e.preventDefault(); commit() }
    if (e.key === 'Escape') { e.preventDefault(); cancel() }
  })
}

async function deleteIssue(ref, el) {
  if (!confirm(`Delete ${ref}?`)) return
  try {
    await api(`/projects/${PROJECT}/issues/${ref}`, { method: 'DELETE' })
    issues = issues.filter(i => i.ref !== ref)
    render()
  } catch (e) {
    alert('Delete failed: ' + e.message)
  }
}

// ── add lane ──────────────────────────────────────────────────────────────────
async function addLane() {
  const color = LANE_COLORS[swimlanes.length % LANE_COLORS.length]
  try {
    const sw = await api(`/projects/${PROJECT}/swimlanes`, {
      method: 'POST',
      body: { name: 'New lane' },
    })
    swimlanes.push(sw)
    render()
    requestAnimationFrame(() => {
      const lanes = document.querySelectorAll('.lane')
      const last  = lanes[lanes.length - 1]
      if (last) last.querySelector('.lane-title')?.focus()
      document.getElementById('boardWrap')?.scrollTo({ left: 1e6, behavior: 'smooth' })
    })
  } catch (e) {
    alert('Failed to add lane: ' + e.message)
  }
}

// ── lane menu ─────────────────────────────────────────────────────────────────
function openLaneMenu(laneEl, laneId) {
  closePopovers()
  const lane = swimlanes.find(l => l.id === laneId)
  if (!lane) return

  const pop = document.createElement('div')
  pop.className = 'popover'
  pop.innerHTML = `
    <div class="popover__swatches">
      ${LANE_COLORS.map(c => `<button class="popover__swatch" data-c="${c}" style="background:${c}" aria-label="${c}"></button>`).join('')}
    </div>
    <button class="popover__btn" data-act="move-left">← Move left</button>
    <button class="popover__btn" data-act="move-right">→ Move right</button>
    <button class="popover__btn popover__danger" data-act="delete">Delete lane</button>
  `
  laneEl.appendChild(pop)

  pop.querySelectorAll('.popover__swatch').forEach(b => {
    b.addEventListener('click', async () => {
      lane.color = b.dataset.c
      render()
      closePopovers()
    })
  })
  pop.querySelector('[data-act="move-left"]').addEventListener('click', () => {
    const i = swimlanes.findIndex(l => l.id === laneId)
    if (i > 0) { [swimlanes[i-1], swimlanes[i]] = [swimlanes[i], swimlanes[i-1]]; render() }
    closePopovers()
  })
  pop.querySelector('[data-act="move-right"]').addEventListener('click', () => {
    const i = swimlanes.findIndex(l => l.id === laneId)
    if (i < swimlanes.length - 1) { [swimlanes[i+1], swimlanes[i]] = [swimlanes[i], swimlanes[i+1]]; render() }
    closePopovers()
  })
  pop.querySelector('[data-act="delete"]').addEventListener('click', async () => {
    const count = issues.filter(i => i.swimlane_id === laneId).length
    const msg = count > 0
      ? `Delete "${lane.name}" and its ${count} issue${count === 1 ? '' : 's'}?`
      : `Delete "${lane.name}"?`
    if (!confirm(msg)) return
    try {
      await api(`/projects/${PROJECT}/swimlanes/${laneId}`, { method: 'DELETE' })
      swimlanes = swimlanes.filter(l => l.id !== laneId)
      issues    = issues.filter(i => i.swimlane_id !== laneId)
      render()
    } catch (e) {
      alert('Delete failed: ' + e.message)
    }
    closePopovers()
  })
}

function closePopovers() {
  document.querySelectorAll('.popover').forEach(p => p.remove())
}
document.addEventListener('click', e => {
  if (!e.target.closest('.popover') && !e.target.closest('.lane-menu')) closePopovers()
})

// ── wire handlers (called after each render) ──────────────────────────────────
function wireHandlers() {
  document.querySelectorAll('.lane-title').forEach(input => {
    const laneId = Number(input.closest('.lane').dataset.laneId)
    input.addEventListener('mousedown', e => e.stopPropagation())
    input.addEventListener('change', async () => {
      const lane = swimlanes.find(l => l.id === laneId)
      if (!lane) return
      const name = input.value.trim() || 'Untitled'
      lane.name = name
    })
    input.addEventListener('keydown', e => { if (e.key === 'Enter') { e.preventDefault(); input.blur() } })
  })

  document.querySelectorAll('.lane-menu').forEach(btn => {
    btn.addEventListener('click', e => {
      e.stopPropagation()
      const laneEl = btn.closest('.lane')
      openLaneMenu(laneEl, Number(laneEl.dataset.laneId))
    })
  })
}

// ── epics panel ───────────────────────────────────────────────────────────────
const epicPanel     = document.getElementById('epicPanel')
const epicList      = document.getElementById('epicList')
const epicToggleBtn = document.getElementById('epicToggleBtn')
const addEpicBtn    = document.getElementById('addEpicBtn')

function renderEpicPanel() {
  if (!epicList) return
  if (epics.length === 0) {
    epicList.innerHTML = '<span class="epic-panel__empty">No epics yet</span>'
    return
  }
  epicList.innerHTML = ''
  epics.forEach(epic => {
    const count = issues.filter(i => i.epic_id === epic.id).length
    const card = document.createElement('div')
    card.className = 'epic-card'
    card.title = epic.name
    card.innerHTML = `
      <span class="epic-card__ref">${escHtml(epic.ref)}</span>
      <span class="epic-card__name">${escHtml(epic.name)}</span>
      <span class="epic-card__count">${count}</span>
    `
    epicList.appendChild(card)
  })
}

epicToggleBtn?.addEventListener('click', () => {
  const open = epicToggleBtn.getAttribute('aria-pressed') === 'true'
  epicToggleBtn.setAttribute('aria-pressed', String(!open))
  epicPanel.hidden = open
  if (!open) renderEpicPanel()
})

addEpicBtn?.addEventListener('click', () => {
  const name = prompt('Epic name:')
  if (!name) return
  api(`/projects/${PROJECT}/epics`, { method: 'POST', body: { name } })
    .then(epic => { epics.push(epic); renderEpicPanel() })
    .catch(e => alert('Failed to create epic: ' + e.message))
})

// ── toolbar filters ───────────────────────────────────────────────────────────
document.getElementById('searchInput')?.addEventListener('input', e => {
  filters.q = e.target.value
  render()
})

document.querySelectorAll('.filter-chip').forEach(chip => {
  chip.addEventListener('click', () => {
    const filter = chip.dataset.filter
    const val    = chip.dataset.val
    const pressed = chip.getAttribute('aria-pressed') === 'true'
    const set = filters[filter]
    if (!set) return
    if (pressed) { set.delete(val); chip.setAttribute('aria-pressed', 'false') }
    else         { set.add(val);    chip.setAttribute('aria-pressed', 'true') }
    render()
  })
})

document.getElementById('addIssueBtn')?.addEventListener('click', () => {
  if (swimlanes.length === 0) { addLane(); return }
  const lane   = swimlanes[0]
  const laneEl = document.querySelector(`.lane[data-lane-id="${lane.id}"]`)
  if (!laneEl) return
  beginAddCard(lane.id, laneEl.querySelector('.lane-body'), laneEl.querySelector('.lane-footer'))
})

// ── drag & drop ───────────────────────────────────────────────────────────────
let drag = null

document.addEventListener('dragstart', e => {
  const cardEl = e.target.closest('[data-drag-kind="card"]')
  const laneHdr = e.target.closest('[data-drag-kind="lane"]')

  if (cardEl) {
    drag = { kind: 'card', ref: cardEl.dataset.cardRef, fromLaneId: Number(cardEl.dataset.laneId) }
    cardEl.classList.add('dragging')
    e.dataTransfer.effectAllowed = 'move'
  } else if (laneHdr) {
    drag = { kind: 'lane', id: Number(laneHdr.dataset.laneId) }
    laneHdr.closest('.lane').classList.add('dragging')
    e.dataTransfer.effectAllowed = 'move'
  }
})

document.addEventListener('dragend', () => {
  document.querySelectorAll('.dragging').forEach(n => n.classList.remove('dragging'))
  document.querySelectorAll('.drop-slot').forEach(n => n.remove())
  document.querySelectorAll('.drop-target-lane').forEach(n => n.classList.remove('drop-target-lane'))
  drag = null
})

document.addEventListener('dragover', e => {
  if (!drag) return

  if (drag.kind === 'card') {
    const body = e.target.closest('.lane-body')
    if (!body) return
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
    const cards = [...body.querySelectorAll('.card:not(.dragging)')]
    let before = null
    for (const c of cards) {
      if (e.clientY < c.getBoundingClientRect().top + c.getBoundingClientRect().height / 2) { before = c; break }
    }
    body.querySelectorAll('.drop-slot').forEach(s => s.remove())
    const slot = document.createElement('div')
    slot.className = 'drop-slot active'
    if (before) body.insertBefore(slot, before)
    else body.appendChild(slot)
    document.querySelectorAll('.lane-body').forEach(b => {
      if (b !== body) b.querySelectorAll('.drop-slot').forEach(s => s.remove())
    })
  }

  if (drag.kind === 'lane') {
    const laneEl = e.target.closest('.lane')
    if (!laneEl) return
    e.preventDefault()
    document.querySelectorAll('.drop-target-lane').forEach(n => n.classList.remove('drop-target-lane'))
    if (Number(laneEl.dataset.laneId) !== drag.id) laneEl.classList.add('drop-target-lane')
  }
})

document.addEventListener('drop', async e => {
  if (!drag) return

  if (drag.kind === 'card') {
    const body = e.target.closest('.lane-body')
    if (!body) return
    e.preventDefault()
    const toLaneId = Number(body.dataset.laneId)
    const issue    = issues.find(i => i.ref === drag.ref)
    if (!issue) return
    try {
      await api(`/projects/${PROJECT}/issues/${issue.ref}/move`, {
        method: 'POST',
        body: { swimlane_id: toLaneId },
      })
      issue.swimlane_id = toLaneId
      render()
    } catch (err) {
      alert('Move failed: ' + err.message)
    }
  }

  if (drag.kind === 'lane') {
    const laneEl = e.target.closest('.lane')
    if (!laneEl) return
    e.preventDefault()
    const targetId = Number(laneEl.dataset.laneId)
    if (targetId === drag.id) return
    const fromIdx = swimlanes.findIndex(l => l.id === drag.id)
    const toIdx   = swimlanes.findIndex(l => l.id === targetId)
    if (fromIdx < 0 || toIdx < 0) return
    const r = laneEl.getBoundingClientRect()
    const insertAfter = e.clientX > r.left + r.width / 2
    const [moved] = swimlanes.splice(fromIdx, 1)
    let newIdx = swimlanes.findIndex(l => l.id === targetId)
    if (insertAfter) newIdx += 1
    swimlanes.splice(newIdx, 0, moved)
    render()
  }
})

// ── boot ──────────────────────────────────────────────────────────────────────
loadBoard().then(render).catch(err => {
  boardEl.innerHTML = `<p style="color:var(--acc-red);padding:24px">Failed to load board: ${escHtml(err.message)}</p>`
})
