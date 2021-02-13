package shell

const (
	// Shell constants
	commandPrompt = "C:\\Windows\\System32\\cmd.exe"
	powerShell    = "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
)

func (s Shell) getSystemShellPath() string {
	if exists(powerShell) {
		return powerShell
	}
	return commandPrompt
}
