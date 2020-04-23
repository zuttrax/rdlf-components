package resty

type options struct {
	path   string
	params Parameters
	body   []byte
}

type RestyOptions func(*options)

func SetQueryParams(p Parameters) RestyOptions {
	return func(opt *options) {
		opt.params = p
	}
}

func SetBody(body []byte) RestyOptions {
	return func(opt *options) {
		opt.body = body
	}
}

func SetPath(path string) RestyOptions {
	return func(opt *options) {
		opt.path = path
	}
}
