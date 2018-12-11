package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	defaultCommand = []string{"go", "test", "./..."}
)

type options struct {
	verbose bool
	quiet   bool
	safe    bool
	timeout time.Duration
	command []string
}

func setupFlags(args []string) *options {
	flags := pflag.NewFlagSet(args[0], pflag.ContinueOnError)
	flags.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	opts := options{}
	flags.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose")
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Quiet")
	flags.BoolVarP(&opts.safe, "safe", "s", true, "Safe mode, aka do not try to update master branches")
	flags.DurationVarP(&opts.timeout, "timeout", "t", 5*time.Minute, "comman timeout")
	// FIXME(vdemeester) git commit template message

	flags.SetInterspersed(false)
	flags.Usage = func() {
		out := os.Stderr
		fmt.Fprintf(out, "Usage:\n  %s [OPTIONS] COMMAND ARGS... \n\n", os.Args[0])
		fmt.Fprint(out, "Options:\n")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(2)
	}
	opts.command = flags.Args()
	return &opts
}

func main() {
	opts := setupFlags(os.Args)
	setupLogging(opts)

	run(opts)
}

func run(opts *options) {
	command := opts.command
	if len(command) == 0 {
		command = defaultCommand
	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	// Extract dependencies
	deps, err := extractDependencies(cwd, opts)
	if err != nil {
		log.Fatal(err)
	}
	errs := []error{}
	log.Debugf("Dependencies to update: %v", deps)
	for _, dep := range deps {
		updated, err := updateDependency(dep, opts)
		if err != nil {
			log.Fatal(err)
		}
		if len(updated) == 0 {
			log.Info("… no updates")
			continue
		} else {
			log.Info("… updates found")
			for _, d := range updated {
				log.Infof("… - %s : %s -> %s", d.Name, d.OldVersion, d.NewVersion)
			}
		}
		if err := runCommand(command, opts); err != nil {
			errs = append(errs, err)
			if ierr := cleanRepository(opts); ierr != nil {
				log.Fatal(ierr)
			}
			continue
		}
		if err := gitCommit(dep, updated, opts); err != nil {
			log.Fatal(err)
		}
	}
	if len(errs) != 0 {
		log.Error("There has been some errors while upgrading")
		for _, e := range errs {
			log.Error(e)
		}
		log.Fatal("§")
	}
}

func cleanRepository(opts *options) error {
	rmCommand := []string{"rm", "-fRv", "vendor/"}
	if err := execute(context.Background(), rmCommand, ioutil.Discard, opts); err != nil {
		return err
	}
	commitCommand := []string{"git", "reset", "--hard", "HEAD"}
	return execute(context.Background(), commitCommand, ioutil.Discard, opts)
}

func updateDependency(dep string, opts *options) ([]dependency, error) {
	log.Infof("Update dependency %s…", dep)
	updated := []dependency{}
	// (3/4) Wrote github.com/gorilla/websocket@v1.4.0: version changed (was v1.2.0)
	var buf = new(bytes.Buffer)
	if err := execute(context.Background(), []string{"dep", "ensure", "-v", "-update", dep}, buf, opts); err != nil {
		return updated, err
	}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile("^.*Wrote (.*)@(.*): version changed \\(was (.*)\\)$")
		if re.MatchString(line) {
			// Do sthg
			subs := re.FindStringSubmatch(line)
			if len(subs) != 4 {
				log.Warnf("subs, %v", subs)
				return []dependency{}, errors.New(fmt.Sprintf("Error while parsing: %s", line))
			}
			updated = append(updated, dependency{Name: subs[1], OldVersion: subs[3], NewVersion: subs[2]})
		}
	}
	return updated, nil
}

type dependency struct {
	Name       string
	OldVersion string
	NewVersion string
}

func runCommand(command []string, opts *options) error {
	log.Infof("… execute command `%s` (timeout %v)… ", strings.Join(command, " "), opts.timeout)
	ctx, cancel := context.WithTimeout(context.Background(), opts.timeout)
	defer cancel()

	return execute(ctx, command, ioutil.Discard, opts)
}

func execute(ctx context.Context, command []string, w io.Writer, opts *options) error {
	var buf io.Writer = os.Stderr
	if !opts.verbose {
		buf = bytes.NewBuffer([]byte{})
	}
	writer := io.MultiWriter(w, buf)
	c := exec.CommandContext(ctx, command[0], command[1:]...)
	c.Stderr = writer
	c.Stdout = writer
	c.Env = os.Environ()
	if err := c.Run(); err != nil {
		if !opts.verbose {
			fmt.Fprint(os.Stderr, buf)
		}
		log.Error(err)
		return err
	}
	return nil
}

func gitCommit(dep string, updated []dependency, opts *options) error {
	log.Infof("… commit dependency changes for %s…", dep)
	var oldversion, newversion string
	deps := []dependency{}
	for _, d := range updated {
		if d.Name == dep {
			oldversion = d.OldVersion
			newversion = d.NewVersion
			continue
		}
		deps = append(deps, d)
	}

	var message strings.Builder
	fmt.Fprintf(&message, "%s: %s -> %s\n", dep, oldversion, newversion)
	fmt.Fprintln(&message, "")
	if len(deps) > 0 {
		fmt.Fprintln(&message, "Additionnal updates…")
	}
	for _, d := range deps {
		fmt.Fprintf(&message, "\t%s: %s -> %s\n", d.Name, d.OldVersion, d.NewVersion)
	}
	addCommand := []string{"git", "add", "Gopkg.lock", "vendor"}
	if err := execute(context.Background(), addCommand, ioutil.Discard, opts); err != nil {
		return err
	}
	commitCommand := []string{"git", "commit", "-sS", "-m", message.String()}
	return execute(context.Background(), commitCommand, ioutil.Discard, opts)
}

func extractDependencies(path string, opts *options) ([]string, error) {
	log.Infof("Gather dependencies to update…")
	buf := bytes.NewBuffer([]byte{})
	deps := []string{}
	c := exec.Command("dep", "status", "-f", "{{.ProjectRoot}};{{.Constraint}}\n")
	c.Stderr = os.Stderr
	c.Stdout = buf
	c.Env = os.Environ()
	if err := c.Run(); err != nil {
		return deps, err
	}
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		log.Debugf("Process dependency", scanner.Text())
		components := strings.SplitN(scanner.Text(), ";", 2)
		dep := components[0]
		constraint := components[1]
		if opts.safe {
			if strings.HasPrefix(constraint, "branch") {
				log.Warnf("Skip dependency %s : constraint %s", dep, constraint)
				continue
			}
		}
		if strings.HasSuffix(constraint, "(override)") {
			log.Warnf("Skip dependency %s : has an override %s", dep, constraint)
			continue
		}
		deps = append(deps, dep)
	}
	return deps, nil
}

func setupLogging(opts *options) {
	formatter := new(prefixed.TextFormatter)
	formatter.DisableTimestamp = true
	log.SetFormatter(formatter)
	if opts.verbose {
		log.SetLevel(log.DebugLevel)
	}
	if opts.quiet {
		log.SetLevel(log.WarnLevel)
	}
}

// getEnv returns the last instance of an environment variable.
func getEnv(env []string, key string) string {
	for i := len(env) - 1; i >= 0; i-- {
		v := env[i]
		kv := strings.SplitN(v, "=", 2)
		if kv[0] == key {
			if len(kv) > 1 {
				return kv[1]
			}
			return ""
		}
	}
	return ""
}
