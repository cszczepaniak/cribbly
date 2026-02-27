package rodtests

import (
	"testing"
	"time"

	"github.com/cszczepaniak/cribbly/internal/assert"
)

func TestValidRoomCodeUnlocksHome(t *testing.T) {
	app := newTestApp(t)

	// Seed a valid room code directly via the repository.
	const code = "ABC123"
	expires := time.Now().Add(time.Hour)
	assert.NoError(t, app.RoomCodeRepo.Create(t.Context(), code, expires))

	page := app.Browser.MustPage(app.BaseURL + "/").MustWaitLoad()
	t.Cleanup(page.MustClose)
	page = page.Timeout(20 * time.Second)

	roomCodeInput := page.MustElement("input[name='room_code']")
	roomCodeInput.MustWaitVisible().MustInput(code)

	// Submit the form and wait for the POST redirect to complete.
	waitNav := page.MustWaitNavigation()
	page.MustElementR("button", "Continue").MustClick()
	waitNav()

	// After the cookie is set, the welcome content should be visible.
	h1 := page.MustElementR("h1", "Welcome to Cribbly")
	assert.Equal(t, "Welcome to Cribbly", h1.MustText())

	// Assert the room_code cookie was set by the server.
	cookies := page.MustCookies(app.BaseURL + "/")
	for _, c := range cookies {
		if c.Name == "room_code" {
			assert.Equal(t, c.Value, code)
			break
		}
	}

	// The primary navigation buttons should also be present and navigable.
	page.MustElementR("button,a", "View Divisions").MustClick()
	page.MustWaitLoad()

	page.MustNavigateBack()
	page.MustWaitLoad()

	page.MustElementR("button,a", "View Standings").MustClick()
	page.MustWaitLoad()
}
