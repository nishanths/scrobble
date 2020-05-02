export const Colors = [
	"red",
	"orange",
	"brown",
	"yellow",
	"green",
	"blue",
	"violet",
	"pink",
	"black",
	"gray",
	"white",
] as const;

export type Color = typeof Colors[number]
