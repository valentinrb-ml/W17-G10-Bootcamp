# RUN PROJECT
# Runs the project from cmd/main.go
.PHONY: run
run:
	go run cmd/main.go

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

# CARRY MODULE COVERAGE
# Runs tests and shows coverage for carry module (service, handler, repository)
.PHONY: cover-carry
cover-carry:
	go test ./internal/service/carry/... ./internal/handler/carry/... ./internal/repository/carry/... -coverprofile=carry_coverage.out && \
	go tool cover -func=carry_coverage.out

# WAREHOUSE MODULE COVERAGE
# Runs tests and shows coverage for warehouse module (service, handler, repository)
.PHONY: cover-warehouse
cover-warehouse:
	go test ./internal/service/warehouse/... ./internal/handler/warehouse/... ./internal/repository/warehouse/... -coverprofile=warehouse_coverage.out && \
	go tool cover -func=warehouse_coverage.out

# SELLER MODULE COVERAGE
# Runs tests and shows coverage for seller module (service, handler, repository)
.PHONY: cover-seller
cover-seller:
	go test ./internal/service/seller/... ./internal/handler/seller/... ./internal/repository/seller/... -coverprofile=seller_coverage.out && \
	go tool cover -func=seller_coverage.out

# GEOGRAPHY MODULE COVERAGE
# Runs tests and shows coverage for geography module (service, handler, repository)
.PHONY: cover-geography
cover-geography:
	go test ./internal/service/geography/... ./internal/handler/geography/... ./internal/repository/geography/... -coverprofile=geography_coverage.out && \
	go tool cover -func=geography_coverage.out

# SECTION MODULE COVERAGE
# Runs tests and shows coverage for section module (service, handler, repository)
.PHONY: cover-section
cover-section:
	go test ./internal/service/section/... ./internal/handler/section/... ./internal/repository/section/... -coverprofile=section_coverage.out && \
	go tool cover -func=section_coverage.out
	go tool cover -html=section_coverage.out

# PRODUCT BATCH MODULE COVERAGE
# Runs tests and shows coverage for product batch module (service, handler, repository)
.PHONY: cover-product-batch
cover-product-batch:
	go test ./internal/service/product_batch/... ./internal/handler/product_batch/... ./internal/repository/product_batch/... -coverprofile=product_batch_coverage.out && \
	go tool cover -func=product_batch_coverage.out
	go tool cover -html=product_batch_coverage.out

# EMPLOYEE MODULE COVERAGE
# Runs tests and shows coverage for employee module (service, handler, repository)
.PHONY: cover-employee
cover-employee:
	go test ./internal/service/employee/... ./internal/handler/employee/... ./internal/repository/employee/... -coverprofile=employee_coverage.out && \
	go tool cover -func=employee_coverage.out

# INBOUND_ORDER MODULE COVERAGE
# Runs tests and shows coverage for inbound_order module (service, handler, repository)
.PHONY: cover-inbound_order
cover-inbound_order:
	go test ./internal/service/inbound_order/... ./internal/handler/inbound_order/... ./internal/repository/inbound_order/... -coverprofile=inbound_order_coverage.out && \
	go tool cover -func=inbound_order_coverage.out

# BUYER MODULE COVERAGE
# Runs tests and shows coverage for buyuer module (service, handler, repository)
.PHONY: cover-carry
cover-carry:
	go test ./internal/service/buyer/... ./internal/handler/buyer/... ./internal/repository/buyer/... -coverprofile=buyer_coverage.out && \
	go tool cover -func=buyer_coverage.out

# PURCHASE_ORDER MODULE COVERAGE
# Runs tests and shows coverage for purchase_order module (service, handler, repository)
.PHONY: cover-warehouse
cover-warehouse:
	go test ./internal/service/purchase_order/... ./internal/handler/purchase_order/... ./internal/repository/purchase_order/... -coverprofile=purchase_order_coverage.out && \
	go tool cover -func=purchase_order_coverage.out
