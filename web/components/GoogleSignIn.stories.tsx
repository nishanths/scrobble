import React from 'react';
import { storiesOf } from '@storybook/react';
import { GoogleSignIn } from './GoogleSignIn';

const s = storiesOf("GoogleSignIn", module)

s.add("default", () => {
  return <GoogleSignIn loginURL="/fakeloginurl" />
})
