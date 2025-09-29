package localizations

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// Global Vars
var (
	Localization *Localizer
)

// Errors
var (
	FailedGetSystemLocalErr = fmt.Errorf("cannot determine locale")
)

const (
	getLocaleErrText = "failed to get system locale, expected in env variable LANG for UNIX or Get-Culture | select -exp Name for windows "

	fallbackLocale = "en"
)

func init() {
	systemLang, _ := GetLocale()
	Localization = New(systemLang, fallbackLocale)
}

// GetLocale source:https://stackoverflow.com/questions/51829386/golang-get-system-language
func GetLocale() (string, string) {
	osHost := runtime.GOOS
	defaultLang := "en"
	defaultLoc := "US"
	switch osHost {
	case "windows":
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "-")
			lang := langLoc[0]
			loc := langLoc[1]
			return lang, loc
		}
	case "darwin":
		cmd := exec.Command("sh", "osascript -e 'user locale of (get system info)'")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := langLoc[1]
			return lang, loc
		}
	case "linux":
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.TrimSpace(envlang)
			langLocRaw = strings.Split(envlang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := langLoc[1]
			return lang, loc
		}
	}
	return defaultLang, defaultLoc
}
