import React, { useState } from "react";
import "../scss/segmented-control.scss";

interface SegmentedControlProps {
  afterChange: (value: string) => void
  initialValue: string
  values: string[]
}

// TODO: this only properly supports two item controls in terms of styling.
export const SegmentedControl: React.FC<SegmentedControlProps> = ({
  afterChange,
  initialValue,
  values,
}) => {
  const [selected, setSelected] = useState(initialValue);
  const className = (selected: boolean, idx: number) => selected ? `c c${idx} selected` : `c c${idx}`
  const onControlClick = (v: string) => {
    const old = selected
    setSelected(v)
    if (v != old) {
      afterChange(v)
    }
  }

  return <div className="SegmentedControl">
    {values.map((v, i) => {
      return <Item key={v} className={className(selected === v, i)} onClick={() => onControlClick(v)} content={v} />
    })}
  </div>
}

const Item: React.StatelessComponent<{ content: string, className: string, onClick: () => void }> = ({
  content,
  className,
  onClick,
}) => {
  return <div className={className} onClick={() => { onClick() }}>{content}</div>
}
