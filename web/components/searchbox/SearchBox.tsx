import React from "react"
import { SearchIcon } from "./SearchIcon"
import "../../scss/searchbox/search-box.scss"

type SearchBoxProps = {
	onChange: (v: string) => void
	value: string
	placeholder?: string
}

// SearchBox provides a component that can be used as a search box.
// The ref provides access to the underlying input element.
export const SearchBox = React.forwardRef<HTMLInputElement, SearchBoxProps>(
	function SearchBox({ onChange, value, placeholder }, ref) {
		const onInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
			onChange(e.target.value)
		}

		return <div className="SearchBox">
			<div className="icon">
				{SearchIcon}
			</div>
			<input ref={ref} type="text" spellCheck="false" value={value} onChange={onInputChange} placeholder={placeholder} />
		</div>
	}
)
