export const colors = [
	"red",
	"orange",
	"blue",
	"yellow",
	"gray",
	"black",
	"white",
	"violet",
	"pink",
	"green",
	"brown",
] as const;

export type Color = typeof colors[number]
