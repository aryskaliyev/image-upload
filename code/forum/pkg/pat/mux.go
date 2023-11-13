package pat

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

type PatternServeMux struct {
	NotFound http.Handler
	handlers map[string][]*patHandler
}

type PageData struct {
	Title        string
	ErrorMessage string
}

func New() *PatternServeMux {
	return &PatternServeMux{handlers: make(map[string][]*patHandler)}
}

func (p *PatternServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, ph := range p.handlers[r.Method] {
		if params, ok := ph.try(r.URL.EscapedPath()); ok {
			if len(params) > 0 && !ph.redirect {
				r.URL.RawQuery = url.Values(params).Encode() + "&" + r.URL.RawQuery
			}
			ph.ServeHTTP(w, r)
			return
		}
	}

	if p.NotFound != nil {
		p.NotFound.ServeHTTP(w, r)
		return
	}

	allowed := make([]string, 0, len(p.handlers))
	for meth, handlers := range p.handlers {
		if meth == r.Method {
			continue
		}

		for _, ph := range handlers {
			if _, ok := ph.try(r.URL.EscapedPath()); ok {
				allowed = append(allowed, meth)
			}
		}
	}

	if len(allowed) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		data := PageData{
			Title:        "Error Page",
			ErrorMessage: "Something went wrong.",
		}
		tmpl, err := template.New("error").Parse(`
					<!doctype html>
					<html>
						<head>
							<meta charset='utf-8'>
							<title>{{.Title}} - Forum</title>
							<link rel='stylesheet' href='/static/css/main.css'>
							<link rel='shortcut icon' href='/static/img/favicon.ico' type='image/x-icon'>
							<link rel='stylesheet' href='https://fonts.googleapis.com/css?family=Ubuntu'>
						</head>
						<body>
							<header>
								<h1><a href='/'>Forum</a></h1>
							</header>
						<nav>
						<div>
								<a href='/'>Home</a>
						</div>
						</nav>
						<section>
							<h1>Oops! {{.ErrorMessage}}</h1>
						</section>
							<footer>Powered by <a href='https://golang.org'>Go</a> in 2023</footer>
							<script src="/static/js/main.js" type="text/javascript"></script>
						</body>
					</html>
		`)
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			return
		}
		return
	}

	w.Header().Add("Allow", strings.Join(allowed, ", "))
	http.Error(w, "Method Not Allowed", 405)
}

func (p *PatternServeMux) Head(pat string, h http.Handler) {
	p.Add("HEAD", pat, h)
}

func (p *PatternServeMux) Get(pat string, h http.Handler) {
	p.Add("HEAD", pat, h)
	p.Add("GET", pat, h)
}

func (p *PatternServeMux) Post(pat string, h http.Handler) {
	p.Add("POST", pat, h)
}

func (p *PatternServeMux) Put(pat string, h http.Handler) {
	p.Add("PUT", pat, h)
}

func (p *PatternServeMux) Del(pat string, h http.Handler) {
	p.Add("DELETE", pat, h)
}

func (p *PatternServeMux) Options(pat string, h http.Handler) {
	p.Add("OPTIONS", pat, h)
}

func (p *PatternServeMux) Patch(pat string, h http.Handler) {
	p.Add("PATCH", pat, h)
}

func (p *PatternServeMux) Add(meth, pat string, h http.Handler) {
	p.add(meth, pat, h, false)
}

func (p *PatternServeMux) add(meth, pat string, h http.Handler, redirect bool) {
	handlers := p.handlers[meth]
	for _, p1 := range handlers {
		if p1.pat == pat {
			return
		}
	}
	handler := &patHandler{
		pat:      pat,
		Handler:  h,
		redirect: redirect,
	}
	p.handlers[meth] = append(handlers, handler)

	n := len(pat)
	if n > 0 && pat[n-1] == '/' {
		p.add(meth, pat[:n-1], http.HandlerFunc(addSlashRedirect), true)
	}
}

func addSlashRedirect(w http.ResponseWriter, r *http.Request) {
	u := *r.URL
	u.Path += "/"
	http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
}

func Tail(pat, path string) string {
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(pat):
			if pat[len(pat)-1] == '/' {
				return path[i:]
			}
			return ""
		case pat[j] == ':':
			var nextc byte
			_, nextc, j = match(pat, isAlnum, j+1)
			_, _, i = match(path, matchPart(nextc), i)
		case path[i] == pat[j]:
			i++
			j++
		default:
			return ""
		}
	}
	return ""
}

type patHandler struct {
	pat string
	http.Handler
	redirect bool
}

func (ph *patHandler) try(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(ph.pat):
			if ph.pat != "/" && len(ph.pat) > 0 && ph.pat[len(ph.pat)-1] == '/' {
				return p, true
			}
			return nil, false
		case ph.pat[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(ph.pat, isAlnum, j+1)
			val, _, i = match(path, matchPart(nextc), i)
			escval, err := url.QueryUnescape(val)
			if err != nil {
				return nil, false
			}
			p.Add(":"+name, escval)
		case path[i] == ph.pat[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(ph.pat) {
		return nil, false
	}
	return p, true
}

func matchPart(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}
