import * as React from "react";

export interface SetUsernameProps {
  usernameChange: (u: string) => void
  host: string
}

export class SetUsername extends React.Component<SetUsernameProps, {username: string}> {
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
