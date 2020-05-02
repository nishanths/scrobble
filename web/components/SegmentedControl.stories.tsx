import React from 'react';
import { storiesOf } from '@storybook/react';
import { SegmentedControl } from './SegmentedControl';

const s = storiesOf("SegmentedControl", module)

s.add("with text", () => {
	const props = {
		values: ["All", "Loved", "By color"] as const,
		initialValue: "All" as const,
		afterChange: () => {},
	}
	return <SegmentedControl {...props}></SegmentedControl>
})
