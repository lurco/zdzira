let container = null

function getContainer() {
  if (!container) {
    container = document.createElement('div')
    container.className = 'toast-container'
    document.body.appendChild(container)
  }
  return container
}

export function showToast(message, type = 'error') {
  const el = document.createElement('div')
  el.className = `toast toast--${type}`
  el.textContent = message

  const close = document.createElement('button')
  close.className = 'toast__close'
  close.textContent = '×'
  close.setAttribute('aria-label', 'Dismiss')
  close.addEventListener('click', () => dismiss(el))
  el.appendChild(close)

  getContainer().appendChild(el)
  setTimeout(() => dismiss(el), 5000)
}

function dismiss(el) {
  el.classList.add('toast--out')
  el.addEventListener('animationend', () => el.remove(), { once: true })
}

window.showToast = showToast
