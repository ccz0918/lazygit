package updates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/jesseduffield/lazygit/pkg/config"
)

// Update checks for updates and does updates
type Updater struct {
	LastChecked string
	Log         *logrus.Logger
	Config      config.AppConfigurer
	NewVersion  string
}

// Updater implements the check and update methods
type Updaterer interface {
	CheckForNewUpdate()
	Update()
}

var (
	projectUrl = "https://github.com/jesseduffield/lazygit"
)

// NewUpdater creates a new updater
func NewUpdater(log *logrus.Logger, config config.AppConfigurer) (*Updater, error) {

	updater := &Updater{
		LastChecked: "today",
		Log:         log,
		Config:      config,
	}
	return updater, nil
}

func (u *Updater) getLatestVersionNumber() (string, error) {
	req, err := http.NewRequest("GET", projectUrl+"/releases/latest", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	byt := []byte(body)
	var dat map[string]interface{}
	if err := json.Unmarshal(byt, &dat); err != nil {
		return "", err
	}
	return dat["tag_name"].(string), nil
}

// CheckForNewUpdate checks if there is an available update
func (u *Updater) CheckForNewUpdate() (string, error) {
	u.Log.Info("Checking for an updated version")
	if u.Config.GetVersion() == "unversioned" {
		u.Log.Info("Current version is not built from an official release so we won't check for an update")
		return "", nil
	}
	newVersion, err := u.getLatestVersionNumber()
	if err != nil {
		return "", err
	}
	u.NewVersion = newVersion
	u.Log.Info("Current version is " + u.Config.GetVersion())
	u.Log.Info("New version is " + newVersion)
	if newVersion != u.Config.GetVersion() {
		u.getBinaryUrl()
		return newVersion, nil
	}
	return "", nil
}

func (u *Updater) mappedOs(os string) string {
	osMap := map[string]string{
		"darwin":  "Darwin",
		"linux":   "Linux",
		"windows": "Windows",
	}
	result, found := osMap[os]
	if found {
		return result
	}
	return os
}

func (u *Updater) mappedArch(arch string) string {
	archMap := map[string]string{
		"386":   "32-bit",
		"amd64": "x86_64",
	}
	result, found := archMap[arch]
	if found {
		return result
	}
	return arch
}

// example: https://github.com/jesseduffield/lazygit/releases/download/v0.1.73/lazygit_0.1.73_Darwin_x86_64.tar.gz
func (u *Updater) getBinaryUrl() (string, error) {
	extension := "tar.gz"
	if runtime.GOOS == "windows" {
		extension = "zip"
	}
	url := fmt.Sprintf(
		"%s/releases/download/%s/lazygit_%s_%s_%s.%s",
		projectUrl,
		u.NewVersion,
		u.NewVersion[1:],
		u.mappedOs(runtime.GOOS),
		u.mappedArch(runtime.GOARCH),
		extension,
	)
	u.Log.Info("url for latest release is " + url)
	return url, nil
}

// worry about what happens if we fail halfway through or if the user closes the app halfway through

func (u *Updater) downloadLatestBinary() error {
	// attempt to download the binary with the url
	// store in our config folder or other appropriate folder
	return nil
}

func (u *Updater) installLatestBinary() error {
	// unzip/untar the binary
	// copy old binary as lazygit_old
	// swap out existing binary for the new one
	return nil
}
