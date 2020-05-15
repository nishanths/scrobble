type DefaultQuery = Map<string, any[]>
type DefaultBody = any
type DefaultHeader = Map<string, any[]>

type Request<H = DefaultHeader, Q = DefaultQuery, B = DefaultBody> = {
	notes?: string
	contentType?: string[]
	header?: H
	query?: Q
	body?: B
}

type BodyByStatusCode<B = DefaultBody> = Map<number, B>

type Response<H = DefaultHeader, B = DefaultBody> = {
	notes?: string
	successStatusCodes: number[]
	failureStatusCodes?: number[]
	successBody?: B | BodyByStatusCode<B>
}

type Endpoint<ReqH = DefaultHeader, ReqQ = DefaultQuery, ReqB = DefaultBody, ResH = DefaultHeader, ResB = DefaultBody> = {
	path: string
	notes?: string
	request?: Request<ReqH, RequQ, ReqB>
	response?: Response<ResH, ResB>
}

type API = {
	baseURL: string
	endpoints: Endpoint[]
}


