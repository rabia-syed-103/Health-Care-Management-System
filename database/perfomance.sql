
-- INDEXES
CREATE INDEX IF NOT EXISTS idx_er_doctor_availability
    ON ER_Shift_Doctor(Doctor_ID, ER_ID);

CREATE INDEX IF NOT EXISTS idx_blood_type_status
    ON Blood(B_Gr, Status);

CREATE INDEX IF NOT EXISTS idx_blood_request_status
    ON Blood_Request(Status);

CREATE INDEX IF NOT EXISTS idx_medicine_name
    ON Medicine(Name);

--  QUERY 1
-- Before
EXPLAIN ANALYZE
SELECT
    esd.Doctor_ID,
    d.Name              AS Doctor_Name,
    d.Specialization,
    ers.Date            AS Shift_Date,
    ers.Time            AS Shift_Time
FROM       ER_Shift_Doctor  esd
JOIN       Doctor           d   ON d.ID   = esd.Doctor_ID
JOIN       ER_Shift         ers ON ers.ID = esd.ER_ID
WHERE esd.ER_ID = 3
ORDER BY d.Name;

-- After
EXPLAIN ANALYZE
SELECT
    esd.Doctor_ID,
    d.Name              AS Doctor_Name,
    d.Specialization,
    ers.Date            AS Shift_Date,
    ers.Time            AS Shift_Time
FROM       ER_Shift_Doctor  esd
JOIN       Doctor           d   ON d.ID   = esd.Doctor_ID
JOIN       ER_Shift         ers ON ers.ID = esd.ER_ID
WHERE esd.ER_ID = 3
ORDER BY d.Name;


--  QUERY 2
-- Before
EXPLAIN ANALYZE
SELECT
    b.ID,
    b.B_Gr,
    b.Unit                                  AS Units_Available,
    b.Collected_Date,
    b.Expiry_Date,
    (b.Expiry_Date - CURRENT_DATE)          AS Days_Until_Expiry
FROM  Blood b
WHERE b.B_Gr        = 'O+'
  AND b.Status      = 'available'
  AND b.Expiry_Date > CURRENT_DATE
  AND b.Unit        > 0
ORDER BY b.Expiry_Date ASC;

-- After
EXPLAIN ANALYZE
SELECT
    b.ID,
    b.B_Gr,
    b.Unit                                  AS Units_Available,
    b.Collected_Date,
    b.Expiry_Date,
    (b.Expiry_Date - CURRENT_DATE)          AS Days_Until_Expiry
FROM  Blood b
WHERE b.B_Gr        = 'O+'
  AND b.Status      = 'available'
  AND b.Expiry_Date > CURRENT_DATE
  AND b.Unit        > 0
ORDER BY b.Expiry_Date ASC;

--  QUERY 3
-- Before
EXPLAIN ANALYZE
SELECT
    br.ID                   AS Request_ID,
    br.Request_Date,
    br.Quantity_Needed,
    br.Status,
    p.Name                  AS Patient_Name,
    p.B_Gr                  AS Patient_Blood_Group,
    p.P_No                  AS Patient_Phone,
    d.Name                  AS Requesting_Doctor,
    d.Specialization
FROM       Blood_Request br
JOIN       Patient       p  ON p.ID = br.Patient_ID
JOIN       Doctor        d  ON d.ID = br.Doctor_ID
WHERE br.Status = 'pending'
ORDER BY br.Request_Date ASC;

-- After
EXPLAIN ANALYZE
SELECT
    br.ID                   AS Request_ID,
    br.Request_Date,
    br.Quantity_Needed,
    br.Status,
    p.Name                  AS Patient_Name,
    p.B_Gr                  AS Patient_Blood_Group,
    p.P_No                  AS Patient_Phone,
    d.Name                  AS Requesting_Doctor,
    d.Specialization
FROM       Blood_Request br
JOIN       Patient       p  ON p.ID = br.Patient_ID
JOIN       Doctor        d  ON d.ID = br.Doctor_ID
WHERE br.Status = 'pending'
ORDER BY br.Request_Date ASC;

--  QUERY 4
-- Before
EXPLAIN ANALYZE
SELECT
    m.ID,
    m.Name              AS Medicine_Name,
    m.Batch_No,
    m.Stock,
    m.Expiry_Date,
    (m.Expiry_Date - CURRENT_DATE)  AS Days_Until_Expiry,
    CASE
        WHEN m.Stock = 0                        THEN 'OUT OF STOCK'
        WHEN m.Expiry_Date <= CURRENT_DATE      THEN 'EXPIRED'
        WHEN m.Expiry_Date <= CURRENT_DATE + 30 THEN 'EXPIRING SOON'
        ELSE                                         'OK'
    END                 AS Stock_Status
FROM  Medicine m
WHERE m.Name = 'Paracetamol'
ORDER BY m.Expiry_Date ASC;

-- After
EXPLAIN ANALYZE
SELECT
    m.ID,
    m.Name              AS Medicine_Name,
    m.Batch_No,
    m.Stock,
    m.Expiry_Date,
    (m.Expiry_Date - CURRENT_DATE)  AS Days_Until_Expiry,
    CASE
        WHEN m.Stock = 0                        THEN 'OUT OF STOCK'
        WHEN m.Expiry_Date <= CURRENT_DATE      THEN 'EXPIRED'
        WHEN m.Expiry_Date <= CURRENT_DATE + 30 THEN 'EXPIRING SOON'
        ELSE                                         'OK'
    END                 AS Stock_Status
FROM  Medicine m
WHERE m.Name = 'Paracetamol'
ORDER BY m.Expiry_Date ASC;