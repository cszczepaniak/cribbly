package divisions

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/jaswdr/faker/v2"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/persistence/divisions"
	"github.com/cszczepaniak/cribbly/internal/persistence/teams"
	divisionservice "github.com/cszczepaniak/cribbly/internal/service/divisions"
	"github.com/cszczepaniak/cribbly/internal/ui/components"
)

type DivisionsHandler struct {
	TeamRepo        teams.Repository
	DivisionRepo    divisions.Repository
	DivisionService divisionservice.Service
}

func (h DivisionsHandler) Index(w http.ResponseWriter, r *http.Request) error {
	divisions, err := h.DivisionRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	teamsByDivision := make(map[string][]teams.Team, len(divisions))
	for _, division := range divisions {
		teams, err := h.TeamRepo.GetForDivision(r.Context(), division.ID)
		if err != nil {
			return err
		}
		teamsByDivision[division.ID] = teams
	}

	return index(divisions, teamsByDivision).Render(r.Context(), w)
}

func (h DivisionsHandler) EditPage(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	division, err := h.DivisionService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	availableTeams, err := h.TeamRepo.GetWithoutDivision(r.Context())
	if err != nil {
		return err
	}

	return editDivision(division, availableTeams).Render(r.Context(), w)
}

func (h DivisionsHandler) EditName(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	division, err := h.DivisionService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(editDivisionInput(division, ""))
}

func (h DivisionsHandler) SaveName(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Name string `json:"name"`
	}
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	division, err := h.DivisionService.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}

	if signals.Name == "" {
		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(editDivisionInput(division, "Division name can't be empty."))
	}

	err = h.DivisionRepo.Rename(r.Context(), division.ID, signals.Name)
	if err != nil {
		return err
	}
	division.Name = signals.Name

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(editDivisionTitle(division))
}

func (h DivisionsHandler) SaveSize(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Size string `json:"size"`
	}
	err := datastar.ReadSignals(r, &signals)
	if err != nil {
		return err
	}

	division, err := h.DivisionService.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		return err
	}

	switch signals.Size {
	case "4":
		if division.Size == 4 {
			return nil
		}
		if len(division.Teams) > 4 {
			return errors.New("cannot decrease division size below the number of teams it has")
		}
		division.Size = 4
	case "6":
		if division.Size == 6 {
			return nil
		}
		division.Size = 6
	default:
		return errors.New("invalid division size")
	}

	err = h.DivisionRepo.UpdateSize(r.Context(), r.PathValue("id"), division.Size)
	if err != nil {
		return nil
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(divisionTeamList(division))
}

func (h DivisionsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	division, err := h.DivisionRepo.Create(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.Redirectf("/admin/divisions/%s", division.ID)
}

func (h DivisionsHandler) Generate(w http.ResponseWriter, r *http.Request) error {
	allTeams, err := h.TeamRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	if len(allTeams)%2 != 0 {
		return components.ShowErrorToast(w, r, "Cannot generate divisions with an odd number of teams.")
	}

	fake := faker.New()
	for len(allTeams) > 0 {
		division, err := h.DivisionRepo.Create(r.Context())
		if err != nil {
			return err
		}

		err = h.DivisionRepo.Rename(r.Context(), division.ID, fake.ProgrammingLanguage().Name())
		if err != nil {
			return err
		}

		nForThisDivision := 4
		if len(allTeams) == 6 {
			nForThisDivision = 6
		}

		forDivision := allTeams[:nForThisDivision]
		allTeams = allTeams[nForThisDivision:]

		for _, team := range forDivision {
			err := h.TeamRepo.AssignToDivision(r.Context(), team.ID, division.ID)
			if err != nil {
				return err
			}
		}
	}

	sse := datastar.NewSSE(w, r)
	return sse.Redirect("/admin/divisions")
}

func (h DivisionsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	division, err := h.DivisionRepo.Get(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	err = sse.PatchElementTempl(confirmDeleteTitle(division.Name))
	if err != nil {
		return err
	}
	return sse.PatchElementTempl(confirmDeleteButton(division.ID))
}

func (h DivisionsHandler) ConfirmDelete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := h.DivisionRepo.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.Redirect("/admin/divisions")
}

func (h DivisionsHandler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	// TODO: do this in a transaction
	divisions, err := h.DivisionRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	for _, d := range divisions {
		division, err := h.DivisionService.Get(r.Context(), d.ID)
		if err != nil {
			return err
		}

		for _, t := range division.Teams {
			err := h.TeamRepo.UnassignFromDivision(r.Context(), t.ID)
			if err != nil {
				return err
			}
		}

		err = h.DivisionRepo.Delete(r.Context(), d.ID)
		if err != nil {
			return err
		}
	}

	sse := datastar.NewSSE(w, r)
	return sse.Redirect("/admin/divisions")
}

type signals struct {
	Name string `json:"name"`
}

func (h DivisionsHandler) Save(w http.ResponseWriter, r *http.Request) error {
	divisionID := r.PathValue("id")

	assign := r.URL.Query().Get("assign")
	unassign := r.URL.Query().Get("unassign")
	if assign != "" || unassign != "" {
		var err error
		if assign != "" {
			err = h.TeamRepo.AssignToDivision(r.Context(), assign, divisionID)
		} else {
			err = h.TeamRepo.UnassignFromDivision(r.Context(), unassign)
		}
		if err != nil {
			return err
		}

		division, err := h.DivisionService.Get(r.Context(), divisionID)
		if err != nil {
			return err
		}

		available, err := h.TeamRepo.GetWithoutDivision(r.Context())
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)
		return sse.PatchElementTempl(editDivisionDetails(division, available))
	}

	var sigs signals
	err := datastar.ReadSignals(r, &sigs)
	if err != nil {
		return err
	}

	if sigs.Name != "" {
		err := h.DivisionRepo.Rename(r.Context(), divisionID, sigs.Name)
		if err != nil {
			return err
		}

		sse := datastar.NewSSE(w, r)

		err = sse.PatchElementTempl(divisionName(divisionID, sigs.Name))
		if err != nil {
			return err
		}
		return sse.MarshalAndPatchSignals(signals{Name: ""})
	}

	return nil
}

type divisionQR struct {
	img          []byte
	divisionName string
}

func (h DivisionsHandler) GenerateQRCodes(w http.ResponseWriter, r *http.Request) error {
	divs, err := h.DivisionRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	var host string
	ref, err := url.Parse(r.Header.Get("Referer"))
	if err != nil {
		slog.Error("could not parse Referer header", "err", err)
		host = r.Host
	} else {
		host = ref.Scheme + "://" + ref.Host
	}

	qrs := make([]divisionQR, 0, len(divs))
	for _, div := range divs {
		q, err := qrcode.New(fmt.Sprintf("%s/divisions/%s", host, div.ID), qrcode.Medium)
		if err != nil {
			return err
		}
		q.DisableBorder = true

		img := q.Image(256)
		encoder := png.Encoder{CompressionLevel: png.BestCompression}

		buf := bytes.NewBuffer(nil)
		w := base64.NewEncoder(base64.StdEncoding, buf)

		err = encoder.Encode(w, img)
		if err != nil {
			return err
		}

		qrs = append(qrs, divisionQR{img: buf.Bytes(), divisionName: div.Name})
	}

	return qrPage(qrs).Render(r.Context(), w)
}
