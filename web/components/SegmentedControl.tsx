import React, { useState, useEffect, ReactNode } from "react";
import "../scss/segmented-control.scss";
import classNames from "classnames"

interface SegmentedControlProps<V extends readonly string[]> {
	afterChange: (value: V[number]) => void
	initialValue: V[number]
	values: V
	// add a new badge to these
	// TODO: make this more generic in the future if the need arises
	newBadges: Set<V[number]>
}

// TODO: support for defining React.FC with generic props?

export const SegmentedControl = <V extends readonly string[]>(props: (SegmentedControlProps<V> & { children?: ReactNode })) => {
	const { afterChange, initialValue, values, newBadges } = props;

	const [selected, setSelected] = useState(initialValue);
	useEffect(() => {
		setSelected(initialValue)
	}, [initialValue])

	const onControlClick = (v: string) => {
		const old = selected
		if (v != old) {
			setSelected(v)
			afterChange(v)
		}
	}

	return <div className="SegmentedControl">
		{values.map((v) => {
			return <Item key={v} className={classNames("c", {
				"selected": selected === v,
			})} onClick={() => onControlClick(v)} content={v} newBadge={newBadges.has(v)} />
		})}
	</div>
}

const Item: React.StatelessComponent<{ content: string, className: string, onClick: () => void, newBadge: boolean }> = ({
	content,
	className,
	onClick,
	newBadge,
}) => {
	return <div className={classNames(className, { "has-badge": newBadge })} onClick={() => { onClick() }}>
		{content}
		{newBadge && <span className="new-badge">new</span>}
	</div>
}
