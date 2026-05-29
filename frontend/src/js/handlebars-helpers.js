import Handlebars from 'handlebars'

Handlebars.registerHelper('selectedIfCurrentProject', slug => {
  const params = new URLSearchParams(location.search)
  return slug === params.get('project') ? 'selected' : ''
})

Handlebars.registerHelper('eq', (a, b) => a === b)
