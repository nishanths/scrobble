import React from 'react';
import { storiesOf } from '@storybook/react';
import { SegmentedControl } from './SegmentedControl';

const s = storiesOf("SegmentedControl", module)

s.add("default", () => {
	const values = ["All", "Loved", "By color"] as const
	type V = typeof values[number]

	const props = {
		values,
		initialValue: "All" as const,
		afterChange: (v: V) => {
			console.log(v)
		},
	}

	return <SegmentedControl {...props} />
})
