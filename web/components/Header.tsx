import * as React from "react"

export class Header extends React.Component<{username: string, logoutURL: string}, {}> {
  render() {
    return <div className="header">
      <span className="username"><span className="emph">{this.props.username}</span>'s scrobbles</span>
      {this.props.logoutURL && <span className="logout"><a href={this.props.logoutURL}>Sign out</a></span>}
    </div>
  }
}
