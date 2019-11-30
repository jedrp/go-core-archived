package pllog

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	ElasticHostURL string `long:"log-host-url" description:"the url of elastichsearch database url" env:"LOG_HOST_URL"`
	Sniff          bool   `long:"log-enable-sniff" description:"Enable or disable sniff" env:"LOG_SNIFF"`
	LogIndexPrefix string `long:"log-prefix" description:"the prefix of index name" env:"LOG_INDEX_PREFIX"`
	LogHostName    string `long:"log-host-name" description:"the prefix of index name" env:"LOG_HOST_NAME"`
	Enable         bool   `long:"log-enable" description:"the prefix of index name" env:"LOG_ENABLE"`
	LogLevel       string `long:"log-level" description:"the prefix of index name" env:"LOG_LEVEL"`
	IndexNameFunc  func() string
	*logrus.Logger
}

func New() PlLogger {
	logrusLogger := &LogrusLogger{
		LogLevel: "debug",
	}

	logrusLogger.IndexNameFunc = func() string {
		dt := time.Now()
		return fmt.Sprintf("%s-%s", logrusLogger.LogIndexPrefix, dt.Format("2006-01-02"))
	}
	parser := flags.NewParser(logrusLogger, flags.IgnoreUnknown)
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
	if !logrusLogger.Enable {
		return &DefaultLogger{}
	}
	return NewWithRef(logrusLogger)
}

func NewWithRef(logrusLogger *LogrusLogger) PlLogger {
	if !logrusLogger.Enable {
		return &DefaultLogger{}
	}
	log := logrus.New()
	level, err := logrus.ParseLevel(logrusLogger.LogLevel)

	if err != nil {
		log.Panic(err)
	}
	log.Level = level

	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(logrusLogger.ElasticHostURL))
	if err != nil {
		log.Panic(err)
	}

	hook, err := NewElasticHookWithFunc(client, logrusLogger.LogHostName, level, logrusLogger.IndexNameFunc)
	if err != nil {
		log.Panic(err)
	}
	log.Hooks.Add(hook)

	logrusLogger.Logger = log
	log.Printf("%+v\n", logrusLogger)
	return logrusLogger
}

func (logrusLogger *LogrusLogger) WithFields(fields map[string]interface{}) PlLogentry {
	return logrusLogger.Logger.WithFields(fields)
}

func NewEntry(logger *LogrusLogger) *logrus.Entry {
	return &logrus.Entry{
		Logger: logger.Logger,
		// Default is three fields, plus one optional.  Give a little extra room.
		Data: make(logrus.Fields, 6),
	}
}
