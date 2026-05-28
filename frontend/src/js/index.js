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

document.addEventListener('DOMContentLoaded', loadProjects)
