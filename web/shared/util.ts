import * as base64js from "base64-js"

export const cookieAuthErrorMessage = "Cookie authentication failed. Try signing in again."

export function capitalize(s: string): string {
	return s.charAt(0).toUpperCase() + s.slice(1);
}

export function pluralize(w: string, n: number): string {
	if (n === 1) {
		return w
	}
	return w + "s"
}

export function hasPrefix(s: string, prefix: string): boolean {
	return s.length >= prefix.length && s.slice(0, prefix.length) == prefix
}

export function trimPrefix(s: string, prefix: string): string {
	if (hasPrefix(s, prefix)) {
		return s.slice(prefix.length)
	}
	return s
}

export function assertExhaustive(value: never, message?: string): never {
	if (message === undefined) {
		message = "unexpected value: " + value
	}
	throw new Error(message);
}

// base64Encode encodes the string using base64 standard encoding, i.e.,
// the encoding corresponding to base64.StdEncoding in Go.
//
// Verified using https://runkit.com/nishanths/5cd735892538b9001a7e08d5
// and https://gobyexample.com/base64-encoding.
export function base64Encode(s: string): string {
	const bytes = new TextEncoder().encode(s);
	return base64js.fromByteArray(bytes);
}

// base64Decode is the inverse of base64Encode.
export function base64Decode(s: string): string {
	const bytes = base64js.toByteArray(s);
	return new TextDecoder("utf-8", { fatal: true }).decode(bytes);
}

// hex encoding and decoding adapted from Go package encoding/hex.

const hextable = "0123456789abcdef"

export function hexEncode(s: string): string {
	let out = "";
	for (let x = 0; x < s.length; x++) {
		out += hextable[s[x].charCodeAt(0) >> 4]
		out += hextable[s[x].charCodeAt(0) & 0x0f]
	}
	return out
}

export function hexDecode(s: string): string {
	let j = 1
	let out = ""
	for (; j < s.length; j += 2) {
		const a = fromHexChar(s[j - 1])
		const b = fromHexChar(s[j])
		if (a === null || b === null) {
			throw "invalid byte in input: " + s
		}
		out += String.fromCharCode((a.charCodeAt(0) << 4) | b.charCodeAt(0))
	}
	if (s.length % 2 === 1) {
		if (fromHexChar(s[j - 1]) === null) {
			throw "invalid byte in input: " + s
		}
		throw "bad length: " + s
	}
	return out
}

function fromHexChar(c: string): string | null {
	if ('0' <= c && c <= '9') {
		return String.fromCharCode(c.charCodeAt(0) - '0'.charCodeAt(0))
	}
	if ('a' <= c && c <= 'f') {
		return String.fromCharCode(c.charCodeAt(0) - 'a'.charCodeAt(0) + 10)
	}
	return null
}

export function pathComponents(p: string): string[] {
	return p.split("/").filter(s => s != "")
}

export function assert(cond: boolean, message = "assertion failed"): asserts cond {
	if (!cond) {
		throw new Error(message)
	}
}

export function copyMap<K, V>(m: Map<K, V>): Map<K, V> {
	const n = new Map<K, V>()
	for (const [key, value] of m.entries()) {
		n.set(key, value)
	}
	return n
}

export class MapDefault<K, V> {
	constructor(private readonly def: () => V, private readonly m: Map<K, V> = new Map()) { }

	getOrDefault(key: K): V {
		const v = this.m.get(key)
		if (v === undefined) {
			return this.def()
		}
		return v
	}

	has(key: K): boolean {
		return this.m.has(key)
	}

	set(key: K, value: V): this {
		this.m.set(key, value)
		return this
	}

	copy(): MapDefault<K, V> {
		return new MapDefault(this.def, copyMap(this.m))
	}
}

// Returns a random integer in [min, max).
// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Math/random#Examples
export function randInt(min: number, max: number): number {
	min = Math.ceil(min);
	max = Math.floor(max);
	return Math.floor(Math.random() * (max - min)) + min; // The maximum is exclusive and the minimum is inclusive
}
