package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/jessevdk/go-flags"
	"gopkg.in/flosch/pongo2.v3"
)

// A Var is a key matched with its value
type Var struct {
	Key   string
	Value interface{}
}

// A Printer is a function that takes a variable and prints it to the given writer
type Printer func(io.Writer, Var) (int, error)

// Vars is a set of variable
type Vars []Var

// Print prints all values to a given writer using a supplied printer function
func (vars Vars) Print(w io.Writer, p Printer) (int, error) {
	var wr int
	for _, v := range vars {
		n, err := p(w, v)
		wr += n
		if err != nil {
			return wr, err
		}
	}

	return wr, nil
}

// EnvVars generats a set of environment variables for all information contained within a URL
func EnvVars(u *url.URL) Vars {
	pw, _ := u.User.Password()

	vars := Vars{
		{Key: "URL_SCHEME", Value: u.Scheme},
		{Key: "URL_HOST", Value: u.Host},
		{Key: "URL_HOSTNAME", Value: u.Hostname()},
		{Key: "URL_PORT", Value: u.Port()},
		{Key: "URL_USERNAME", Value: u.User.Username()},
		{Key: "URL_PASSWORD", Value: pw},
		{Key: "URL_URI", Value: u.RequestURI()},
		{Key: "URL_PATH", Value: u.Path},
		{Key: "URL_ESCAPED_PATH", Value: u.EscapedPath()},
		{Key: "URL_QUERY", Value: u.Query().Encode()},
		{Key: "URL_FRAGMENT", Value: u.Fragment},
	}

	// Append all query parameters as URL_QUERY_x
	for k := range u.Query() {
		vars = append(vars, Var{Key: fmt.Sprintf("URL_QUERY_%s", k), Value: u.Query().Get(k)})
	}

	return vars
}

// Options defines the command-line options for urlsplit
type Options struct {
	Export bool    `short:"e" long:"export" description:"Print URL variables as a set of export statements"`
	Key    *string `short:"k" long:"key" description:"Print the value of this key"`
	Format *string `short:"f" long:"format" description:"Render a template"`
	URL    struct {
		URL string `description:"The URL to parse"`
	} `positional-args:"yes" required:"yes"`
}

func Main() error {
	var options = new(Options)
	_, err := flags.Parse(options)
	if err != nil {
		return err
	}

	u, err := url.Parse(options.URL.URL)
	if err != nil {
		return err
	}

	var (
		w    = os.Stdout
		vars = EnvVars(u)
	)

	// When --export is enabled then print each variable as a set of export statements
	if options.Export {
		_, err = vars.Print(w, func(w io.Writer, v Var) (int, error) {
			return fmt.Fprintf(w, "export %q\n", fmt.Sprintf("%s=%s", v.Key, v.Value))
		})
		return err
	}

	// When --key is enabled then get a specific key from the variable set and print it.
	// When no key matches the request key an error is raised.
	if options.Key != nil {
		for _, v := range vars {
			if v.Key == *options.Key {
				_, err = fmt.Fprint(w, v.Value)
				return err
			}
		}

		return fmt.Errorf("No such URL component named %q", *options.Key)
	}

	// When --format is enabled then render a template with Pongo2
	if options.Format != nil {
		t, err := pongo2.FromString(*options.Format)
		if err != nil {
			return err
		}

		ctx := make(pongo2.Context, len(vars))
		for _, v := range vars {
			ctx[v.Key] = v.Value
		}

		return t.ExecuteWriter(ctx, w)
	}

	return fmt.Errorf("One of ('-e', '-k', '-f') must be specified")
}

func main() {
	err := Main()
	if err != nil {
		log.Fatal(err)
	}
}
