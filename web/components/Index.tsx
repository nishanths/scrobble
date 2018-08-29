import * as React from "react";
import { BootstrapArgs } from "../src/shared";

interface UserDetailProps {
  username: string
  usernameChange: (u: string) => void
  host: string
}


class UserDetail extends React.Component<UserDetailProps, {apiKey: string, error: any, fetched: boolean}> {
  constructor(props: UserDetailProps) {
    super(props)
    this.state = {
      apiKey: "",
      error: null,
      fetched: false
    }
  }

  componentDidMount() {
    this.fetchInitial()
  }

  private fetchInitial() {
    if (!this.props.username) {
      return;
    }
    fetch(`https://${this.props.host}/account`)
      .then(res => res.json())
      .then(
        r => {
          this.setState({apiKey: r.apiKey})
        },
        error => {
          this.setState({error})
        }
      )
  }

  private setUsername() {
    fetch(`https://${this.props.host}/setUsername`, {method: "POST"})
      .then(res => res.json())
      .then(
        response => {
        },
        error => {
        }
      )
  }

  private generateAPIKey() {
    fetch(`https://${this.props.host}/generateAPIKey`, {method: "POST"})
      .then(res => res.json())
      .then(
        response => {
          this.setState({apiKey: response})
        },
        error => {
        }
      )
  }

  render() {
    return <div>
      <div><span>Username</span><span>{this.props.username}</span></div>
    </div>
  }
}


interface SetUsernameProps {
  usernameChange: (u: string) => void
  host: string
}

class SetUsername extends React.Component<SetUsernameProps, {username: string}> {
  private input: HTMLInputElement|null = null;

  constructor(props: SetUsernameProps) {
    super(props)
    this.state = {
      username: ""
    }
  }

  componentDidMount() {
    this.input!.focus();
  }

  private setUsername(u: string) {
    fetch(`https://${this.props.host}/setUsername`, {method: "POST"})
      .then(res => {
        if (res.ok) {
          this.props.usernameChange(u)
          return
        }
        return res.json()
      })
      .then(
        r => {
          console.log(r)
        },
        error => {
          console.log(error)
        }
      )
  }

  private onChange() {
    this.setState({username: this.input!.value})
  }

  private handleSubmit(e: any) {
    e.preventDefault()
    this.setUsername(this.state.username)
  }

  render() {
    return <div>
      <form onSubmit={this.handleSubmit.bind(this)}>
        <label>
          Set your username
          <input type="text" value={this.state.username} onChange={this.onChange.bind(this)}
            ref={r => { this.input = r }}></input>
        </label>
        <input type="submit" value="OK" />
      </form>
    </div>
  }
}

type IndexProps = BootstrapArgs

export class Index extends React.Component<{p: IndexProps}, {username?: string}> {
  constructor(props: {p: IndexProps}) {
    super(props)
    this.state = {
      username: this.props.p.account && this.props.p.account.username
    }
  }

  updateUsername(u: string) {
    this.setState({username: u})
  }

  private signIn() {
    if (this.props.p.email) {
      return <p><a href={this.props.p.logoutURL}>Sign out</a> ({this.props.p.email})</p>
    }
    return <p>To get started, <a href={this.props.p.loginURL}>sign in with Google</a></p>
  }

  private visit() {
    return <p>To see a user's scrobbled songs, visit {this.props.p.host}/u/<i>username</i></p>
  }

  private profile() {
    return this.state.username && <p><a href={"https://" + this.props.p.host + "/u/" + this.state.username}>Your profile</a></p>
  }

  private download() {
    return <p><a href={this.props.p.downloadURL}>Download</a> menu bar client for macOS (v1.0, updated Aug 28 2018)</p>
  }

  render() {
    return [
      <h1>{this.props.p.host}</h1>,
      this.props.p.email && !this.state.username &&
        <SetUsername host={this.props.p.host} usernameChange={this.updateUsername.bind(this)}/>,
      this.props.p.email && this.state.username &&
        <UserDetail username="foo" host={this.props.p.host} usernameChange={this.updateUsername.bind(this)}/>,
      this.signIn(),
      this.profile(),
      this.download(),
      this.visit(),
    ]
  }
}

