--  Hospital, Pharmacy & Blood Bank Management System
-- Doctor
INSERT INTO Doctor (ID, Name, Email, P_No, Specialization) VALUES
(1,  'Dr. Ahmed Raza',       'ahmed.raza@hospital.pk',       '0300-1111111', 'Cardiology'),
(2,  'Dr. Sara Khan',        'sara.khan@hospital.pk',         '0300-2222222', 'Neurology'),
(3,  'Dr. Bilal Awan',       'bilal.awan@hospital.pk',        '0300-3333333', 'Orthopedics'),
(4,  'Dr. Ayesha Malik',     'ayesha.malik@hospital.pk',      '0300-4444444', 'Pediatrics'),
(5,  'Dr. Usman Tariq',      'usman.tariq@hospital.pk',       '0300-5555555', 'General Surgery'),
(6,  'Dr. Hina Javed',       'hina.javed@hospital.pk',        '0300-6666666', 'Dermatology'),
(7,  'Dr. Zeeshan Qureshi',  'zeeshan.qureshi@hospital.pk',  '0300-7777777', 'Psychiatry'),
(8,  'Dr. Nadia Hussain',    'nadia.hussain@hospital.pk',     '0300-8888888', 'Oncology'),
(9,  'Dr. Farhan Siddiqui',  'farhan.siddiqui@hospital.pk',  '0300-9999999', 'Emergency Medicine'),
(10, 'Dr. Rabia Noor',       'rabia.noor@hospital.pk',        '0300-1010101', 'Endocrinology'),
(11, 'Dr. Kamran Baig',      'kamran.baig@hospital.pk',       '0300-1111112', 'Nephrology'),
(12, 'Dr. Sadia Iqbal',      'sadia.iqbal@hospital.pk',       '0300-1212121', 'Gynecology');


--  RECEPTIONIST
INSERT INTO Receptionist (ID, Name, Email, P_No) VALUES
(1, 'Sana Mehmood',   'sana.mehmood@hospital.pk',    '0311-1111111'),
(2, 'Kamran Ali',     'kamran.ali@hospital.pk',       '0311-2222222'),
(3, 'Noor Fatima',    'noor.fatima@hospital.pk',      '0311-3333333'),
(4, 'Tariq Bashir',   'tariq.bashir@hospital.pk',     '0311-4444444'),
(5, 'Amna Shahid',    'amna.shahid@hospital.pk',      '0311-5555555'),
(6, 'Bilal Yousaf',   'bilal.yousaf@hospital.pk',     '0311-6666666');

--   PHARMACIST

INSERT INTO Pharmacist (ID, Name, Email, P_No) VALUES
(1, 'Asif Iqbal',      'asif.iqbal@hospital.pk',      '0322-1111111'),
(2, 'Mehwish Karim',   'mehwish.karim@hospital.pk',   '0322-2222222'),
(3, 'Omer Farooq',     'omer.farooq@hospital.pk',     '0322-3333333'),
(4, 'Lubna Sheikh',    'lubna.sheikh@hospital.pk',    '0322-4444444'),
(5, 'Rizwan Baig',     'rizwan.baig@hospital.pk',     '0322-5555555'),
(6, 'Huma Zafar',      'huma.zafar@hospital.pk',      '0322-6666666');

--  BLOOD_MANAGER 

INSERT INTO Blood_Manager (ID, Name, Email, P_No) VALUES
(1, 'Zahid Mehmood',   'zahid.mehmood@hospital.pk',   '0333-1111111'),
(2, 'Farzana Bibi',    'farzana.bibi@hospital.pk',    '0333-2222222'),
(3, 'Imran Chaudhry',  'imran.chaudhry@hospital.pk',  '0333-3333333'),
(4, 'Saima Rehman',    'saima.rehman@hospital.pk',    '0333-4444444');

--   PATIENT  

