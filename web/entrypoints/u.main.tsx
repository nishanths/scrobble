import React from "react";
import thunk from 'redux-thunk'
import * as ReactDOM from "react-dom";
import { Router } from "react-router"
import { createStore, applyMiddleware } from "redux";
import { Provider } from "react-redux";
import { U } from "../components/Ux"
import { UArgs } from "../shared/types";
import reducer from "../redux/reducers/u"

declare var uargs: UArgs;

const store = createStore(reducer, { uargs }, applyMiddleware(thunk));

ReactDOM.render(
  <Provider store={store}>
    <Router history={}>
      <U {...uargs} wnd={window} />
    </Router>
  </Provider>,
  document.querySelector("#app")
);
