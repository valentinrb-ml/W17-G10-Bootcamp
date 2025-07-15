INSERT INTO countries (id, country_name) VALUES
                                             (1, 'Argentina'), (2, 'Brasil'), (3, 'Chile'), (4, 'Uruguay'), (5, 'Paraguay'),
                                             (6, 'Bolivia'), (7, 'Perú'), (8, 'Ecuador'), (9, 'Colombia'), (10, 'Venezuela');
INSERT INTO provinces (id, province_name, id_country_fk) VALUES
                                                             (1, 'Buenos Aires', 1), (2, 'Córdoba', 1), (3, 'Santa Fe', 1), (4, 'Mendoza', 1),
                                                             (5, 'San Pablo', 2), (6, 'Rio de Janeiro', 2), (7, 'Antofagasta', 3), (8, 'Montevideo', 4),
                                                             (9, 'Asunción', 5), (10, 'La Paz', 6);
INSERT INTO localities (id, locality_name, province_id) VALUES
                                                            (1, 'La Plata', 1), (2, 'Córdoba Capital', 2), (3, 'Rosario', 3), (4, 'Godoy Cruz', 4),
                                                            (5, 'Campinas', 5), (6, 'Niterói', 6), (7, 'Calama', 7), (8, 'Centro', 8),
                                                            (9, 'Lambaré', 9), (10, 'El Alto', 10);
INSERT INTO sellers (id, cid, company_name, address, telephone, locality_id) VALUES
                                                                                 (1, 101, 'Frutas del Sur', 'Calle 1', '221-111', 1),
                                                                                 (2, 102, 'Verdulería Norte', 'Calle 2', '221-112', 2),
                                                                                 (3, 103, 'Carnes Argentinas', 'Calle 3', '221-113', 3),
                                                                                 (4, 104, 'Almacén Cordobés', 'Calle 4', '221-114', 4),
                                                                                 (5, 105, 'Exportadora Brasil', 'Calle 5', '11-221', 5),
                                                                                 (6, 106, 'Café do Brasil', 'Calle 6', '21-222', 6),
                                                                                 (7, 107, 'Viña Andina', 'Calle 7', '32-333', 7),
                                                                                 (8, 108, 'Delicatessen Uy', 'Calle 8', '34-444', 8),
                                                                                 (9, 109, 'Maíz PY', 'Calle 9', '41-555', 9),
                                                                                 (10, 110, 'Bolivian Imports', 'Calle 10', '51-666', 10);
INSERT INTO carriers (id, cid, company_name, address, telephone, locality_id) VALUES
                                                                                  (1, 'C001', 'Transporte Sureño', 'Av 10', '421-001', 1),
                                                                                  (2, 'C002', 'Logística Pampeana', 'Av 2', '421-002', 2),
                                                                                  (3, 'C003', 'Carga Express', 'Av 3', '421-003', 3),
                                                                                  (4, 'C004', 'Rápido Mendoza', 'Av 4', '421-004', 4),
                                                                                  (5, 'C005', 'Transportes Brasil', 'Rua 5', '31-101', 5),
                                                                                  (6, 'C006', 'Rio Cargo', 'Rua 6', '33-202', 6),
                                                                                  (7, 'C007', 'Chile Express', 'Av Chile', '56-303', 7),
                                                                                  (8, 'C008', 'Uy Delivery', 'Av U', '12-904', 8),
                                                                                  (9, 'C009', 'PY Truck', 'Ruta PY', '37-105', 9),
                                                                                  (10, 'C010', 'Cargas Bolívar', 'Av Bolivia', '53-207', 10);
INSERT INTO buyers (id, id_card_number, first_name, last_name) VALUES
                                                                   (1, '4001', 'Ana', 'Pérez'), (2, '4002', 'Bernardo', 'Gómez'), (3, '4003', 'Camila', 'Ríos'),
                                                                   (4, '4004', 'David', 'Silva'), (5, '4005', 'Esteban', 'Arce'), (6, '4006', 'Felipe', 'Sosa'),
                                                                   (7, '4007', 'Gabriela', 'Campos'), (8, '4008', 'Hugo', 'Castro'), (9, '4009', 'Irene', 'Fernández'),
                                                                   (10, '4010', 'Joaquín', 'de la Vega');
INSERT INTO warehouse (id, address, telephone, warehouse_code, minimum_capacity, minimum_temperature, locality_id) VALUES
                                                                                (1, 'Depósito Sur', '155-201', 'WS-001', 100, -10.5, 1),
                                                                                (2, 'Bodega Central', '155-202', 'WS-002', 150, -8.0, 2),
                                                                                (3, 'Frigorifico Norte', '155-203', 'WS-003', 200, -18.0, 3),
                                                                                (4, 'Almacén Andes', '155-204', 'WS-004', 80, -5.0, 4),
                                                                                (5, 'Galpón SP', '155-205', 'WS-005', 120, -12.0, 5),
                                                                                (6, 'Almacen Rj', '155-206', 'WS-006', 90, -7.5, 6),
                                                                                (7, 'Depósito Calama', '155-207', 'WS-007', 110, -15.0, 7),
                                                                                (8, 'Almacen CentroUy', '155-208', 'WS-008', 95, -6.0, 8),
                                                                                (9, 'Depósito Lambaré', '155-209', 'WS-009', 130, -9.0, 9),
                                                                                (10, 'Bodega Alto', '155-210', 'WS-010', 140, -11.0, 10);
