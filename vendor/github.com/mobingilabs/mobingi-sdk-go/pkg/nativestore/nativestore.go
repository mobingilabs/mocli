package nativestore

import (
	"fmt"

	dcred "github.com/docker/docker-credential-helpers/credentials"
	"github.com/pkg/errors"
)

func Set(lbl, url, user, secret string) error {
	switch fmt.Sprintf("%T", ns) {
	case "osxkeychain.Osxkeychain", "wincred.Wincred":
	default:
		return errors.New("native store not supported yet")
	}

	cr := &dcred.Credentials{
		ServerURL: url,
		Username:  user,
		Secret:    secret,
	}

	dcred.SetCredsLabel(lbl)
	return ns.Add(cr)
}

func Get(url string) (string, string, error) {
	return ns.Get(url)
}