INSERT INTO Patient (ID, Name, Email, B_Gr, P_No) VALUES
(1,  'Ali Hassan',        'ali.hassan@gmail.com',        'A+',  '0345-1111111'),  -- A+
(2,  'Fatima Zahra',      'fatima.zahra@gmail.com',      'B+',  '0345-2222222'),  -- B+
(3,  'Muhammad Usman',    'usman.m@gmail.com',            'O+',  '0345-3333333'),  -- O+
(4,  'Zainab Mirza',      'zainab.mirza@gmail.com',      'AB+', '0345-4444444'),  -- AB+
(5,  'Hassan Riaz',       'hassan.riaz@gmail.com',        'A-',  '0345-5555555'),  -- A-
(6,  'Mariam Aslam',      'mariam.aslam@gmail.com',      'B-',  '0345-6666666'),  -- B-
(7,  'Saad Nawaz',        'saad.nawaz@gmail.com',         'O-',  '0345-7777777'),  -- O-
(8,  'Hira Yousaf',       'hira.yousaf@gmail.com',        'A+',  '0345-8888888'),  -- A+
(9,  'Junaid Bhatti',     'junaid.bhatti@gmail.com',      'B+',  '0345-9999999'),  -- B+
(10, 'Sobia Anwar',       'sobia.anwar@gmail.com',        'O+',  '0345-1010101'),  -- O+
(11, 'Adeel Farhan',      'adeel.farhan@gmail.com',       'AB-', '0345-1111112'),  -- AB-
(12, 'Raheela Pervaiz',   'raheela.pervaiz@gmail.com',   'A+',  '0345-1212121'),  -- A+
(13, 'Kashif Mahmood',    'kashif.mahmood@gmail.com',    'B+',  '0345-1313131'),  -- B+
(14, 'Iqra Siddiqui',     'iqra.siddiqui@gmail.com',     'O+',  '0345-1414141'),  -- O+
(15, 'Naveed Sultan',     'naveed.sultan@gmail.com',      'A-',  '0345-1515151'),  -- A-
(16, 'Amna Butt',         'amna.butt@gmail.com',          'B-',  '0345-1616161'),  -- B-
(17, 'Fahad Rehman',      'fahad.rehman@gmail.com',       'O-',  '0345-1717171'),  -- O-
(18, 'Sana Rashid',       'sana.rashid@gmail.com',        'AB+', '0345-1818181'),  -- AB+
(19, 'Umar Hayat',        'umar.hayat@gmail.com',         'A+',  '0345-1919191'),  -- A+
(20, 'Nadia Waheed',      'nadia.waheed@gmail.com',       'B+',  '0345-2020202');  -- B+



--   ER_PATIENT  

INSERT INTO ER_Patient (ID, Name, Age, P_No, Arrival_Time) VALUES
(1,  'Shahid Latif',    45, '0301-1111111', '08:30:00'),
(2,  'Rukhsana Bibi',   62, '0301-2222222', '09:15:00'),
(3,  'Hamza Akbar',     28, '0301-3333333', '10:00:00'),
(4,  'Perveen Akhtar',  55, '0301-4444444', '11:45:00'),
(5,  'Babar Zaman',     33, '0301-5555555', '13:20:00'),
(6,  'Samina Gul',      41, '0301-6666666', '14:05:00'),
(7,  'Waqas Rehman',    19, '0301-7777777', '15:30:00'),
(8,  'Nasreen Begum',   70, '0301-8888888', '16:50:00'),
(9,  'Tariq Mehmood',   37, '0301-9999999', '07:45:00'),
(10, 'Zara Qureshi',    24, '0301-1010101', '12:00:00'),
(11, 'Imran Ghani',     52, '0301-1111112', '17:30:00'),
(12, 'Bushra Naz',      66, '0301-1212121', '20:00:00');



--   DONOR  

