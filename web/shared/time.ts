export function dateDisplay(d: Date, now: Date): string {
  let copy = new Date(now.getTime())

  // is it today?
  if (sameDate(d, copy)) {
    return "Today"
  }

  let original = copy.getDate() // save before modification

  // is it yesterday?
  copy.setDate(original - 1)
  if (sameDate(d, copy)) {
    return "Yesterday"
  }

  for (let c of [2, 3, 4]) {
    copy.setDate(original - c)
    if (sameDate(d, copy)) {
      return c + " days ago"
    }
  }

  copy.setDate(original) // restore

  // use a regular looking date string
  return d.getFullYear() != copy.getFullYear() ? `${shortMonth(d)} ${d.getDate()} ${d.getFullYear()}` :
    `${shortMonth(d)} ${d.getDate()}`
}

function sameDate(a: Date, b: Date): boolean {
  return a.getDate() == b.getDate() &&
    a.getMonth() == b.getMonth() &&
    a.getFullYear() == b.getFullYear()
}

const shortMonthNames = [
  "Jan", "Feb", "Mar", "Apr", "May", "Jun",
  "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
]

function shortMonth(d: Date) {
  return shortMonthNames[d.getMonth()]
}
