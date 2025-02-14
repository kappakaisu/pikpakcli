package url

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/52funny/pikpakcli/conf"
	"github.com/52funny/pikpakcli/internal/pikpak"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var NewUrlCommand = &cobra.Command{
	Use:   "url",
	Short: `Create a file according to url`,
	Run: func(cmd *cobra.Command, args []string) {
		p := pikpak.NewPikPak(conf.Config.Username, conf.Config.Password)
		err := p.Login()
		if err != nil {
			logrus.Errorln("Login Failed:", err)
		}
		// input mode
		if strings.TrimSpace(input) != "" {
			f, err := os.OpenFile(input, os.O_RDONLY, 0666)
			if err != nil {
				logrus.Errorln("Open file %s failed:", input, err)
				return
			}
			reader := bufio.NewReader(f)
			shas := make([]string, 0)
			for {
				lineBytes, _, err := reader.ReadLine()
				if err == io.EOF {
					break
				}
				shas = append(shas, string(lineBytes))
			}
			handleNewUrl(&p, shas)
			return
		}

		// args mode
		if len(args) > 0 {
			handleNewUrl(&p, args)
		} else {
			logrus.Errorln("Please input the folder name")
		}
	},
}

var path string

var input string

func init() {
	NewUrlCommand.Flags().StringVarP(&path, "path", "p", "/", "The path of the folder")
	NewUrlCommand.Flags().StringVarP(&input, "input", "i", "", "The input of the sha file")
}

// new folder
func handleNewUrl(p *pikpak.PikPak, shas []string) {
	parentId, err := p.GetPathFolderId(path)
	if err != nil {
		logrus.Errorf("Get parent id failed: %s\n", err)
		return
	}

	for _, url := range shas {
		err := p.CreateUrlFile(parentId, url)
		if err != nil {
			logrus.Errorln("Create url file failed: ", err)
			continue
		}
		logrus.Infoln("Create url file success: ", url)
	}
}
