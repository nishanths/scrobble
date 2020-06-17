import FlexSearch, { CreateOptions, Index } from "flexsearch"
import { hashCode } from "../../shared/util"

type DocFields = (keyof Omit<Doc, "id" | "songIdent">)[]

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

	// TODO: make these options work
	//
	// async: true,
	// worker: 4,
}

export const createIndex = (): Index<Doc> => {
	return new FlexSearch(createOptions) as Index<Doc> // bleh, flexsearch has poor type definition files
}

export const toIndexID = hashCode

export const hasActiveSearch = (s: string): boolean => {
	return s.trim().toLowerCase().length != 0
}
