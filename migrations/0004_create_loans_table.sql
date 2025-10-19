
CREATE TABLE IF NOT EXISTS loans (
  id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id INT UNSIGNED NOT NULL,
  book_id INT UNSIGNED NOT NULL,
  status ENUM('reserved','borrowed','returned','cancelled') NOT NULL DEFAULT 'reserved',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,

  reserved_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  borrowed_at DATETIME NULL,
  due_date    DATETIME NULL,
  returned_at DATETIME NULL,
  cancelled_at DATETIME NULL,
  notes VARCHAR(255) NULL,

  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  CONSTRAINT fk_loans_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_loans_book FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE INDEX idx_loans_user    ON loans (user_id);
CREATE INDEX idx_loans_book    ON loans (book_id);
CREATE INDEX idx_loans_status  ON loans (status);
CREATE INDEX idx_loans_active  ON loans (is_active);
