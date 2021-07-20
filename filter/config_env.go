package filter

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Make os.LookupEnv mockable for unit test.
var _LookupEnv = os.LookupEnv //nolint: gochecknoglobals

//NewFromEnv instantiate a Filter object from the environment variable configuration.
//nolint: funlen
func NewFromEnv(upLog logrus.FieldLogger) (string, uint8, *Filter) {
	var ok bool

	var c Config

	var from, to, restricteds string

	urls := []string{}

	url, ok := _LookupEnv("VILLIP_URL")
	if !ok {
		upLog.Fatal("Missing VILLIP_URL environment variable")
	}

	c.URL = url

	if villipPriority, ok := _LookupEnv("VILLIP_PRIORITY"); ok {
		priority, err := strconv.Atoi(villipPriority)
		if err != nil {
			upLog.Fatalf("%s is not a valid priority", villipPriority)
		}

		if priority < 0 || priority > 255 {
			upLog.Fatalf("%s is not a valid priority", villipPriority)
		}

		c.Priority = uint8(priority)
	}

	villipPort, _ := _LookupEnv("VILLIP_PORT")

	if villipPort == "" {
		villipPort = "8080"
	}

	port, err := strconv.Atoi(villipPort)
	if err != nil {
		upLog.Fatalf("%s is not a valid TCP port", villipPort)
	}

	c.Port = port

	c.Force = false
	if _, ok := _LookupEnv("VILLIP_FORCE"); ok {
		c.Force = true
	}

	if _, ok := _LookupEnv("VILLIP_INSECURE"); ok {
		c.Insecure = true
	}

	if dumpFolder, ok := _LookupEnv("VILLIP_DUMPFOLDER"); ok {
		c.Dump.Folder = dumpFolder
	}

	c.Replace = make([]Creplacement, 0)
	c.Request.Header = make([]Cheader, 0)
	c.Response.Header = make([]Cheader, 0)
	c.Request.Replace = make([]Creplacement, 0)
	c.Response.Replace = make([]Creplacement, 0)

	if from, ok = _LookupEnv("VILLIP_FROM"); ok {
		if to, ok = _LookupEnv("VILLIP_TO"); !ok {
			upLog.Fatal("Missing VILLIP_TO environment variable")
		}

		if urlList, ok := _LookupEnv("VILLIP_FOR"); ok {
			urls = strings.Split(strings.Replace(urlList, " ", "", -1), ",")
		}

		c.Response.Replace = append(c.Response.Replace, Creplacement{From: from, To: to, Urls: urls})
	}

	if restricteds, ok = _LookupEnv("VILLIP_RESTRICTED"); ok {
		c.Restricted = strings.Split(strings.Replace(restricteds, " ", "", -1), ",")
	}

	i := 1

	for {
		from, ok = _LookupEnv(fmt.Sprintf("VILLIP_FROM_%d", i))
		if !ok {
			break
		}

		to, ok = _LookupEnv(fmt.Sprintf("VILLIP_TO_%d", i))
		if !ok {
			upLog.Fatalf("Missing VILLIP_TO_%d environment variable", i)
		}

		urls = []string{}
		if urlList, ok := _LookupEnv(fmt.Sprintf("VILLIP_FOR_%d", i)); ok {
			urls = strings.Split(strings.Replace(urlList, " ", "", -1), ",")
		}

		c.Response.Replace = append(c.Response.Replace, Creplacement{From: from, To: to, Urls: urls})
		i++
	}

	if contenttypes, ok := _LookupEnv("VILLIP_TYPES"); ok {
		c.ContentTypes = strings.Split(strings.Replace(contenttypes, " ", "", -1), ",")
	}

	if dumpURLs, ok := _LookupEnv("VILLIP_DUMPURLS"); ok {
		c.Dump.URLs = strings.Split(strings.Replace(dumpURLs, " ", "", -1), ",")
	}

	return _newFromConfig(upLog, c)
}
