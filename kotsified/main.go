package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var clientset *kubernetes.Clientset

var namespace string

var pool *pgxpool.Pool

type data struct {
	Broken      bool
	Clicks      int
	Env         string
	LicenseData license
	LicenseText string
	Logo        string
}

type license struct {
	Assignee       string         `json:"assignee"`
	Fields         []licenseField `json:"fields"`
	InstallationID string         `json:"installation_id"`
	LicenseID      string         `json:"license_id"`
	ReleaseChannel string         `json:"release_channel"`
}

type licenseField struct {
	Field            string      `json:"field"`
	HideFromCustomer bool        `json:"hide_from_customer"`
	Title            string      `json:"title"`
	Value            interface{} `json:"value"`
}

func (d data) GetPods() string {
	var output string

	if clientset != nil {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{
			LabelSelector: "app=multiprem",
		})
		if err != nil {
			log.Print("unable to get pods: ", err)
		}

		for _, pod := range pods.Items {
			output += fmt.Sprintf("%s: %s\n", pod.Name, pod.Status.Phase)
		}
	}

	return output
}

func (d data) GetStatefulSet() string {
	var output string

	if clientset != nil {
		sts, err := clientset.AppsV1().StatefulSets(namespace).Get(context.Background(), "multiprem-postgres", v1.GetOptions{})
		if err != nil {
			log.Print("unable to get statefulset: ", err)
		}

		output = fmt.Sprintf("%s: %d", sts.Name, *sts.Spec.Replicas)
	}

	return output
}

func main() {
	var err error

	ctx := context.Background()

	d := data{
		Broken: false,
		Env:    strings.Join(os.Environ(), "\n"),
		LicenseData: license{
			Assignee: "Replicated",
			Fields: []licenseField{
				{
					Value: 2,
				},
			},
		},
		LicenseText: "Unable to retrieve license data",
		Logo:        "https://www.replicated.com/wp-content/uploads/2020/12/logo-mark.png",
	}

	if os.Getenv("MULTIPREM_LOGO") != "" {
		d.Logo = os.Getenv("MULTIPREM_LOGO")
		log.Print("found custom logo")
	}

	log.Print("super secret info: HIDEME")

	// Connect to database
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOSTNAME"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DATABASE"))

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Print("unable to parse database config: ", err)
	}

	pool, err = pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Print("unable to connect to database: ", err)
	}

	log.Print("connected to postgres")

	if _, err = pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS click (
	  button text primary key
	, clicks integer not null default 0
);
`); err != nil {
		log.Print("unable to create database schema: ", err)
	}

	log.Print("created database schema")

	// Get license info
	r, err := http.Get("http://kotsadm:3000/license/v1/license")
	if err == nil {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			d.LicenseText = string(body)

			err := json.Unmarshal([]byte(d.LicenseText), &d.LicenseData)
			if err != nil {
				log.Print("couldn't parse license: ", err)
			} else {
				log.Print("found license data")
			}
		} else {
			log.Print("couldn't parse license: ", err)
		}
	} else {
		log.Print("couldn't retrieve license: ", err)
	}

	// Get kubernetes info
	kubeconfig, err := rest.InClusterConfig()
	if err != nil {
		log.Print("unable to config from service account: ", err)
	} else {
		log.Print("found service account config")

		clientset, err = kubernetes.NewForConfig(kubeconfig)
		if err != nil {
			log.Print("unable to create clientset: ", err)
		} else {
			log.Print("created clientset")

			ns, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
			if err != nil {
				log.Print("unable to read namespace name: ", err)
			} else {
				namespace = string(ns)
			}
		}
	}

	tmpl := template.Must(template.New("template").Parse(`
<!DOCTYPE html>
<html lang="en">
	<head>
		<link rel="icon" href="{{ .Logo }}">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.9.3/css/bulma.min.css">
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Multi-Prem</title>
	</head>
	<body>
		<section class="section">
			<div class="is-flex is-align-items-center is-justify-content-space-between">
				<div class="is-flex is-align-items-center">
					<figure class="pr-3">
		  			<img src="{{ .Logo }}" width="100">
					</figure>
					<div class="is-flex is-flex-direction-column">
						<h1 class="title">Multi-Prem</h1>
						<h2 class="subtitle">Hello {{ .LicenseData.Assignee }}!</h1>
						<h2 class="subtitle">You are licensed for {{ (index .LicenseData.Fields 1).Value }} replicas</h1>
					</div>
				</div>
				<div>
					<button class="button is-info" onclick="fetch('/increase').then(() => { window.location.reload() })">Click Me!</button>
					{{ if .Broken }}
					<p>Something isn't working!!</p>
					{{ else }}
					<p>{{ .Clicks }} Click{{if ne .Clicks 1}}s{{end}}</p>
					{{ end }}
				</div>
				<div>
					<button class="button is-danger" onclick="fetch('/breakapp')">Don't Click Me!</button>
				</div>
			</div>
			<h2 class="title is-5 pt-5">License Details</h2>
			<pre>{{ .LicenseText }}</pre>
			<h2 class="title is-5 pt-5">Pod Status</h2>
			<pre>{{ .GetPods }}</pre>
			<h2 class="title is-5 pt-5">Stateful Set Replicas</h2>
			<pre>{{ .GetStatefulSet }}</pre>
		</section>
	</body>
</html>
`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := d

		if err := pool.QueryRow(r.Context(), "SELECT clicks FROM click WHERE button = 'me'").Scan(&t.Clicks); err != nil {
			if err != pgx.ErrNoRows {
				t.Broken = true
				log.Print("cannot query postgres: ", err)
			}
		}

		if err := tmpl.Execute(w, t); err != nil {
			log.Print("cannot render template: ", err)
		}
	})

	http.HandleFunc("/breakapp", func(w http.ResponseWriter, r *http.Request) {
		if clientset != nil {
			_, err := clientset.AppsV1().StatefulSets(string(namespace)).Patch(context.Background(), "multiprem-postgres", types.MergePatchType, []byte(`{"spec": {"replicas": 0}}`), v1.PatchOptions{})
			if err != nil {
				log.Print("unable to delete database: ", err)
				return
			}

			log.Print("oh no, customer broke the app")
		}
	})

	http.HandleFunc("/increase", func(w http.ResponseWriter, r *http.Request) {
		if _, err := pool.Exec(r.Context(), `
INSERT INTO click AS original (
	  button
	, clicks
) VALUES (
	  'me'
	, 1
) ON CONFLICT (button)
DO
	UPDATE SET clicks = original.clicks + 1 WHERE original.button = 'me'
`); err != nil {
			d.Broken = true
			log.Print("cannot update postgres: ", err)
		}

		log.Print("incremented counter")
	})

	http.ListenAndServe(":3000", nil) // nolint errcheck
	log.Print("started http listener")
}
