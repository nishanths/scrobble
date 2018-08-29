import * as React from "react";
import { BootstrapArgs, Account } from "../src/shared";
import { SetUsername, SetUsernameProps } from "./SetUsername";
import { AccountDetail, AccountDetailProps } from "./AccountDetail";
import "../scss/index.scss";

type IndexProps = BootstrapArgs

export class Index extends React.Component<IndexProps, {account: Account}> {
  private static readonly downloadURL = "https://github.com/nishanths/scrobble/releases/latest"

  constructor(props: IndexProps) {
    super(props)
    this.state = {
      account: this.props.account
    }
  }

  updateAccount(a: Account) {
    this.setState({account: a})
  }

  private signIn() {
    if (this.props.email && this.props.logoutURL) {
      return <p><a href={this.props.logoutURL}>Sign out</a> ({this.props.email})</p>
    }
    if (this.props.loginURL) {
      return <p><a href={this.props.loginURL}>Sign in with Google</a> to get started</p>
    }
    return null;
  }

  private visit() {
    return <p>To see a user's scrobbled songs, go to <a href="">https://{this.props.host}/u/&lt;username&gt;</a></p>
  }

  private profile() {
    return this.state.account && <p><a href={"https://" + this.props.host + "/u/" + this.state.account.username}>Your scrobbles</a></p>
  }

  private download() {
    return <p><a href={Index.downloadURL}>Download</a> menu bar client for iTunes (macOS 10.13+)</p>
  }

  render() {
    return <div>
      <h1>{this.props.host}</h1>
      {this.props.email && !this.state.account.username &&
        <SetUsername host={this.props.host} accountChange={this.updateAccount.bind(this)}/>}
      {this.props.email && this.state.account.username &&
        <AccountDetail account={this.state.account} host={this.props.host}/>}
      {this.signIn()}
      {this.state.account.username ? this.profile() : this.visit()}
      {this.download()}
    </div>
  }
}

