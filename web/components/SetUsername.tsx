import * as React from "react";
import { Account } from "../src/shared"
import "../scss/username.scss";

export interface SetUsernameProps {
  accountChange: (a: Account) => void
  host: string
}

export class SetUsername extends React.Component<SetUsernameProps, {username: string, error: string}> {
  private input: HTMLInputElement|null = null;
  private static usernameRe = /^[a-z0-9]*$/

  constructor(props: SetUsernameProps) {
    super(props)
    this.state = {
      username: "",
      error: "",
    }
  }

  componentDidMount() {
    this.input!.focus();
  }

  private setUsername(u: string) {
     let {reason, ok} = SetUsername.validate(u)
     if (!ok) {
       this.setState({error: reason})
       return
     }

    this.setState({error: ""})
    const genericError = "Something went wrong? Try again"
    let success = false;

    fetch(`https://${this.props.host}/setUsername?username=${u}`, {method: "POST"})
      .then(res => {
        if (res.status == 200) {
          success = true
          return res.json()
        }
        if (res.status == 406) {
          this.setState({error: "That username is already taken."})
          return res.text()
        }
        this.setState({error: genericError})
        return res.blob()
      })
      .then(r => {
        if (!success) { return }
        this.props.accountChange(r as Account)
      }, err => {
        console.error(err)
        this.setState({error: genericError})
      })
  }

  private static validate(u: string): {reason: string, ok: boolean} {
    if (u.length <= 1) {
      return {reason: "Username must be longer than 2 characters", ok: false}
    }
    if (u.length >= 31) {
      return {reason: "Username must be shorter than 32 characters", ok: false}
    }
    if (!SetUsername.usernameRe.test(u)) {
      return {reason: "Username must contain only a-z and 0-9", ok: false}
    }
    if (u.indexOf("scrobble") != -1) {
      return {reason: "Username cannot contain 'scrobble'", ok: false};
    }
    if (u == "username") {
      return {reason: "Username cannot be 'username'", ok: false};
    }
    return {reason: "", ok: true}
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
          Set your username:
          <input type="text" value={this.state.username} onChange={this.onChange.bind(this)}
            ref={r => { this.input = r }}></input>
        </label>
        <input type="submit" value="OK" />
        {this.state.error && <span className="error">{this.state.error}</span>}
      </form>
    </div>
  }
}
