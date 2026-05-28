const currentProject = new URLSearchParams(location.search).get('project')

export async function initTopbar() {
  const el = document.getElementById('projectSwitcher')
  if (!el) return

  try {
    const res = await fetch('/api/v1/projects')
    const projects = await res.json()

    if (projects.length === 0) {
      el.innerHTML = '<option value="">No projects</option>'
      return
    }

    el.innerHTML = projects
      .map(p => `<option value="${p.slug}"${p.slug === currentProject ? ' selected' : ''}>${p.name}</option>`)
      .join('')
  } catch {
    el.innerHTML = '<option value="">Error loading projects</option>'
  }

  el.addEventListener('change', () => {
    if (el.value) location.href = `/board.html?project=${el.value}`
  })
}

document.addEventListener('DOMContentLoaded', initTopbar)
