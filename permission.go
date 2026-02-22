package main

import (
	"encoding/csv"
	"os"
	"strings"
)

type Permission struct {
	Includes []string `json:"includes"`
	Excludes []string `json:"excludes"`
}

type Distributor struct {
	Name       string     `json:"name"`
	ParentName string     `json:"parent_name"`
	Permission Permission `json:"permission"`
}

type PermissionService struct {
	distributors map[string]*Distributor
	validRegions map[string]bool
}

func NewPermissionService() *PermissionService {

	ps := &PermissionService{
		distributors: make(map[string]*Distributor),
		validRegions: make(map[string]bool),
	}

	ps.loadCities("cities.csv")

	return ps
}

func (ps *PermissionService) loadCities(filePath string) {

	file, err := os.Open(filePath)
	if err != nil {
		panic("Cannot open cities.csv")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic("Cannot read cities.csv")
	}

	for i, row := range records {
		if i == 0 {
			continue // skip header
		}
		if len(row) < 6 {
			continue
		}

		city := normalize(row[3])
		state := normalize(row[4])
		country := normalize(row[5])

		full := city + "-" + state + "-" + country
		stateLevel := state + "-" + country

		ps.validRegions[full] = true
		ps.validRegions[stateLevel] = true
		ps.validRegions[country] = true
	}
}

func normalize(value string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(value), " ", ""))
}

func (ps *PermissionService) AddDistributor(d *Distributor) {
	ps.distributors[d.Name] = d
}

func (ps *PermissionService) CheckPermission(name string, region string) bool {

	region = normalize(region)

	if !ps.validRegions[region] {
		return false
	}

	dist, ok := ps.distributors[name]
	if !ok {
		return false
	}

	// Parent restriction
	if dist.ParentName != "" {
		if !ps.CheckPermission(dist.ParentName, region) {
			return false
		}
	}

	// Exclude
	for _, ex := range dist.Permission.Excludes {
		if match(region, normalize(ex)) {
			return false
		}
	}

	// Include
	for _, in := range dist.Permission.Includes {
		if match(region, normalize(in)) {
			return true
		}
	}

	return false
}

func match(region string, rule string) bool {

	regionParts := strings.Split(region, "-")
	ruleParts := strings.Split(rule, "-")

	if len(ruleParts) > len(regionParts) {
		return false
	}

	for i := 1; i <= len(ruleParts); i++ {
		if regionParts[len(regionParts)-i] != ruleParts[len(ruleParts)-i] {
			return false
		}
	}

	return true
}
