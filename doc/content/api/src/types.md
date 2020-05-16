[Back to index](/)

# Types

* [Song](./#type-song)
* [SongResponse](./#type-songresponse)
* [SongIdentifier](./#type-songidentifier)
* [ArtworkHash](./#type-artworkhash)
* [MediaItem](./#type-mediaitem)

## Type: Song

A `Song` represents a scrobbled song.

```
type Song = {
  albumTitle:   string
  artistName:   string
  title:        string
  totalTime:    number // length of song in nanoseconds
  year:         number
  releaseDate:  number // timestamp in unix seconds
  lastPlayed:   number // timestamp in unix seconds
  playCount:    number
  artworkHash:  ArtworkHash
  loved:        boolean
  previewURL:   string
  trackViewURL: string
}
```

All properties will be present in a resonse. Values that are not available
are represented by the Go zero value representation for the type. For example,
if the `year` property is unavailable for a song, the `year` key will be present
but will have value `0`.

## Type: SongResponse

```
type SongResponse = Song & { ident: SongIdentifier }
```

## Type: SongIdentifer

A SongIdentifier represents the id for a song.

```
type SongIdentifier = string
```

The format is

```
sprintf("%s|%s|%s|%s", base64Encode(albumTitle), base64Encode(artistName), base64Encode(title), base64Encode(itoa(year)))
```

where `albumTitle`, `artistName`, etc. are properties of the `Song` type, and
`base64Encode()` represents the base64-encoding behavior of [`base64.StdEncoding`](https://golang.org/pkg/encoding/base64/#pkg-variables) in the Go standard library.

As an example, the `SongIdentifier` for a `Song` with the following properties
```
albumTitle: "Trick of the Light - Single"
artistName: "La Mar"
title: "Trick of the Light"
year: 2015
```
is
```
VHJpY2sgb2YgdGhlIExpZ2h0IC0gU2luZ2xl|TGEgTWFy|VHJpY2sgb2YgdGhlIExpZ2h0|MjAxNQ==
```

## Type: ArtworkHash

```
type ArtworkHash = string
```

The format is

```
decimalEncode(sha1(<artworkImageData> + '|' + <artworkFormat>))
```
An example Go implementation to compute the artwork hash given the artwork image
data and the artwork image format is

```go
import (
    "bytes"
    "crypto/sha1"
    "fmt"
)

func artworkHash(artwork []byte, format string) string {
    h := sha1.New()
    h.Write(artwork)
    h.Write([]byte("|"))
    h.Write([]byte(format))
    sum := h.Sum(nil)

    var buf bytes.Buffer
    for _, b := range sum {
        buf.WriteString(fmt.Sprintf("%d", b))
    }
    return buf.String()
}
```

Valid artwork format values are
```
type ArtworkFormat =
  | "GIF"
  | "JPEG"
  | "JPEG2000"
  | "PNG"
```

## Type: MediaItem

A `MediaItem` represents a song in requests.

The properties correspond to those in the [`ITLibMediaItem`](https://developer.apple.com/documentation/ituneslibrary/itlibmediaitem)
class in the `iTunesLibrary` framework.

```
type MediaItem struct {
  added:          number // timestamp in unix seconds
  albumTitle:     string
  sortAlbumTitle: string
  artistName:     string
  sortArtistName: string
  genre:          string
  hasArtwork:     boolean
  kind:           string
  lastPlayed:     number // timestamp in unix seconds
  playCount:      number
  releaseDate:    number // timestamp in unix seconds
  sortTitle:      string
  title:          string
  totalTime:      number // milliseconds
  year:           number
  persistentID:   string
  artworkHash:    string
  loved:          boolean
}
```
