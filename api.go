package main

type DemoGreeting struct {
	Message  string
	Greeting Greeting
}

type Greeting struct {
	Message string
}

type HealthStatus struct {
	Status string `json:"status"`
}
