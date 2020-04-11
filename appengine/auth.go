package main

type Secret struct {
	TasksSecret string `datastore:",noindex"`
}
