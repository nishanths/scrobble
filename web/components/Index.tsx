import * as React from "react";
import { BootstrapArgs } from "../src/shared";
import { SetUsername, SetUsernameProps } from "./SetUsername";
import { Account, AccountProps } from "./Account";

type IndexProps = BootstrapArgs

export class Index extends React.Component<{p: IndexProps}, {username?: string}> {
  constructor(props: {p: IndexProps}) {
    super(props)
    this.state = {
      username: this.props.p.account.username
    }
  }

  updateUsername(u: string) {
    this.setState({username: u})
  }

  private signIn() {
    if (this.props.p.email && this.props.p.logoutURL) {
      return <p><a href={this.props.p.logoutURL}>Sign out</a> ({this.props.p.email})</p>
    }
    if (this.props.p.loginURL) {
      return <p><a href={this.props.p.loginURL}>Sign in with Google</a> to get started</p>
    }
    return null;
  }

  private visit() {
    return <p>To see a user's scrobbled songs, visit {this.props.p.host}/u/&lt;username&gt;</p>
  }

  private profile() {
    return this.state.username && <p><a href={"https://" + this.props.p.host + "/u/" + this.state.username}>Your profile</a></p>
  }

  private download() {
    return <p><a href={this.props.p.downloadURL} title="v1.0, last updated Aug 28 2018">Download</a> menu bar client (macOS 10.13+)</p>
  }

  render() {
    return [
      <h1>{this.props.p.host}</h1>,
      this.props.p.email && !this.state.username &&
        <SetUsername host={this.props.p.host} usernameChange={this.updateUsername.bind(this)}/>,
      this.props.p.email && this.state.username &&
        <Account username={this.state.username} host={this.props.p.host} usernameChange={this.updateUsername.bind(this)}/>,
      this.signIn(),
      this.profile(),
      this.download(),
      this.visit(),
    ]
  }
}

