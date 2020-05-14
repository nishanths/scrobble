import React from "react"

type BaseProps = {}

export class Base extends React.Component<BaseProps> {
  constructor(props: BaseProps) {
    super(props)
  }

  render() {
    return <div className="Base"></div>
  }
}
