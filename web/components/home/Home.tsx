import React from "react"
import { GoogleSignIn } from "../GoogleSignIn"
import { randInt } from "../../shared/util"
import { macOSAppLink, guideLink, apiDocLink } from "../../shared/const"

interface HomeProps {
    loginURL: string
}

const artworks = [
    { img: "sonne", caption: "Artwork: Odyssey / Sonne – Rival Consoles" },
    { img: "paspatou", caption: "Artwork: Paspatou – Parra for Cuva" },
    { img: "onclejazz", caption: "Artwork: Oncle Jazz – Men I Trust" },
]

const taglines = [
    "an Apple Music scrobbling service.",
    "showcase your Apple Music library like never before.",
    "a beautiful way to scrobble your Apple Music songs.",
]

export const Home: React.FC<HomeProps> = ({ loginURL }) => {
    const artwork = artworks[randInt(0, artworks.length)]
    const tagline = taglines[randInt(0, taglines.length)]

    return <div className="Home">
        <div className="start">
            <div className="imgContainer">
                <img src={`/static/img/home/${artwork.img}.jpeg`} />
                <div className="caption">{artwork.caption}</div>
            </div>
            <div className="line"><b>scrobble</b> — {tagline}</div>
        </div>
        <div className="sign-in">
            <GoogleSignIn loginURL={loginURL} />
        </div>
        <div className="footer">
            <div>
                <a href={guideLink}><div className="item">How do I use Scrobble?</div></a>
                <a href={"https://scrobble.littleroot.org/u/nishanth"}><div className="item">Example profile</div></a>
                <a href={apiDocLink}><div className="item">API doc</div></a>
                <a href={macOSAppLink}><div className="item">macOS app</div></a>
            </div>
        </div>
    </div>
}
