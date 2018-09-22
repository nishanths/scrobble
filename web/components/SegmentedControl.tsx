import * as React from "react";

interface SegmentedControlProps {
  afterChange: (idx: number) => void
}

// TODO: this only supports two item controls, and it does not
// support configuring labels, etc.
export class SegmentedControl extends React.Component<SegmentedControlProps, {selectedIdx: number}> {
  constructor(props: SegmentedControlProps) {
    super(props)
    this.state = {
      selectedIdx: 0
    }
  }

  private static className(selected: boolean, idx: number): string {
    return selected ? `c c${idx} selected` : `c c${idx}`
  }

  private onControlClicked(idx: number) {
    let current = this.state.selectedIdx
    this.setState({selectedIdx: idx})
    if (idx != current) {
      this.props.afterChange(idx)
    }
  }

  render() {
    let c0 = this.state.selectedIdx == 0

    return <div className="SegmentedControl">
      <div className={SegmentedControl.className(c0, 0)} onClick={ () => this.onControlClicked(0) }>All</div>
      <div className={SegmentedControl.className(!c0, 1)} onClick={ () => this.onControlClicked(1) }>Loved</div>
    </div>
  }
}
