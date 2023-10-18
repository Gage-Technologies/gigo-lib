package logging

//func TestCreateAlexisLogger(t *testing.T) {
//	logger, err := CreateAlexisLogger("log-test.log")
//	if err != nil {
//		t.Errorf("\nCreate Alexis logger failed\n    Error: %v", err)
//		return
//	}
//
//	if logger.Debug() == nil {
//		t.Error("\nCreate Alexis logger failed\n    Error: failed to create debug logger")
//		return
//	}
//
//	if logger.Info() == nil {
//		t.Error("\nCreate Alexis logger failed\n    Error: failed to create info logger")
//		return
//	}
//
//	if logger.Warn() == nil {
//		t.Error("\nCreate Alexis logger failed\n    Error: failed to create warn logger")
//		return
//	}
//
//	if logger.Error() == nil {
//		t.Error("\nCreate Alexis logger failed\n    Error: failed to create error logger")
//		return
//	}
//
//	// delete log file
//	err = os.Remove("log-test.log")
//	if err != nil {
//		t.Errorf("\nCreate Alexis logger failed\n    Error: %v", err)
//		return
//	}
//
//	t.Log("\nCreate Alexis logger succeeded")
//}
