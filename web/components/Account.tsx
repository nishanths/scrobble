import * as React from "react";

export interface AccountProps {
  username: string
  usernameChange: (u: string) => void
  host: string
}


export class Account extends React.Component<AccountProps, {apiKey: string, error: any, fetched: boolean}> {
  constructor(props: AccountProps) {
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
