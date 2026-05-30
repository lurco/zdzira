function autoResize(el) {
  el.style.height = 'inherit'
  el.style.height = `min(${el.scrollHeight}px, 60vh)`
}

document.body.addEventListener('input', event => {
  if (event.target.tagName === 'TEXTAREA') autoResize(event.target)
})

export function initTextareas(container) {
  container.querySelectorAll('textarea').forEach(autoResize)
}
