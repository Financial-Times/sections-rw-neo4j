package sections

import (
	"os"
	"testing"

	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

const (
	sectionUUID          = "12345"
	newSectionUUID       = "123456"
	tmeID                = "TME_ID"
	newTmeID             = "NEW_TME_ID"
	fsetID               = "fset_ID"
	leiCodeID            = "leiCode"
	prefLabel            = "Test"
	specialCharPrefLabel = "Test 'special chars"
)

var defaultTypes = []string{"Thing", "Concept", "Classification", "Section"}

func TestConnectivityCheck(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)
	err := sectionsDriver.Check()
	assert.NoError(err, "Unexpected error on connectivity check")
}

func TestPrefLabelIsCorrectlyWritten(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{UUIDS: []string{sectionUUID}}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	err := sectionsDriver.Write(sectionToWrite)
	assert.NoError(err, "ERROR happened during write time")

	storedSection, found, err := sectionsDriver.Read(sectionUUID)
	assert.NoError(err, "ERROR happened during read time")
	assert.Equal(true, found)
	assert.NotEmpty(storedSection)

	assert.Equal(prefLabel, storedSection.(Section).PrefLabel, "PrefLabel should be "+prefLabel)
	cleanUp(assert, sectionUUID, sectionsDriver)
}

func TestPrefLabelSpecialCharactersAreHandledByCreate(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{sectionUUID}}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")

	//add default types that will be automatically added by the writer
	sectionToWrite.Types = defaultTypes
	//check if sectionToWrite is the same with the one inside the DB
	readSectionForUUIDAndCheckFieldsMatch(assert, sectionsDriver, sectionUUID, sectionToWrite)
	cleanUp(assert, sectionUUID, sectionsDriver)
}

func TestCreateCompleteSectionWithPropsAndIdentifiers(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}, FactsetIdentifier: fsetID, LeiCode: leiCodeID}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")

	//add default types that will be automatically added by the writer
	sectionToWrite.Types = defaultTypes
	//check if sectionToWrite is the same with the one inside the DB
	readSectionForUUIDAndCheckFieldsMatch(assert, sectionsDriver, sectionUUID, sectionToWrite)
	cleanUp(assert, sectionUUID, sectionsDriver)
}

func TestUpdateWillRemovePropertiesAndIdentifiersNoLongerPresent(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)

	allAlternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{sectionUUID}, FactsetIdentifier: fsetID, LeiCode: leiCodeID}
	sectionToWrite := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: allAlternativeIdentifiers}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")
	//add default types that will be automatically added by the writer
	sectionToWrite.Types = defaultTypes
	readSectionForUUIDAndCheckFieldsMatch(assert, sectionsDriver, sectionUUID, sectionToWrite)

	tmeAlternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}}
	updatedSection := Section{UUID: sectionUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: tmeAlternativeIdentifiers}

	assert.NoError(sectionsDriver.Write(updatedSection), "Failed to write updated section")
	//add default types that will be automatically added by the writer
	updatedSection.Types = defaultTypes
	readSectionForUUIDAndCheckFieldsMatch(assert, sectionsDriver, sectionUUID, updatedSection)

	cleanUp(assert, sectionUUID, sectionsDriver)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}}
	sectionToDelete := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(sectionsDriver.Write(sectionToDelete), "Failed to write section")

	found, err := sectionsDriver.Delete(sectionUUID)
	assert.True(found, "Didn't manage to delete section for uuid %", sectionUUID)
	assert.NoError(err, "Error deleting section for uuid %s", sectionUUID)

	p, found, err := sectionsDriver.Read(sectionUUID)

	assert.Equal(Section{}, p, "Found section %s who should have been deleted", p)
	assert.False(found, "Found section for uuid %s who should have been deleted", sectionUUID)
	assert.NoError(err, "Error trying to find section for uuid %s", sectionUUID)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)

	alternativeIds := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{sectionUUID}}
	sectionOneToCount := Section{UUID: sectionUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIds}

	assert.NoError(sectionsDriver.Write(sectionOneToCount), "Failed to write section")

	nr, err := sectionsDriver.Count()
	assert.Equal(1, nr, "Should be 1 sections in DB - count differs")
	assert.NoError(err, "An unexpected error occurred during count")

	newAlternativeIds := alternativeIdentifiers{TME: []string{newTmeID}, UUIDS: []string{newSectionUUID}}
	sectionTwoToCount := Section{UUID: newSectionUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: newAlternativeIds}

	assert.NoError(sectionsDriver.Write(sectionTwoToCount), "Failed to write section")

	nr, err = sectionsDriver.Count()
	assert.Equal(2, nr, "Should be 2 sections in DB - count differs")
	assert.NoError(err, "An unexpected error occurred during count")

	cleanUp(assert, sectionUUID, sectionsDriver)
	cleanUp(assert, newSectionUUID, sectionsDriver)
}

func readSectionForUUIDAndCheckFieldsMatch(assert *assert.Assertions, sectionsDriver service, uuid string, expectedSection Section) {

	storedSection, found, err := sectionsDriver.Read(uuid)

	assert.NoError(err, "Error finding section for uuid %s", uuid)
	assert.True(found, "Didn't find section for uuid %s", uuid)
	assert.Equal(expectedSection, storedSection, "sections should be the same")
}

func getSectionsCypherDriver(t *testing.T) service {
	assert := assert.New(t)
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	db, err := neoism.Connect(url)
	assert.NoError(err, "Failed to connect to Neo4j")
	return NewCypherSectionsService(neoutils.StringerDb{db}, db)
}

func cleanUp(assert *assert.Assertions, uuid string, sectionsDriver service) {
	found, err := sectionsDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete section for uuid %", uuid)
	assert.NoError(err, "Error deleting section for uuid %s", uuid)
}