INSERT INTO Donor (ID, Name, P_No, Email, B_Gr, Last_Donate) VALUES
(1,  'Khalid Mehmood',   '0312-1111111', 'khalid.m@gmail.com',     'A+',  '2025-10-01'),
(2,  'Shazia Parveen',   '0312-2222222', 'shazia.p@gmail.com',     'B+',  '2025-09-15'),
(3,  'Tariq Javed',      '0312-3333333', 'tariq.j@gmail.com',      'O+',  '2025-08-20'),
(4,  'Maryam Riaz',      '0312-4444444', 'maryam.r@gmail.com',     'AB+', '2025-07-10'),
(5,  'Zubair Ahmed',     '0312-5555555', 'zubair.a@gmail.com',     'A-',  '2025-10-05'),
(6,  'Faisal Iqbal',     '0312-6666666', 'faisal.i@gmail.com',     'B-',  '2025-09-01'),
(7,  'Amna Rauf',        '0312-7777777', 'amna.rauf@gmail.com',    'O-',  '2025-08-11'),
(8,  'Sajid Nawaz',      '0312-8888888', 'sajid.n@gmail.com',      'O+',  '2025-07-25'),
(9,  'Kiran Shaheen',    '0312-9999999', 'kiran.s@gmail.com',      'A+',  '2025-10-10'),
(10, 'Waseem Akram',     '0312-1010101', 'waseem.ak@gmail.com',    'B+',  '2025-09-20'),
(11, 'Rabia Tahir',      '0312-1111112', 'rabia.t@gmail.com',      'AB-', '2025-08-01'),
(12, 'Usman Ghani',      '0312-1212121', 'usman.g@gmail.com',      'O-',  '2025-07-15'),
(13, 'Nazia Saleem',     '0312-1313131', 'nazia.s@gmail.com',      'A+',  '2025-09-05'),
(14, 'Arif Hussain',     '0312-1414141', 'arif.h@gmail.com',       'B-',  '2025-10-15'),
(15, 'Saira Batool',     '0312-1515151', 'saira.b@gmail.com',      'O+',  '2025-08-30');



--   BLOOD  

INSERT INTO Blood (ID, B_Gr, Collected_Date, Status, Expiry_Date, Unit, Donor_ID) VALUES
(1,  'A+',  '2025-11-01', 'available', '2026-05-01', 8,  1),
(2,  'B+',  '2025-11-03', 'available', '2026-05-03', 6,  2),
(3,  'O+',  '2025-11-05', 'available', '2026-05-05', 12, 3),
(4,  'AB+', '2025-11-07', 'available', '2026-05-07', 4,  4),
(5,  'A-',  '2025-11-09', 'available', '2026-05-09', 5,  5),
(6,  'B-',  '2025-11-11', 'available', '2026-05-11', 7,  6),
(7,  'O-',  '2025-11-13', 'available', '2026-05-13', 10, 7),
(8,  'O+',  '2025-11-15', 'available', '2026-05-15', 9,  8),
(9,  'A+',  '2025-11-17', 'available', '2026-05-17', 6,  9),
(10, 'B+',  '2025-11-19', 'available', '2026-05-19', 5,  10),
(11, 'AB-', '2025-11-21', 'available', '2026-05-21', 3,  11),
(12, 'O-',  '2025-11-23', 'available', '2026-05-23', 8,  12),
(13, 'A-',  '2025-11-25', 'available', '2026-05-25', 4,  13),
(14, 'B-',  '2025-11-27', 'available', '2026-05-27', 6,  14),
(15, 'O-',  '2025-11-29', 'available', '2026-05-29', 9,  15),
(16, 'AB+', '2025-12-01', 'available', '2026-06-01', 3,  NULL),
(17, 'O+',  '2025-12-03', 'available', '2026-06-03', 7,  NULL),
(18, 'A+',  '2025-12-05', 'available', '2026-06-05', 5,  NULL);

SELECT SETVAL('blood_id_seq', 18);

--   DONATION 

