package async

type SyncFunc func(done chan error)

func SyncParallel(funcs ...SyncFunc) error {
	cpt := len(funcs)
	if cpt == 0 {
		return nil
	}
	done := make(chan error, cpt)
	for _, f := range funcs {
		go f(done)
	}

	var err error
	for ; cpt > 0; cpt-- {
		if e := <-done; e != nil {
			err = e
		}
	}

	close(done)
	done = nil

	return err
}
