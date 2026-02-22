DROP TABLE IF EXISTS Blood_Request_Fulfillment   CASCADE;
DROP TABLE IF EXISTS Blood_Request               CASCADE;
DROP TABLE IF EXISTS Donation                    CASCADE;
DROP TABLE IF EXISTS Blood                       CASCADE;
DROP TABLE IF EXISTS ER_Patient_Entry            CASCADE;
DROP TABLE IF EXISTS ER_Shift_Doctor             CASCADE;
DROP TABLE IF EXISTS ER_Shift                    CASCADE;
DROP TABLE IF EXISTS ER_Patient                  CASCADE;
DROP TABLE IF EXISTS Dispensing                  CASCADE;
DROP TABLE IF EXISTS Prescription_Medicine       CASCADE;
DROP TABLE IF EXISTS Prescription                CASCADE;
DROP TABLE IF EXISTS Appointment                 CASCADE;
DROP TABLE IF EXISTS OT                          CASCADE;
DROP TABLE IF EXISTS Medicine                    CASCADE;
DROP TABLE IF EXISTS Pharmacist                  CASCADE;
DROP TABLE IF EXISTS Donor                       CASCADE;
DROP TABLE IF EXISTS Blood_Manager               CASCADE;
DROP TABLE IF EXISTS Patient                     CASCADE;
DROP TABLE IF EXISTS Doctor                      CASCADE;
DROP TABLE IF EXISTS Receptionist                CASCADE;

DROP FUNCTION IF EXISTS fn_donor_management_before();
DROP FUNCTION IF EXISTS fn_donor_management_after();
DROP FUNCTION IF EXISTS fn_blood_before();
DROP FUNCTION IF EXISTS fn_blood_fulfillment_before();
DROP FUNCTION IF EXISTS fn_blood_fulfillment_after();
DROP FUNCTION IF EXISTS fn_prescription_before();
DROP FUNCTION IF EXISTS fn_pm_before();
DROP FUNCTION IF EXISTS fn_dispensing_before();
DROP FUNCTION IF EXISTS fn_dispensing_after_insert();
DROP FUNCTION IF EXISTS fn_dispensing_after_delete();
DROP FUNCTION IF EXISTS fn_appointment_management();
DROP FUNCTION IF EXISTS fn_er_shift_before();


DROP VIEW IF EXISTS Active_Blood_Inventory;
DROP VIEW IF EXISTS Donor_Donation_History;
DROP VIEW IF EXISTS ER_Shift_Staffing;
DROP VIEW IF EXISTS Medicine_Expiry_Alert;

DROP INDEX IF EXISTS idx_prescription_patient_id;
DROP INDEX IF EXISTS idx_blood_type_status;
DROP INDEX IF EXISTS idx_blood_request_status;
DROP INDEX IF EXISTS idx_er_entry_patient;


--  TABLES

CREATE TABLE Doctor (
    ID              INT             NOT NULL,
    Name            VARCHAR(100)    NOT NULL,
    Email           VARCHAR(150)    NOT NULL,
    P_No            VARCHAR(20)     NOT NULL,
    Specialization  VARCHAR(100)    NOT NULL,
    CONSTRAINT PK_Doctor        PRIMARY KEY (ID),
    CONSTRAINT UQ_Doctor_Email  UNIQUE (Email)
);

CREATE TABLE Receptionist (
    ID      INT             NOT NULL,
    Name    VARCHAR(100)    NOT NULL,
    Email   VARCHAR(100)    NOT NULL,
    P_No    VARCHAR(20)     NOT NULL,
    CONSTRAINT PK_Receptionist          PRIMARY KEY (ID),
    CONSTRAINT UQ_Receptionist_Email    UNIQUE (Email)
);

CREATE TABLE Pharmacist (
    ID      INT             NOT NULL,
    Name    VARCHAR(100)    NOT NULL,
    Email   VARCHAR(150)    NOT NULL,
    P_No    VARCHAR(20)     NOT NULL,
    CONSTRAINT PK_Pharmacist        PRIMARY KEY (ID),
    CONSTRAINT UQ_Pharmacist_Email  UNIQUE (Email)
);

CREATE TABLE Blood_Manager (
    ID      INT             NOT NULL,
    Name    VARCHAR(100)    NOT NULL,
    Email   VARCHAR(150)    NOT NULL,
    P_No    VARCHAR(20)     NOT NULL,
    CONSTRAINT PK_Blood_Manager         PRIMARY KEY (ID),
    CONSTRAINT UQ_Blood_Manager_Email   UNIQUE (Email)
);

