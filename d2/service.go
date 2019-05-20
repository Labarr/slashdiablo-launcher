package d2

import (
	"fmt"
	"io"
	"os"

	"github.com/nokka/slash-launcher/github"
)

// Service is responible for all things related to Diablo II.
type Service struct {
	path          string
	githubService github.Service
	//store         storage.Store
}

// Exec will exec the Diablo 2.
func (s *Service) Exec() error {
	/*config, err := s.store.Read()
	if err != nil {
		return err
	}

	fmt.Println("LAUNCH on path", s.path)
	fmt.Println("WITH INSTANCES", config.D2Instances)*/

	return nil
}

// Patch will check for updates and if found, patch the game.
func (s *Service) Patch() <-chan float32 {
	progress := make(chan float32)

	go func() {
		defer close(progress)

		// Create the file, but give it a tmp file extension, this means we won't overwrite a
		// file until it's downloaded, but we'll remove the tmp extension once downloaded.
		out, err := os.Create(s.path + "/test.tmp")
		if err != nil {
			fmt.Println(err)
		}

		defer out.Close()

		contents, err := s.githubService.GetFile("Patch_D2.mpq")
		if err != nil {
			fmt.Println(err)
			return
		}

		// Create our progress reporter and pass it to be used alongside our writer
		counter := &WriteCounter{
			Total:    float32(2108703),
			progress: progress,
		}

		_, err = io.Copy(out, io.TeeReader(contents, counter))
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	return progress

}

// NewService returns a service with all the dependencies.
func NewService(path string, githubService github.Service) *Service {
	return &Service{
		path:          path,
		githubService: githubService,
	}
}