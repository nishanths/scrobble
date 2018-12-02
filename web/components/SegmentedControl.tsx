import * as React from "react";

interface SegmentedControlProps {
  afterChange: (value: string) => void
  initial: string
  values: string[]
}

// TODO: this only properly supports two item controls in terms of styling.
export class SegmentedControl extends React.Component<SegmentedControlProps, {selected: string}> {
  constructor(props: SegmentedControlProps) {
    super(props)
    this.state = {
      selected: this.props.initial
    }
  }

  private static className(selected: boolean, idx: number): string {
    return selected ? `c c${idx} selected` : `c c${idx}`
  }

  private onControlClicked(v: string) {
    let current = this.state.selected
    this.setState({selected: v})
    if (v != current) {
      this.props.afterChange(v)
    }
  }

  private items(): JSX.Element[] {
    return this.props.values.map((v, i) => {
      let selected = this.state.selected == v;
      return <div key={v} className={SegmentedControl.className(selected, i)} onClick={ () => this.onControlClicked(v) }>{v}</div>
    })
  }

  render() {
    return <div className="SegmentedControl">
      {this.items()}
    </div>
  }
}
