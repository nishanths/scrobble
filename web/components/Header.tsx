import * as React from "react"

export class Header extends React.Component<{username: string, logoutURL: string}, {}> {
  render() {
    return <div>
      <span className="username">{this.props.username}'s scrobbles</span>
      {this.props.logoutURL && <span className="logout"></span>}
    </div>
  }
}
