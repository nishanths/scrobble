import FlexSearch, { CreateOptions, Index } from "flexsearch"
import { hashCode, OmitStrict } from "../../shared/util"

type DocFields = (keyof OmitStrict<Doc, "id" | "songIdent">)[]

export type Doc = {
	id: number
	songIdent: string
	title: string
	artist: string
	album: string
}

const indexFields: DocFields = [
	"title",
	"artist",
	"album",
]

const createOptions: CreateOptions = {
	encode: "icase",
	tokenize: "full",
	doc: {
		id: "id",
		field: indexFields,
	},
	threshold: 0,
	resolution: 9,
	depth: 4,

	// TODO: make these options work
	//
	// async: true,
	// worker: 4,
}

export const createIndex = (): Index<Doc> => {
	return FlexSearch.create(createOptions)
}

export const indexID = hashCode

export const hasActiveSearch = (s: string): boolean => {
	return s.trim().toLowerCase().length != 0
}
