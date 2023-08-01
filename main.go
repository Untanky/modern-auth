package main

type Person struct {
	Name string
}

func (p Person) String() string {
	return p.Name
}

func main() {
	app := App{}
	app.Start()
	defer app.Stop()
}
