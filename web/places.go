package web

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

// ErrPlacesNotFound is returned by GetPlaces when the job's CSV output does not
// exist. Callers use it to distinguish a missing job (404) from other errors.
var ErrPlacesNotFound = errors.New("places not found")

// Place is a single map-able result extracted from a job's CSV output.
type Place struct {
	Title        string  `json:"title"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Link         string  `json:"link"`
	Category     string  `json:"category"`
	Phone        string  `json:"phone"`
	Website      string  `json:"website"`
	ReviewRating float64 `json:"review_rating"`
	// Distance is the great-circle distance in meters from the job's center
	// point. It is 0 when no center point was provided for the job.
	Distance float64 `json:"distance"`
	// DistanceText is Distance formatted for humans with km/m units, e.g.
	// "450 m" or "1.20 km". Empty when no center point was provided.
	DistanceText string `json:"distance_text"`
}

// haversine returns the great-circle distance in meters between two lat/lon
// points using the haversine formula.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusM = 6371000.0

	rad := math.Pi / 180

	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)

	return 2 * earthRadiusM * math.Asin(math.Min(1, math.Sqrt(a)))
}

// formatDistance renders a distance in meters as a human string, switching from
// meters to kilometers at 1 km.
func formatDistance(m float64) string {
	if m < 1000 {
		return fmt.Sprintf("%.0f m", m)
	}

	return fmt.Sprintf("%.2f km", m/1000)
}

// withDistances fills in each place's Distance/DistanceText relative to the
// given center point and returns the slice sorted nearest-first. When hasCenter
// is false the places are returned unchanged (no distance, original order).
func withDistances(places []Place, centerLat, centerLon float64, hasCenter bool) []Place {
	if !hasCenter {
		return places
	}

	for i := range places {
		d := haversine(centerLat, centerLon, places[i].Latitude, places[i].Longitude)
		places[i].Distance = d
		places[i].DistanceText = formatDistance(d)
	}

	sort.SliceStable(places, func(i, j int) bool {
		return places[i].Distance < places[j].Distance
	})

	return places
}

// GetPlaces locates the job's CSV output and parses it into mappable places.
// In web mode each job writes exactly one {id}.csv, so that file is the single
// source of truth for the map.
func (s *Service) GetPlaces(_ context.Context, id string) ([]Place, error) {
	path, err := s.csvPath(id)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("csv file not found for job %s: %w", id, ErrPlacesNotFound)
		}

		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	return parsePlaces(f)
}

// parsePlaces reads scraped results from a CSV stream and returns the places
// that have valid coordinates. Columns are resolved by header name so the
// parser tolerates reordering; the names mirror gmaps.Entry.CsvHeaders().
func parsePlaces(r io.Reader) ([]Place, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return []Place{}, nil
		}

		return nil, err
	}

	col := make(map[string]int, len(header))
	for i, name := range header {
		col[name] = i
	}

	get := func(row []string, name string) string {
		idx, ok := col[name]
		if !ok || idx >= len(row) {
			return ""
		}

		return row[idx]
	}

	places := []Place{}

	// seen deduplicates places so repeated rows (same listing scraped more than
	// once) do not pile up in the map/list view. The key prefers a stable Google
	// identifier and falls back to the maps link, then rounded coordinates.
	seen := make(map[string]struct{})

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		lat, errLat := strconv.ParseFloat(get(row, "latitude"), 64)
		lon, errLon := strconv.ParseFloat(get(row, "longitude"), 64)

		if errLat != nil || errLon != nil {
			continue
		}

		if !finite(lat) || !finite(lon) {
			continue
		}

		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			continue
		}

		if lat == 0 && lon == 0 {
			continue
		}

		key := get(row, "place_id")
		if key == "" {
			key = get(row, "cid")
		}

		if key == "" {
			key = get(row, "link")
		}

		if key == "" {
			key = fmt.Sprintf("%s|%.6f|%.6f", get(row, "title"), lat, lon)
		}

		if _, dup := seen[key]; dup {
			continue
		}

		seen[key] = struct{}{}

		rating, _ := strconv.ParseFloat(get(row, "review_rating"), 64)
		if !finite(rating) {
			rating = 0
		}

		places = append(places, Place{
			Title:        get(row, "title"),
			Address:      get(row, "address"),
			Latitude:     lat,
			Longitude:    lon,
			Link:         get(row, "link"),
			Category:     get(row, "category"),
			Phone:        get(row, "phone"),
			Website:      get(row, "website"),
			ReviewRating: rating,
		})
	}

	return places, nil
}

// finite reports whether f is a usable, real number (not NaN or ±Inf).
func finite(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0)
}
