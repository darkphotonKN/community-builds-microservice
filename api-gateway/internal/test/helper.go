package testsuite

import (
	"log"
)

/**
* Helper functions for testing.
**/

/**
* Helps extracting metadata for your specific test.
**/
func (t *TestSuite) GetMetaData(testName string) MetaData {

	// assert meta data
	metaData, metaDataOk := t.MetaData[testName].(MetaData)

	if !metaDataOk {
		log.Fatalf("Failed to retrieve metaData: metaDataOk=%v", metaDataOk)
	}

	return metaData
}

/**
* Helps set metadata for a specific test.
**/
func (t *TestSuite) setMetaData(testName string, data interface{}) {
	t.MetaData[testName] = data
}
