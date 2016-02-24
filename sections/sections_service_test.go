package sections

import (
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

var sectionsDriver baseftrwapp.Service

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"

	sectionsDriver = getSectionsCypherDriver(t)

	sectionToDelete := Section{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(sectionsDriver.Write(sectionToDelete), "Failed to write section")

	found, err := sectionsDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete section for uuid %", uuid)
	assert.NoError(err, "Error deleting section for uuid %s", uuid)

	p, found, err := sectionsDriver.Read(uuid)

	assert.Equal(Section{}, p, "Found section %s who should have been deleted", p)
	assert.False(found, "Found section for uuid %s who should have been deleted", uuid)
	assert.NoError(err, "Error trying to find section for uuid %s", uuid)
}

func TestCreateAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	sectionsDriver = getSectionsCypherDriver(t)

	sectionToWrite := Section{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")

	readSectionForUUIDAndCheckFieldsMatch(t, uuid, sectionToWrite)

	cleanUp(t, uuid)
}

func TestCreateHandlesSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	sectionsDriver = getSectionsCypherDriver(t)

	sectionToWrite := Section{UUID: uuid, CanonicalName: "Test 'special chars", TmeIdentifier: "TME_ID"}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")

	readSectionForUUIDAndCheckFieldsMatch(t, uuid, sectionToWrite)

	cleanUp(t, uuid)
}

func TestCreateNotAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	sectionsDriver = getSectionsCypherDriver(t)

	sectionToWrite := Section{UUID: uuid, CanonicalName: "Test"}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")

	readSectionForUUIDAndCheckFieldsMatch(t, uuid, sectionToWrite)

	cleanUp(t, uuid)
}

func TestUpdateWillRemovePropertiesNoLongerPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	sectionsDriver = getSectionsCypherDriver(t)

	sectionToWrite := Section{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")
	readSectionForUUIDAndCheckFieldsMatch(t, uuid, sectionToWrite)

	updatedSection := Section{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(sectionsDriver.Write(updatedSection), "Failed to write updated section")
	readSectionForUUIDAndCheckFieldsMatch(t, uuid, updatedSection)

	cleanUp(t, uuid)
}

func TestConnectivityCheck(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver = getSectionsCypherDriver(t)
	err := sectionsDriver.Check()
	assert.NoError(err, "Unexpected error on connectivity check")
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

func readSectionForUUIDAndCheckFieldsMatch(t *testing.T, uuid string, expectedSection Section) {
	assert := assert.New(t)
	storedSection, found, err := sectionsDriver.Read(uuid)

	assert.NoError(err, "Error finding section for uuid %s", uuid)
	assert.True(found, "Didn't find section for uuid %s", uuid)
	assert.Equal(expectedSection, storedSection, "sections should be the same")
}

func TestWritePrefLabelIsAlsoWrittenAndIsEqualToName(t *testing.T) {
	assert := assert.New(t)
	sectionsDriver := getSectionsCypherDriver(t)
	uuid := "12345"
	sectionToWrite := Section{UUID: uuid, CanonicalName: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(sectionsDriver.Write(sectionToWrite), "Failed to write section")

	result := []struct {
		PrefLabel string `json:"t.prefLabel"`
	}{}

	getPrefLabelQuery := &neoism.CypherQuery{
		Statement: `
				MATCH (t:Section {uuid:"12345"}) RETURN t.prefLabel
				`,
		Result: &result,
	}

	err := sectionsDriver.cypherRunner.CypherBatch([]*neoism.CypherQuery{getPrefLabelQuery})
	assert.NoError(err)
	assert.Equal("Test", result[0].PrefLabel, "PrefLabel should be 'Test")
	cleanUp(t, uuid)
}

func cleanUp(t *testing.T, uuid string) {
	assert := assert.New(t)
	found, err := sectionsDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete section for uuid %", uuid)
	assert.NoError(err, "Error deleting section for uuid %s", uuid)
}
