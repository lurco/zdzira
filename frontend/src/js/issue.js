import '../styles/main.sass'
import './mode'
import './topbar'

const params = new URLSearchParams(location.search)
const PROJECT = params.get('project') || ''
const REF = params.get('ref') || ''

function esc(s) {
  return String(s).replace(/[&<>"']/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]))
}

function badge(cls, text) {
  return `<span class="badge badge--${cls.toLowerCase()}">${esc(text)}</span>`
}

async function loadIssue() {
  const main = document.querySelector('.issue-detail')
  if (!main) return
  if (!PROJECT || !REF) {
    main.innerHTML = '<p class="error">Missing project or issue reference in URL.</p>'
    return
  }
  try {
    const res = await fetch(`/api/v1/projects/${esc(PROJECT)}/issues/${esc(REF)}`)
    if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
    const issue = await res.json()

    document.title = `Zdzira — ${esc(issue.name)}`

    main.innerHTML = `
      <header class="issue-header">
        <span class="issue-ref">${esc(REF.toUpperCase())}</span>
        ${badge('type-' + issue.type, issue.type)}
        ${badge('priority-' + issue.priority, issue.priority)}
      </header>
      <h1 class="issue-title">${esc(issue.name)}</h1>
      <div class="issue-body">
        ${issue.description ? `<p class="issue-description">${esc(issue.description)}</p>` : '<p class="issue-description issue-description--empty">No description.</p>'}
      </div>
    `
  } catch (e) {
    main.innerHTML = `<p class="error">Failed to load issue: ${esc(e.message)}</p>`
  }
}

document.addEventListener('DOMContentLoaded', loadIssue)
