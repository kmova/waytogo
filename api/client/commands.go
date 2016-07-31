package client

// Command returns a cli command handler if one exists
func (cli *WaytogoCli) Command(name string) func(...string) error {
	return map[string]func(...string) error{
	}[name]
}
