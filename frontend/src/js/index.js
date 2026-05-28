import 'htmx.org'
import '../styles/main.sass'
import './mode'
import './topbar'

function esc(s) {
  return String(s).replace(/[&<>"']/g, c => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]))
}

async function loadProjects() {
  const grid = document.getElementById('projectGrid')
  if (!grid) return

  try {
    const res = await fetch('/api/v1/projects')
    const projects = await res.json()

    if (projects.length === 0) {
      grid.innerHTML = '<p class="project-list__empty">No projects yet.</p>'
      return
    }

    grid.innerHTML = projects.map(p => `
      <a class="project-card" href="/board.html?project=${esc(p.slug)}">
        <span class="project-card__shortcut">${esc(p.shortcut)}</span>
        <span class="project-card__name">${esc(p.name)}</span>
      </a>
    `).join('')
  } catch (e) {
    grid.innerHTML = `<p class="project-list__empty">Failed to load: ${esc(e.message)}</p>`
  }
}

async function createProject() {
  const name = prompt('Project name:')
  if (!name) return
  const shortcut = prompt('Issue prefix (e.g. PROJ):')
  if (!shortcut) return
  try {
    const res = await fetch('/api/v1/projects', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name, shortcut }),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.detail || `${res.status} ${res.statusText}`)
    }
    const p = await res.json()
    location.href = `/board.html?project=${esc(p.slug)}`
  } catch (e) {
    alert('Failed to create project: ' + e.message)
  }
}

document.addEventListener('DOMContentLoaded', () => {
  loadProjects()
  const btn = document.getElementById('newProjectBtn')
  if (btn) btn.addEventListener('click', createProject)
})
