
ALTER TABLE books
  ADD COLUMN stock INT NOT NULL DEFAULT 1 AFTER description;

UPDATE books
SET reservation_status = CASE WHEN stock > 0 THEN 'available' ELSE 'reserved' END;

