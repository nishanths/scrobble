import { configure } from '@storybook/react'

const req = require.context('../components', true, /.stories.tsx$/)

const loadStories = () => {
  req.keys().forEach(req)
}

configure(loadStories, module)