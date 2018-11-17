package security

import (
	"testing"
)

func TestJwtTokenParsing(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI3NzAyNmYxMWM2ODQ0YzVlOWEyYmM1YWNkNGQxYjljNyIsImlhdCI6MTU0MDcyODAxMiwiaXNzIjoibG9naW4uYmluZ2dsLm5ldCIsImV4cCI6MTc0MTMzMjgxMiwic3ViIjoiaGVucmlrLmJpbmdnbEBnbWFpbC5jb20iLCJUeXBlIjoibG9naW4uVXNlciIsIlVzZXJOYW1lIjoiaGVucmlrLmJpbmdnbCIsIkVtYWlsIjoiaGVucmlrLmJpbmdnbEBnbWFpbC5jb20iLCJDbGFpbXMiOlsibXlkbXN8aHR0cHM6Ly9teWRtcy5iaW5nZ2wubmV0L3xVc2VyIiwibG9naW4uYmluZ2dsLm5ldHxodHRwczovL2xvZ2luLmJpbmdnbC5uZXR8VXNlcjtBZG1pbiIsInRlc3RzaXRlfGh0dHBzOi8vd3d3LmJpbmdnbC5uZXQvfFJvbGUxIl0sIlVzZXJJZCI6IjExODAwNDU5MzUwMzk2MzkwMDc5NCIsIkRpc3BsYXlOYW1lIjoiSGVucmlrIEJpbmdnbCIsIlN1cm5hbWUiOiJCaW5nZ2wiLCJHaXZlbk5hbWUiOiJIZW5yaWsifQ.aI5tN5tAku4ZEA5iOlS6VypAS-a4QGxwIawB9dtMu3U"
	secret := "secret"

	_, err := ParseJwtToken(token, secret, "login.binggl.net")
	if err != nil {
		t.Errorf("Could not parse jwt token %s", err)
	}

	_, err = ParseJwtToken(token, secret, "wrong issued")
	if err == nil {
		t.Error("The issuer check should work!", err)
	}
}
