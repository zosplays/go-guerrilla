package backends

import (
	"github.com/flashmob/go-guerrilla/mail"
	"io"
	"strings"
	"time"
)

func init() {
	streamers["header"] = func() StreamDecorator {
		return *StreamHeader()
	}
}

type streamHeader struct {
	addHead []byte
	w       io.Writer
	i       int
}

func newStreamHeader(w io.Writer) *streamHeader {
	sc := new(streamHeader)
	sc.w = w
	return sc
}

func (sh *streamHeader) addHeader(e *mail.Envelope, config HeaderConfig) {
	to := strings.TrimSpace(e.RcptTo[0].User) + "@" + config.PrimaryHost
	hash := "unknown"
	if len(e.Hashes) > 0 {
		hash = e.Hashes[0]
	}
	var addHead string
	addHead += "Delivered-To: " + to + "\n"
	addHead += "Received: from " + e.Helo + " (" + e.Helo + "  [" + e.RemoteIP + "])\n"
	if len(e.RcptTo) > 0 {
		addHead += "	by " + e.RcptTo[0].Host + " with SMTP id " + hash + "@" + e.RcptTo[0].Host + ";\n"
	}
	addHead += "	" + time.Now().Format(time.RFC1123Z) + "\n"

	sh.addHead = []byte(addHead)
}

func (sh *streamHeader) Write(p []byte) (n int, err error) {
	if sh.i < len(sh.addHead) {
		for {
			if N, err := sh.w.Write(sh.addHead[sh.i:]); err != nil {
				return N, err
			} else {
				sh.i += N
				if sh.i >= len(sh.addHead) {
					break
				}
			}
		}
	}
	return sh.w.Write(p)
}

func StreamHeader() *StreamDecorator {
	sd := &StreamDecorator{}
	sd.p =

		func(sp StreamProcessor) StreamProcessor {
			var dc *streamHeader
			x := 1 + 5
			_ = x
			sd.Open = func(e *mail.Envelope) error {
				dc = newStreamHeader(sp)
				hc := HeaderConfig{"sharklasers.com"}
				dc.addHeader(e, hc)
				return nil
			}
			return StreamProcessWith(func(p []byte) (int, error) {

				return sp.Write(p)
			})
		}

		/*
			func(sp StreamProcessor) StreamProcessor {
				var dc *streamHeader

				sd.Open = func(e *mail.Envelope) error {
					dc = newStreamHeader(sp)
					hc := HeaderConfig{"sharklasers.com"}
					dc.addHeader(e, hc)
					return nil
				}

				return StreamProcessWith(dc.Write)


			}
		*/
	return sd
}
