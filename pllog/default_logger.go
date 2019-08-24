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

//WithFields ...
func (log *DefaultLogger) WithFields(map[string]interface{}) PlLogentry {
	return log
}
