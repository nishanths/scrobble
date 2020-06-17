import React, { useState } from 'react';
import { storiesOf } from '@storybook/react';
import { SearchBox } from './SearchBox';

const s = storiesOf("SearchBox", module)

s.add("default", () => {
	const [value, setValue] = useState("")

	const props = {
		value,
		placeholder: "Filter by album, artist, or song title...",
		onChange: (v: string) => {
			console.log(v)
			setValue(v)
		},
	}

	return <SearchBox {...props} />
})
