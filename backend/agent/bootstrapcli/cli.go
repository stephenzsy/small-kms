package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/stephenzsy/small-kms/backend/agent/bootstrap"
	"github.com/urfave/cli/v2"
)

var BuildID = "dev"

func main() {
	app := &cli.App{
		Name: "agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env-file",
				Usage: "path to the .env file",
				Action: func(ctx *cli.Context, s string) error {
					if s != "" {
						return godotenv.Load(s)
					}
					return nil
				},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "bootstrap",
				Usage: "bootstrap the agent",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "token-cache-file",
						Usage: "path to the token cache file",
						Value: "./token-cache.json",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "login",
						Usage: "Admin login to initialize agent",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "device-code",
								Usage: "use device code to login",
								Value: false,
							},
						},
						Action: func(c *cli.Context) error {
							cacheFilePath := c.String("token-cache-file")
							return bootstrap.NewServicePrincipalBootstraper().Login(c.Context, cacheFilePath, c.Bool("device-code"))
						},
					},
					{
						Name:  "service-principal",
						Usage: "bootstrap service principal for agent",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "cert-path",
								Usage:   "path to the client cert",
								EnvVars: []string{"AZURE_CLIENT_CERTIFICATE_PATH", "CLIENT_CERTIFICATE_PATH"},
								Value:   "./sp-client-cert.pem",
							},
							&cli.StringFlag{
								Name:     "cert-policy-id",
								Usage:    "policy identifier for the client cert",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return bootstrap.NewServicePrincipalBootstraper().Bootstrap(c.Context,
								c.String("cert-policy-id"),
								c.String("cert-path"), c.String("token-cache-file"))
						},
					},
					{
						Name:  "active-server",
						Usage: "bootstrap active server for agent",
						Action: func(c *cli.Context) error {
							return bootstrap.NewServicePrincipalBootstraper().BootstarpActiveServer(c.Context)
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
