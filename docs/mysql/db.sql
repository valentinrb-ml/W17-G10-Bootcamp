-- Creación de la Base de Datos
CREATE DATABASE IF NOT EXISTS db_warehouse
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;
USE db_warehouse;
-- Tabla: countries
CREATE TABLE countries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    country_name VARCHAR(255) NOT NULL
);
-- Tabla: provinces
CREATE TABLE provinces (
    id INT AUTO_INCREMENT PRIMARY KEY,
    province_name VARCHAR(255) NOT NULL,
    id_country_fk INT NOT NULL
);
-- Tabla: localities
CREATE TABLE localities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    locality_name VARCHAR(255) NOT NULL,
    province_id INT NOT NULL
);
-- Tabla: sellers
CREATE TABLE sellers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid INT NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    telephone VARCHAR(255),
    locality_id INT NOT NULL
);
-- Tabla: carriers
CREATE TABLE carriers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    telephone VARCHAR(255),
    locality_id INT NOT NULL
);
-- Tabla: buyers
CREATE TABLE buyers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_card_number VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL
);
-- Tabla: warehouse
CREATE TABLE warehouse (
    id INT AUTO_INCREMENT PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    telephone VARCHAR(255),
    warehouse_code VARCHAR(255) NOT NULL,
    locality_id INT NOT NULL
);
-- Tabla: employees
CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    id_card_number VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    warehouse_id INT NOT NULL
);
-- Tabla: products_types
CREATE TABLE products_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(255) NOT NULL
);
-- Tabla: products
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(255) NOT NULL,
    expiration_rate DECIMAL(19,2),
    freezing_rate DECIMAL(19,2),
    height DECIMAL(19,2),
    length DECIMAL(19,2),
    net_weight DECIMAL(19,2),
    product_code VARCHAR(255) NOT NULL,
    recommended_freezing_temperature DECIMAL(19,2),
    width DECIMAL(19,2),
    product_type_id INT NOT NULL,
    seller_id INT NOT NULL
);
-- Tabla: sections
CREATE TABLE sections (
    id INT AUTO_INCREMENT PRIMARY KEY,
    section_number VARCHAR(255) NOT NULL,
    current_capacity INT,
    current_temperature DECIMAL(19,2),
    maximum_capacity INT,
    minimum_capacity INT,
    minimum_temperature DECIMAL(19,2),
    product_type_id INT NOT NULL,
    warehouse_id INT NOT NULL
);
-- Tabla: order_status
CREATE TABLE order_status (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(255) NOT NULL
);
-- Tabla: purchase_orders
CREATE TABLE purchase_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_number VARCHAR(255) NOT NULL,
    order_date DATETIME(6) NOT NULL,
    tracking_code VARCHAR(255),
    buyer_id INT NOT NULL,
    carrier_id INT NOT NULL,
    order_status_id INT NOT NULL,
    warehouse_id INT NOT NULL
);
-- Tabla: product_batches
CREATE TABLE product_batches (
    id INT AUTO_INCREMENT PRIMARY KEY,
    batch_number VARCHAR(255) NOT NULL,
    current_quantity INT,
    current_temperature DECIMAL(19,2),
    due_date DATETIME(6),
    initial_quantity INT,
    manufacturing_date DATETIME(6),
    manufacturing_hour DATETIME(6),
    minimum_temperature DECIMAL(19,2),
    product_id INT NOT NULL,
    section_id INT NOT NULL
);
-- Tabla: inbound_orders
CREATE TABLE inbound_orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_date DATETIME(6) NOT NULL,
    order_number VARCHAR(255) NOT NULL,
    employee_id INT NOT NULL,
    product_batch_id INT NOT NULL,
    warehouse_id INT NOT NULL
);
-- Tabla: product_records
CREATE TABLE product_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    last_update_date DATETIME(6),
    purchase_price DECIMAL(19,2),
    sale_price DECIMAL(19,2),
    product_id INT NOT NULL
);
-- Tabla: order_details
CREATE TABLE order_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    clean_liness_status VARCHAR(255),
    quantity INT,
    temperature DECIMAL(19,2),
    product_record_id INT NOT NULL,
    purchase_order_id INT NOT NULL
);

