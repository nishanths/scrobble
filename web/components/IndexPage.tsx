import * as React from "react";
import { BootstrapArgs, Account } from "../src/shared";
import { SetUsername, SetUsernameProps } from "./SetUsername";
import { AccountDetail, AccountDetailProps } from "./AccountDetail";
import "../scss/index.scss";

type IndexPageProps = BootstrapArgs

export class IndexPage extends React.Component<IndexPageProps, {account: Account}> {
  private static readonly downloadURL = "https://github.com/nishanths/scrobble/releases/latest"

  constructor(props: IndexPageProps) {
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
    return <p>Profiles can be found at /u/&lt;username&gt;, e.g., <a href="/u/nishanth">/u/nishanth</a>, <a href="/u/888">/u/888</a></p>
  }

  private profile() {
    return this.state.account && <p><a href={"https://" + this.props.host + "/u/" + this.state.account.username}>Your scrobbles</a></p>
  }

  private download() {
    return <p><a href={IndexPage.downloadURL}>Download</a> menu bar client for iTunes (macOS 10.13+)</p>
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

