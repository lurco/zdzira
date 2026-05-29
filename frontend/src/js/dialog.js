import Handlebars from 'handlebars'

const MODAL_ID = 'sharedModal'
const BODY_ID = 'sharedModalBody'

const templateCache = new Map()

function getCompiled(templateId) {
  if (templateCache.has(templateId)) return templateCache.get(templateId)
  const el = document.getElementById(templateId)
  if (!el) throw new Error(`template not found: ${templateId}`)
  const compiled = Handlebars.compile(el.innerHTML)
  templateCache.set(templateId, compiled)
  return compiled
}

export function renderTemplate(templateId, data = {}) {
  return getCompiled(templateId)(data)
}

function getModal() {
  const modal = document.getElementById(MODAL_ID)
  if (!modal) throw new Error(`shared dialog not found in DOM: #${MODAL_ID}`)
  return modal
}

export function openDialog(templateId, data = {}) {
  const modal = getModal()
  const body = document.getElementById(BODY_ID)
  body.innerHTML = getCompiled(templateId)(data)
  window.htmx.process(body)
  modal.showModal()
}

export function closeDialog() {
  getModal().close()
}

document.addEventListener('click', event => {
  const trigger = event.target.closest('[data-dialog-open]')
  if (trigger) {
    event.preventDefault()
    const dataAttr = trigger.getAttribute('data-dialog-data')
    const data = dataAttr ? JSON.parse(dataAttr) : {}
    openDialog(trigger.getAttribute('data-dialog-open'), data)
    return
  }

  const closer = event.target.closest('[data-dialog-close]')
  if (closer) {
    event.preventDefault()
    closeDialog()
  }
})

document.body.addEventListener('htmx:afterRequest', event => {
  if (!event.detail.successful) return
  const modal = document.getElementById(MODAL_ID)
  if (!modal || !modal.open) return
  if (!modal.contains(event.detail.elt)) return
  closeDialog()
})

window.openDialog = openDialog
window.closeDialog = closeDialog