INSERT INTO employees (id, id_card_number, first_name, last_name, warehouse_id) VALUES
                                                                                    (1, 'E001', 'Lucas', 'Martínez', 1),
                                                                                    (2, 'E002', 'Martina', 'García', 2),
                                                                                    (3, 'E003', 'Pedro', 'Suárez', 3),
                                                                                    (4, 'E004', 'Raúl', 'Pereyra', 4),
                                                                                    (5, 'E005', 'Nadia', 'Lagos', 5),
                                                                                    (6, 'E006', 'Carlos', 'Montero', 6),
                                                                                    (7, 'E007', 'Rosa', 'Duarte', 7),
                                                                                    (8, 'E008', 'Sergio', 'Paz', 8),
                                                                                    (9, 'E009', 'Tomás', 'Yáñez', 9),
                                                                                    (10, 'E010', 'Viviana', 'Reyes', 10);
INSERT INTO products_types (id, description) VALUES
                                                 (1, 'Frutas'), (2, 'Carnes'), (3, 'Verduras'), (4, 'Bebidas'), (5, 'Lácteos'),
                                                 (6, 'Panadería'), (7, 'Enlatados'), (8, 'Limpieza'), (9, 'Snacks'), (10, 'Granos');

INSERT INTO products (id, description, expiration_rate, freezing_rate, height, length, net_weight, product_code, recommended_freezing_temperature, width, product_type_id, seller_id)
VALUES
    (1, 'Manzanas', 30, 0, 10, 15, 5.0, 'P001', 0, 10, 1, 1),
    (2, 'Carne Res', 10, -18, 12, 32, 20.5, 'P002', -18, 18, 2, 3),
    (3, 'Espinaca', 20, 0, 4, 10, 1.5, 'P003', 0, 8, 3, 2),
    (4, 'Jugo Naranja', 180, 2, 23, 9, 2.5, 'P004', 2, 9, 4, 5),
    (5, 'Leche Entera', 15, 4, 23, 28, 15, 'P005', 4, 16, 5, 6),
    (6, 'Pan Baguette', 5, 0, 42, 10, 0.9, 'P006', 0, 7, 6, 8),
    (7, 'Atún Lata', 365, 0, 7, 15, 0.4, 'P007', 0, 4, 7, 4),
    (8, 'Detergente', 730, 0, 25, 8, 2, 'P008', 0, 8, 8, 7),
    (9, 'Almendras', 365, 0, 2, 3, 0.5, 'P009', 0, 4, 9, 9),
    (10, 'Arroz', 730, 0, 8, 14, 9.0, 'P010', 0, 10, 10, 10);
INSERT INTO sections (id, section_number, current_capacity, current_temperature, maximum_capacity, minimum_capacity, minimum_temperature, product_type_id, warehouse_id)
VALUES
    (1, 1, 100, 5, 200, 50, 4, 1, 1),
    (2, 2, 120, 1, 200, 60, 0, 2, 2),
    (3, 3, 70, 5, 100, 30, 5, 3, 3),
    (4, 4, 90, 8, 120, 40, 6, 4, 4),
    (5, 5, 50, 4, 80, 20, 4, 5, 5),
    (6, 6, 200, 6, 400, 100, 5, 6, 6),
    (7, 7, 80, 12, 150, 30, 10, 7, 7),
    (8, 8, 60, 15, 200, 40, 12, 8, 8),
    (9, 9, 40, 3, 70, 10, 3, 9, 9),
    (10, 10, 180, 7, 250, 70, 6, 10, 10);
INSERT INTO order_status (id, description) VALUES
                                               (1, 'Pendiente'), (2, 'Confirmada'), (3, 'Cancelada'), (4, 'En reparto'),
                                               (5, 'Entregada'), (6, 'En preparación'), (7, 'Facturada'), (8, 'Devuelta'),
                                               (9, 'Revisión'), (10, 'Cerrada');
