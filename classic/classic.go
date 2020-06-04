package niniteclassic

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
)

// Refer to https://ninite.com/help/features/switches.html for any missing details.

// Classic uses the Ninite Pro Classic interface on the local computer.
type Classic struct {
	// attributes
	path   string
	locale string
	proxy  struct {
		server string
		port   int
	}
	proxyAuth struct {
		username string
		password string
	}
	selectedApps []string
	excludedApps []string
	remote       []string
	remoteAuth   struct {
		username string
		password string
	}
	disableShortcuts  bool
	disableAutoUpdate bool
	allUsers          bool
	cachePath         string
	noCache           bool
	cleanCache        bool

	// verbs
	updateOnly bool
	uninstall  bool
	freeze     struct {
		outputFilename string
		locales        []string
	}
	list bool
}

// Status is a status
type Status struct {
	App     string
	Status  string
	Reason  string
	Version string
}

// AppVersion is an available app version
type AppVersion struct {
	App              string
	Version          string
	CurrentVersion   bool // ???: I don't think this indicates that it is installed nor the version that is presently installed
	AlternateVersion bool // this indicates that this version has to be selected explicitly in order to be installed
}

// AppAudit is an app that may or may not be installed
type AppAudit struct {
	App       string
	Version   string
	Status    string
	Installed bool
}

var statusMatch = regexp.MustCompile(`^\s*(?P<app>[^:\r\n]+)\s+:\s+(?P<status>[^\r\n\(\)]+)(?:\s+\((?P<reason>.+)\))?$`)
var versionMatch = regexp.MustCompile(`^\s*(?P<app>[^:\r\n]+)\s+:\s+(?P<type>[\*\(])?(?P<version>[^\r\n\(\)]+)\)?$`)
var auditMatch = regexp.MustCompile(`^\s*(?P<app>[^:\r\n]+)\s+:\s+(?P<status>[^\r\n\(\)\-]+)(?:\s+-\s+(?P<version>.+))?$`)

func (c Classic) composeArgs() []string {
	args := []string{"/silent", "."}

	//
	// attributes

	if c.locale != "" {
		args = append(args, "/locale", c.locale)
	}

	if c.proxy.server != "" && c.proxy.port != 0 {
		args = append(args, "/proxy", c.proxy.server, strconv.Itoa(c.proxy.port))
	}

	if c.proxyAuth.username != "" && c.proxyAuth.password != "" {
		args = append(args, "/proxyauth", c.proxyAuth.username, c.proxyAuth.password)
	}

	if len(c.selectedApps) > 0 {
		args = append(append(args, "/select"), c.selectedApps...)
	}

	if len(c.excludedApps) > 0 {
		args = append(append(args, "/exclude"), c.excludedApps...)
	}

	if len(c.remote) > 0 {
		args = append(append(args, "/remote"), c.remote...)
	}

	if c.remoteAuth.username != "" && c.remoteAuth.password != "" {
		args = append(args, "/remoteauth", c.remoteAuth.username, c.remoteAuth.password)
	}

	if c.disableShortcuts {
		args = append(args, "/disableshortcuts")
	}

	if c.disableAutoUpdate {
		args = append(args, "/disableautoupdate")
	}

	if c.allUsers {
		args = append(args, "/allusers")
	}

	if c.cachePath != "" {
		args = append(args, "/cachepath", c.cachePath)
	}

	if c.noCache {
		args = append(args, "/nocache")
	}

	if c.cleanCache {
		args = append(args, "/cleancache")
	}

	//
	// verbs

	if c.updateOnly {
		args = append(args, "/updateonly")
	}

	if c.freeze.outputFilename != "" {
		args = append(args, "/freeze")
		if len(c.freeze.locales) > 0 {
			args = append(args, c.freeze.locales...)
		}
		args = append(args, c.freeze.outputFilename)
	}

	if c.list {
		args = append(args, "/list", "versions")
	}

	return args
}

func (c Classic) start() (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	args := c.composeArgs()
	cmd := exec.Command(c.path, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, nil, nil, err
	}

	return cmd, stdout, stderr, nil
}

