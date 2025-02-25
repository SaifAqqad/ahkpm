package cmd

import (
	"ahkpm/src/core"
	"ahkpm/src/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs either the specified package or all packages listed in ahkpm.json",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			utils.Exit("Error getting current directory")
		}

		ahkpmFileExists, err := utils.FileExists(cwd + `\ahkpm.json`)
		if err != nil {
			utils.Exit("Error checking if ahkpm.json exists")
		}

		if !ahkpmFileExists {
			fmt.Println("ahkpm.json not found in current directory. Run `ahkpm init` to create one.")
			os.Exit(1)
		}

		installer := core.Installer{}

		if len(args) == 0 {
			fmt.Println("Installing all dependencies")
			dependencies := core.AhkpmJson{}.ReadFromFile().Dependencies
			for dep, ver := range dependencies() {
				v, err := core.Version{}.FromString(ver)
				if err != nil {
					utils.Exit(err.Error())
				}
				installer.InstallSinglePackage(dep, v)
			}
			return
		}

		if len(args) > 1 {
			// TODO: support specifying multiple packages
			fmt.Println("Please specify only one package to install")
			return
		}

		packageToInstall := args[0]
		var versionSpecifier string
		if strings.Contains(packageToInstall, "@") {
			splitArg := strings.SplitN(packageToInstall, "@", 2)
			packageToInstall = splitArg[0]
			versionSpecifier = splitArg[1]
		}

		version, err := core.Version{}.FromString(versionSpecifier)
		if err != nil {
			utils.Exit(err.Error())
		}

		installer.InstallSinglePackage(packageToInstall, version)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
