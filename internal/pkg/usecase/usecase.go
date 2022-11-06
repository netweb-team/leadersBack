package usecase

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"leaders_apartments/internal/pkg/config"
	"leaders_apartments/internal/pkg/correct"
	"leaders_apartments/internal/pkg/domain"
	"leaders_apartments/internal/pkg/parser/html"
	"leaders_apartments/internal/pkg/parser/xslx"
	"leaders_apartments/internal/pkg/utils"
	"net/http"
	"os"

	"github.com/ddulesov/gogost/gost34112012512"
	"github.com/labstack/gommon/log"
)

type serverUsecases struct {
	repo domain.Repository
}

func New(repo domain.Repository) domain.Usecase {
	return &serverUsecases{repo}
}

func (u *serverUsecases) ImportXlsx(f io.Reader, user int) *domain.Table {
	data, err := xslx.Parse(f)
	if err != nil {
		return nil
	}

	if data.ID, err = u.repo.SaveTable(data.Path, user, data.Flat); err != nil {
		os.Remove(data.Path)
		return nil
	}
	return data
}

func (u *serverUsecases) ExportXlsx(id int) string {
	filename, err := u.repo.GetTableName(id)
	if err != nil {
		return ""
	}
	return filename
}

func (u *serverUsecases) FindAnalogs(id, ptnIndex int) *domain.PatternAnalogs {
	filename, err := u.repo.GetTableName(id)
	if err != nil {
		return nil
	}
	row := xslx.ReadRow(filename, ptnIndex)
	if row == nil {
		return nil
	}
	row.ID = ptnIndex
	cfg := config.New()
	resp, err := http.Get(fmt.Sprintf(cfg.MapApi, row.Address, cfg.ApiKey))
	if err != nil {
		log.Info("Cannot get request for coordinates: ", err)
		return nil
	}
	defer resp.Body.Close()
	coords := new(domain.MapResponse)
	json.NewDecoder(resp.Body).Decode(coords)
	row.Longitude, row.Latitude = coords.Features[0].Center[0], coords.Features[0].Center[1]
	analogs := html.FindAnalogs(row)
	if len(analogs) == 0 {
		return nil
	}
	for _, anal := range analogs {
		anal.Rooms, anal.Floors, anal.Walls, anal.Good, anal.AvgCost = row.Rooms, row.Floors, row.Walls, 1, float64(anal.Cost)/anal.Total
	}
	coefs := correct.Do(row, analogs)
	sum := 0.0
	for _, anal := range analogs {
		sum += anal.AvgCost
	}
	row.AvgCost = sum / float64(len(analogs))
	row.Cost = int(row.Total * row.AvgCost)
	if u.repo.SaveAnalogs(id, row, analogs, coefs) != nil {
		return nil
	}
	return &domain.PatternAnalogs{Pattern: row, Analogs: analogs, Correct: coefs}
}

func (u *serverUsecases) CalcPool(id int) []*domain.Row {
	filename, err := u.repo.GetTableName(id)
	if err != nil {
		return nil
	}
	table := xslx.ReadTable(filename)
	if table == nil {
		return nil
	}
	data, err := u.repo.GetPatternAnalogs(id)
	if err != nil {
		return nil
	}
	for i := range data {
		idx := data[i].Pattern.ID - 1
		table[idx].AvgCost = data[i].Pattern.AvgCost
		data[i].Pattern = table[idx]
	}
	groups := make([][]*domain.Row, len(data))
	for i := range groups {
		groups[i] = make([]*domain.Row, 0)
	}
	for j, row := range table {
		for i, ptn := range data {
			p := ptn.Pattern
			if j == p.ID-1 || row.Rooms == p.Rooms && row.Segment == p.Segment && row.Floors == p.Floors && row.Walls == p.Walls {
				row.AvgCost = p.AvgCost
				groups[i] = append(groups[i], row)
				break
			}
		}
	}
	for i := range groups {
		correct.Do(data[i].Pattern, groups[i])
	}
	//path := filename
	filename = xslx.SavePrice(filename, table)
	// u.repo.ChangePath(id, filename)
	// os.Remove(path)
	return table
}

func (u *serverUsecases) CreateUser(user *domain.User) string {
	hasher.Reset()
	hasher.Write([]byte(user.Password))
	user.HashPass = hasher.Sum(nil)
	if u.repo.CreateUser(user) != nil {
		return ""
	}
	cookie := utils.RandString(32)
	if u.repo.CreateCookie(user.ID, cookie) != nil {
		return ""
	}
	return cookie
}

func (u *serverUsecases) CreateAuth(user *domain.User) string {
	hasher.Reset()
	hasher.Write([]byte(user.Password))
	user.HashPass = hasher.Sum(nil)
	dbUser := u.repo.GetUser(user.Login)
	if dbUser == nil || hex.EncodeToString(dbUser.HashPass) != hex.EncodeToString(user.HashPass) {
		log.Info(dbUser)
		return ""
	}
	cookie := utils.RandString(32)
	if u.repo.CreateCookie(dbUser.ID, cookie) != nil {
		log.Info(cookie)
		return ""
	}
	return cookie
}

func (u *serverUsecases) DeleteAuth(cookie string) {
	u.repo.DeleteCookie(cookie)
}

func (u *serverUsecases) CheckAuth(cookie string, pool int) int {
	if pool > 0 {
		return u.repo.CheckPool(cookie, pool)
	}
	return u.repo.CheckCookie(cookie)
}

var hasher = gost34112012512.New()

func (u *serverUsecases) ChangeCorrect(pool, id int, coefs *domain.CorrectCoefs) *domain.PatternAnalogs {
	ptn := u.repo.ChangeCorrect(pool, id, coefs)
	if ptn == 0 {
		return nil
	}
	analogs, corrects := u.repo.GetAnalogs(pool, ptn)
	sum := 0.0
	for _, a := range analogs {
		sum += a.AvgCost
	}
	sum /= float64(len(analogs))
	if u.repo.SavePatternPrice(pool, ptn, sum) != nil {
		return nil
	}
	return &domain.PatternAnalogs{Pattern: &domain.Row{ID: ptn, AvgCost: sum}, Analogs: analogs, Correct: corrects}
}

func (u *serverUsecases) ChangeAnalog(pool, id int) *domain.PatternAnalogs {
	ptn := u.repo.ChangeAnalog(pool, id)
	if ptn == 0 {
		return nil
	}
	analogs, corrects := u.repo.GetAnalogs(pool, ptn)
	sum := 0.0
	for _, a := range analogs {
		sum += a.AvgCost
	}
	sum /= float64(len(analogs))
	if u.repo.SavePatternPrice(pool, ptn, sum) != nil {
		return nil
	}
	return &domain.PatternAnalogs{Pattern: &domain.Row{ID: ptn, AvgCost: sum}, Analogs: analogs, Correct: corrects}
}

func (u *serverUsecases) GetUserArchives(user int) []*domain.Table {
	return u.repo.GetUserArchives(user)
}

func (u *serverUsecases) GetArchive(id int) *domain.Table {
	archive := u.repo.GetArchive(id)
	if archive == nil {
		return nil
	}
	archive.Rows = xslx.ReadTable(archive.Path)
	if err := xslx.ReadPrice(archive.Path, archive.Rows); err != nil {
		return nil
	}
	archive.PA, _ = u.repo.GetPatternAnalogs(id)
	return archive
}