CREATE TABLE Patient (
    ID      INT             NOT NULL,
    Name    VARCHAR(100)    NOT NULL,
    Email   VARCHAR(150)    NOT NULL,
    B_Gr    CHAR(3)         NOT NULL,
    P_No    VARCHAR(20)     NOT NULL,
    CONSTRAINT PK_Patient       PRIMARY KEY (ID),
    CONSTRAINT UQ_Patient_Email UNIQUE (Email),
    CONSTRAINT CK_Patient_BGr   CHECK (B_Gr IN ('A+','A-','B+','B-','AB+','AB-','O+','O-'))
);

CREATE TABLE ER_Patient (
    ID              INT             NOT NULL,
    Name            VARCHAR(100)    NOT NULL,
    Age             INT             NOT NULL,
    P_No            VARCHAR(20)     NOT NULL,
    Arrival_Time    TIME            NOT NULL,
    CONSTRAINT PK_ER_Patient    PRIMARY KEY (ID),
    CONSTRAINT CK_ER_Age        CHECK (Age >= 0 AND Age <= 150)
);

CREATE TABLE Donor (
    ID          INT             NOT NULL,
    Name        VARCHAR(100)    NOT NULL,
    P_No        VARCHAR(20)     NOT NULL,
    Email       VARCHAR(100)    NOT NULL,
    B_Gr        VARCHAR(5)      NOT NULL,
    Last_Donate DATE            NOT NULL,
    CONSTRAINT PK_Donor         PRIMARY KEY (ID),
    CONSTRAINT UQ_Donor_Email   UNIQUE (Email),
    CONSTRAINT CK_Donor_BGr     CHECK (B_Gr IN ('A+','A-','B+','B-','AB+','AB-','O+','O-'))
);

CREATE TABLE Blood (
    ID              SERIAL          NOT NULL,
    B_Gr            CHAR(3)         NOT NULL,
    Collected_Date  DATE            NOT NULL,
    Status          VARCHAR(20)     NOT NULL    DEFAULT 'available',
    Expiry_Date     DATE            NOT NULL,
    Unit            INT             NOT NULL    DEFAULT 0,
    Donor_ID        INT             NULL,
    CONSTRAINT PK_Blood             PRIMARY KEY (ID),
    CONSTRAINT CK_Blood_BGr         CHECK (B_Gr IN ('A+','A-','B+','B-','AB+','AB-','O+','O-')),
    CONSTRAINT CK_Blood_Status      CHECK (Status IN ('available','reserved','used','expired','discarded')),
    CONSTRAINT CK_Blood_Expiry      CHECK (Expiry_Date > Collected_Date),
    CONSTRAINT CK_Blood_Unit        CHECK (Unit >= 0),
    CONSTRAINT FK_Blood_Donor       FOREIGN KEY (Donor_ID)
                                        REFERENCES Donor(ID)
                                        ON DELETE SET NULL
                                        ON UPDATE CASCADE
);