// NewClassic returns a Classic, which uses the Ninite Pro Classic interface running on the local computer.
func NewClassic(path string) (Classic, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return Classic{}, err
	}
	if fi.IsDir() {
		return NewClassic(filepath.Join(path, "NinitePro.exe"))
	}
	return Classic{
		path: path,
	}, nil
}

//
// Attributes

// Locale sets the locale for this instance of Classic.
func (c Classic) Locale(code string) Classic {
	c.locale = code
	return c
}

// Proxy sets the proxy server and port for this instance of Classic.
func (c Classic) Proxy(server string, port int) Classic {
	c.proxy.server = server
	c.proxy.port = port
	return c
}

// ProxyAuth sets the proxy username and password for this instance of Classic.
func (c Classic) ProxyAuth(username, password string) Classic {
	c.proxyAuth.username = username
	c.proxyAuth.password = password
	return c
}

// Select selects specific Ninite-managed apps.
func (c Classic) Select(apps ...string) Classic {
	c.selectedApps = apps
	return c
}

// Exclude excludes specific Ninite-managed apps.
func (c Classic) Exclude(apps ...string) Classic {
	c.excludedApps = apps
	return c
}

// Prefer sets version preferences per Ninite-managed app.
func (c Classic) Prefer() Classic {
	if true {
		panic("Unimplemented") // TODO: implement this
	}
	return Classic{}
}

// Remote identifies remote computers to manage.  This can be either machine addresses or filenames prefixed with `file:`.
func (c Classic) Remote(machines ...string) Classic {
	c.remote = machines
	return c
}

// RemoteAuth sets the remote machine username and password for this instance of Classic.
func (c Classic) RemoteAuth(username, password string) Classic {
	c.remoteAuth.username = username
	c.remoteAuth.password = password
	return c
}

// DisableShortcuts disables creation of shortcuts for installed apps for this instance of Classic.
func (c Classic) DisableShortcuts() Classic {
	c.disableShortcuts = true
	return c
}

// DisableAutoUpdate disables auto-update of installed apps for this instance of Classic.
func (c Classic) DisableAutoUpdate() Classic {
	c.disableAutoUpdate = true
	return c
}

// AllUsers will force some apps to install for all users for this instance of Classic.
func (c Classic) AllUsers() Classic {
	c.allUsers = true
	return c
}

// CachePath sets the cache path for this instance of Classic.
func (c Classic) CachePath(path string) Classic {
	c.cachePath = path
	return c
}

// NoCache disables use or creation of a local download cache for this instance of Classic.
func (c Classic) NoCache() Classic {
	c.noCache = true
	return c
}

// CleanCache cleans up the local download cache by deleting older unused files for this instance of Classic.
func (c Classic) CleanCache() Classic {
	c.cleanCache = true
	return c
}

//
// Verbs

