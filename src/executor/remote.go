package executor

type Remote struct {
	Name string
	Url  string
}

func ListRemote() []Remote {
	return []Remote{
		{
			Name: "origin",
			Url:  "https://github.com/julien040/gut-essai.git",
		},
	}
}

func AddRemote() {
}

func RemoveRemote() {
}
