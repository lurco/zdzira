import Handlebars from 'handlebars'

Handlebars.registerHelper('selectedIfCurrentProject', slug => {
  const params = new URLSearchParams(location.search)
  return slug === params.get('project') ? 'selected' : ''
})

Handlebars.registerHelper('eq', (a, b) => a === b)

Handlebars.registerHelper('laneColor', color => color || '#6B655A')

Handlebars.registerHelper('priorityClass', priority => {
  const map = { IMMEDIATE: 'p0', HIGH: 'p1', LOW: 'p3' }
  return map[priority] || 'p3'
})

Handlebars.registerHelper('lower', value => String(value || '').toLowerCase())

Handlebars.registerHelper('projectSlug', () => {
  return new URLSearchParams(location.search).get('project') || 'main'
})

const LANE_COLORS = [
  '#2A6FDB', '#F5D547', '#7A4FD6',
  '#E63946', '#2A9D5A', '#6B655A', '#0A0908',
]
Handlebars.registerHelper('laneColors', () => LANE_COLORS)

Handlebars.registerHelper('json', value => JSON.stringify(value))
