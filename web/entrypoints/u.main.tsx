import React from "react";
import thunk from 'redux-thunk'
import * as ReactDOM from "react-dom";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom"
import { createStore, applyMiddleware } from "redux";
import { Provider } from "react-redux";
import { U, Mode } from "../components/U"
import { colors } from "../components/colorpicker"
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
        <Route exact path="/u/:username/color" render={p => <U {...uargs} wnd={window} mode={Mode.Color} {...p} />} />
        {
          colors.map(c => {
            return <Route exact path={`/u/:username/color/${c}`} render={p => <U {...uargs} wnd={window} mode={Mode.Color} color={c} {...p} />} />
          })
        }
      </Switch>
    </Router>
  </Provider>,
  document.querySelector("#app")
);
