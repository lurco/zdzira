const STORAGE_KEY = 'zdzira-mode'
const root = document.documentElement

function applyMode(mode) {
  root.dataset.mode = mode
  document.querySelectorAll('[data-set]').forEach(btn => {
    btn.setAttribute('aria-pressed', String(btn.dataset.set === mode))
  })
  localStorage.setItem(STORAGE_KEY, mode)
}

document.addEventListener('DOMContentLoaded', () => {
  const saved = localStorage.getItem(STORAGE_KEY) || 'light'
  applyMode(saved)

  document.querySelectorAll('[data-set]').forEach(btn => {
    btn.addEventListener('click', () => applyMode(btn.dataset.set))
  })
})
