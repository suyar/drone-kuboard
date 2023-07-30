package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type payload struct {
	Kind      string            `json:"kind"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Images    map[string]string `json:"images"`
}

type errorResponse struct {
	Reason string
}

func main() {
	app := &cli.App{
		Name:  "drone-kuboard",
		Usage: "Update Kuboard Workloads Image Tag",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "suyar",
				Email: "su@zorzz.com",
			},
		},
		Action: runPlugin,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cluster",
				Usage:   "cluster name",
				EnvVars: []string{"PLUGIN_CLUSTER"},
			},
			&cli.StringFlag{
				Name:    "kind",
				Usage:   "workload type",
				EnvVars: []string{"PLUGIN_KIND"},
			},
			&cli.StringFlag{
				Name:    "name",
				Usage:   "workload name",
				EnvVars: []string{"PLUGIN_NAME"},
			},
			&cli.StringFlag{
				Name:    "namespace",
				Usage:   "workload namespace",
				EnvVars: []string{"PLUGIN_NAMESPACE"},
			},
			&cli.StringFlag{
				Name:    "image",
				Usage:   "image uri",
				EnvVars: []string{"PLUGIN_IMAGE"},
			},
			&cli.StringFlag{
				Name:    "tag",
				Usage:   "image tag",
				EnvVars: []string{"PLUGIN_TAG"},
			},
			&cli.StringFlag{
				Name:    "kuboard_uri",
				Usage:   "kuboard uri",
				EnvVars: []string{"PLUGIN_KUBOARD_URI"},
			},
			&cli.StringFlag{
				Name:    "kuboard_username",
				Usage:   "kuboard username",
				EnvVars: []string{"PLUGIN_KUBOARD_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "kuboard_key",
				Usage:   "kuboard access key",
				EnvVars: []string{"PLUGIN_KUBOARD_KEY"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	log.Println("The image version was successfully updated.")
}

func runPlugin(ctx *cli.Context) error {
	cluster := ctx.Value("cluster").(string)
	uri := strings.Trim(ctx.Value("kuboard_uri").(string), "/")
	uri = uri + "/kuboard-api/cluster/" + cluster + "/kind/CICDApi/admin/resource/updateImageTag"
	image := strings.Trim(ctx.Value("image").(string), "/")
	tag := ctx.Value("tag").(string)

	data, err := json.Marshal(payload{
		Name:      ctx.Value("name").(string),
		Namespace: ctx.Value("namespace").(string),
		Kind:      ctx.Value("kind").(string),
		Images: map[string]string{
			image: image + ":" + tag,
		},
	})

	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", uri, bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{
		Name:  "KuboardUsername",
		Value: ctx.Value("kuboard_username").(string),
	})
	request.AddCookie(&http.Cookie{
		Name:  "KuboardAccessKey",
		Value: ctx.Value("kuboard_key").(string),
	})

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		errResp := &errorResponse{}
		err := json.Unmarshal(body, errResp)

		if err != nil {
			return errors.New(response.Status)
		} else {
			return errors.New(errResp.Reason)
		}
	}

	return err
}
