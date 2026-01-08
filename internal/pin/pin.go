package pin

import "github.com/kernelle-soft/gimme/internal/log"

func Pin(repo string) {
	log.Print("Pinning repository \"{}\".", repo)
}

func Unpin(repo string) {
	log.Print("Unpinning repository \"{}\".", repo)
}