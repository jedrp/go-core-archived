package pllog

type DefaultLogger struct {
}

func (log *DefaultLogger) Trace(args ...interface{}) {

}
func (log *DefaultLogger) Debug(args ...interface{}) {

}
func (log *DefaultLogger) Info(args ...interface{}) {

}
func (log *DefaultLogger) Warn(args ...interface{}) {

}
func (log *DefaultLogger) Error(args ...interface{}) {

}
func (log *DefaultLogger) Fatal(args ...interface{}) {

}
func (log *DefaultLogger) Panic(args ...interface{}) {

}
func (log *DefaultLogger) Tracef(format string, args ...interface{}) {

}
func (log *DefaultLogger) Debugf(format string, args ...interface{}) {

}
func (log *DefaultLogger) Infof(format string, args ...interface{}) {

}
func (log *DefaultLogger) Warnf(format string, args ...interface{}) {

}
func (log *DefaultLogger) Errorf(format string, args ...interface{}) {

}
func (log *DefaultLogger) Fatalf(format string, args ...interface{}) {

}
func (log *DefaultLogger) Panicf(format string, args ...interface{}) {

}

//WithFields ...
func (log *DefaultLogger) WithFields(map[string]interface{}) PlLogentry {
	return log
}
