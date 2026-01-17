CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    status TEXT NOT NULL,
    priority TEXT NOT NULL
);

-- Insert initial data
INSERT INTO tasks (title, description, status, priority) VALUES
    ('Task One', 'First task description', 'Pending', 'High'),
    ('Task Two', 'Second task description', 'In Progress', 'Medium'),
    ('Task Three', 'Third task description', 'Completed', 'Low');
