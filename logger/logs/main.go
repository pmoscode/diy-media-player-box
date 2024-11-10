package logs

func New(filename string, filesize int) *Log {
	logger := &Log{
		logRotationPeriodAfterBytes: filesize,
		filename:                    filename,
		logFile:                     nil,
		byteCounter:                 0,
	}

	logger.init()

	return logger
}
