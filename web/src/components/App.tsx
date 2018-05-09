import * as React from "react";
import { Playback } from "./Playback";

export interface AppProps {
  compiler: string;
  framework: string;
}

export class App extends React.Component<AppProps, {}> {
	componentDidMount() {
		fetch("http://localhost:8080/api/scrobbles?username=nishanths")
      .then(r => {
        console.log(r);
        return r.json();
      })
      .then(r => this.setState(r));
	}

  render() {
    console.log(this.state);
    return <p>
      Hello from {this.props.compiler} and {this.props.framework}!
    </p>;
  }
}
