import React from 'react';
import { storiesOf } from '@storybook/react';
import { ColorPicker } from './ColorPicker';

const s = storiesOf("ColorPicker", module)

s.add("default", () => {
	return <ColorPicker prompt="Pick a color." />
})
