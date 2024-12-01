package testsuite

import "fmt"

/**
* Helper functions for testing.
**/

/**
* Helps extracting metadata for your specific test.
**/
func (t *TestSuite) getMetaData(testName string) MetaData {

	// assert meta data
	metaData, metaDataOk := t.MetaData[testName].(MetaData)

	if !metaDataOk {
		fmt.Printf("Failed to retrieve metaData: metaDataOk=%v", metaDataOk)
	}

	return metaData
}

/**
* Helps set metadata for a specific test.
**/
func (t *TestSuite) setMetaData(testName string, data interface{}) {

	// -- TestAddSkillLinksToBuildService_Success METADATA --
	t.MetaData[testName] = data
}
