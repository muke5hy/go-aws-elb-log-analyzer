package main

import (
	"net"
	"regexp"
	"strconv"
	"time"
)

var matcherClasic []string
var matcherApp []string
var reClasic *regexp.Regexp
var reApp *regexp.Regexp

func init() {
	reClasic = regexp.MustCompile(`(?P<date>[^Z]+Z) (?P<elb>[^\s]+) (?P<ipclient>[^:]+?):[0-9]+ (?P<ipnode>[^:]+?):[0-9]+ (?P<reqtime>[0-9\.]+) (?P<backtime>[0-9\.]+) (?P<restime>[0-9\.]+) (?P<elbcode>[0-9]{3}) (?P<backcode>[0-9]{3}) (?P<lenght1>[0-9]+) (?P<lenght2>[0-9]+) "(?P<Method>[A-Z]+) (?P<URL>[^"]+) HTTP/[0-9\.]+".*`)
	matcherClasic = reClasic.SubexpNames()

	reApp = regexp.MustCompile(`(?P<proto>[^\s]+) (?P<date>[^Z]+Z) (?P<target>[^\s]+) (?P<ipclient>[^:]+?):[0-9]+ (?P<ipnode>[^:]+?):[0-9]+ (?P<reqtime>[0-9\.]+) (?P<backtime>[-0-9\.]+) (?P<restime>[-0-9\.]+) (?P<elbcode>[0-9]{3}) (?P<backcode>([0-9]{3}|-)) (?P<lenght1>[0-9]+) (?P<lenght2>[0-9]+) "(?P<Method>[A-Z]+) (?P<URL>[^"]+) HTTP/[0-9\.]+".*`)
	matcherApp = reApp.SubexpNames()
}

// LineLog is the struct to analyce and store a line
type LineLog struct {
	URL      string
	Time     time.Time
	Elapsed  float64
	Method   string
	IPClient net.IP
	Seek     int64
	Filelog  *FileLog
	Len      int
}

// NewLineLog create an structure of anlyzce record
func NewLineLog(raw string, filelog *FileLog, seek int64) *LineLog {

	line := &LineLog{
		Filelog: filelog,
		Seek:    seek,
		Len:     len([]byte(raw)),
	}

	line.parse(raw)

	return line
}

// parse the raw record with regular expresion and store in the struct
func (l *LineLog) parse(raw string) {

	var r [][]string
	var matcher []string

	if l.Filelog.ElbType == "classic" {
		r = reClasic.FindAllStringSubmatch(raw, -1)
		matcher = matcherClasic
	} else {
		r = reApp.FindAllStringSubmatch(raw, -1)
		matcher = matcherApp
	}

	if r == nil {
		return
	}

	for i, n := range r[0] {
		switch matcher[i] {
		case "date":
			l.Time, _ = time.Parse(time.RFC3339Nano, n)
			break
		case "ipclient":
			if analyze {
				l.IPClient = net.ParseIP(n)
			}
			break
		case "reqtime":
			if analyze {
				_time, _ := strconv.ParseFloat(n, 64)
				l.Elapsed = l.Elapsed + _time
			}
			break
		case "backtime":
			if analyze {
				_time, _ := strconv.ParseFloat(n, 64)
				l.Elapsed = l.Elapsed + _time
			}
			break
		case "restime":
			if analyze {
				_time, _ := strconv.ParseFloat(n, 64)
				l.Elapsed = l.Elapsed + _time
			}
			break
		}

	}

}
