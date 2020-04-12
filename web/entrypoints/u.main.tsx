import React from "react";
import thunk from 'redux-thunk'
import * as ReactDOM from "react-dom";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom"
import { createStore, applyMiddleware } from "redux";
import { Provider } from "react-redux";
import { U, Mode } from "../components/U"
import { UArgs } from "../shared/types";
import reducer from "../redux/reducers/u"

declare var uargs: UArgs;

const store = createStore(reducer, { uargs }, applyMiddleware(thunk));

ReactDOM.render(
  <Provider store={store}>
    <Router>
      <Switch>
        <Route exact path="/u/:username" render={p => <U {...uargs} wnd={window} mode={Mode.All} {...p} />} />
        <Route exact path="/u/:username/all" render={p => <U {...uargs} wnd={window} mode={Mode.All} {...p} />} />
        <Route exact path="/u/:username/loved" render={p => <U {...uargs} wnd={window} mode={Mode.Loved} {...p} />} />
      </Switch>
    </Router>
  </Provider>,
  document.querySelector("#app")
);
