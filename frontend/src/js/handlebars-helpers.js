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
