package main

import (
	"context"
	"testing"

	"github.com/ahmetb/grpcoin/server/firestoretestutil"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testUser struct {
	id   string
	name string
}

func (t testUser) DBKey() string       { return t.id }
func (t testUser) DisplayName() string { return t.name }
func (t testUser) ProfileURL() string  { return "https://" + t.name }

func TestGetUser_notFound(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	udb := &userDB{fs: firestoretestutil.StartEmulator(t, ctx)}
	tu := testUser{id: "foo"}

	u, ok, err := udb.get(ctx, tu)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("was not expecting to find user: %#v", u)
	}
}

func TestNewUser(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	udb := &userDB{fs: firestoretestutil.StartEmulator(t, ctx)}
	tu := testUser{id: "foobar", name: "ab"}

	err := udb.create(ctx, tu)
	if err != nil {
		t.Fatal(err)
	}

	uv, ok, err := udb.get(ctx, tu)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("not found created user")
	}

	expected := User{
		ID:          "foobar",
		DisplayName: "ab",
		ProfileURL:  "https://ab",
	}
	if diff := cmp.Diff(uv, expected,
		cmpopts.IgnoreFields(User{}, "CreatedAt")); diff != "" {
		t.Fatal(diff)
	}
}

func TestEnsureAccountExists(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	ctx := context.Background()
	udb := &userDB{fs: firestoretestutil.StartEmulator(t, ctx)}
	tu := testUser{id: "testuser", name: "abc"}

	u, err := udb.ensureAccountExists(ctx, tu)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == "" {
		t.Fatal("id should not be empty")
	}
	u2, err := udb.ensureAccountExists(ctx, tu)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(u, u2); diff != "" {
		t.Fatal(diff)
	}
}
