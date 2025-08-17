package util

import "go.uber.org/zap"

func ErrorFromChannel(chErr chan error) (err error) {
	close(chErr)
	select {
	case err, ok := <-chErr: // Error occured
		if ok {
			zap.S().Error(err)
			return err
		}
	default: // No error, moving on
	}

	return nil
}
