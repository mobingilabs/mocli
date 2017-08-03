package cmd

import (
	"encoding/json"

	"github.com/mobingilabs/mocli/client"
	"github.com/mobingilabs/mocli/pkg/check"
	"github.com/mobingilabs/mocli/pkg/cli"
	"github.com/mobingilabs/mocli/pkg/credentials"
	d "github.com/mobingilabs/mocli/pkg/debug"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type authPayload struct {
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	GrantType    string `json:"grant_type,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to Mobingi API",
		Long: `Login to Mobingi API server. If 'grant_type' is set to 'password', you will be prompted to
enter your username and password. Token will be saved in $HOME/.` + cli.BinName() + `/credentials.

Valid 'grant-type' values: client_credentials, password

Examples:

  $ ` + cli.BinName() + ` login --client-id=foo --client-secret=bar`,
		Run: login,
	}

	cmd.Flags().StringP("client-id", "i", "", "client id (required)")
	cmd.Flags().StringP("client-secret", "s", "", "client secret (required)")
	cmd.Flags().StringP("grant-type", "g", "client_credentials", "grant type")
	cmd.Flags().StringP("username", "u", "", "user name")
	cmd.Flags().StringP("password", "p", "", "password")
	return cmd
}

func login(cmd *cobra.Command, args []string) {
	idsec := &credentials.ClientIdSecret{
		Id:     cli.GetCliStringFlag(cmd, "client-id"),
		Secret: cli.GetCliStringFlag(cmd, "client-secret"),
	}

	err := idsec.EnsureInput(false)
	if err != nil {
		check.ErrorExit(err, 1)
	}

	grant := cli.GetCliStringFlag(cmd, "grant-type")

	var p *authPayload
	if grant == "client_credentials" {
		p = &authPayload{
			ClientId:     idsec.Id,
			ClientSecret: idsec.Secret,
			GrantType:    grant,
		}
	}

	if grant == "password" {
		userpass := userPass(cmd)
		p = &authPayload{
			ClientId:     idsec.Id,
			ClientSecret: idsec.Secret,
			Username:     userpass.Username,
			Password:     userpass.Password,
			GrantType:    grant,
		}
	}

	// should not be nil when `grant_type` is valid
	if p == nil {
		check.ErrorExit("Invalid argument(s). See `help` for more information.", 1)
	}

	cnf := cli.ReadCliConfig()
	if cnf == nil {
		check.ErrorExit("read config failed", 1)
	}

	// we always overwrite 'runenv' during login
	runenv := cli.GetCliStringFlag(cmd, "runenv")
	if runenv == "" {
		runenv = "prod"
	}

	viper.Set("runenv", runenv)
	payload, err := json.Marshal(p)
	check.ErrorExit(err, 1)

	c := client.NewClient(client.NewApiConfig(cmd))
	token, err := c.GetAccessToken(payload)
	check.ErrorExit(err, 1)

	cnf.AccessToken = token
	cnf.RunEnv = runenv
	err = cnf.WriteToConfig()
	check.ErrorExit(err, 1)

	err = viper.ReadInConfig()
	check.ErrorExit(err, 1)
	d.Info("Login successful.")
}
