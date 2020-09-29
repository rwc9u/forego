package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestParseConcurrencyFlagEmpty(t *testing.T) {
	c, err := parseConcurrency("")
	if err != nil {
		t.Fatal(err)
	}
	if len(c) > 0 {
		t.Fatal("expected no concurrency settings with ''")
	}
}

func TestParseConcurrencyFlagSimle(t *testing.T) {
	c, err := parseConcurrency("foo=2")
	if err != nil {
		t.Fatal(err)
	}

	if len(c) != 1 {
		t.Fatal("expected 1 concurrency settings with 'foo=2'")
	}

	if c["foo"] != 2 {
		t.Fatal("expected concurrency settings of 2 with 'foo=2'")
	}
}

func TestParseConcurrencyFlagMultiple(t *testing.T) {
	c, err := parseConcurrency("foo=2,bar=3")
	if err != nil {
		t.Fatal(err)
	}

	if len(c) != 2 {
		t.Fatal("expected 1 concurrency settings with 'foo=2'")
	}

	if c["foo"] != 2 {
		t.Fatal("expected concurrency settings of 2 with 'foo=2'")
	}

	if c["bar"] != 3 {
		t.Fatal("expected concurrency settings of 3 with 'bar=3'")
	}
}

func TestParseConcurrencyFlagNonInt(t *testing.T) {
	_, err := parseConcurrency("foo=x")
	if err == nil {
		t.Fatal("foo=x should fail")
	}
}

func TestParseConcurrencyFlagWhitespace(t *testing.T) {
	c, err := parseConcurrency("foo   =   2, bar = 3")
	if err != nil {
		t.Fatalf("foo   =   2, bar = 4 should not fail:%s", err)
	}

	if len(c) != 2 {
		t.Fatal("expected 1 concurrency settings with 'foo=2'")
	}

	if c["foo"] != 2 {
		t.Fatal("expected concurrency settings of 2 with 'foo=2'")
	}

	if c["bar"] != 3 {
		t.Fatal("expected concurrency settings of 3 with 'bar=3'")
	}
}

func TestParseConcurrencyFlagMultipleEquals(t *testing.T) {
	_, err := parseConcurrency("foo===2")
	if err == nil {
		t.Fatalf("foo===2 should fail: %s", err)
	}
}

func TestParseConcurrencyFlagNoValue(t *testing.T) {
	_, err := parseConcurrency("foo=")
	if err == nil {
		t.Fatalf("foo= should fail: %s", err)
	}

	_, err = parseConcurrency("=")
	if err == nil {
		t.Fatalf("= should fail: %s", err)
	}

	_, err = parseConcurrency("=1")
	if err == nil {
		t.Fatalf("= should fail: %s", err)
	}

	_, err = parseConcurrency(",")
	if err == nil {
		t.Fatalf(", should fail: %s", err)
	}

	_, err = parseConcurrency(",,,")
	if err == nil {
		t.Fatalf(",,, should fail: %s", err)
	}

}

func TestPortFromEnv(t *testing.T) {
	env := make(Env)
	port, err := basePort(env)
	if err != nil {
		t.Fatalf("Can not get base port: %s", err)
	}
	if port != 5000 {
		t.Fatal("Base port should be 5000")
	}

	os.Setenv("PORT", "4000")
	port, err = basePort(env)
	if err != nil {
		t.Fatalf("Can not get port: %s", err)
	}
	if port != 4000 {
		t.Fatal("Base port should be 4000")
	}

	env["PORT"] = "6000"
	port, err = basePort(env)
	if err != nil {
		t.Fatalf("Can not get base port: %s", err)
	}
	if port != 6000 {
		t.Fatal("Base port should be 6000")
	}

	env["PORT"] = "forego"
	port, err = basePort(env)
	if err == nil {
		t.Fatalf("Port 'forego' should fail: %s", err)
	}

}

func TestConfigBeOverrideByForegoFile(t *testing.T) {
	var procfile = "Profile"
	var port = 5000
	var concurrency string = "web=2"
	var gracetime int = 3
	err := readConfigFile("./fixtures/configs/.forego", &procfile, &port, &concurrency, &gracetime)

	if err != nil {
		t.Fatalf("Cannot set default values from forego config file")
	}

	if procfile != "Procfile.dev" {
		t.Fatal("Procfile should be Procfile.dev")
	}

	if port != 15000 {
		t.Fatalf("port should be 15000, got %d", port)
	}

	if concurrency != "foo=2,bar=3" {
		t.Fatalf("concurrency should be 'foo=2,bar=3', got %s", concurrency)
	}

	if gracetime != 30 {
		t.Fatalf("gracetime should be 3, got %d", gracetime)
	}
}

func TestStartProcess(t *testing.T) {
	of := NewOutletFactory()
	of.LeftFormatter = fmt.Sprintf("%%-%ds | ", 20)
	procFileEntry := ProcfileEntry{Name: "testproc", Command: "sleep 1"}
	env := Env{}
	ctx, cancel := context.WithCancel(context.Background())

	f := &Forego{
		outletFactory:  of,
		teardown:       ctx,
		teardownCancel: cancel,
	}

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	go f.monitorInterrupt()

	f.startProcess(5000, 0, 0, procFileEntry, env, of)
	f.startProcess(5000, 0, 1, procFileEntry, env, of)

	<-f.teardown.Done()
	f.wg.Wait()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	match, _ := regexp.MatchString("(?is)port 5000.*port 5001", string(out))

	if match != true {
		t.Fatalf("Did not see consecutive ports for concurrent processes, got %s", out)
	}
}
