export function dateDisplay(d: Date, now: Date): string {
    return dateDisplayDesc(d, now)[0]
}

// dateDisplayDesc returns a human-readable display string for the given
// date, and a boolean indicating whether the display string is of the
// form e.g. "today", "3 days ago" (true) or "28 Dec" (false).
//
// All output strings are safe to be capitalized.
export function dateDisplayDesc(d: Date, now: Date): [string, boolean] {
    const copy = new Date(now.getTime())

    // is it today?
    if (sameDate(d, copy)) {
        return ["today", true]
    }

    const original = copy.getDate() // save before modification

    // is it yesterday?
    copy.setDate(original - 1)
    if (sameDate(d, copy)) {
        return ["yesterday", true]
    }

    for (const c of [2, 3, 4]) {
        copy.setDate(original - c)
        if (sameDate(d, copy)) {
            return [c + " days ago", true]
        }
    }

    copy.setDate(original) // restore

    // use a regular looking date string
    const out = d.getFullYear() != copy.getFullYear() ?
        `${d.getDate()} ${shortMonth(d)} ${d.getFullYear()}` :
        `${d.getDate()} ${shortMonth(d)}`

    return [out, false]
}

// sameDate returns whether the given dates have the same date,
// month, and year.
export function sameDate(a: Date, b: Date): boolean {
    return a.getDate() == b.getDate() &&
        a.getMonth() == b.getMonth() &&
        a.getFullYear() == b.getFullYear()
}

const shortMonthNames = [
    "Jan", "Feb", "Mar", "Apr", "May", "Jun",
    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
]

export function shortMonth(d: Date) {
    return shortMonthNames[d.getMonth()]
}
