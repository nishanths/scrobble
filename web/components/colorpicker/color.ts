export const colors = [
  "gray",
  "red",
  "orange",
  "blue",
  "yellow",
  "white",
  "green",
  "violet",
  "pink",
  "brown",
  "black",
] as const;

export type Color = typeof colors[number]
