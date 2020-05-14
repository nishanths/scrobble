import React from "react"
import * as ReactDOM from "react-dom"
import { BrowserRouter as Router, Switch, Route, Redirect } from "react-router-dom"
import { Dashboard, Mode } from "../components/dashboard"
import { BootstrapArgs, Notie } from "../shared/types"
import "../scss/dashboard/dashboard.scss"

declare const bootstrap: BootstrapArgs
declare const notie: Notie

ReactDOM.render(
  <Router>
    <Switch>
      <Route exact path="/" render={(p) => <Dashboard history={p.history} mode={Mode.Base} {...bootstrap} notie={notie} />} />
      <Route exact path="/dashboard"><Redirect to="/" /></Route>
      <Route exact path="/dashboard/privacy" render={(p) => <Dashboard history={p.history} mode={Mode.Privacy} {...bootstrap} notie={notie} />} />
      <Route exact path="/dashboard/api-key" render={(p) => <Dashboard history={p.history} mode={Mode.APIKey} {...bootstrap} notie={notie} />} />
    </Switch>
  </Router>,
  document.querySelector("#app")
);
