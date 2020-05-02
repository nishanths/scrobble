import { configure } from '@storybook/react'

import '!style-loader!css-loader!sass-loader!./base.scss';

const req = require.context('../components', true, /.stories.tsx$/)
const loadStories = () => { req.keys().forEach(req) }
configure(loadStories, module)