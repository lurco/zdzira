import { PROJECT_API, refreshBoard } from './project'

const SLOT_CLASS = 'drop-slot'
const DRAGGING_CLASS = 'dragging'
const LANE_DROP_CLASS = 'drop-target-lane'

let dragRef = null // card being dragged
let dragLaneId = null // lane being dragged

function clearDropSlots() {
  document.querySelectorAll(`.${SLOT_CLASS}`).forEach(slot => slot.remove())
}

function clearLaneDropTargets() {
  document.querySelectorAll(`.${LANE_DROP_CLASS}`).forEach(lane => lane.classList.remove(LANE_DROP_CLASS))
}

function pickInsertBefore(laneBody, clientY) {
  const cards = [...laneBody.querySelectorAll(`.card:not(.${DRAGGING_CLASS})`)]
  for (const card of cards) {
    const rect = card.getBoundingClientRect()
    if (clientY < rect.top + rect.height / 2) return card
  }
  return null
}

function currentLaneOrder() {
  return [...document.querySelectorAll('.lane[data-lane-id]')].map(lane => Number(lane.dataset.laneId))
}

function reorderLanes(ids) {
  fetch(`${PROJECT_API}/swimlanes/reorder`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ids }),
  }).then(refreshBoard)
}

document.addEventListener('dragstart', event => {
  const laneGrip = event.target.closest('[data-drag-kind="lane"]')
  if (laneGrip) {
    dragLaneId = Number(laneGrip.dataset.laneId)
    const lane = laneGrip.closest('.lane')
    if (lane) {
      lane.classList.add(DRAGGING_CLASS)
      event.dataTransfer.setDragImage(lane, 0, 0)
    }
    event.dataTransfer.effectAllowed = 'move'
    return
  }

  const card = event.target.closest('.card[data-card-ref]')
  if (!card) return
  dragRef = card.dataset.cardRef
  card.classList.add(DRAGGING_CLASS)
  event.dataTransfer.effectAllowed = 'move'
})

document.addEventListener('dragover', event => {
  if (dragLaneId !== null) {
    const lane = event.target.closest('.lane[data-lane-id]')
    if (!lane) return
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
    clearLaneDropTargets()
    if (Number(lane.dataset.laneId) !== dragLaneId) lane.classList.add(LANE_DROP_CLASS)
    return
  }

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
  if (dragLaneId !== null) {
    const targetLane = event.target.closest('.lane[data-lane-id]')
    if (!targetLane) return
    event.preventDefault()

    const movingId = dragLaneId
    const targetId = Number(targetLane.dataset.laneId)
    if (movingId === targetId) return

    const rect = targetLane.getBoundingClientRect()
    const dropAfter = event.clientX > rect.left + rect.width / 2
    const ids = currentLaneOrder().filter(id => id !== movingId)
    const targetIndex = ids.indexOf(targetId)
    ids.splice(dropAfter ? targetIndex + 1 : targetIndex, 0, movingId)
    reorderLanes(ids)
    return
  }

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
  clearLaneDropTargets()
  dragRef = null
  dragLaneId = null
})
