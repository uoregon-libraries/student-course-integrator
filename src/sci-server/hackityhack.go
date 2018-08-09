// hackityhack.go (not to be confused with donttalkback.go) is absolutely not a
// terrible hack just to try and reduce the slowness on the first hit to the
// lambda-like person lookup service.

package sciserver

import (
	"time"

	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/service"
)

var lastWarmed time.Time

func warmCache() {
	if time.Since(lastWarmed) < time.Hour {
		return
	}
	go hackityhackWarmIt()
}

func hackityhackWarmIt() {
	logger.Debugf("Warming the person lookup cache")
	lastWarmed = time.Now()
	service.DuckID("nobody").Call()
	logger.Debugf("Warmed")
}
