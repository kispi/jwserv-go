package services

import (
	"os"

	logging "github.com/op/go-logging"
)

var Log = logging.MustGetLogger("kispi")

func InitLogger() {
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{color}%{time:2006/1/2 15:04:05} %{level:.1s} %{message} %{color:reset}@%{shortfile}`,
	)

	backend_formatter := logging.NewBackendFormatter(backend, format)
	backend_leveled := logging.AddModuleLevel(backend_formatter)

	// if SConfig.Debug {
	backend_leveled.SetLevel(logging.DEBUG, "")
	// } else {
	// backend_leveled.SetLevel(logging.INFO, "")
	// }

	logging.SetBackend(backend_leveled)
}
