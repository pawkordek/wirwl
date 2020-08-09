package wirwl

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io"
	"os"
	"path/filepath"
	"testing"
	"wirwl/internal/data"
	"wirwl/internal/log"
)

/*File containing various functions allowing to setup/clean up the test environment.
TestMain is run before all tests which setups the test directory and test db.
After all tests run it removes the test directory with all of it's data.
Therefore all tests making use of the filesystem should run in the test directory.
If they need a test db, they should make a copy of it.
If they need to store data persisting between the tests they should use the persistent test data directory, otherwise they
are free to do as they wish but for app config dir/data directory, the standard paths defined below should be used.
After a test is run, after test cleanup function should be executed to ensure that test environment (mostly test directories)
are restored back to the state before test run.
*/

/*All path variables below should be treated as constants. They cannot be made const as they need to have slashes adapted
for cross platform compatibility.
*/

//Folder that should be used for storing any temporary data when testing and for storing directories used as paths in
//passed into application's Config file.
//It's contents are cleared after every test, except for folder 'persistent' and it's contents.
//The folder itself is created at the beginning before first test is run and removed after all tests are run
var testDataDirPath = filepath.FromSlash("../testdata/")

//Should be used for storing data that must persist between tests as it's the only folder in testdata directory
//which is not removed after each test's execution
var persistentTestDataDirPath = filepath.FromSlash(testDataDirPath + "persistent/")

//Used as application's data directory path when testing
var testAppDataDirPath = filepath.FromSlash(testDataDirPath + "app_data/")

//Used as application's config directory path when testing
var testConfigDirPath = filepath.FromSlash(testDataDirPath + "config/")

/* It's a path to a database file which is generated every time tests are run but before any test executes.
If shown in the application, the data would look as follows:
comics:
	some comic1
	some comic2
music:
	some music1
	some music2
videos:
	some video1
	some video2
The file should be copied to perform any tests that require an existing data.
*/
var testDbPath = filepath.FromSlash(persistentTestDataDirPath + "data.db")

//Path used by tests to store a copy of test database file so that they don't affect the original file
var testDbCopyPath = filepath.FromSlash(testAppDataDirPath + "data.db")

func TestMain(m *testing.M) {
	setupTestEnvironment()
	exitCode := m.Run()
	cleanupTestEnvironment()
	os.Exit(exitCode)
}

func setupTestEnvironment() {
	/*Cleanup is run in the case that a test crashed in the previous run and couldn't run it's cleanup functions, leaving
	potentially unwanted files and directories*/
	cleanupTestEnvironment()
	err := data.CreateDirIfNotExist(testDataDirPath)
	if err != nil {
		log.Fatal(err)
	}
	createTestDb()
}

func createTestDb() {
	err := data.CreateDirIfNotExist(persistentTestDataDirPath)
	if err != nil {
		log.Fatal(err)
	}
	dataProvider := data.NewBoltProvider(testDbPath)
	err = dataProvider.SaveEntries(data.GetTestEntries())
	if err != nil {
		log.Fatal(err)
	}
}

func cleanupTestEnvironment() {
	err := data.DeleteDirWithContents(testDataDirPath)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanupAfterTestRun() {
	removeAllNonPersistentFilesInTestDataDir()
}

func removeAllNonPersistentFilesInTestDataDir() {
	err := data.DeleteAllInDirExceptForDirs(testDataDirPath, "persistent")
	if err != nil {
		log.Fatal(err)
	}
}

func areFilesInPathsTheSame(filePath1 string, filePath2 string) bool {
	file1, err := os.Open(filePath1)
	if err != nil {
		log.Fatal(err)
	}
	defer file1.Close()
	file2, err := os.Open(filePath2)
	if err != nil {
		log.Fatal(err)
	}
	defer file2.Close()
	return areFilesTheSame(file1, file2)
}

func areFilesTheSame(file1 *os.File, file2 *os.File) bool {
	const chunkSize = 4000
	for {
		bytesOfFile1 := make([]byte, chunkSize)
		_, err1 := file1.Read(bytesOfFile1)
		bytesOfFile2 := make([]byte, chunkSize)
		_, err2 := file2.Read(bytesOfFile2)
		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Error(err1)
				log.Fatal(err2)
			}
		}
		if !bytes.Equal(bytesOfFile1, bytesOfFile2) {
			return false
		}
	}
	return true
}

func createCorrectSavedWirwlConfigFileInPath(path string) {
	err := data.CreateDirIfNotExist(path)
	if err != nil {
		log.Fatal(err)
	}
	config := Config{
		AppDataDirPath: testAppDataDirPath,
		ConfigDirPath:  testConfigDirPath,
	}
	config.saveConfigIn(path + "wirwl_correct.cfg")
}

func createCorrectWirwlConfigFileForLoadingInPath(path string) {
	err := data.CreateDirIfNotExist(path)
	if err != nil {
		log.Fatal(err)
	}
	config := Config{
		AppDataDirPath: "some db path",
		ConfigDirPath:  testConfigDirPath,
	}
	config.saveConfigIn(path + "wirwl.cfg")
}

func (config *Config) saveConfigIn(configFilePath string) {
	configFile, err := os.OpenFile(configFilePath, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatal(err)
	}
	err = toml.NewEncoder(configFile).Encode(config)
	if err != nil {
		log.Fatal(err)
	}
	err = configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
}
