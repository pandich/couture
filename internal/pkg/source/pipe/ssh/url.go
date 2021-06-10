package ssh

import (
	"couture/internal/pkg/model"
	"github.com/melbahja/goph"
	"net/url"
	"os"
	"os/user"
	"path"
	"strconv"
)

type sshURL model.SourceURL

// FEATURE rsa key configurable and ability to use PEM key
// see: client.Config.Auth, err = goph.Key() etc
func (u sshURL) getClient() (*goph.Client, error) {
	auth, err := u.getAuth()
	if err != nil {
		return nil, err
	}

	usr, err := u.getUser()
	if err != nil {
		return nil, err
	}

	client, err := goph.New(usr, u.Host, auth)
	if err != nil {
		return nil, err
	}
	var port uint = 22
	rawURL := url.URL(u)
	if s := rawURL.Port(); s != "" {
		i, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, err
		}
		port = uint(i)
	}
	client.Config.Port = port
	return client, err
}

func (u sshURL) getUser() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	var usr = currentUser.Username
	if userInfo := u.User; userInfo != nil {
		usr = userInfo.Username()
	}
	return usr, nil
}

func (u sshURL) getPassphrase() string {
	if u.User != nil {
		if password, ok := u.User.Password(); ok {
			return password
		}
	}
	return ""
}

func (u sshURL) getAuth() (goph.Auth, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	idPath := path.Join(homeDir, ".ssh", "id_rsa")
	return goph.Key(idPath, u.getPassphrase())
}
