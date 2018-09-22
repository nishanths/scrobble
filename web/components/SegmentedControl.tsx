import * as React from "react";

export class SegmentedControl extends React.Component<{}, {}> {
  constructor(props: {}) {
    super(props)
  }

  render() {
    return <div className="SegmentedControl">
      <div className="c c0 sel">All</div>
      <div className="c c1">Loved</div>
    </div>
  }
}
