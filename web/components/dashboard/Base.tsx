import React from "react"
import "../../scss/dashboard/base.scss"

type BaseProps = {
  username: string
}

export class Base extends React.Component<BaseProps> {
  constructor(props: BaseProps) {
    super(props)
  }

  render() {
    return <div className="Base">
      <div className="welcome">Welcome, {this.props.username}!</div>
    </div>
  }
}
