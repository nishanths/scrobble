import React from 'react';
import { storiesOf } from '@storybook/react';
import { Swatch } from './Swatch';

const s = storiesOf("Swatch", module)

s.add("default", () => {
    return <Swatch selected={false} color="orange" />
})

s.add("selected", () => {
    return <Swatch selected={true} color="violet" />
})
