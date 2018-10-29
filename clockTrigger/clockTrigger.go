package clockTrigger

import (
	"net/http"
	"time"
)

func doEvery(d time.Duration, f func())  {
	for range time.Tick(d) {
		f()
	}
}


func triggerWebhook() {
	http.Get("https://igcapi.herokuapp.com/paragliding/admin/api/webhooks")
}


func main() {

	doEvery(10*time.Minute, triggerWebhook)

}
