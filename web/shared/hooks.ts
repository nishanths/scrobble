import { useState, useEffect, useRef } from "react"

export const useStateRef = <S>(initialState: S) => {
	const [v, setv] = useState(initialState)
	const ref = useRef(v)
	useEffect(() => {
		ref.current = v
	}, [v])
	return [v, ref, setv] as const
}
