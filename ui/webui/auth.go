package webui

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
)

// authCookie carries the browser's proof of having presented the token. It
// holds a digest, so the token itself is not kept in the browser's cookie jar.
const authCookie = "shellgym_auth"

// requireAuth is the optional shared-token check (Options.AuthToken), for a
// daemon reachable by more people than its student. Every route is covered,
// the static assets included. A request passes with any of:
//
//   - `Authorization: Bearer <token>` - non-browser clients (`shellgym solve`);
//   - the auth cookie - the browser, on every request after the first;
//   - `?token=<token>` on a GET - the browser's first visit: the cookie is set
//     and the request is redirected to the same URL without the token.
//
// Without a token configured, the check is off.
func (s *Server) requireAuth(next http.Handler) http.Handler {
	token := s.Opts.AuthToken
	if token == "" {
		return next
	}
	digest := cookieDigest(token)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if bearer, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer "); ok && sameSecret(bearer, token) {
			next.ServeHTTP(w, r)
			return
		}
		if c, err := r.Cookie(authCookie); err == nil && sameSecret(c.Value, digest) {
			next.ServeHTTP(w, r)
			return
		}
		q := r.URL.Query()
		if r.Method == http.MethodGet && q.Has("token") && sameSecret(q.Get("token"), token) {
			http.SetCookie(w, &http.Cookie{
				Name:     authCookie,
				Value:    digest,
				Path:     "/",
				HttpOnly: true,
				// Lax, not Strict: the cookie must survive the redirect below
				// when the link was clicked on another site. It is still never
				// sent with a cross-site POST or WebSocket handshake.
				SameSite: http.SameSiteLaxMode,
				Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
			})
			q.Del("token")
			u := *r.URL
			u.RawQuery = q.Encode()
			http.Redirect(w, r, u.RequestURI(), http.StatusSeeOther)
			return
		}
		w.Header().Set("WWW-Authenticate", `Bearer realm="shellgym"`)
		http.Error(w, "This gym requires an auth token. Open the page as /?token=<token>, "+
			"or send an `Authorization: Bearer <token>` header.", http.StatusUnauthorized)
	})
}

func cookieDigest(token string) string {
	sum := sha256.Sum256([]byte("shellgym-auth-cookie:" + token))
	return hex.EncodeToString(sum[:])
}

func sameSecret(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
