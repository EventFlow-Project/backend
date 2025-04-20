CREATE TABLE IF NOT EXISTS friend_requests (
    id VARCHAR(255) PRIMARY KEY,
    from_id VARCHAR(255) NOT NULL,
    to_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_friend_requests_from_user FOREIGN KEY (from_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_friend_requests_to_user FOREIGN KEY (to_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_friend_request_status CHECK (status IN ('pending', 'accepted', 'rejected'))
);

CREATE INDEX idx_friend_requests_from_id ON friend_requests(from_id);
CREATE INDEX idx_friend_requests_to_id ON friend_requests(to_id);
CREATE INDEX idx_friend_requests_status ON friend_requests(status);
CREATE INDEX idx_friend_requests_deleted_at ON friend_requests(deleted_at); 