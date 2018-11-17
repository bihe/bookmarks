package security

import (
	"testing"
)

func TestJwtTokenParsing(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiI3NzAyNmYxMWM2ODQ0YzVlOWEyYmM1YWNkNGQxYjljNyIsImlhdCI6MTU0MDcyODAxMiwiaXNzIjoibG9naW4uYS5iLmMuZGUiLCJleHAiOjE3NDEzMzI4MTIsInN1YiI6ImhlbnJpa0BhLmIuYy5kZSIsIlR5cGUiOiJsb2dpbi5Vc2VyIiwiVXNlck5hbWUiOiJ1c2VyLm5hbWUiLCJFbWFpbCI6ImhlbnJpa0BhLmIuYy5kZSIsIkNsYWltcyI6WyJib29rbWFya3N8aHR0cDovL2xvY2FsaG9zdDozMDAwfEFkbWluIl0sIlVzZXJJZCI6IjExMSIsIkRpc3BsYXlOYW1lIjoiVXNlcjEgTmFtZSIsIlN1cm5hbWUiOiJOYW1lIiwiR2l2ZW5OYW1lIjoiVXNlcjEifQ.GzR2RAEYyTazO6zaDcCWdSlvaso3TUW-Tem3-l0qnnw"
	secret := "secret"

	_, err := ParseJwtToken(token, secret, "login.a.b.c.de")
	if err != nil {
		t.Errorf("Could not parse jwt token %s", err)
	}

	_, err = ParseJwtToken(token, secret, "wrong issued")
	if err == nil {
		t.Error("The issuer check should work!", err)
	}
}