INSERT INTO Donation (Donor_ID, Manager_ID, Blood_ID, Donation_Date, Status) VALUES
(1,  1, 1,  '2025-10-01', 'completed'),
(2,  1, 2,  '2025-09-15', 'completed'),
(3,  2, 3,  '2025-08-20', 'completed'),
(4,  2, 4,  '2025-07-10', 'completed'),
(5,  3, 5,  '2025-10-05', 'completed'),
(6,  3, 6,  '2025-09-01', 'completed'),
(7,  4, 7,  '2025-08-11', 'completed'),
(8,  4, 8,  '2025-07-25', 'completed'),
(9,  1, 9,  '2025-10-10', 'completed'),
(10, 2, 10, '2025-09-20', 'completed'),
(11, 3, 11, '2025-08-01', 'completed'),
(12, 4, 12, '2025-07-15', 'completed'),
(13, 1, 13, '2025-09-05', 'completed'),
(14, 2, 14, '2025-10-15', 'completed'),
(15, 3, 15, '2025-08-30', 'completed');

SELECT SETVAL('donation_id_seq', 15);

--   OT  

INSERT INTO OT (ID, Is_Available) VALUES
(1, TRUE),
(2, TRUE),
(3, FALSE),
(4, TRUE),
(5, FALSE),
(6, TRUE);



--   APPOINTMENT  

INSERT INTO Appointment (ID, Receptionist_ID, Patient_ID, Doctor_ID, Date, Time, Status, OT_ID) VALUES
(1,  1, 1,  1,  '2026-01-05', '09:00:00', 'completed', NULL),
(2,  1, 2,  2,  '2026-01-05', '10:00:00', 'completed', NULL),
(3,  2, 3,  3,  '2026-01-06', '09:30:00', 'completed', 1),
(4,  2, 4,  4,  '2026-01-06', '11:00:00', 'completed', NULL),
(5,  3, 5,  5,  '2026-01-07', '08:00:00', 'completed', 2),
(6,  3, 6,  6,  '2026-01-07', '09:00:00', 'cancelled', NULL),
(7,  4, 7,  7,  '2026-01-08', '10:00:00', 'completed', NULL),
(8,  4, 8,  8,  '2026-01-08', '11:30:00', 'completed', NULL),
(9,  5, 9,  9,  '2026-01-09', '08:30:00', 'completed', NULL),
(10, 5, 10, 10, '2026-01-09', '09:30:00', 'completed', NULL),
(11, 1, 11, 11, '2026-01-10', '10:00:00', 'completed', NULL),
(12, 2, 12, 12, '2026-01-10', '11:00:00', 'cancelled', NULL),
(13, 3, 13, 1,  '2026-01-12', '09:00:00', 'completed', NULL),
(14, 4, 14, 2,  '2026-01-12', '10:30:00', 'completed', NULL),
(15, 5, 15, 3,  '2026-01-13', '08:00:00', 'completed', NULL),
(16, 6, 16, 4,  '2026-01-13', '09:00:00', 'completed', NULL),
(17, 1, 17, 5,  '2026-01-14', '10:00:00', 'cancelled', NULL),
(18, 2, 18, 6,  '2026-01-14', '11:00:00', 'completed', 4),
(19, 3, 19, 7,  '2026-01-15', '09:00:00', 'pending',   NULL),
(20, 4, 20, 8,  '2026-01-15', '10:00:00', 'pending',   NULL);


--  12. MEDICINE  