// UpdateOnly performs an update on software that is already installed and does not cause any new software to become installed.
func (c Classic) UpdateOnly(statusChan chan<- Status) error {
	c.updateOnly = true

	cmd, stdout, stderr, err := c.start()
	if err != nil {
		return err
	}
	defer cmd.Wait() // ???: is this necessary? it is possible to return before cmd.Wait is run without this.

	b := bufio.NewReader(stdout)
	for {
		line, err := b.ReadString('\n')
		if err == io.EOF {
			close(statusChan)
			break
		} else if err != nil {
			return err
		}

		if m := statusMatch.FindStringSubmatch(line); len(m) > 0 {
			statusChan <- Status{
				App:    m[1],
				Status: m[2],
				Reason: m[3],
			}
		}
	}

	var stderrResult error
	if se, err := ioutil.ReadAll(stderr); err == nil {
		if len(se) > 0 {
			stderrResult = errors.New(string(se)) // FIXME: this is naive
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if stderrResult != nil { // if all is apparently well but there was text in stderr, use that as an error
		return stderrResult
	}

	return nil
}

// Uninstall performs an uninstall on selected Ninite-managed apps.
func (c Classic) Uninstall(statusChan chan<- Status) error {
	c.uninstall = true

	cmd, stdout, stderr, err := c.start()
	if err != nil {
		return err
	}
	defer cmd.Wait() // ???: is this necessary? it is possible to return before cmd.Wait is run without this.

	b := bufio.NewReader(stdout)
	for {
		line, err := b.ReadString('\n')
		if err == io.EOF {
			close(statusChan)
			break
		} else if err != nil {
			return err
		}

		if m := statusMatch.FindStringSubmatch(line); len(m) > 0 {
			statusChan <- Status{
				App:    m[1],
				Status: m[2],
				Reason: m[3],
			}
		}
	}

	var stderrResult error
	if se, err := ioutil.ReadAll(stderr); err == nil {
		if len(se) > 0 {
			stderrResult = errors.New(string(se)) // FIXME: this is naive
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if stderrResult != nil { // if all is apparently well but there was text in stderr, use that as an error
		return stderrResult
	}

	return nil
}

// Freeze creates an offline installer for the selected Ninite-managed apps.
func (c Classic) Freeze(statusChan chan<- Status, output string, locales ...string) error {
	c.freeze.outputFilename = output
	c.freeze.locales = locales

	cmd, stdout, stderr, err := c.start()
	if err != nil {
		return err
	}
	defer cmd.Wait() // ???: is this necessary? it is possible to return before cmd.Wait is run without this.

	b := bufio.NewReader(stdout)
	for {
		line, err := b.ReadString('\n')
		if err == io.EOF {
			close(statusChan)
			break
		} else if err != nil {
			return err
		}

		if m := statusMatch.FindStringSubmatch(line); len(m) > 0 {
			statusChan <- Status{
				App:     m[1],
				Version: m[2],
			}
		}
	}

	var stderrResult error
	if se, err := ioutil.ReadAll(stderr); err == nil {
		if len(se) > 0 {
			stderrResult = errors.New(string(se)) // FIXME: this is naive
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if stderrResult != nil { // if all is apparently well but there was text in stderr, use that as an error
		return stderrResult
	}

	return nil
}

// List lists all (or selected) Ninite-managed apps available for install, including their versions.
func (c Classic) List(versionChan chan<- AppVersion) error {
	c.list = true

	cmd, stdout, stderr, err := c.start()
	if err != nil {
		return err
	}
	defer cmd.Wait() // ???: is this necessary? it is possible to return before cmd.Wait is run without this.

	b := bufio.NewReader(stdout)
	for {
		line, err := b.ReadString('\n')
		if err == io.EOF {
			close(versionChan)
			break
		} else if err != nil {
			return err
		}

		if m := versionMatch.FindStringSubmatch(line); len(m) > 0 {
			var currentVersion, alternateVersion bool
			if m[2] == "*" {
				currentVersion = true
			} else if m[2] == "(" {
				alternateVersion = true
			}
			versionChan <- AppVersion{
				App:              m[1],
				Version:          m[3],
				CurrentVersion:   currentVersion,
				AlternateVersion: alternateVersion,
			}
		}
	}

	var stderrResult error
	if se, err := ioutil.ReadAll(stderr); err == nil {
		if len(se) > 0 {
			stderrResult = errors.New(string(se)) // FIXME: this is naive
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if stderrResult != nil { // if all is apparently well but there was text in stderr, use that as an error
		return stderrResult
	}

	return nil
}

// Audit lists all (or selected) Ninite-managed apps, including their versions and whether they are installed.
func (c Classic) Audit(auditChan chan<- AppAudit) error {
	c.list = true

	cmd, stdout, stderr, err := c.start()
	if err != nil {
		return err
	}
	defer cmd.Wait() // ???: is this necessary? it is possible to return before cmd.Wait is run without this.

	b := bufio.NewReader(stdout)
	for {
		line, err := b.ReadString('\n')
		if err == io.EOF {
			close(auditChan)
			break
		} else if err != nil {
			return err
		}

		if m := auditMatch.FindStringSubmatch(line); len(m) > 0 {
			var installed bool
			if len(m[3]) > 0 {
				installed = true
			}
			auditChan <- AppAudit{
				App:       m[1],
				Status:    m[2],
				Version:   m[3],
				Installed: installed,
			}
		}
	}

	var stderrResult error
	if se, err := ioutil.ReadAll(stderr); err == nil {
		if len(se) > 0 {
			stderrResult = errors.New(string(se)) // FIXME: this is naive
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	if stderrResult != nil { // if all is apparently well but there was text in stderr, use that as an error
		return stderrResult
	}

	return nil
}
