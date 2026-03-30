package players

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"codeberg.org/tealeg/xlsx/v4"
	"github.com/jaswdr/faker/v2"
	"github.com/starfederation/datastar-go/datastar"

	"github.com/cszczepaniak/cribbly/internal/moreiter"
	"github.com/cszczepaniak/cribbly/internal/persistence/players"
)

type PlayersHandler struct {
	PlayerRepo players.Repository
}

type excelSheetData struct {
	Name string     `json:"name"`
	Rows [][]string `json:"rows"`
}

type excelPreviewSheet struct {
	Name    string
	Rows    [][]string
	MaxCols int
}

type excelPreviewPageData struct {
	WorkbookJSON string
	Sheets       []excelPreviewSheet
}

func (h PlayersHandler) RegistrationPage(w http.ResponseWriter, r *http.Request) error {
	players, err := h.PlayerRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	tm := playerRegistrationPage(players)
	return tm.Render(r.Context(), w)
}

func (h PlayersHandler) PostPlayer(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	if err := datastar.ReadSignals(r, &signals); err != nil {
		return err
	}

	_, err := h.PlayerRepo.Create(r.Context(), signals.FirstName, signals.LastName)
	if err != nil {
		return err
	}

	players, err := h.PlayerRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)

	signals.FirstName = ""
	signals.LastName = ""
	err = sse.MarshalAndPatchSignals(signals)
	if err != nil {
		return err
	}

	return sse.PatchElementTempl(playerTable(players))
}

func (h PlayersHandler) GenerateRandomPlayers(w http.ResponseWriter, r *http.Request) error {
	var signals struct {
		Num int `json:"num"`
	}
	if err := datastar.ReadSignals(r, &signals); err != nil {
		return err
	}

	fake := faker.New()

	for range signals.Num {
		firstName := fake.Person().FirstName()
		lastName := fake.Person().LastName()
		_, err := h.PlayerRepo.Create(r.Context(), firstName, lastName)
		if err != nil {
			return err
		}
	}

	players, err := h.PlayerRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(playerTable(players))
}

func (h PlayersHandler) DeleteAllPlayers(w http.ResponseWriter, r *http.Request) error {
	players, err := h.PlayerRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	for _, p := range players {
		if p.TeamID != "" {
			err := h.PlayerRepo.UnassignFromTeam(r.Context(), p.TeamID, moreiter.Of(p.ID))
			if err != nil {
				return err
			}
		}

		err = h.PlayerRepo.Delete(r.Context(), p.ID)
		if err != nil {
			return err
		}
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(playerTable(nil))
}

func (h PlayersHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	err := h.PlayerRepo.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	players, err := h.PlayerRepo.GetAll(r.Context())
	if err != nil {
		return err
	}

	sse := datastar.NewSSE(w, r)
	return sse.PatchElementTempl(playerTable(players))
}

func (h PlayersHandler) UploadExcel(w http.ResponseWriter, r *http.Request) error {
	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	workbook, err := xlsx.OpenReaderAt(file, header.Size)
	if err != nil {
		return err
	}

	fullSheets, previewSheets, err := workbookData(workbook)
	if err != nil {
		return err
	}
	if len(previewSheets) == 0 {
		return fmt.Errorf("excel file has no sheets")
	}

	fullSheetsJSON, err := json.Marshal(fullSheets)
	if err != nil {
		return err
	}

	tm := excelImportPreviewPage(excelPreviewPageData{
		WorkbookJSON: string(fullSheetsJSON),
		Sheets:       previewSheets,
	})
	return tm.Render(r.Context(), w)
}

func (h PlayersHandler) ImportExcel(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	workbookJSON := r.FormValue("workbook_json")
	if workbookJSON == "" {
		return fmt.Errorf("missing workbook payload")
	}

	var sheets []excelSheetData
	if err := json.Unmarshal([]byte(workbookJSON), &sheets); err != nil {
		return err
	}
	if len(sheets) == 0 {
		return fmt.Errorf("workbook has no sheets")
	}

	sheetIndex, err := strconv.Atoi(r.FormValue("sheet_index"))
	if err != nil {
		return err
	}
	if sheetIndex < 0 || sheetIndex >= len(sheets) {
		return fmt.Errorf("sheet index out of range: %d", sheetIndex)
	}

	nameCol1Based, err := strconv.Atoi(r.FormValue("name_col"))
	if err != nil {
		return err
	}
	if nameCol1Based < 1 {
		return fmt.Errorf("name column must be 1-based and >= 1")
	}
	nameCol := nameCol1Based - 1

	startRow := 0
	if r.FormValue("skip_header") == "on" {
		startRow = 1
	}

	// TODO: transaction
	for rowIdx := startRow; rowIdx < len(sheets[sheetIndex].Rows); rowIdx++ {
		row := sheets[sheetIndex].Rows[rowIdx]
		if nameCol >= len(row) {
			continue
		}

		fullName := strings.TrimSpace(row[nameCol])
		if fullName == "" {
			continue
		}

		firstName, lastName, _ := strings.Cut(fullName, " ")
		if _, err := h.PlayerRepo.Create(r.Context(), firstName, lastName); err != nil {
			return err
		}
	}

	http.Redirect(w, r, "/admin/players", http.StatusSeeOther)
	return nil
}

func workbookData(workbook *xlsx.File) ([]excelSheetData, []excelPreviewSheet, error) {
	fullSheets := make([]excelSheetData, 0, len(workbook.Sheets))
	previewSheets := make([]excelPreviewSheet, 0, len(workbook.Sheets))

	for _, sheet := range workbook.Sheets {
		fullRows := make([][]string, 0, sheet.MaxRow)
		maxCols := 0

		err := sheet.ForEachRow(func(row *xlsx.Row) error {
			values := []string{}
			err := row.ForEachCell(func(cell *xlsx.Cell) error {
				col, _ := cell.GetCoordinates()
				for len(values) <= col {
					values = append(values, "")
				}
				values[col] = strings.TrimSpace(cell.String())
				return nil
			})
			if err != nil {
				return err
			}

			fullRows = append(fullRows, values)
			if len(values) > maxCols {
				maxCols = len(values)
			}
			return nil
		}, xlsx.SkipEmptyRows)
		if err != nil {
			return nil, nil, err
		}

		fullSheets = append(fullSheets, excelSheetData{
			Name: sheet.Name,
			Rows: fullRows,
		})

		previewSheets = append(previewSheets, excelPreviewSheet{
			Name:    sheet.Name,
			Rows:    fullRows,
			MaxCols: maxCols,
		})
	}

	return fullSheets, previewSheets, nil
}
