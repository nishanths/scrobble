package server

import "net/http"

func devScrobblesHandler(w http.ResponseWriter, r *http.Request) {
	const b = `
[
  {
    "Song": {
      "Duration": 257800,
      "Genre": "Rock",
      "Name": "Dreams",
      "Artist": "Fleetwood Mac",
      "Album": "Rumours (Deluxe)",
      "Year": 1977,
      "Urlp": "651871544",
      "Urli": "651871679"
    },
    "StartTimes": [
      1525850554
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 264906,
      "Genre": "Rock",
      "Name": "Gypsy",
      "Artist": "Fleetwood Mac",
      "Album": "Greatest Hits",
      "Year": 1982,
      "Urlp": "202271826",
      "Urli": "202272422"
    },
    "StartTimes": [
      1525850378
    ],
    "Artwork": "https://is5-ssl.mzstatic.com/image/thumb/Music/fa/e8/30/mzi.izeorbmm.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 174933,
      "Genre": "Dance",
      "Name": "Solo Dance",
      "Artist": "Martin Jensen",
      "Album": "Ultra 2018",
      "Year": 2017,
      "Urlp": "1309963000",
      "Urli": "1309963896"
    },
    "StartTimes": [
      1525850282,
      1525850274
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music118/v4/f7/1b/6d/f71b6d28-cfa2-accb-84ee-fbe68ef913d2/0617465919053.png/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 264906,
      "Genre": "Rock",
      "Name": "Gypsy",
      "Artist": "Fleetwood Mac",
      "Album": "Greatest Hits",
      "Year": 1982,
      "Urlp": "202271826",
      "Urli": "202272422"
    },
    "StartTimes": [
      1525850235,
      1525850202
    ],
    "Artwork": "https://is5-ssl.mzstatic.com/image/thumb/Music/fa/e8/30/mzi.izeorbmm.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 174933,
      "Genre": "Dance",
      "Name": "Solo Dance",
      "Artist": "Martin Jensen",
      "Album": "Ultra 2018",
      "Year": 2017,
      "Urlp": "1309963000",
      "Urli": "1309963896"
    },
    "StartTimes": [
      1525850138
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music118/v4/f7/1b/6d/f71b6d28-cfa2-accb-84ee-fbe68ef913d2/0617465919053.png/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 257800,
      "Genre": "Rock",
      "Name": "Dreams",
      "Artist": "Fleetwood Mac",
      "Album": "Rumours (Deluxe)",
      "Year": 1977,
      "Urlp": "651871544",
      "Urli": "651871679"
    },
    "StartTimes": [
      1525850134,
      1525834082,
      1525833824,
      1525833566,
      1525833308,
      1525833050,
      1525832792,
      1525832535,
      1525832277,
      1525832019,
      1525831761,
      1525831502
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 289747,
      "Genre": "Alternative",
      "Name": "Shadow of the Day",
      "Artist": "LINKIN PARK",
      "Album": "Minutes to Midnight",
      "Year": 2007,
      "Urlp": "528975362",
      "Urli": "528975367"
    },
    "StartTimes": [
      1525823374
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Features/v4/9a/c2/68/9ac26832-eeb4-3854-0c47-33f3e5c11cc8/dj.xkwgleci.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 164457,
      "Genre": "Alternative",
      "Name": "Bleed It Out",
      "Artist": "LINKIN PARK",
      "Album": "Minutes to Midnight",
      "Year": 2007,
      "Urlp": "528975362",
      "Urli": "528975366"
    },
    "StartTimes": [
      1525823248
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Features/v4/9a/c2/68/9ac26832-eeb4-3854-0c47-33f3e5c11cc8/dj.xkwgleci.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 178197,
      "Genre": "Alternative",
      "Name": "Blood To Gold",
      "Artist": "slenderbodies",
      "Album": "Fabulist: Extended - EP",
      "Year": 2018,
      "Urlp": "1327615188",
      "Urli": "1327615251"
    },
    "StartTimes": [
      1525821721
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music118/v4/cd/7f/7a/cd7f7a41-deb4-6d55-e706-93e36807c10e/content_art_2F63vrfvh5QiCIINUkfUlZ_IMG_1030.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 243216,
      "Genre": "Alternative",
      "Name": "Heaven's Only Wishful",
      "Artist": "MorMor",
      "Album": "Heaven's Only Wishful - Single",
      "Year": 2018,
      "Urlp": "1365980325",
      "Urli": "1365980333"
    },
    "StartTimes": [
      1525821713
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music128/v4/a2/3c/a9/a23ca98c-5e9e-bc7a-1f15-2e24bac4e864/5054526784345_1.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 257800,
      "Genre": "Rock",
      "Name": "Dreams",
      "Artist": "Fleetwood Mac",
      "Album": "Rumours (Deluxe)",
      "Year": 1977,
      "Urlp": "651871544",
      "Urli": "651871679"
    },
    "StartTimes": [
      1525813132
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 227000,
      "Genre": "",
      "Name": "Party Police",
      "Artist": "Alvvays",
      "Album": "Alvvays",
      "Year": 2014,
      "Urlp": "877969472",
      "Urli": "877969479"
    },
    "StartTimes": [
      1525812957
    ],
    "Artwork": "https://is5-ssl.mzstatic.com/image/thumb/Music/v4/66/6c/b0/666cb02b-7ddb-d592-58d5-83c9e6e3f6c7/644110028297.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 257800,
      "Genre": "Rock",
      "Name": "Dreams",
      "Artist": "Fleetwood Mac",
      "Album": "Rumours (Deluxe)",
      "Year": 1977,
      "Urlp": "651871544",
      "Urli": "651871679"
    },
    "StartTimes": [
      1525808302,
      1525808045,
      1525807787,
      1525807529,
      1525807271,
      1525807013,
      1525806756,
      1525806498,
      1525806240,
      1525805982
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music2/v4/69/ef/ca/69efca19-5e7a-67a1-a9c4-a748fa8b3db6/603497925759.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 205593,
      "Genre": "Pop",
      "Name": "Alarm",
      "Artist": "Anne-Marie",
      "Album": "Speak Your Mind (Deluxe)",
      "Year": 2016,
      "Urlp": "1374052241",
      "Urli": "1374052255"
    },
    "StartTimes": [
      1525732507
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music128/v4/2c/67/ed/2c67ed61-ae1f-9fd3-d9b4-35805eefe6c6/190295633509.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 221900,
      "Genre": "Dance",
      "Name": "This Is What You Came For (feat. Rihanna)",
      "Artist": "Calvin Harris",
      "Album": "This Is What You Came For (feat. Rihanna) - Single",
      "Year": 2016,
      "Urlp": "1108212521",
      "Urli": "1108212668"
    },
    "StartTimes": [
      1525732285
    ],
    "Artwork": "https://is5-ssl.mzstatic.com/image/thumb/Music60/v4/36/28/5a/36285aca-a659-aab6-a1f8-8d4b4485aa98/886445857290.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 230461,
      "Genre": "R\u0026B/Soul",
      "Name": "Starboy (feat. Daft Punk)",
      "Artist": "The Weeknd",
      "Album": "Starboy",
      "Year": 2016,
      "Urlp": "1156172520",
      "Urli": "1156172683"
    },
    "StartTimes": [
      1525732055
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music62/v4/d1/8c/44/d18c44ef-a19b-1c33-47c2-0ae58c5acad3/UMG_cvrart_00602557212396_01_RGB72_1800x1800_16UMGIM67863.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 252933,
      "Genre": "Dance",
      "Name": "Something Just Like This (Alesso Remix)",
      "Artist": "The Chainsmokers \u0026 Coldplay",
      "Album": "Something Just Like This (Remix Pack) - EP",
      "Year": 2017,
      "Urlp": "1229287358",
      "Urli": "1229287791"
    },
    "StartTimes": [
      1525731802
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music91/v4/3a/31/ae/3a31ae26-2795-b43b-834a-4a013f12e96a/886446461090.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 204502,
      "Genre": "Pop",
      "Name": "Strip That Down (feat. Quavo)",
      "Artist": "Liam Payne",
      "Album": "Strip That Down (feat. Quavo) - Single",
      "Year": 2017,
      "Urlp": "1236272435",
      "Urli": "1236272436"
    },
    "StartTimes": [
      1525731597
    ],
    "Artwork": "https://is2-ssl.mzstatic.com/image/thumb/Music127/v4/e6/62/c6/e662c69e-593a-53a2-c614-9da9ed8df392/UMG_cvrart_00602557629576_01_RGB72_1800x1800_17UMGIM89908.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 292586,
      "Genre": "Dance",
      "Name": "Waiting All Night (feat. Ella Eyre)",
      "Artist": "Rudimental",
      "Album": "Home (Deluxe Edition)",
      "Year": 2013,
      "Urlp": "666989860",
      "Urli": "666989917"
    },
    "StartTimes": [
      1525731305
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music6/v4/8f/67/0d/8f670db5-3c7d-c56b-2171-0bb3c6f9eb8a/825646544738.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 189196,
      "Genre": "Dance",
      "Name": "No Money",
      "Artist": "Galantis",
      "Album": "The Aviary",
      "Year": 2016,
      "Urlp": "1257258777",
      "Urli": "1257259293"
    },
    "StartTimes": [
      1525731116
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music128/v4/c2/96/bf/c296bfde-fce6-3e38-1362-266787bc5abe/075679894755.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 208826,
      "Genre": "Pop",
      "Name": "New Rules",
      "Artist": "Dua Lipa",
      "Album": "Dua Lipa (Deluxe)",
      "Year": 2017,
      "Urlp": "1228739599",
      "Urli": "1228739609"
    },
    "StartTimes": [
      1525730907
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music122/v4/15/1e/32/151e323b-a12a-4059-44ce-16a6735b382e/190295807870.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 195678,
      "Genre": "Pop",
      "Name": "Chasing Highs",
      "Artist": "ALMA",
      "Album": "Chasing Highs - Single",
      "Year": 2017,
      "Urlp": "1216772739",
      "Urli": "1216772784"
    },
    "StartTimes": [
      1525730711
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music111/v4/01/72/2f/01722f7b-7e73-1374-234a-e06318007ce2/UMG_cvrart_00602557429022_01_RGB72_1800x1800_17UMGIM03037.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 173984,
      "Genre": "Hip-Hop/Rap",
      "Name": "One Dance (feat. Wizkid \u0026 Kyla)",
      "Artist": "Drake",
      "Album": "Views",
      "Year": 2016,
      "Urlp": "1109766593",
      "Urli": "1109766881"
    },
    "StartTimes": [
      1525730537
    ],
    "Artwork": "https://is2-ssl.mzstatic.com/image/thumb/Music60/v4/31/b5/7d/31b57d0c-4881-1563-8aa4-cdd0d1a32e81/UMG_cvrart_00602547943507_01_RGB72_1800x1800_16UMGIM27643.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 201933,
      "Genre": "Dance",
      "Name": "Levels (Radio Edit)",
      "Artist": "Avicii",
      "Album": "20 #1's: Workout",
      "Year": 2011,
      "Urlp": "1185110292",
      "Urli": "1185110484"
    },
    "StartTimes": [
      1525730335
    ],
    "Artwork": "https://is2-ssl.mzstatic.com/image/thumb/Music122/v4/d1/3e/7a/d13e7a7d-54c2-6ddc-2d70-3da9845cf9a7/UMG_cvrart_00602557244373_01_RGB72_1800x1800_16UMGIM74016.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 174933,
      "Genre": "Dance",
      "Name": "Solo Dance",
      "Artist": "Martin Jensen",
      "Album": "Ultra 2018",
      "Year": 2016,
      "Urlp": "1309963000",
      "Urli": "1309963896"
    },
    "StartTimes": [
      1525730160
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music118/v4/f7/1b/6d/f71b6d28-cfa2-accb-84ee-fbe68ef913d2/0617465919053.png/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 212026,
      "Genre": "Dance",
      "Name": "Don't You Worry Child (feat. John Martin)",
      "Artist": "Swedish House Mafia",
      "Album": "Aftercluv Dancelab the Drop, Vol. 1",
      "Year": 2012,
      "Urlp": "1018117789",
      "Urli": "1018117798"
    },
    "StartTimes": [
      1525729948
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music7/v4/37/3c/94/373c94df-364b-a485-9c9d-3ae50aa530e9/UMG_cvrart_00600753611319_01_RGB72_1500x1500_15UMGIM27846.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 267080,
      "Genre": "R\u0026B/Soul",
      "Name": "Don't Stop the Music",
      "Artist": "Rihanna",
      "Album": "Good Girl Gone Bad: Reloaded",
      "Year": 2007,
      "Urlp": "1168770543",
      "Urli": "1168770971"
    },
    "StartTimes": [
      1525729681
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music71/v4/50/7c/ee/507cee76-58ce-f8df-579f-15eb6aac14de/UMG_cvrart_00602557132335_01_RGB72_1800x1800_16UMGIM59202.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 252000,
      "Genre": "Pop",
      "Name": "Move Your Body (Single Mix)",
      "Artist": "Sia",
      "Album": "Move Your Body (Single Mix) - Single",
      "Year": 2017,
      "Urlp": "1187387847",
      "Urli": "1187387952"
    },
    "StartTimes": [
      1525729651
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music122/v4/1e/19/75/1e1975ef-7161-06c4-d458-787429318eee/886446298153.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 131064,
      "Genre": "Pop",
      "Name": "Mine",
      "Artist": "Bazzi",
      "Album": "COSMIC",
      "Year": 2017,
      "Urlp": "1369438774",
      "Urli": "1369439270"
    },
    "StartTimes": [
      1525729563,
      1525729432,
      1525729301,
      1525729194
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music128/v4/9c/70/fd/9c70fd08-9adc-3a57-4cb6-945fa4c54127/075679874580.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 243216,
      "Genre": "Alternative",
      "Name": "Heaven's Only Wishful",
      "Artist": "MorMor",
      "Album": "Heaven's Only Wishful - Single",
      "Year": 2018,
      "Urlp": "1365980325",
      "Urli": "1365980333"
    },
    "StartTimes": [
      1525729187
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music128/v4/a2/3c/a9/a23ca98c-5e9e-bc7a-1f15-2e24bac4e864/5054526784345_1.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 178197,
      "Genre": "Alternative",
      "Name": "Blood To Gold",
      "Artist": "slenderbodies",
      "Album": "Fabulist: Extended - EP",
      "Year": 2018,
      "Urlp": "1327615188",
      "Urli": "1327615251"
    },
    "StartTimes": [
      1525729072,
      1525728894,
      1525728715,
      1525728537,
      1525728359,
      1525728181,
      1525728003,
      1525727824,
      1525727646,
      1525727468,
      1525727290,
      1525727112,
      1525726979,
      1525726801,
      1525726622,
      1525726444,
      1525726266,
      1525726084
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music118/v4/cd/7f/7a/cd7f7a41-deb4-6d55-e706-93e36807c10e/content_art_2F63vrfvh5QiCIINUkfUlZ_IMG_1030.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 322266,
      "Genre": "Rock",
      "Name": "Short Change Hero",
      "Artist": "The Heavy",
      "Album": "The House That Dirt Built",
      "Year": 2009,
      "Urlp": "1202348085",
      "Urli": "1202348740"
    },
    "StartTimes": [
      1525676771
    ],
    "Artwork": "https://is2-ssl.mzstatic.com/image/thumb/Music111/v4/c2/94/de/c294de9a-c763-4047-51e4-3ea54e1cdc9e/5054429118537.png/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 251480,
      "Genre": "Alternative",
      "Name": "Next Year",
      "Artist": "Two Door Cinema Club",
      "Album": "Beacon",
      "Year": 2012,
      "Urlp": "544390690",
      "Urli": "544390803"
    },
    "StartTimes": [
      1525676732
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music/v4/7b/2d/1a/7b2d1a06-cebf-b431-c4dc-85499fbe0f8a/UMG_cvrart_00602537115150_01_RGB72_1200x1200_12UMGIM37743.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 236440,
      "Genre": "Alternative",
      "Name": "Sleep Alone",
      "Artist": "Two Door Cinema Club",
      "Album": "Beacon",
      "Year": 2012,
      "Urlp": "544390690",
      "Urli": "544390861"
    },
    "StartTimes": [
      1525676495
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music/v4/7b/2d/1a/7b2d1a06-cebf-b431-c4dc-85499fbe0f8a/UMG_cvrart_00602537115150_01_RGB72_1200x1200_12UMGIM37743.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 147026,
      "Genre": "Alternative",
      "Name": "White Winter Hymnal",
      "Artist": "Fleet Foxes",
      "Album": "Fleet Foxes",
      "Year": 2008,
      "Urlp": "281086394",
      "Urli": "281086428"
    },
    "StartTimes": [
      1525676163,
      1525675869,
      1525675722,
      1525675574
    ],
    "Artwork": "https://is2-ssl.mzstatic.com/image/thumb/Music/f5/f0/37/mzi.smnpawlm.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 292093,
      "Genre": "Hip Hop/Rap",
      "Name": "Power",
      "Artist": "Kanye West",
      "Album": "My Beautiful Dark Twisted Fantasy",
      "Year": 2010,
      "Urlp": "403822142",
      "Urli": "403822299"
    },
    "StartTimes": [
      1525675570
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music/81/08/42/mzi.jpzmjowd.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 272842,
      "Genre": "Electronic",
      "Name": "The Animals",
      "Artist": "Ladytron",
      "Album": "The Animals - Single",
      "Year": 2018,
      "Urlp": "1357762655",
      "Urli": "1357762656"
    },
    "StartTimes": [
      1525675479,
      1525675206,
      1525671832,
      1525671559
    ],
    "Artwork": "https://is1-ssl.mzstatic.com/image/thumb/Music128/v4/70/9b/77/709b7746-6465-84cb-5d86-57c869a4ab7d/5054526684454_1.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 178197,
      "Genre": "Alternative",
      "Name": "Blood To Gold",
      "Artist": "slenderbodies",
      "Album": "Fabulist: Extended - EP",
      "Year": 2018,
      "Urlp": "1327615188",
      "Urli": "1327615251"
    },
    "StartTimes": [
      1525671526,
      1525671348,
      1525668756,
      1525668497
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music118/v4/cd/7f/7a/cd7f7a41-deb4-6d55-e706-93e36807c10e/content_art_2F63vrfvh5QiCIINUkfUlZ_IMG_1030.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 243216,
      "Genre": "Alternative",
      "Name": "Heaven's Only Wishful",
      "Artist": "MorMor",
      "Album": "Heaven's Only Wishful - Single",
      "Year": 2018,
      "Urlp": "1365980325",
      "Urli": "1365980333"
    },
    "StartTimes": [
      1525668495
    ],
    "Artwork": "https://is3-ssl.mzstatic.com/image/thumb/Music128/v4/a2/3c/a9/a23ca98c-5e9e-bc7a-1f15-2e24bac4e864/5054526784345_1.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 239200,
      "Genre": "Alternative",
      "Name": "Deep Forest Green",
      "Artist": "Husky Rescue",
      "Album": "The Long Lost Friend (Special Edition)",
      "Year": 2015,
      "Urlp": "1052912154",
      "Urli": "1052912459"
    },
    "StartTimes": [
      1525663454
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music19/v4/b6/d8/e0/b6d8e0c3-a374-e664-df72-c5c7d2f4676e/mzm.uqbftnkm.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 178197,
      "Genre": "Alternative",
      "Name": "Blood To Gold",
      "Artist": "slenderbodies",
      "Album": "Fabulist: Extended - EP",
      "Year": 2018,
      "Urlp": "1327615188",
      "Urli": "1327615251"
    },
    "StartTimes": [
      1525663300,
      1525663122
    ],
    "Artwork": "https://is4-ssl.mzstatic.com/image/thumb/Music118/v4/cd/7f/7a/cd7f7a41-deb4-6d55-e706-93e36807c10e/content_art_2F63vrfvh5QiCIINUkfUlZ_IMG_1030.jpg/1200x630bb.jpg"
  },
  {
    "Song": {
      "Duration": 237985,
      "Genre": "Alternative",
      "Name": "Wildwind",
      "Artist": "Young Dreams",
      "Album": "Waves 2 You",
      "Year": 2018,
      "Urlp": "1300583115",
      "Urli": "1300583184"
    },
    "StartTimes": [
      1525663060
    ],
    "Artwork": "https://is5-ssl.mzstatic.com/image/thumb/Music118/v4/0a/53/83/0a538394-bdd9-4d8d-ee9c-550f992f45de/cover.jpg/1200x630bb.jpg"
  }
]
`
	w.Write([]byte(b))
}
