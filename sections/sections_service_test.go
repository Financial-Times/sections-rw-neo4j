package sections

import (
	"os"
	"testing"

	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/stretchr/testify/assert"
	"sort"
)

const (
	sectionUUID          = "12345"
	newSectionUUID       = "123456"
	tmeID                = "TME_ID"
	newTmeID             = "NEW_TME_ID"
	prefLabel            = "Test"
	specialCharPrefLabel = "Test 'special chars"
)

var defaultTypes = []string{"Thing", "Concept", "Classification", "Section"}

func TestConnectivityCheck(t *testing.T) {

	sectionsDriver := getSectionsCypherDriver(t)
	err := sectionsDriver.Check()
	assert.NoError(t, err, "Unexpected error on connectivity check")
}

func TestPrefLabelIsCorrectlyWritten(t *testing.T) {

	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{UUIDS: []string{sectionUUID}}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	err := sectionsDriver.Write(sectionToWrite)
	assert.NoError(t, err, "ERROR happened during write time")

	storedSection, found, err := sectionsDriver.Read(sectionUUID)
	assert.NoError(t, err, "ERROR happened during read time")
	assert.Equal(t, true, found)
	assert.NotEmpty(t, storedSection)

	assert.Equal(t, prefLabel, storedSection.(Section).PrefLabel, "PrefLabel should be "+prefLabel)
	cleanUp(t, sectionUUID, sectionsDriver)
}

func TestPrefLabelSpecialCharactersAreHandledByCreate(t *testing.T) {

	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{sectionUUID}}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(t, sectionsDriver.Write(sectionToWrite), "Failed to write section")

	//add default types that will be automatically added by the writer
	sectionToWrite.Types = defaultTypes
	//check if sectionToWrite is the same with the one inside the DB
	readSectionForUUIDAndCheckFieldsMatch(t, sectionsDriver, sectionUUID, sectionToWrite)
	cleanUp(t, sectionUUID, sectionsDriver)
}

func TestCreateCompleteSectionWithPropsAndIdentifiers(t *testing.T) {

	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(t, sectionsDriver.Write(sectionToWrite), "Failed to write section")

	//add default types that will be automatically added by the writer
	sectionToWrite.Types = defaultTypes
	//check if sectionToWrite is the same with the one inside the DB
	readSectionForUUIDAndCheckFieldsMatch(t, sectionsDriver, sectionUUID, sectionToWrite)
	cleanUp(t, sectionUUID, sectionsDriver)
}

func TestUpdateWillRemovePropertiesAndIdentifiersNoLongerPresent(t *testing.T) {

	sectionsDriver := getSectionsCypherDriver(t)

	allAlternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{sectionUUID}}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: allAlternativeIdentifiers}

	assert.NoError(t, sectionsDriver.Write(sectionToWrite), "Failed to write section")
	//add default types that will be automatically added by the writer
	sectionToWrite.Types = defaultTypes
	readSectionForUUIDAndCheckFieldsMatch(t, sectionsDriver, sectionUUID, sectionToWrite)

	tmeAlternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}}
	updatedSection := Section{UUID: sectionUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: tmeAlternativeIdentifiers}

	assert.NoError(t, sectionsDriver.Write(updatedSection), "Failed to write updated section")
	//add default types that will be automatically added by the writer
	updatedSection.Types = defaultTypes
	readSectionForUUIDAndCheckFieldsMatch(t, sectionsDriver, sectionUUID, updatedSection)

	cleanUp(t, sectionUUID, sectionsDriver)
}

func TestDelete(t *testing.T) {

	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}}
	sectionToDelete := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(t, sectionsDriver.Write(sectionToDelete), "Failed to write section")

	found, err := sectionsDriver.Delete(sectionUUID)
	assert.True(t, found, "Didn't manage to delete section for uuid %", sectionUUID)
	assert.NoError(t, err, "Error deleting section for uuid %s", sectionUUID)

	p, found, err := sectionsDriver.Read(sectionUUID)

	assert.Equal(t, Section{}, p, "Found section %s who should have been deleted", p)
	assert.False(t, found, "Found section for uuid %s who should have been deleted", sectionUUID)
	assert.NoError(t, err, "Error trying to find section for uuid %s", sectionUUID)
}

func TestCount(t *testing.T) {

	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIds := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}}
	sectionOneToCount := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIds}

	assert.NoError(t, sectionsDriver.Write(sectionOneToCount), "Failed to write section")

	nr, err := sectionsDriver.Count()
	assert.Equal(t, 1, nr, "Should be 1 sections in DB - count differs")
	assert.NoError(t, err, "An unexpected error occurred during count")

	newAlternativeIds := alternativeIdentifiers{TME: []string{newTmeID}, UUIDS: []string{newSectionUUID}}
	sectionTwoToCount := Section{UUID: newSectionUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: newAlternativeIds}

	assert.NoError(t, sectionsDriver.Write(sectionTwoToCount), "Failed to write section")

	nr, err = sectionsDriver.Count()
	assert.Equal(t, 2, nr, "Should be 2 sections in DB - count differs")
	assert.NoError(t, err, "An unexpected error occurred during count")

	cleanUp(t, sectionUUID, sectionsDriver)
	cleanUp(t, newSectionUUID, sectionsDriver)
}

func readSectionForUUIDAndCheckFieldsMatch(t *testing.T, sectionsDriver service, uuid string, expectedSection Section) {
	section, found, err := sectionsDriver.Read(uuid)
	sort.Strings(expectedSection.Types)
	sort.Strings(expectedSection.AlternativeIdentifiers.TME)
	sort.Strings(expectedSection.AlternativeIdentifiers.UUIDS)

	storedSection := section.(Section)
	sort.Strings(storedSection.Types)
	sort.Strings(storedSection.AlternativeIdentifiers.TME)
	sort.Strings(storedSection.AlternativeIdentifiers.UUIDS)

	assert.NoError(t, err, "Error finding section for uuid %s", uuid)
	assert.True(t, found, "Didn't find section for uuid %s", uuid)
	assert.Equal(t, expectedSection, storedSection, "sections should be the same")
}

func getSectionsCypherDriver(t *testing.T) service {
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	db, err := neoutils.Connect(url, conf)
	assert.NoError(t, err, "Failed to connect to Neo4j")
	service := NewCypherSectionsService(db)
	service.Initialise()
	return service
}

func cleanUp(t *testing.T, uuid string, sectionsDriver service) {
	found, err := sectionsDriver.Delete(uuid)
	assert.True(t, found, "Didn't manage to delete section for uuid %", uuid)
	assert.NoError(t, err, "Error deleting section for uuid %s", uuid)
}
