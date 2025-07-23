# RUN ALL TESTS
# Runs all tests in the current directory and subdirectories, showing detailed output.
.PHONY: test
test:
	go test -v ./...

# GET COVERAGE FILE
# Runs all tests, calculates code coverage, and creates a 'coverage.out' report file.
.PHONY: cover
cover:
	go test -cover -coverprofile=coverage.out ./... && \
	grep -v '_stub.go' coverage.out | grep -v '_mock.go' > coverage.tmp && \
	mv coverage.tmp coverage.out

# OPEN COVERAGE FILE IN BROWSER
# Shows the coverage report in HTML format in your browser.
# This target depends on 'cover', so it first generates the coverage file.
.PHONY: cover-html
cover-html: cover
	go tool cover -html=coverage.out

# SHOW COVERAGE SUMMARY
# Shows a summary of code coverage (per function and total percentage) in the terminal.
# Depends on 'cover', so it generates the coverage file first.
.PHONY: cover-summary
cover-summary: cover
	go tool cover -func=coverage.out
