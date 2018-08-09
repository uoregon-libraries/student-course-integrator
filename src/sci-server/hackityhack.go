// hackityhack.go (not to be confused with donttalkback.go) is absolutely not a
// terrible hack just to try and reduce the slowness on the first hit to the
// lambda-like person lookup service.

package sciserver

import (
	"sync"
	"time"

	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/service"
)

var lastWarmed time.Time
var m sync.Mutex

func warmCache() {
	m.Lock()
	defer m.Unlock()

	if time.Since(lastWarmed) < time.Hour {
		return
	}
	lastWarmed = time.Now()
	go hackityhackWarmIt()
}

func hackityhackWarmIt() {
	logger.Debugf("Warming the person lookup cache")
	service.DuckID("nobody").Call()
	logger.Debugf("Warmed")
}