INSERT INTO Medicine (ID, Batch_No, Name, Stock, Expiry_Date) VALUES
(1,  'BATCH-001', 'Paracetamol',        500, '2028-06-01'),
(2,  'BATCH-002', 'Amoxicillin',        300, '2027-03-15'),
(3,  'BATCH-003', 'Metformin',          400, '2028-08-20'),
(4,  'BATCH-004', 'Atorvastatin',       250, '2028-07-10'),
(5,  'BATCH-005', 'Omeprazole',         350, '2027-05-25'),
(6,  'BATCH-006', 'Ciprofloxacin',      200, '2027-04-30'),
(7,  'BATCH-007', 'Aspirin',            600, '2029-01-01'),
(8,  'BATCH-008', 'Ibuprofen',          450, '2028-09-15'),
(9,  'BATCH-009', 'Lisinopril',         180, '2027-06-20'),
(10, 'BATCH-010', 'Azithromycin',       150, '2027-02-28'),
(11, 'BATCH-011', 'Insulin Glargine',   100, '2027-12-01'),
(12, 'BATCH-012', 'Salbutamol',         220, '2028-10-10'),
(13, 'BATCH-013', 'Diazepam',           80,  '2027-03-05'),
(14, 'BATCH-014', 'Ceftriaxone',        120, '2027-11-20'),
(15, 'BATCH-015', 'Dexamethasone',      160, '2028-01-15'),
(16, 'BATCH-016', 'Metronidazole',      200, '2027-08-10'),
(17, 'BATCH-017', 'Ranitidine',         180, '2028-04-20'),
(18, 'BATCH-018', 'Amlodipine',         140, '2027-09-30'),
(19, 'BATCH-019', 'Prednisolone',       110, '2028-02-28'),
(20, 'BATCH-020', 'Folic Acid',         300, '2029-05-01');


--  13. PRESCRIPTION  

INSERT INTO Prescription (ID, Doctor_ID, Patient_ID, Date) VALUES
(1,  1,  1,  '2026-01-05'),
(2,  2,  2,  '2026-01-05'),
(3,  3,  3,  '2026-01-06'),
(4,  4,  4,  '2026-01-06'),
(5,  5,  5,  '2026-01-07'),
(6,  6,  6,  '2026-01-07'),
(7,  7,  7,  '2026-01-08'),
(8,  8,  8,  '2026-01-08'),
(9,  9,  9,  '2026-01-09'),
(10, 10, 10, '2026-01-09'),
(11, 11, 11, '2026-01-10'),
(12, 12, 12, '2026-01-10'),
(13, 1,  13, '2026-01-12'),
(14, 2,  14, '2026-01-12'),
(15, 3,  15, '2026-01-13'),
(16, 4,  16, '2026-01-13'),
(17, 5,  17, '2026-01-14'),
(18, 6,  18, '2026-01-14');



--   PRESCRIPTION_MEDICINE 

INSERT INTO Prescription_Medicine (ID, Prescription_ID, Medicine_ID, Quantity) VALUES
(1,  1,  1,  2),
(2,  1,  7,  1),
(3,  2,  2,  1),
(4,  2,  5,  2),
(5,  3,  3,  3),
(6,  3,  4,  1),
(7,  4,  6,  2),
(8,  4,  8,  1),
(9,  5,  9,  1),
(10, 5,  13, 2),
(11, 6,  10, 2),
(12, 6,  12, 1),
(13, 7,  11, 1),
(14, 7,  15, 2),
(15, 8,  14, 2),
(16, 8,  16, 1),
(17, 9,  1,  3),
(18, 9,  7,  1),
(19, 10, 3,  2),
(20, 10, 17, 1),
(21, 11, 2,  1),
(22, 11, 18, 2),
(23, 12, 5,  3),
(24, 12, 19, 1),
(25, 13, 6,  2),
(26, 13, 20, 1),
(27, 14, 8,  1),
(28, 14, 9,  2),
(29, 15, 1,  2),
(30, 15, 10, 1);



--   DISPENSING  

INSERT INTO Dispensing (ID, Medicine_ID, Pharmacist_ID, Prescription_ID, Quantity) VALUES
(1,  1,  1, 1,  2),
(2,  7,  1, 1,  1),
(3,  2,  2, 2,  1),
(4,  5,  2, 2,  2),
(5,  3,  3, 3,  3),
(6,  4,  3, 3,  1),
(7,  6,  4, 4,  2),
(8,  8,  4, 4,  1),
(9,  9,  5, 5,  1),
(10, 13, 5, 5,  2),
(11, 10, 6, 6,  2),
(12, 12, 6, 6,  1),
(13, 11, 1, 7,  1),
(14, 15, 1, 7,  2),
(15, 14, 2, 8,  2),
(16, 16, 2, 8,  1),
(17, 1,  3, 9,  3),
(18, 7,  3, 9,  1);

