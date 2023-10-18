package logging

import (
	"testing"
	"time"
)

func TestCreateESLogger(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	logger, err := CreateESLogger([]string{"http://localhost:9200"}, "alexis_system", "ff1cba5992265a8950164ece6b7c72ceb543188664597bf22d845750465e4d4c", "dragonfly-db-test", "test")
	if err != nil {
		t.Errorf("\nCreate Alexis logger failed\n    Error: %v", err)
		return
	}

	if logger.DebugLogger() == nil {
		t.Error("\nCreate Alexis logger failed\n    Error: failed to create debug logger")
		return
	}

	logger.Debug("logger test debug")

	if logger.InfoLogger() == nil {
		t.Error("\nCreate Alexis logger failed\n    Error: failed to create info logger")
		return
	}

	logger.Info("logger test info")

	if logger.WarnLogger() == nil {
		t.Error("\nCreate Alexis logger failed\n    Error: failed to create warn logger")
		return
	}

	logger.Warn("logger test warn")

	if logger.ErrorLogger() == nil {
		t.Error("\nCreate Alexis logger failed\n    Error: failed to create error logger")
		return
	}

	logger.Error("logger test error")

	// delete log file
	//err = os.Remove("log-test.log")
	//if err != nil {
	//	t.Errorf("\nCreate Alexis logger failed\n    Error: %v", err)
	//	return
	//}

	time.Sleep(time.Second * 20)

	t.Log("\nCreate Alexis logger succeeded")
}
