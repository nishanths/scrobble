import React, { useState, useEffect, ReactNode } from "react";
import "../scss/segmented-control.scss";

interface SegmentedControlProps<V extends readonly string[]> {
  afterChange: (value: V[number]) => void
  initialValue: V[number]
  values: V
}

// TODO: support for React.FC with generic props?

export const SegmentedControl = <V extends readonly string[]>(props: (SegmentedControlProps<V> & { children?: ReactNode })) => {
  const { afterChange, initialValue, values } = props;

  const [selected, setSelected] = useState(initialValue);
  useEffect(() => {
    setSelected(initialValue)
  }, [initialValue])

  const className = (selected: boolean) => selected ? `c selected` : `c`
  const onControlClick = (v: string) => {
    const old = selected
    if (v != old) {
      setSelected(v)
      afterChange(v)
    }
  }

  return <div className="SegmentedControl">
    {values.map((v, i) => {
      return <Item key={v} className={className(selected === v)} onClick={() => onControlClick(v)} content={v} />
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
