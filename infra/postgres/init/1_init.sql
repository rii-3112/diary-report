-- ユーザーテーブル
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 公開設定テーブル（どの項目を他人に見せるか）
CREATE TABLE user_public_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    show_title BOOLEAN DEFAULT true,
    show_content BOOLEAN DEFAULT false,
    show_learning_notes BOOLEAN DEFAULT false,
    show_habit_status BOOLEAN DEFAULT false
);

-- 日報テーブル
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    learning_notes TEXT,
    is_habit_done BOOLEAN DEFAULT false,
    is_public BOOLEAN DEFAULT false,
    public_token VARCHAR(255) UNIQUE NOT NULL, -- 共有URL用のトークン
    submitted_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);