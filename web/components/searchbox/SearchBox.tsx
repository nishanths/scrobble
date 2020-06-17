import React from "react"
import "../../scss/searchbox/search-box.scss"

type SearchBoxProps = {
	onChange: (v: string) => void
	value: string
	placeholder?: string
}

export const SearchBox: React.FC<SearchBoxProps> = ({ value, onChange, placeholder }) => {
	const onInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		onChange(e.target.value)
	}

	return <div className="SearchBox">
		<input value={value} onChange={onInputChange} placeholder={placeholder} />
	</div>
}
