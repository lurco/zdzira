import { PROJECT_API, refreshBoard } from './project'

const SLOT_CLASS = 'drop-slot'
const DRAGGING_CLASS = 'dragging'

let dragRef = null

function clearDropSlots() {
  document.querySelectorAll(`.${SLOT_CLASS}`).forEach(slot => slot.remove())
}

function pickInsertBefore(laneBody, clientY) {
  const cards = [...laneBody.querySelectorAll(`.card:not(.${DRAGGING_CLASS})`)]
  for (const card of cards) {
    const rect = card.getBoundingClientRect()
    if (clientY < rect.top + rect.height / 2) return card
  }
  return null
}

document.addEventListener('dragstart', event => {
  const card = event.target.closest('.card[data-card-ref]')
  if (!card) return
  dragRef = card.dataset.cardRef
  card.classList.add(DRAGGING_CLASS)
  event.dataTransfer.effectAllowed = 'move'
})

document.addEventListener('dragover', event => {
  if (!dragRef) return
  const laneBody = event.target.closest('.lane-body')
  if (!laneBody) return

  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'

  clearDropSlots()
  const slot = document.createElement('div')
  slot.className = `${SLOT_CLASS} active`
  const before = pickInsertBefore(laneBody, event.clientY)
  if (before) laneBody.insertBefore(slot, before)
  else laneBody.appendChild(slot)
})

document.addEventListener('drop', event => {
  if (!dragRef) return
  const laneBody = event.target.closest('.lane-body')
  if (!laneBody) return

  event.preventDefault()
  const toLaneId = Number(laneBody.dataset.laneId)
  const movingRef = dragRef

  window.htmx.ajax('POST', `${PROJECT_API}/issues/${movingRef}/move`, {
    values: { swimlane_id: toLaneId },
    swap: 'none',
  }).then(refreshBoard)
})

document.addEventListener('dragend', () => {
  document.querySelectorAll(`.${DRAGGING_CLASS}`).forEach(el => el.classList.remove(DRAGGING_CLASS))
  clearDropSlots()
  dragRef = null
})
