import React from "react"
import { GoogleSignIn } from "../GoogleSignIn"
import { randInt } from "../../shared/util"

interface HomeProps {
  loginURL: string
}

const artworks = [
  { img: "sonne", caption: "Artwork: Odyssey / Sonne – Rival Consoles" },
  { img: "paspatou", caption: "Artwork: Paspatou – Parra for Cuva" },
  { img: "onclejazz", caption: "Artwork: Oncle Jazz – Men I Trust" },
]

export const Home: React.FC<HomeProps> = ({ loginURL }) => {
  const artwork = artworks[randInt(0, artworks.length)]

  return <div className="Home">
    <div className="start">
      <img src={`/static/img/home/${artwork.img}.jpeg`} title={artwork.caption} />
      <div className="line"><b>scrobble</b>, an Apple Music scrobbling service.</div>
    </div>
    <div className="sign-in">
      <GoogleSignIn loginURL={loginURL} />
    </div>
    <div className="footer">
      <a href="/guide"><div className="item">Guide</div></a> ·
      <a href="https://github.com/nishanths/scrobble/releases/latest"><div className="item">macOS app</div></a> ·
      <a href=""><div className="item">API doc (upcoming)</div></a>
    </div>
  </div>
}
