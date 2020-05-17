import React from "react";
import thunk from 'redux-thunk'
import * as ReactDOM from "react-dom";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom"
import { createStore, applyMiddleware } from "redux";
import { Provider } from "react-redux";
import { U, Mode, DetailKind } from "../components/u"
import { colors } from "../components/colorpicker"
import { UArgs, NProgress } from "../shared/types";
import reducer from "../redux/reducers/u"
import { Loupe } from 'loupe-js'

declare const uargs: UArgs;
declare const NProgress: NProgress

const store = createStore(reducer, { uargs }, applyMiddleware(thunk))
const loupe = new Loupe({
  magnification: 1.5,
  style: { boxShadow: "4px 5px 5px 4px rgba(0,0,0,0.5)" },
})

// configure NProgress globally
NProgress.configure({ showSpinner: false, minimum: 0.1, trickleSpeed: 150, speed: 500 })

ReactDOM.render(
  <Provider store={store}>
    <Router>
      <Switch>
        <Route exact path="/u/:username" render={p => <U {...uargs} wnd={window} mode={Mode.All} nProgress={NProgress} loupe={loupe} {...p} />} />
        <Route exact path="/u/:username/song/:hexSongIdent" render={p => <U {...uargs} wnd={window} mode={Mode.All} detail={{ kind: DetailKind.Song, hexIdent: p.match.params["hexSongIdent"] }} nProgress={NProgress} loupe={loupe} {...p} />} />
        <Route exact path="/u/:username/album/:hexSongIdent" render={p => <U {...uargs} wnd={window} mode={Mode.All} detail={{ kind: DetailKind.Album, hexIdent: p.match.params["hexSongIdent"] }} nProgress={NProgress} loupe={loupe} {...p} />} />

        <Route exact path="/u/:username/all" render={p => <U {...uargs} wnd={window} mode={Mode.All} nProgress={NProgress} loupe={loupe} {...p} />} />
        <Route exact path="/u/:username/all/song/:hexSongIdent" render={p => <U {...uargs} wnd={window} mode={Mode.All} detail={{ kind: DetailKind.Song, hexIdent: p.match.params["hexSongIdent"] }} nProgress={NProgress} loupe={loupe} {...p} />} />
        <Route exact path="/u/:username/all/album/:hexSongIdent" render={p => <U {...uargs} wnd={window} mode={Mode.All} detail={{ kind: DetailKind.Album, hexIdent: p.match.params["hexSongIdent"] }} nProgress={NProgress} loupe={loupe} {...p} />} />

        <Route exact path="/u/:username/loved" render={p => <U {...uargs} wnd={window} mode={Mode.Loved} nProgress={NProgress} loupe={loupe} {...p} />} />
        <Route exact path="/u/:username/loved/song/:hexSongIdent" render={p => <U {...uargs} wnd={window} mode={Mode.Loved} detail={{ kind: DetailKind.Song, hexIdent: p.match.params["hexSongIdent"] }} nProgress={NProgress} loupe={loupe} {...p} />} />
        <Route exact path="/u/:username/loved/album/:hexSongIdent" render={p => <U {...uargs} wnd={window} mode={Mode.Loved} detail={{ kind: DetailKind.Album, hexIdent: p.match.params["hexSongIdent"] }} nProgress={NProgress} loupe={loupe} {...p} />} />

        <Route exact path="/u/:username/color" render={p => <U {...uargs} wnd={window} mode={Mode.Color} nProgress={NProgress} loupe={loupe} {...p} />} />

        {colors.map(c => <Route key={c} exact path={`/u/:username/color/${c}`} render={p => <U {...uargs} wnd={window} mode={Mode.Color} color={c} nProgress={NProgress} loupe={loupe} {...p} />} />)}
        {colors.map(c => <Route key={c + "a"} exact path={`/u/:username/color/${c}/album/:hexSongIdent`} render={p => <U {...uargs} wnd={window} mode={Mode.Color} color={c} detail={{ kind: DetailKind.Album, hexIdent: p.match.params["hexSongIdent"] }} nProgress={NProgress} loupe={loupe} {...p} />} />)}
      </Switch>
    </Router>
  </Provider>,
  document.querySelector("#app")
);
