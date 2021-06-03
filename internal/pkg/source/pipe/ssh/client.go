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

func getClient(sourceURL model.SourceURL) (*goph.Client, error) {
	auth, err := getAuth(sourceURL)
	if err != nil {
		return nil, err
	}

	usr, err := getUser(sourceURL)
	if err != nil {
		return nil, err
	}

	client, err := goph.New(usr, sourceURL.Host, auth)
	if err != nil {
		return nil, err
	}
	var port uint = 22
	u := url.URL(sourceURL)
	if s := u.Port(); s != "" {
		i, err := strconv.ParseUint(s, 10, 16)
		if err != nil {
			return nil, err
		}
		port = uint(i)
	}
	// FEATURE rsa key configurable and ability to use PEM key @Jim
	// client.Config.Auth, err = goph.Key()
	client.Config.Port = port
	return client, err
}

func getUser(sourceURL model.SourceURL) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	var usr = currentUser.Username
	if userInfo := sourceURL.User; userInfo != nil {
		usr = userInfo.Username()
	}
	return usr, nil
}

func getPassphrase(sourceURL model.SourceURL) string {
	if sourceURL.User != nil {
		if password, ok := sourceURL.User.Password(); ok {
			return password
		}
	}
	return ""
}

func getAuth(sourceURL model.SourceURL) (goph.Auth, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	idPath := path.Join(homeDir, ".ssh", "id_rsa")
	return goph.Key(idPath, getPassphrase(sourceURL))
}
