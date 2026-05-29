// Client-side filtering of board cards by priority, type, and free text.
// Chips toggle on click; within a group active values OR together, and the
// groups (plus the search box) AND together. Re-applied whenever the board
// re-renders, since htmx replaces #board's contents on every update.

const activeByGroup = { priority: new Set(), type: new Set() }
let searchText = ''

function toggleChip(chip) {
  const group = activeByGroup[chip.dataset.filter]
  if (!group) return
  const wasPressed = chip.getAttribute('aria-pressed') === 'true'
  chip.setAttribute('aria-pressed', String(!wasPressed))
  if (wasPressed) group.delete(chip.dataset.val)
  else group.add(chip.dataset.val)
  applyFilters()
}

function cardMatches(card) {
  const { priority, type } = activeByGroup
  if (priority.size && !priority.has(card.dataset.priority)) return false
  if (type.size && !type.has(card.dataset.type)) return false
  if (searchText) {
    const title = (card.querySelector('.card-title')?.textContent || '').toLowerCase()
    const ref = (card.dataset.cardRef || '').toLowerCase()
    if (!title.includes(searchText) && !ref.includes(searchText)) return false
  }
  return true
}

function applyFilters() {
  document.querySelectorAll('.lane').forEach(lane => {
    let visible = 0
    lane.querySelectorAll('.card').forEach(card => {
      const show = cardMatches(card)
      card.hidden = !show
      if (show) visible += 1
    })
    const count = lane.querySelector('.lane-count')
    if (count) count.textContent = visible
  })
}

document.addEventListener('click', event => {
  const chip = event.target.closest('.filter-chip')
  if (chip) toggleChip(chip)
})

document.addEventListener('input', event => {
  if (event.target.id !== 'searchInput') return
  searchText = event.target.value.trim().toLowerCase()
  applyFilters()
})

// The board's innerHTML is swapped on every update, so re-apply afterwards.
document.body.addEventListener('htmx:afterSettle', event => {
  if (event.detail.target?.id === 'board') applyFilters()
})