INSERT INTO purchase_orders (id, order_number, order_date, tracking_code, buyer_id, carrier_id, order_status_id, warehouse_id)
VALUES
    (1, 'PO001', '2024-01-01 10:00:00', 'TK001', 1, 1, 1, 1),
    (2, 'PO002', '2024-01-02 11:30:00', 'TK002', 2, 2, 2, 2),
    (3, 'PO003', '2024-01-03 09:15:00', 'TK003', 3, 3, 1, 3),
    (4, 'PO004', '2024-01-04 15:20:00', 'TK004', 4, 4, 4, 4),
    (5, 'PO005', '2024-01-05 17:45:00', 'TK005', 5, 5, 5, 5),
    (6, 'PO006', '2024-01-06 13:10:00', 'TK006', 6, 6, 1, 6),
    (7, 'PO007', '2024-01-07 14:00:00', 'TK007', 7, 7, 2, 7),
    (8, 'PO008', '2024-01-08 08:30:00', 'TK008', 8, 8, 6, 8),
    (9, 'PO009', '2024-01-09 16:10:00', 'TK009', 9, 9, 4, 9),
    (10, 'PO010', '2024-01-10 10:50:00', 'TK010', 10, 10, 1, 10);
INSERT INTO product_batches (id, batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id)
VALUES
    (1, 1, 50, 3, '2024-06-10 00:00:00', 70, '2024-05-10 00:00:00', '2024-05-10 08:00:00', 2, 1, 1),
    (2, 2, 20, -18, '2024-08-01 00:00:00', 25, '2024-06-15 00:00:00', '2024-06-15 06:00:00', -15, 2, 2),
    (3, 3, 32, 5, '2024-06-20 00:00:00', 40, '2024-05-20 00:00:00', '2024-05-20 12:00:00', 4, 3, 3),
    (4, 4, 10, 2, '2024-10-10 00:00:00', 15, '2024-07-15 00:00:00', '2024-07-15 11:00:00', 2, 4, 4),
    (5, 5, 25, 4, '2024-07-15 00:00:00', 30, '2024-06-10 00:00:00', '2024-06-10 09:30:00', 4, 5, 5),
    (6, 6, 200, 6, '2024-12-20 00:00:00', 210, '2024-06-25 00:00:00', '2024-06-25 15:25:00', 5, 6, 6),
    (7, 7, 75, 11, '2025-01-15 00:00:00', 80, '2024-06-10 00:00:00', '2024-06-10 13:10:00', 8, 7, 7),
    (8, 8, 56, 16, '2024-12-01 00:00:00', 57, '2024-06-03 00:00:00', '2024-06-03 10:30:00', 14, 8, 8),
    (9, 9, 39, 3, '2024-06-05 00:00:00', 44, '2024-05-01 00:00:00', '2024-05-01 07:10:00', 2, 9, 9),
    (10, 10, 88, 7, '2024-11-21 00:00:00', 90, '2024-05-25 00:00:00', '2024-05-25 08:15:00', 7, 10, 10);

INSERT INTO inbound_orders (id, order_date, order_number, employee_id, product_batch_id, warehouse_id)
VALUES
    (1, '2024-05-10 09:00:00', 'IN001', 1, 1, 1),
    (2, '2024-06-15 07:30:00', 'IN002', 2, 2, 2),
    (3, '2024-05-20 13:15:00', 'IN003', 3, 3, 3),
    (4, '2024-07-15 12:00:00', 'IN004', 4, 4, 4),
    (5, '2024-06-10 11:00:00', 'IN005', 5, 5, 5),
    (6, '2024-06-25 16:00:00', 'IN006', 6, 6, 6),
    (7, '2024-06-10 14:00:00', 'IN007', 7, 7, 7),
    (8, '2024-06-03 11:10:00', 'IN008', 8, 8, 8),
    (9, '2024-05-01 08:00:00', 'IN009', 9, 9, 9),
    (10, '2024-05-25 10:30:00', 'IN010', 10, 10, 10);

INSERT INTO product_records (id, last_update_date, purchase_price, sale_price, product_id)
VALUES
    (1, '2024-05-10 09:30:00', 10, 14, 1),
    (2, '2024-06-15 08:00:00', 32, 42, 2),
    (3, '2024-05-20 15:00:00', 3, 4, 3),
    (4, '2024-07-15 14:00:00', 8, 12, 4),
    (5, '2024-06-10 12:00:00', 12, 15, 5),
    (6, '2024-06-25 17:00:00', 1, 3, 6),
    (7, '2024-06-10 15:00:00', 25, 31, 7),
    (8, '2024-06-03 12:00:00', 10, 14, 8),
    (9, '2024-05-01 09:00:00', 23, 28, 9),
    (10, '2024-05-25 11:30:00', 3, 5, 10);

INSERT INTO order_details (id, clean_liness_status, quantity, temperature, product_record_id, purchase_order_id)
VALUES
    (1, 'Limpio', 5, 3, 1, 1),
    (2, 'Satisfactorio', 3, -18, 2, 2),
    (3, 'Bueno', 2, 5, 3, 3),
    (4, 'Excelente', 4, 2, 4, 4),
    (5, 'Correcto', 6, 4, 5, 5),
    (6, 'Normal', 9, 6, 6, 6),
    (7, 'Aceptable', 8, 11, 7, 7),
    (8, 'Perfecto', 12, 16, 8, 8),
    (9, 'Bueno', 4, 3, 9, 9),
    (10, 'Limpio', 10, 7, 10, 10);