CREATE TABLE Donation (
    ID              SERIAL          NOT NULL,
    Donor_ID        INT             NOT NULL,
    Manager_ID      INT             NOT NULL,
    Blood_ID        INT             NOT NULL,
    Donation_Date   DATE            NOT NULL,
    Status          VARCHAR(20)     NOT NULL    DEFAULT 'completed',
    CONSTRAINT PK_Donation          PRIMARY KEY (ID),
    CONSTRAINT CK_Donation_Status   CHECK (Status IN ('completed','cancelled','deferred')),
    CONSTRAINT FK_Donation_Donor    FOREIGN KEY (Donor_ID)
                                        REFERENCES Donor(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Donation_Manager  FOREIGN KEY (Manager_ID)
                                        REFERENCES Blood_Manager(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Donation_Blood    FOREIGN KEY (Blood_ID)
                                        REFERENCES Blood(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE
);

CREATE TABLE OT (
    ID              INT     NOT NULL,
    Is_Available    BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT PK_OT PRIMARY KEY (ID)
);

CREATE TABLE Appointment (
    ID                  INT             NOT NULL,
    Receptionist_ID     INT             NOT NULL,
    Patient_ID          INT             NOT NULL,
    Doctor_ID           INT             NOT NULL,
    Date                DATE            NOT NULL,
    Time                TIME            NOT NULL,
    Status              VARCHAR(50)     NOT NULL    DEFAULT 'pending',
    OT_ID               INT             NULL,
    CONSTRAINT PK_Appointment       PRIMARY KEY (ID),
    CONSTRAINT CK_Appt_Status       CHECK (Status IN ('pending','cancelled','completed')),
    CONSTRAINT FK_Appt_Receptionist FOREIGN KEY (Receptionist_ID)
                                        REFERENCES Receptionist(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Appt_Patient      FOREIGN KEY (Patient_ID)
                                        REFERENCES Patient(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Appt_Doctor       FOREIGN KEY (Doctor_ID)
                                        REFERENCES Doctor(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Appt_OT           FOREIGN KEY (OT_ID)
                                        REFERENCES OT(ID)
                                        ON DELETE SET NULL
                                        ON UPDATE CASCADE
);

CREATE TABLE Medicine (
    ID              INT             NOT NULL,
    Batch_No        VARCHAR(50)     NOT NULL,
    Name            VARCHAR(150)    NOT NULL,
    Stock           INT             NOT NULL,
    Expiry_Date     DATE            NOT NULL,
    CONSTRAINT PK_Medicine          PRIMARY KEY (ID),
    CONSTRAINT UQ_Medicine_Batch    UNIQUE (Batch_No),
    CONSTRAINT CK_Medicine_Stock    CHECK (Stock >= 0)
);

CREATE TABLE Prescription (
    ID          INT     NOT NULL,
    Doctor_ID   INT     NOT NULL,
    Patient_ID  INT     NOT NULL,
    Date        DATE    NOT NULL,
    CONSTRAINT PK_Prescription      PRIMARY KEY (ID),
    CONSTRAINT FK_Presc_Doctor      FOREIGN KEY (Doctor_ID)
                                        REFERENCES Doctor(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Presc_Patient     FOREIGN KEY (Patient_ID)
                                        REFERENCES Patient(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE
);

CREATE TABLE Prescription_Medicine (
    ID                  INT     NOT NULL,
    Prescription_ID     INT     NOT NULL,
    Medicine_ID         INT     NOT NULL,
    Quantity            INT     NOT NULL,
    CONSTRAINT PK_Prescription_Medicine PRIMARY KEY (ID),
    CONSTRAINT CK_PM_Quantity           CHECK (Quantity > 0),
    CONSTRAINT FK_PM_Prescription       FOREIGN KEY (Prescription_ID)
                                            REFERENCES Prescription(ID)
                                            ON DELETE CASCADE
                                            ON UPDATE CASCADE,
    CONSTRAINT FK_PM_Medicine           FOREIGN KEY (Medicine_ID)
                                            REFERENCES Medicine(ID)
                                            ON DELETE RESTRICT
                                            ON UPDATE CASCADE
);

CREATE TABLE Dispensing (
    ID                  INT     NOT NULL,
    Medicine_ID         INT     NOT NULL,
    Pharmacist_ID       INT     NOT NULL,
    Prescription_ID     INT     NOT NULL,
    Quantity            INT     NOT NULL,
    CONSTRAINT PK_Dispensing        PRIMARY KEY (ID),
    CONSTRAINT CK_Disp_Quantity     CHECK (Quantity > 0),
    CONSTRAINT FK_Disp_Medicine     FOREIGN KEY (Medicine_ID)
                                        REFERENCES Medicine(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Disp_Pharmacist   FOREIGN KEY (Pharmacist_ID)
                                        REFERENCES Pharmacist(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_Disp_Presc        FOREIGN KEY (Prescription_ID)
                                        REFERENCES Prescription(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE
);

CREATE TABLE ER_Shift (
    ID                  SERIAL          NOT NULL,
    Receptionist_ID     INT             NULL,
    Time                TIME            NOT NULL,
    Date                DATE            NOT NULL,
    CONSTRAINT PK_ER_Shift          PRIMARY KEY (ID),
    CONSTRAINT UQ_ER_Shift_DateTime UNIQUE (Date, Time),
    CONSTRAINT FK_ER_Shift_Recept   FOREIGN KEY (Receptionist_ID)
                                        REFERENCES Receptionist(ID)
                                        ON DELETE SET NULL
                                        ON UPDATE CASCADE
);

CREATE TABLE ER_Shift_Doctor (
    ID          SERIAL  NOT NULL,
    Doctor_ID   INT     NOT NULL,
    ER_ID       INT     NOT NULL,
    CONSTRAINT PK_ER_Shift_Doctor   PRIMARY KEY (ID),
    CONSTRAINT UQ_ER_Shift_Doctor   UNIQUE (Doctor_ID, ER_ID),
    CONSTRAINT FK_ESD_Doctor        FOREIGN KEY (Doctor_ID)
                                        REFERENCES Doctor(ID)
                                        ON DELETE CASCADE
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_ESD_Shift         FOREIGN KEY (ER_ID)
                                        REFERENCES ER_Shift(ID)
                                        ON DELETE CASCADE
                                        ON UPDATE CASCADE
);

CREATE TABLE ER_Patient_Entry (
    ID              INT             NOT NULL,
    ER_Patient_ID   INT             NOT NULL,
    Doctor_ID       INT             NOT NULL,
    ER_ID           INT             NOT NULL,
    Status          VARCHAR(50)     NOT NULL    DEFAULT 'waiting',
    CONSTRAINT PK_ER_Patient_Entry  PRIMARY KEY (ID),
    CONSTRAINT CK_ERPE_Status       CHECK (Status IN ('waiting','in_treatment','admitted','discharged','transferred')),
    CONSTRAINT FK_ERPE_Patient      FOREIGN KEY (ER_Patient_ID)
                                        REFERENCES ER_Patient(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_ERPE_Doctor       FOREIGN KEY (Doctor_ID)
                                        REFERENCES Doctor(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_ERPE_Shift        FOREIGN KEY (ER_ID)
                                        REFERENCES ER_Shift(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE
);

CREATE TABLE Blood_Request (
    ID                  INT             NOT NULL,
    Doctor_ID           INT             NOT NULL,
    Patient_ID          INT             NOT NULL,
    Quantity_Needed     INT             NOT NULL,
    Status              VARCHAR(50)     NOT NULL    DEFAULT 'pending',
    Request_Date        DATE            NOT NULL,
    CONSTRAINT PK_Blood_Request     PRIMARY KEY (ID),
    CONSTRAINT CK_BR_Status         CHECK (Status IN ('pending','complete','cancelled')),
    CONSTRAINT CK_BR_Quantity       CHECK (Quantity_Needed > 0),
    CONSTRAINT FK_BR_Doctor         FOREIGN KEY (Doctor_ID)
                                        REFERENCES Doctor(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE,
    CONSTRAINT FK_BR_Patient        FOREIGN KEY (Patient_ID)
                                        REFERENCES Patient(ID)
                                        ON DELETE RESTRICT
                                        ON UPDATE CASCADE
);

CREATE TABLE Blood_Request_Fulfillment (
    ID                  SERIAL  NOT NULL,
    Request_ID          INT     NOT NULL,
    Blood_ID            INT     NOT NULL,
    Manager_ID          INT     NOT NULL,
    Quantity_Provided   INT     NOT NULL,
    Fulfillment_Date    DATE    NOT NULL,
    CONSTRAINT PK_BRF           PRIMARY KEY (ID),
    CONSTRAINT CK_BRF_Quantity  CHECK (Quantity_Provided > 0),
    CONSTRAINT FK_BRF_Request   FOREIGN KEY (Request_ID)
                                    REFERENCES Blood_Request(ID)
                                    ON DELETE CASCADE
                                    ON UPDATE CASCADE,
    CONSTRAINT FK_BRF_Blood     FOREIGN KEY (Blood_ID)
                                    REFERENCES Blood(ID)
                                    ON DELETE RESTRICT
                                    ON UPDATE CASCADE,
    CONSTRAINT FK_BRF_Manager   FOREIGN KEY (Manager_ID)
                                    REFERENCES Blood_Manager(ID)
                                    ON DELETE RESTRICT
                                    ON UPDATE CASCADE
);

--TRIGGERS

--TRIGGER 1 — trg_donor_management
--Validates 90-day donation gap and prevents future dates, 
--then updates donor's last donation date and increments blood inventory.

CREATE OR REPLACE FUNCTION fn_donor_management_before()
RETURNS TRIGGER AS $$
DECLARE
    v_last_donate   DATE;
    v_donor_name    VARCHAR(100);
    v_days_since    INT;
BEGIN
    IF NEW.Donation_Date > CURRENT_DATE THEN
        RAISE EXCEPTION
            '[Donor Management] Donation date "%" is in the future. Today is %.',
            NEW.Donation_Date, CURRENT_DATE;
    END IF;

    SELECT Last_Donate, Name
    INTO   v_last_donate, v_donor_name
    FROM   Donor
    WHERE  ID = NEW.Donor_ID;

    v_days_since := (CURRENT_DATE - v_last_donate);

    IF v_days_since < 90 THEN
        RAISE EXCEPTION
            '[Donor Management] Donor "%" (ID: %) donated % days ago on %. '
            'Must wait until % (90-day rule).',
            v_donor_name, NEW.Donor_ID,
            v_days_since, v_last_donate,
            v_last_donate + INTERVAL '90 days';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_donor_management_before
BEFORE INSERT ON Donation
FOR EACH ROW
EXECUTE FUNCTION fn_donor_management_before();


CREATE OR REPLACE FUNCTION fn_donor_management_after()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.Status = 'completed' THEN
        UPDATE Donor
        SET    Last_Donate = NEW.Donation_Date
        WHERE  ID          = NEW.Donor_ID;

        UPDATE Blood
        SET    Unit   = Unit + 1,
               Status = CASE WHEN Status = 'used' THEN 'available' ELSE Status END
        WHERE  ID = NEW.Blood_ID;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_donor_management_after
AFTER INSERT ON Donation
FOR EACH ROW
EXECUTE FUNCTION fn_donor_management_after();

--TRIGGER 2 — trg_blood_inventory_management
--Prevents negative blood units, auto-expires outdated blood, 
--validates fulfillment requests, deducts allocated units, and 
--updates request status

CREATE OR REPLACE FUNCTION fn_blood_before()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.Unit < 0 THEN
        RAISE EXCEPTION
            '[Blood Inventory] Blood ID % cannot have negative units. Value: %.',
            NEW.ID, NEW.Unit;
    END IF;

    IF NEW.Expiry_Date <= CURRENT_DATE
       AND NEW.Status NOT IN ('used', 'discarded', 'expired') THEN
        NEW.Status := 'expired';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_blood_before
BEFORE INSERT OR UPDATE ON Blood
FOR EACH ROW
EXECUTE FUNCTION fn_blood_before();


CREATE OR REPLACE FUNCTION fn_blood_fulfillment_before()
RETURNS TRIGGER AS $$
DECLARE
    v_status        VARCHAR(20);
    v_expiry        DATE;
    v_units         INT;
    v_bgr           CHAR(3);
    v_patient_bgr   CHAR(3);
    v_patient_name  VARCHAR(100);
    v_compatible    BOOLEAN := FALSE;
BEGIN

    SELECT Status, Expiry_Date, Unit, B_Gr
    INTO   v_status, v_expiry, v_units, v_bgr
    FROM   Blood
    WHERE  ID = NEW.Blood_ID;

    
    IF v_expiry <= CURRENT_DATE OR v_status = 'expired' THEN
        RAISE EXCEPTION
            '[Blood Inventory] Blood ID % (%) expired on %. Cannot use for fulfillment.',
            NEW.Blood_ID, v_bgr, v_expiry;
    END IF;


    IF v_status NOT IN ('available', 'reserved') THEN
        RAISE EXCEPTION
            '[Blood Inventory] Blood ID % has status "%" and cannot be used.',
            NEW.Blood_ID, v_status;
    END IF;

    IF NEW.Quantity_Provided < 1 OR NEW.Quantity_Provided > 20 THEN
        RAISE EXCEPTION
            '[Blood Inventory] Fulfillment quantity must be between 1 and 20. Requested: %.',
            NEW.Quantity_Provided;
    END IF;

    IF NEW.Quantity_Provided > v_units THEN
        RAISE EXCEPTION
            '[Blood Inventory] Insufficient units for Blood ID %. Available: %, Requested: %.',
            NEW.Blood_ID, v_units, NEW.Quantity_Provided;
    END IF;

    SELECT p.B_Gr, p.Name
    INTO   v_patient_bgr, v_patient_name
    FROM   Blood_Request br
    JOIN   Patient       p ON p.ID = br.Patient_ID
    WHERE  br.ID = NEW.Request_ID;

   
    v_compatible := CASE v_patient_bgr
        WHEN 'O-'  THEN v_bgr IN ('O-')
        WHEN 'O+'  THEN v_bgr IN ('O-', 'O+')
        WHEN 'A-'  THEN v_bgr IN ('O-', 'A-')
        WHEN 'A+'  THEN v_bgr IN ('O-', 'O+', 'A-', 'A+')
        WHEN 'B-'  THEN v_bgr IN ('O-', 'B-')
        WHEN 'B+'  THEN v_bgr IN ('O-', 'O+', 'B-', 'B+')
        WHEN 'AB-' THEN v_bgr IN ('O-', 'A-', 'B-', 'AB-')
        WHEN 'AB+' THEN v_bgr IN ('O-', 'O+', 'A-', 'A+', 'B-', 'B+', 'AB-', 'AB+')
        ELSE FALSE
    END;

    IF NOT v_compatible THEN
        RAISE EXCEPTION
            '[Blood Inventory] Incompatible blood type for Request ID %. '
            'Patient "%" (%) cannot receive % blood. '
            'Compatible donor types for %: %.',
            NEW.Request_ID,
            v_patient_name, v_patient_bgr,
            v_bgr,
            v_patient_bgr,
            CASE v_patient_bgr
                WHEN 'O-'  THEN 'O-'
                WHEN 'O+'  THEN 'O-, O+'
                WHEN 'A-'  THEN 'O-, A-'
                WHEN 'A+'  THEN 'O-, O+, A-, A+'
                WHEN 'B-'  THEN 'O-, B-'
                WHEN 'B+'  THEN 'O-, O+, B-, B+'
                WHEN 'AB-' THEN 'O-, A-, B-, AB-'
                WHEN 'AB+' THEN 'O-, O+, A-, A+, B-, B+, AB-, AB+'
            END;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_blood_fulfillment_before
BEFORE INSERT ON Blood_Request_Fulfillment
FOR EACH ROW
EXECUTE FUNCTION fn_blood_fulfillment_before();


CREATE OR REPLACE FUNCTION fn_blood_fulfillment_after()
RETURNS TRIGGER AS $$
DECLARE
    v_units_after   INT;
    v_qty_needed    INT;
    v_qty_fulfilled INT;
BEGIN
    UPDATE Blood
    SET    Unit = Unit - NEW.Quantity_Provided
    WHERE  ID   = NEW.Blood_ID
    RETURNING Unit INTO v_units_after;

    UPDATE Blood
    SET    Status = CASE WHEN v_units_after = 0 THEN 'used' ELSE 'reserved' END
    WHERE  ID = NEW.Blood_ID;

    SELECT br.Quantity_Needed,
           COALESCE(SUM(brf.Quantity_Provided), 0)
    INTO   v_qty_needed, v_qty_fulfilled
    FROM   Blood_Request br
    LEFT JOIN Blood_Request_Fulfillment brf ON brf.Request_ID = br.ID
    WHERE  br.ID = NEW.Request_ID
    GROUP BY br.Quantity_Needed;

    IF v_qty_fulfilled >= v_qty_needed THEN
        UPDATE Blood_Request SET Status = 'complete' WHERE ID = NEW.Request_ID;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_blood_fulfillment_after
AFTER INSERT ON Blood_Request_Fulfillment
FOR EACH ROW
EXECUTE FUNCTION fn_blood_fulfillment_after();

--TRIGGER 3 — trg_medicine_dispensing
--Prevents duplicate prescriptions and medicines, validates medicine expiry 
--and stock before dispensing, automatically deducts stock on dispensing, and 
--restores stock on dispensing cancellation

CREATE OR REPLACE FUNCTION fn_prescription_before()
RETURNS TRIGGER AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM   Prescription
    WHERE  Doctor_ID  = NEW.Doctor_ID
      AND  Patient_ID = NEW.Patient_ID
      AND  Date       = NEW.Date;

    IF v_count > 0 THEN
        RAISE EXCEPTION
            '[Medicine Dispensing] A prescription already exists for '
            'Patient ID % from Doctor ID % on %.',
            NEW.Patient_ID, NEW.Doctor_ID, NEW.Date;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prescription_before
BEFORE INSERT ON Prescription
FOR EACH ROW
EXECUTE FUNCTION fn_prescription_before();


CREATE OR REPLACE FUNCTION fn_pm_before()
RETURNS TRIGGER AS $$
DECLARE
    v_count         INT;
    v_medicine_name VARCHAR(150);
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM   Prescription_Medicine
    WHERE  Prescription_ID = NEW.Prescription_ID
      AND  Medicine_ID     = NEW.Medicine_ID;

    IF v_count > 0 THEN
        SELECT Name INTO v_medicine_name FROM Medicine WHERE ID = NEW.Medicine_ID;
        RAISE EXCEPTION
            '[Medicine Dispensing] "%" (ID: %) is already in Prescription ID %. '
            'Update the quantity instead.',
            v_medicine_name, NEW.Medicine_ID, NEW.Prescription_ID;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_pm_before
BEFORE INSERT ON Prescription_Medicine
FOR EACH ROW
EXECUTE FUNCTION fn_pm_before();


CREATE OR REPLACE FUNCTION fn_dispensing_before()
RETURNS TRIGGER AS $$
DECLARE
    v_expiry        DATE;
    v_stock         INT;
    v_medicine_name VARCHAR(150);
BEGIN
    SELECT Expiry_Date, Stock, Name
    INTO   v_expiry, v_stock, v_medicine_name
    FROM   Medicine
    WHERE  ID = NEW.Medicine_ID;

    IF v_expiry <= CURRENT_DATE THEN
        RAISE EXCEPTION
            '[Medicine Dispensing] Cannot dispense "%" (ID: %). Expired on %.',
            v_medicine_name, NEW.Medicine_ID, v_expiry;
    END IF;

    IF v_stock < NEW.Quantity THEN
        RAISE EXCEPTION
            '[Medicine Dispensing] Insufficient stock for "%" (ID: %). '
            'Available: %, Requested: %.',
            v_medicine_name, NEW.Medicine_ID, v_stock, NEW.Quantity;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_dispensing_before
BEFORE INSERT ON Dispensing
FOR EACH ROW
EXECUTE FUNCTION fn_dispensing_before();


CREATE OR REPLACE FUNCTION fn_dispensing_after_insert()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE Medicine SET Stock = Stock - NEW.Quantity WHERE ID = NEW.Medicine_ID;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_dispensing_after_insert
AFTER INSERT ON Dispensing
FOR EACH ROW
EXECUTE FUNCTION fn_dispensing_after_insert();


CREATE OR REPLACE FUNCTION fn_dispensing_after_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE Medicine SET Stock = Stock + OLD.Quantity WHERE ID = OLD.Medicine_ID;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_dispensing_after_delete
AFTER DELETE ON Dispensing
FOR EACH ROW
EXECUTE FUNCTION fn_dispensing_after_delete();

--TRIGGER 4 — trg_appointment_management
--Prevents doctors from being double-booked at the same date and time

CREATE OR REPLACE FUNCTION fn_appointment_management()
RETURNS TRIGGER AS $$
DECLARE
    v_count         INT;
    v_doctor_name   VARCHAR(100);
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM   Appointment
    WHERE  Doctor_ID = NEW.Doctor_ID
      AND  Date      = NEW.Date
      AND  Time      = NEW.Time
      AND  Status    = 'pending';

    IF v_count > 0 THEN
        SELECT Name INTO v_doctor_name FROM Doctor WHERE ID = NEW.Doctor_ID;
        RAISE EXCEPTION
            '[Appointment Management] Dr. "%" (ID: %) already has an appointment '
            'on % at %. Choose a different time slot.',
            v_doctor_name, NEW.Doctor_ID, NEW.Date, NEW.Time;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_appointment_management
BEFORE INSERT ON Appointment
FOR EACH ROW
EXECUTE FUNCTION fn_appointment_management();

--TRIGGER 5 — trg_er_shift_management
--Prevents scheduling ER shifts in the past and limits doctor assignments to maximum 5 per shift

CREATE OR REPLACE FUNCTION fn_er_shift_before()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.Date < CURRENT_DATE THEN
        RAISE EXCEPTION
            '[ER Shift Management] Cannot schedule a shift in the past. '
            'Requested: %, Today: %.',
            NEW.Date, CURRENT_DATE;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_er_shift_before
BEFORE INSERT ON ER_Shift
FOR EACH ROW
EXECUTE FUNCTION fn_er_shift_before();


CREATE OR REPLACE FUNCTION fn_er_doctor_before()
RETURNS TRIGGER AS $$
DECLARE
    v_count         INT;
    v_shift_date    DATE;
    v_shift_time    TIME;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM   ER_Shift_Doctor
    WHERE  ER_ID = NEW.ER_ID;

    IF v_count >= 5 THEN
        SELECT Date, Time INTO v_shift_date, v_shift_time
        FROM   ER_Shift WHERE ID = NEW.ER_ID;

        RAISE EXCEPTION
            '[ER Shift Management] Shift ID % (% at %) already has 5 doctors. '
            'Maximum capacity reached.',
            NEW.ER_ID, v_shift_date, v_shift_time;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_er_doctor_before
BEFORE INSERT ON ER_Shift_Doctor
FOR EACH ROW
EXECUTE FUNCTION fn_er_doctor_before();

--VIEWS

--VIEW 1 — Active_Blood_Inventory
--Displays available and reserved blood units with expiry countdown and donor 
--information for prioritizing usage and preventing expired blood administration

CREATE VIEW Active_Blood_Inventory AS
SELECT
    b.ID                                        AS Blood_ID,
    b.B_Gr                                      AS Blood_Group,
    b.Unit                                      AS Units_Available,
    b.Collected_Date,
    b.Expiry_Date,
    (b.Expiry_Date - CURRENT_DATE)              AS Days_Until_Expiry,
    b.Status,
    CASE
        WHEN (b.Expiry_Date - CURRENT_DATE) <= 7  THEN 'CRITICAL'
        WHEN (b.Expiry_Date - CURRENT_DATE) <= 14 THEN 'WARNING'
        ELSE                                           'OK'
    END                                         AS Expiry_Alert,
    d.ID                                        AS Donor_ID,
    d.Name                                      AS Donor_Name,
    d.B_Gr                                      AS Donor_Blood_Group,
    d.P_No                                      AS Donor_Phone
FROM       Blood b
LEFT JOIN  Donor d ON d.ID = b.Donor_ID
WHERE  b.Status      IN ('available', 'reserved')
  AND  b.Expiry_Date  > CURRENT_DATE
  AND  b.Unit         > 0
ORDER BY
    b.Expiry_Date ASC,
    b.B_Gr        ASC;

--VIEW 2 — Donor_Donation_History
--Shows donor profiles with total donations, last donation date, and current 
--eligibility status based on 90-day rule for outreach and screening purposes

CREATE VIEW Donor_Donation_History AS
SELECT
    d.ID                                                        AS Donor_ID,
    d.Name                                                      AS Donor_Name,
    d.B_Gr                                                      AS Blood_Group,
    d.P_No                                                      AS Phone,
    d.Email,
    d.Last_Donate                                               AS Last_Donation_Date,
    (CURRENT_DATE - d.Last_Donate)                              AS Days_Since_Last_Donation,
    COUNT(don.ID)
        FILTER (WHERE don.Status = 'completed')                 AS Total_Completed_Donations,
    COUNT(don.ID)
        FILTER (WHERE don.Status = 'cancelled')                 AS Total_Cancelled_Donations,
    COUNT(don.ID)
        FILTER (WHERE don.Status = 'deferred')                  AS Total_Deferred_Donations,
    CASE
        WHEN (CURRENT_DATE - d.Last_Donate) >= 90 THEN 'ELIGIBLE'
        ELSE 'NOT ELIGIBLE — Wait until ' ||
             TO_CHAR(d.Last_Donate + INTERVAL '90 days', 'DD Mon YYYY')
    END                                                         AS Eligibility_Status
FROM       Donor     d
LEFT JOIN  Donation  don ON don.Donor_ID = d.ID
GROUP BY
    d.ID, d.Name, d.B_Gr, d.P_No,
    d.Email, d.Last_Donate
ORDER BY
    d.Last_Donate ASC;

--VIEW 3 — ER_Shift_Staffing
--Displays ER shift details with assigned receptionist, doctor count, and doctor 
--names to monitor staffing levels and identify coverage gaps

CREATE VIEW ER_Shift_Staffing AS
SELECT
    ers.ID                                      AS Shift_ID,
    ers.Date                                    AS Shift_Date,
    ers.Time                                    AS Shift_Time,
    r.ID                                        AS Receptionist_ID,
    r.Name                                      AS Receptionist_Name,
    r.P_No                                      AS Receptionist_Phone,
    COUNT(esd.Doctor_ID)                        AS Doctors_Assigned,
    (5 - COUNT(esd.Doctor_ID))                  AS Slots_Remaining,
    CASE
        WHEN COUNT(esd.Doctor_ID) = 0 THEN 'UNSTAFFED'
        WHEN COUNT(esd.Doctor_ID) < 3 THEN 'UNDERSTAFFED'
        WHEN COUNT(esd.Doctor_ID) = 5 THEN 'FULL'
        ELSE                               'ADEQUATE'
    END                                         AS Staffing_Status,
    STRING_AGG(doc.Name, ', ' ORDER BY doc.Name) AS Doctor_Names
FROM            ER_Shift        ers
LEFT JOIN       Receptionist    r   ON  r.ID  = ers.Receptionist_ID
LEFT JOIN       ER_Shift_Doctor esd ON esd.ER_ID    = ers.ID
LEFT JOIN       Doctor          doc ON doc.ID = esd.Doctor_ID
GROUP BY
    ers.ID, ers.Date, ers.Time,
    r.ID,   r.Name,   r.P_No
ORDER BY
    ers.Date ASC,
    ers.Time ASC;

--VIEW 4 — Medicine_Expiry_Alert
--Shows medicines expiring within 30 days or out of stock with alert types and days 
--until expiry for restocking and disposal decisions

CREATE VIEW Medicine_Expiry_Alert AS
SELECT
    m.ID,
    m.Name                                      AS Medicine_Name,
    m.Batch_No,
    m.Stock,
    m.Expiry_Date,
    (m.Expiry_Date - CURRENT_DATE)              AS Days_Until_Expiry,
    CASE
        WHEN m.Stock = 0                                THEN 'OUT OF STOCK'
        WHEN m.Expiry_Date <= CURRENT_DATE              THEN 'EXPIRED'
        WHEN m.Expiry_Date <= CURRENT_DATE + 7          THEN 'CRITICAL — Expires in 7 days'
        WHEN m.Expiry_Date <= CURRENT_DATE + 30         THEN 'EXPIRING SOON'
    END                                         AS Alert_Type
FROM  Medicine m
WHERE m.Stock = 0
   OR m.Expiry_Date <= CURRENT_DATE + 30
ORDER BY
    m.Expiry_Date ASC,
    m.Stock       ASC;

--INDEX 1 — idx_er_doctor_availability
--Table: ER_Shift_Doctor (Doctor_ID, ER_ID)
--Enables instant doctor availability checks during patient triage by avoiding full 
--table scans, ensuring immediate patient assignments and preventing emergency response delays

CREATE INDEX idx_er_doctor_availability
    ON ER_Shift_Doctor(Doctor_ID, ER_ID);

--INDEX 2 — idx_blood_type_status
--Table: Blood (B_Gr, Status)
--Provides instant blood availability lookups by blood type and status during life-critical 
--emergencies where every second counts, directly impacting patient safety

CREATE INDEX idx_blood_type_status
    ON Blood(B_Gr, Status);


--INDEX 3 — idx_blood_request_status
--Table: Blood_Request (Status)
--Eliminates full table scans when filtering blood requests by status, preventing
--performance degradation during high-demand periods in dashboards and fulfillment workflows

CREATE INDEX idx_blood_request_status
    ON Blood_Request(Status);

--INDEX 4 — idx_medicine_name
--Table: Medicine (Name)
--Enables fast medicine name searches during prescriptions, dispensing, 
--and stock checks, preventing delays that could affect patient medication timing as inventory grows

CREATE INDEX idx_medicine_name
    ON Medicine(Name);