-- Índices y Claves Foráneas
-- Provincias -> countries
ALTER TABLE provinces
ADD CONSTRAINT fk_provinces_country
FOREIGN KEY(id_country_fk) REFERENCES countries(id);
-- Localities -> provinces
ALTER TABLE localities
ADD CONSTRAINT fk_localities_province
FOREIGN KEY(province_id) REFERENCES provinces(id);
-- Sellers -> localities
ALTER TABLE sellers
ADD CONSTRAINT fk_sellers_locality
FOREIGN KEY(locality_id) REFERENCES localities(id);
-- Carriers -> localities
ALTER TABLE carriers
ADD CONSTRAINT fk_carriers_locality
FOREIGN KEY(locality_id) REFERENCES localities(id);
-- Warehouse -> localities
ALTER TABLE warehouse
ADD CONSTRAINT fk_warehouse_locality
FOREIGN KEY(locality_id) REFERENCES localities(id);
-- Employees -> warehouse
ALTER TABLE employees
ADD CONSTRAINT fk_employees_warehouse
FOREIGN KEY(warehouse_id) REFERENCES warehouse(id);
-- Products -> products_types
ALTER TABLE products
ADD CONSTRAINT fk_products_type
FOREIGN KEY(product_type_id) REFERENCES products_types(id);
-- Products -> sellers
ALTER TABLE products
ADD CONSTRAINT fk_products_seller
FOREIGN KEY(seller_id) REFERENCES sellers(id);
-- Sections -> products_types
ALTER TABLE sections
ADD CONSTRAINT fk_sections_product_type
FOREIGN KEY(product_type_id) REFERENCES products_types(id);
-- Sections -> warehouse
ALTER TABLE sections
ADD CONSTRAINT fk_sections_warehouse
FOREIGN KEY(warehouse_id) REFERENCES warehouse(id);
-- Purchase_orders -> buyers
ALTER TABLE purchase_orders
ADD CONSTRAINT fk_purchase_orders_buyer
FOREIGN KEY(buyer_id) REFERENCES buyers(id);
-- Purchase_orders -> carriers
ALTER TABLE purchase_orders
ADD CONSTRAINT fk_purchase_orders_carrier
FOREIGN KEY(carrier_id) REFERENCES carriers(id);
-- Purchase_orders -> order_status
ALTER TABLE purchase_orders
ADD CONSTRAINT fk_purchase_orders_status
FOREIGN KEY(order_status_id) REFERENCES order_status(id);
-- Purchase_orders -> warehouse
ALTER TABLE purchase_orders
ADD CONSTRAINT fk_purchase_orders_warehouse
FOREIGN KEY(warehouse_id) REFERENCES warehouse(id);
-- Product_batches -> products
ALTER TABLE product_batches
ADD CONSTRAINT fk_product_batches_product
FOREIGN KEY(product_id) REFERENCES products(id);
-- Product_batches -> sections
ALTER TABLE product_batches
ADD CONSTRAINT fk_product_batches_section
FOREIGN KEY(section_id) REFERENCES sections(id);
-- Inbound_orders -> employees
ALTER TABLE inbound_orders
ADD CONSTRAINT fk_inbound_orders_employee
FOREIGN KEY(employee_id) REFERENCES employees(id);
-- Inbound_orders -> product_batches
ALTER TABLE inbound_orders
ADD CONSTRAINT fk_inbound_orders_batch
FOREIGN KEY(product_batch_id) REFERENCES product_batches(id);
-- Inbound_orders -> warehouse
ALTER TABLE inbound_orders
ADD CONSTRAINT fk_inbound_orders_warehouse
FOREIGN KEY(warehouse_id) REFERENCES warehouse(id);
-- Product_records -> products
ALTER TABLE product_records
ADD CONSTRAINT fk_product_records_product
FOREIGN KEY(product_id) REFERENCES products(id);
-- Order_details -> product_records
ALTER TABLE order_details
ADD CONSTRAINT fk_order_details_product_record
FOREIGN KEY(product_record_id) REFERENCES product_records(id);
-- Order_details -> purchase_orders
ALTER TABLE order_details
ADD CONSTRAINT fk_order_details_purchase_order
FOREIGN KEY(purchase_order_id) REFERENCES purchase_orders(id);
