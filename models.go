package main

type Location struct {
	Name, Host, Port, Basepath string
	Auth                       Auth
}
type Auth struct {
	Username, Password string
}
