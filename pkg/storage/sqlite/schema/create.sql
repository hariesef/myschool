--sqlite does not support BIGINT or SERIAL
CREATE TABLE "students" (
	uid INTEGER PRIMARY KEY AUTOINCREMENT , --sqlite compatibility mode
    created_at INT NOT NULL,
    updated_at INT NOT NULL,
    deleted_at INT NOT NULL,
    name TEXT NOT NULL,
    Gender TEXT NOT NULL
);
CREATE INDEX `idx_students_deleted_at` ON `students`(`deleted_at`);