-- Adjust stock to match dispensed quantities

UPDATE Medicine SET Stock = Stock - 5  WHERE ID = 1;
UPDATE Medicine SET Stock = Stock - 1  WHERE ID = 2;
UPDATE Medicine SET Stock = Stock - 3  WHERE ID = 3;
UPDATE Medicine SET Stock = Stock - 1  WHERE ID = 4;
UPDATE Medicine SET Stock = Stock - 2  WHERE ID = 5;
UPDATE Medicine SET Stock = Stock - 2  WHERE ID = 6;
UPDATE Medicine SET Stock = Stock - 2  WHERE ID = 7;
UPDATE Medicine SET Stock = Stock - 1  WHERE ID = 8;
UPDATE Medicine SET Stock = Stock - 1  WHERE ID = 9;
UPDATE Medicine SET Stock = Stock - 2  WHERE ID = 10;
UPDATE Medicine SET Stock = Stock - 1  WHERE ID = 11;
UPDATE Medicine SET Stock = Stock - 1  WHERE ID = 12;
UPDATE Medicine SET Stock = Stock - 2  WHERE ID = 13;
UPDATE Medicine SET Stock = Stock - 2  WHERE ID = 14;
UPDATE Medicine SET Stock = Stock - 2  WHERE ID = 15;
UPDATE Medicine SET Stock = Stock - 1  WHERE ID = 16;


--   ER_SHIFT  
INSERT INTO ER_Shift (ID, Receptionist_ID, Time, Date) VALUES
(1, 1, '08:00:00', '2026-04-01'),
(2, 2, '16:00:00', '2026-04-01'),
(3, 3, '00:00:00', '2026-04-02'),
(4, 1, '08:00:00', '2026-04-02'),
(5, 4, '16:00:00', '2026-04-03'),
(6, 5, '08:00:00', '2026-04-03'),
(7, 6, '16:00:00', '2026-04-04'),
(8, 2, '08:00:00', '2026-04-04');

SELECT SETVAL('er_shift_id_seq', 8);



--   ER_SHIFT_DOCTOR  
INSERT INTO ER_Shift_Doctor (Doctor_ID, ER_ID) VALUES
(9, 1), (1, 1), (2, 1),
(3, 2), (4, 2), (5, 2),
(9, 3), (6, 3), (7, 3),
(8, 4), (10, 4), (11, 4),
(12, 5), (1, 5), (3, 5),
(5, 6), (7, 6), (9, 6),
(2, 7), (4, 7), (6, 7), (8, 7);

SELECT SETVAL('er_shift_doctor_id_seq', (SELECT MAX(ID) FROM ER_Shift_Doctor));


--   ER_PATIENT_ENTRY  
INSERT INTO ER_Patient_Entry (ID, ER_Patient_ID, Doctor_ID, ER_ID, Status) VALUES
(1,  1,  9, 1, 'discharged'),
(2,  2,  1, 1, 'admitted'),
(3,  3,  2, 1, 'discharged'),
(4,  4,  3, 2, 'in_treatment'),
(5,  5,  4, 2, 'waiting'),
(6,  6,  9, 3, 'discharged'),
(7,  7,  6, 3, 'transferred'),
(8,  8,  8, 4, 'in_treatment'),
(9,  9,  10,4, 'waiting'),
(10, 10, 11,4, 'discharged'),
(11, 11, 12,5, 'admitted'),
(12, 12, 1, 5, 'in_treatment');



--  19. BLOOD_REQUEST 
INSERT INTO Blood_Request (ID, Doctor_ID, Patient_ID, Quantity_Needed, Status, Request_Date) VALUES
(1,  1,  1,  3, 'complete',  '2026-01-06'),  
(2,  2,  2,  2, 'pending',   '2026-01-07'),  
(3,  3,  3,  4, 'complete',  '2026-01-08'),  
(4,  4,  5,  2, 'pending',   '2026-01-08'),  
(5,  5,  7,  1, 'cancelled', '2026-01-09'),  
(6,  6,  10, 3, 'complete',  '2026-01-10'),  
(7,  7,  11, 2, 'pending',   '2026-01-11'),  
(8,  8,  12, 2, 'complete',  '2026-01-12'),  
(9,  9,  9,  1, 'pending',   '2026-01-13'),  
(10, 10, 16, 2, 'pending',   '2026-01-14'),  
(11, 11, 17, 1, 'complete',  '2026-01-15'),  
(12, 12, 18, 3, 'pending',   '2026-01-16'),  
(13, 1,  15, 2, 'complete',  '2026-01-17'),  
(14, 2,  19, 1, 'pending',   '2026-01-18'),  
(15, 3,  20, 2, 'pending',   '2026-01-19');  


--   BLOOD_REQUEST_FULFILLMENT 
INSERT INTO Blood_Request_Fulfillment
    (Request_ID, Blood_ID, Manager_ID, Quantity_Provided, Fulfillment_Date)
VALUES
(1,  1, 1, 3, '2026-01-07'),
(3,  3, 2, 3, '2026-01-09'),
(3,  7, 2, 1, '2026-01-09'),
(6,  8, 3, 3, '2026-01-11'),
(8,  9, 1, 2, '2026-01-13'),
(11, 7, 4, 1, '2026-01-16'),
(13, 5, 2, 2, '2026-01-18');

SELECT SETVAL('blood_request_fulfillment_id_seq',
    (SELECT MAX(ID) FROM Blood_Request_Fulfillment));

UPDATE Blood SET Unit = Unit - 3 WHERE ID = 1;
UPDATE Blood SET Unit = Unit - 3 WHERE ID = 3;
UPDATE Blood SET Unit = Unit - 2 WHERE ID = 5;
UPDATE Blood SET Unit = Unit - 2 WHERE ID = 7;  
UPDATE Blood SET Unit = Unit - 3 WHERE ID = 8;
UPDATE Blood SET Unit = Unit - 2 WHERE ID = 9;


--  VERIFICATION
SELECT table_name, records FROM (
    SELECT 'Doctor'                    AS table_name, COUNT(*) AS records FROM Doctor
    UNION ALL SELECT 'Receptionist',               COUNT(*) FROM Receptionist
    UNION ALL SELECT 'Pharmacist',                 COUNT(*) FROM Pharmacist
    UNION ALL SELECT 'Blood_Manager',              COUNT(*) FROM Blood_Manager
    UNION ALL SELECT 'Patient',                    COUNT(*) FROM Patient
    UNION ALL SELECT 'ER_Patient',                 COUNT(*) FROM ER_Patient
    UNION ALL SELECT 'Donor',                      COUNT(*) FROM Donor
    UNION ALL SELECT 'Blood',                      COUNT(*) FROM Blood
    UNION ALL SELECT 'Donation',                   COUNT(*) FROM Donation
    UNION ALL SELECT 'OT',                         COUNT(*) FROM OT
    UNION ALL SELECT 'Appointment',                COUNT(*) FROM Appointment
    UNION ALL SELECT 'Medicine',                   COUNT(*) FROM Medicine
    UNION ALL SELECT 'Prescription',               COUNT(*) FROM Prescription
    UNION ALL SELECT 'Prescription_Medicine',      COUNT(*) FROM Prescription_Medicine
    UNION ALL SELECT 'Dispensing',                 COUNT(*) FROM Dispensing
    UNION ALL SELECT 'ER_Shift',                   COUNT(*) FROM ER_Shift
    UNION ALL SELECT 'ER_Shift_Doctor',            COUNT(*) FROM ER_Shift_Doctor
    UNION ALL SELECT 'ER_Patient_Entry',           COUNT(*) FROM ER_Patient_Entry
    UNION ALL SELECT 'Blood_Request',              COUNT(*) FROM Blood_Request
    UNION ALL SELECT 'Blood_Request_Fulfillment',  COUNT(*) FROM Blood_Request_Fulfillment
) t
ORDER BY table_name;
